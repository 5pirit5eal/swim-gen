package rag

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
	"github.com/gocolly/colly"
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

func (db *RAGDB) NewCollector(ctx context.Context, visitedURLs *URLMap, syncGroup *sync.WaitGroup, c chan<- models.Document, ec chan<- error) *colly.Collector {
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
		go db.Client.ImprovePlan(ctx, &models.ScrapedPlan{
			PlanID:      GenerateUUID(url),
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
func (db *RAGDB) startScraping(ctx context.Context, alreadyVisited []string, c chan models.Document, ec chan error, url string) {
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
	// Set the embedder to document embedding mode
	db.Client.DocumentMode()
	logger.Info("Starting to scrape")
	// Load urls in the database into the scraper
	alreadyVisited, err := db.GetAlreadyVisitedURLs(ctx)
	if err != nil {
		return fmt.Errorf("db.GetAlreadyVisitedURLs: %w", err)
	}
	logger.Info("Queried database successfully")

	// Scrape the URL
	dc := make(chan models.Document)
	ec := make(chan error)
	go db.startScraping(ctx, alreadyVisited, dc, ec, url)

	documents := make([]models.Document, 0)
	errors := make([]error, 0)

Forloop:
	for {
		select {
		case doc, dcOpen := <-dc:
			if !dcOpen {
				logger.Info("Document channel closed, stopping scraping")
				break Forloop
			}
			// Filter out empty documents
			if len(doc.Plan.Plan().Table) != 0 {
				documents = append(documents, doc)
			}
		case err, ecOpen := <-ec:
			if !ecOpen {
				logger.Info("Error channel closed, stopping scraping")
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
		logger.Warn("Failed to improve plans or errors during scraping", "err", slog.AnyValue(errors))
	}
	logger.Info("Adding documents to the database")

	// Store documents and their embeddings in the database
	langchainDocs := make([]schema.Document, len(documents))
	for i := range documents {
		// Convert the models.Document to a schema.Document
		doc, err := models.PlanToDoc(&documents[i])
		if err != nil {
			logger.Error("Failed to convert plan to document", httplog.ErrAttr(err))
			return fmt.Errorf("PlanToDoc: %w", err)
		}
		// Add the document to the langchainDocs slice
		langchainDocs[i] = doc
	}
	_, err = db.Store.AddDocuments(ctx, langchainDocs)
	if err != nil {
		logger.Error("Failed to add documents to the database", httplog.ErrAttr(err))
		return fmt.Errorf("Store.AddDocuments: %w", err)
	}
	logger.Info("Added documents to the database successfully")

	var collectionUUID string
	err = pgxscan.Get(ctx, db.Conn, &collectionUUID,
		fmt.Sprintf(`SELECT uuid FROM %s WHERE name = $1 ORDER BY name limit 1`, CollectionTableName), db.cfg.Embedding.Model)
	if err != nil {
		logger.Error("Failed to query collection UUID", httplog.ErrAttr(err))
		return fmt.Errorf("failed to query collection UUID: %w", err)
	}

	if err := db.AddScrapedPlans(ctx, collectionUUID, documents); err != nil {
		logger.Error("Failed to add scraped plans to the database", httplog.ErrAttr(err))
		return fmt.Errorf("failed to add scraped plans to the database: %w", err)
	}
	logger.Info("Added scraped plans to the database successfully")
	logger.Info("Scraping and adding data finished successfully", "documents", len(documents), "errors", len(errors))

	return nil
}

func (db *RAGDB) GetAlreadyVisitedURLs(ctx context.Context) ([]string, error) {
	logger := httplog.LogEntry(ctx)
	// Query the database for already visited URLs
	var alreadyVisited []string
	err := pgxscan.Select(ctx, db.Conn, &alreadyVisited, fmt.Sprintf(`
		SELECT url FROM %s
		WHERE (created_at > now() - interval '60 days'
		AND collection_id = (
			SELECT uuid FROM %s WHERE name = $1
			ORDER BY name limit 1
		))`, ScrapedTableName, CollectionTableName), db.cfg.Embedding.Model)
	if err != nil {
		logger.Error("Error querying database", httplog.ErrAttr(err))
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	return alreadyVisited, nil
}

func (db *RAGDB) AddScrapedPlans(ctx context.Context, collectionUUID string, documents []models.Document) error {
	logger := httplog.LogEntry(ctx)
	logger.Info("Adding scraped plans and visited urls to the database")

	// Start transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Failed to begin transaction", httplog.ErrAttr(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer rollback in case of an error
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			logger.Error("Failed to rollback transaction", httplog.ErrAttr(err))
		}
		logger.Info("Transaction rolled back")
	}()
	// Add the scraped urls to the database
	for _, document := range documents {
		plan := document.Plan.Plan()
		planMap := document.Plan.Map()
		// Check if the URL is already in the database
		if _, found := planMap["url"]; !found {
			logger.Warn("No URL found in the plan, skipping")
			continue
		}
		url := planMap["url"].(string)
		if plan.PlanID == "" {
			logger.Warn("PlanID is empty, skipping insertion", "url", url)
			continue
		}

		pseudoTx, err := tx.Begin(ctx)
		if err != nil {
			logger.Error("Failed to begin transaction", httplog.ErrAttr(err))
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		logger.Debug(fmt.Sprintf("Inserting URL %s into database", url), "url", url)
		// If a URL already exists, delete the old entry and insert the new one
		_, err = pseudoTx.Exec(ctx, fmt.Sprintf(`
			WITH deleted AS (
				DELETE FROM %s 
				WHERE uuid = (
					SELECT plan_id 
					FROM %s 
					WHERE url = $1 AND collection_id = $3
				)
				RETURNING uuid
			)
			INSERT INTO %s (url, plan_id, collection_id) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (url, collection_id) 
			DO UPDATE SET plan_id = EXCLUDED.plan_id`, db.cfg.Embedding.Name, ScrapedTableName, ScrapedTableName),
			url, plan.PlanID, collectionUUID)

		if err != nil {
			logger.Error("Failed to insert scraped plan into database", httplog.ErrAttr(err))
			return fmt.Errorf("failed to insert scraped plan into database: %w", err)
		}
		// Add the plan to the plan table
		_, err = pseudoTx.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (plan_id, title, description, plan_table)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (plan_id) DO NOTHING`, PlanTableName),
			plan.PlanID, plan.Title, plan.Description, plan.Table)
		if err != nil {
			logger.Error("Failed to insert plan into database", httplog.ErrAttr(err))
			return fmt.Errorf("failed to insert plan into database: %w", err)
		}

		// Commit the transaction
		if err := pseudoTx.Commit(ctx); err != nil {
			logger.Error("Failed to commit transaction", httplog.ErrAttr(err))
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		logger.Debug("Inserted URL into database successfully", "url", url)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit transaction", httplog.ErrAttr(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Inserted scraped URLs into database successfully", slog.Int("n_urls", len(documents)))

	return nil
}

func (db *RAGDB) GetScrapedPlan(ctx context.Context, planID string) (*models.ScrapedPlan, error) {
	logger := httplog.LogEntry(ctx)
	var plan models.ScrapedPlan
	err := pgxscan.Get(ctx, db.Conn, &plan,
		fmt.Sprintf(`
			SELECT sp.url, sp.plan_id, sp.created_at, p.title, p.description, p.plan_table
			FROM %s sp
			JOIN %s p ON sp.plan_id = p.plan_id
			WHERE sp.plan_id = $1`, ScrapedTableName, PlanTableName), planID)
	if err != nil {
		logger.Error("Error querying scraped plan", httplog.ErrAttr(err))
		return nil, err
	}
	return &plan, nil
}
