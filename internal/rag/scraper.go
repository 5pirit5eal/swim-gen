package rag

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5"
	"github.com/tmc/langchaingo/schema"
)

// URLMap is a thread-safe map to store URLs that have already been visited
// and to prevent duplicate scraping of the same URL.
type URLMap struct {
	mux sync.Mutex
	m   map[string]bool
}

// NewURLMap initializes a new URLMap with the given already visited URLs.
func NewURLMap(alreadyVisited []string) *URLMap {
	m := make(map[string]bool)
	for _, url := range alreadyVisited {
		m[url] = true
	}
	return &URLMap{
		mux: sync.Mutex{},
		m:   m,
	}
}

func (um *URLMap) Store(key string) {
	um.mux.Lock()
	defer um.mux.Unlock()
	um.m[key] = true
}

func (um *URLMap) Load(key string) bool {
	um.mux.Lock()
	defer um.mux.Unlock()
	_, found := um.m[key]
	return found
}

func (um *URLMap) Len() int {
	return len(um.m)
}

func (db *RAGDB) NewCollector(ctx context.Context, visitedURLs *URLMap, syncGroup *sync.WaitGroup, c chan schema.Document, ec chan error) *colly.Collector {
	logger := httplog.LogEntry(ctx)
	tables := 0

	// Create a new scraper
	scraper := colly.NewCollector(
		colly.AllowedDomains("docswim.de"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.MaxDepth(3),
		colly.Async(true),
	)
	scraper.Limit(&colly.LimitRule{
		DomainGlob:  "*docswim.de*",
		Parallelism: 10,
		// Delay:       2 * time.Second,
	})
	scraper.OnError(func(_ *colly.Response, err error) {
		// Handle errors during scraping
		logger.Error("Error occurred while scraping", httplog.ErrAttr(err))
		ec <- err
	})

	scraper.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Check if the link is a valid URL
		if e.Attr("href") == "" || strings.HasPrefix(e.Attr("href"), "#") {
			logger.Debug("Invalid link, skipping")
			return
		}
		// Check if the link is an internal link
		if e.Attr("href")[0] == '/' {
			logger.Debug("Internal link, skipping")
			return
		}
		// Extract the href attribute
		href := e.Attr("href")
		// Check if the href has been visited
		if found := visitedURLs.Load(href); !found {
			logger.Debug("Found new link", "new", e.Attr("href"))
			// Mark the href as visited
			visitedURLs.Store(href)
			// Visit the link
			err := e.Request.Visit(href)
			if err != nil {
				if !strings.Contains(err.Error(), "Max depth limit reached") && !strings.Contains(err.Error(), "Forbidden domain") {
					logger.Error("Error visiting link", httplog.ErrAttr(err))
				}
			}
		}
	})
	scraper.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.ChildText("h1")
		desc := e.ChildText("div.cm-posts > article.post h3")

		if table := e.ChildText("table"); table == "" {
			logger.Debug("No table found, skipping")
			return
		}

		e.ForEach("div.cm-posts > article.post p:not(:has(span), :has(iframe))", func(_ int, p *colly.HTMLElement) {
			if p.Text != "" {
				desc = desc + "\n" + p.Text
			}
		})

		table := make(models.Table, 0)
		// Extract the table content
		e.ForEach("div.cm-posts > article.post tr:not(:has(strong), [colspan], :has(span))", func(_ int, r *colly.HTMLElement) {
			empty := 0
			amount, err := strconv.Atoi(r.ChildText("td:nth-child(1)"))
			if err != nil {
				amount = 0
				empty++
			}
			distance, err := strconv.Atoi(r.ChildText("td:nth-child(3)"))
			if err != nil {
				distance = 0
				empty++
			}
			sum, err := strconv.Atoi(r.ChildText("td:nth-child(7)"))
			if err != nil {
				sum = 0
				empty++
			}

			if empty >= 3 {
				if r.ChildText("td:nth-child(5)") != "" {
					// Append overreaching content to the last row
					if len(table) > 0 {
						table[len(table)-1].Content += " " + r.ChildText("td:nth-child(5)")
					}
				}
				return
			}

			row := models.Row{
				Amount:     amount,
				Multiplier: r.ChildText("td:nth-child(2)"),
				Distance:   distance,
				Break:      r.ChildText("td:nth-child(4)"),
				Content:    r.ChildText("td:nth-child(5)"),
				Intensity:  r.ChildText("td:nth-child(6)"),
				Sum:        sum,
			}

			// Append the row to the table
			table = append(table, row)
		})

		tables++
		logger.Info(fmt.Sprintf("Found table nr. %d", tables), "length", len(table), "title", title)
		logger.Debug("Found description", "description", desc)

		if len(table) > 0 {
			table.AddSum()
		}

		url := e.Request.URL.String()

		visitedURLs.Store(url)

		syncGroup.Add(1)
		go db.Client.ImprovePlan(ctx, models.Plan{
			URL:         url,
			Title:       title,
			Description: desc,
			Table:       table,
		}, syncGroup, c, ec)
	})

	scraper.OnScraped(func(r *colly.Response) {
		logger.Debug("Scraping finished", "sub_url", r.Request.URL)
	})
	return scraper
}

// Scrapes the given URLs and extracts the relevant data from the HTML content.
func (db *RAGDB) scrape(ctx context.Context, alreadyVisited []string, c chan schema.Document, ec chan error, url string) {
	logger := httplog.LogEntry(ctx)
	syncGroup := &sync.WaitGroup{}
	defer close(c)
	defer close(ec)

	// Mark the seed as visited
	// Create a map to track visited URLs and models.plans
	visitedURLs := NewURLMap(alreadyVisited)
	visitedURLs.Store(url)
	collector := db.NewCollector(ctx, visitedURLs, syncGroup, c, ec)

	err := collector.Visit(url)
	if err != nil {
		logger.Error("Error visiting seed URL", httplog.ErrAttr(err))
		ec <- err
	}

	collector.Wait()
	syncGroup.Wait()
}

func (db *RAGDB) ScrapeURL(ctx context.Context, url string) error {
	logger := httplog.LogEntry(ctx)
	logger.Info("Starting to scrape")
	// Load urls in the database into the scraper
	alreadyVisited := make([]string, 0)

	rows, err := db.Conn.Query(ctx, fmt.Sprintf(`
		SELECT url FROM urls
		WHERE (created_at > now() - interval '60 days'
		AND collection_id = (
			SELECT uuid FROM %s WHERE name = $1
			ORDER BY name limit 1
		))`, CollectionTableName), db.cfg.Embedding.Model)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}
	logger.Info("Queried database successfully")
	defer rows.Close()
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		alreadyVisited = append(alreadyVisited, url)
	}

	// Scrape the URL
	dc := make(chan schema.Document)
	ec := make(chan error)
	go db.scrape(ctx, alreadyVisited, dc, ec, url)

	documents := make([]schema.Document, 0)
	errors := make([]error, 0)

Forloop:
	for {
		select {
		case doc, dcOpen := <-dc:
			if !dcOpen {
				logger.Info("Channel closed, stopping scraping")
				break Forloop
			}
			documents = append(documents, doc)
		case err, ecOpen := <-ec:
			if !ecOpen {
				logger.Info("Channel closed, stopping scraping")
				break Forloop
			}
			errors = append(errors, err)
		case <-ctx.Done():
			logger.Info("Context done, stopping scraping")
			break Forloop
		case <-time.After(3600 * time.Second):
			errors = append(errors, fmt.Errorf("timeout while scraping"))
			break Forloop
		}
	}
	logger.Info("Scraping finished", "documents", len(documents), "errors", len(errors))
	if len(errors) > 0 {
		logger.Error("Failed to improve plans or errors during scraping", "err", slog.AnyValue(errors))
		// return fmt.Errorf("failed to improve plans or errors during scraping: %v", errors)
	}
	logger.Info("Adding documents to the database")

	// Store documents and their embeddings in the database
	ids, err := db.Store.AddDocuments(ctx, documents)
	if err != nil {
		logger.Error("Failed to add documents to the database", httplog.ErrAttr(err))
		return fmt.Errorf("failed to add documents to the database: %w", err)
	}
	logger.Info("Added documents to the database successfully")

	var collectionUUID string
	err = db.Conn.QueryRow(ctx, `SELECT uuid FROM documents WHERE name = $1 ORDER BY name limit 1`, db.cfg.Embedding.Model).Scan(&collectionUUID)
	if err != nil {
		logger.Error("Failed to query collection UUID", httplog.ErrAttr(err))
		return fmt.Errorf("failed to query collection UUID: %w", err)
	}

	batch := &pgx.Batch{}
	for i := range documents {
		logger.Debug(fmt.Sprintf("Inserting URL %s into database", documents[i].Metadata["url"]), "url", documents[i].Metadata["url"])
		batch.Queue(fmt.Sprintf(`
			WITH deleted AS (
				DELETE FROM %s 
				WHERE uuid = (
					SELECT document_id 
					FROM urls 
					WHERE url = $1 AND collection_id = $3
				)
				RETURNING uuid
			)
			INSERT INTO urls (url, document_id, collection_id) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (url, collection_id) 
			DO UPDATE SET document_id = EXCLUDED.document_id`, db.cfg.Embedding.Name),
			documents[i].Metadata["url"], ids[i], collectionUUID)
	}
	logger.Info("Inserting URLs into database...")
	// Execute the batch
	br := db.Conn.SendBatch(ctx, batch)

	if err := br.Close(); err != nil {
		return err
	}
	logger.Info("Inserted URLs into database successfully", slog.Int("n_urls", len(documents)))
	return nil
}
