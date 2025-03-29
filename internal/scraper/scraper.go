package scraper

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/5pirit5eal/swim-rag/internal/rag"
	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type URLMap struct {
	mux sync.Mutex
	m   map[string]bool
}

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

func Scrape(alreadyVisited []string, ctx context.Context, client llms.Model, c chan schema.Document, ec chan error, seeds ...string) {
	defer close(c)
	defer close(ec)
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
		log.Println("Error occurred while scraping: ", err)
	})

	// Create a map to track visited URLs and models.plans
	visitedURLs := NewURLMap(alreadyVisited)

	scraper.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Check if the link is a valid URL
		if e.Attr("href") == "" || strings.HasPrefix(e.Attr("href"), "#") {
			log.Println("Invalid link, skipping")
			return
		}
		// Check if the link is an internal link
		if e.Attr("href")[0] == '/' {
			log.Println("Internal link, skipping")
			return
		}
		// Extract the href attribute
		href := e.Attr("href")
		// Check if the href has been visited
		if found := visitedURLs.Load(href); !found {
			log.Println("Found new link:", e.Attr("href"))
			// Mark the href as visited
			visitedURLs.Store(href)
			// Visit the link
			err := e.Request.Visit(href)
			if err != nil {
				if !strings.Contains(err.Error(), "Max depth limit reached") {
					log.Println("Error visiting link:", err)
				}
			}
		}
	})
	scraper.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.ChildText("h1")
		desc := e.ChildText("div.cm-posts > article.post h3")

		if table := e.ChildText("table"); table == "" {
			log.Println("No table found, skipping")
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
					// else {
					// 	// log.Println("No previous row to append content to")
					// }
				}
				// else {
				// 	log.Println("Skipping empty row")
				// }
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
			// log.Println("Added row to table:", row.String())
			table = append(table, row)
		})

		// log.Println("Found description:", desc)
		// log.Println("Found table with", len(table), "rows")
		tables++
		log.Println("Found table nr. ", tables)

		if len(table) > 0 {
			table.AddSum()
		}

		url := e.Request.URL.String()

		visitedURLs.Store(url)

		go ImprovePlan(ctx, client, models.Plan{
			URL:         url,
			Title:       title,
			Description: desc,
			Table:       table,
		}, c, ec)
	})

	scraper.OnScraped(func(r *colly.Response) {
		log.Println("Scraping finished for:", r.Request.URL)
	})

	// Visit each seed URL and scrape the data
	for _, url := range seeds {
		// Mark the seed as visited
		visitedURLs.Store(url)

		err := scraper.Visit(url)
		if err != nil {
			log.Println("Error visiting seed URL:", err)
			ec <- err
		}
	}
	scraper.Wait()
}

func ScrapeURL(db *rag.RAGDB, ctx context.Context, client llms.Model, cfg models.Config, url string) error {
	// Load urls in the database into the scraper
	alreadyVisited := make([]string, 0)

	rows, err := db.Conn.Query(ctx, fmt.Sprintf(`
		SELECT url FROM urls
		WHERE (created_at > now() - interval '60 days'
		AND collection_id = (
			SELECT uuid FROM %s WHERE name = $1
			ORDER BY name limit 1
		))`, rag.CollectionTableName), cfg.Embedding.Model)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}
	log.Println("Queried database successfully")
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
	go Scrape(alreadyVisited, ctx, client, dc, ec, url)

	documents := make([]schema.Document, 0)
	errors := make([]error, 0)

Forloop:
	for {
		select {
		case doc, dcOpen := <-dc:
			if dcOpen {
				documents = append(documents, doc)
			} else {
				log.Println("Channel closed, stopping scraping")
				break Forloop
			}
		case err, ecOpen := <-ec:
			if ecOpen {
				errors = append(errors, err)
			} else {
				log.Println("Channel closed, stopping scraping")
				break Forloop
			}
		case <-ctx.Done():
			log.Println("Context done, stopping scraping")
			break Forloop
		case <-time.After(3600 * time.Second):
			errors = append(errors, fmt.Errorf("timeout while scraping"))
			break Forloop
		}
	}
	log.Println("Scraping finished, received", len(documents), "documents and", len(errors), "errors")
	if len(errors) > 0 {
		log.Println("Failed to improve plans or errors during scraping:", errors)
		// return fmt.Errorf("failed to improve plans or errors during scraping: %v", errors)
	}
	log.Println("Adding documents to the database...")

	// Store documents and their embeddings in the database
	ids, err := db.Store.AddDocuments(ctx, documents)
	if err != nil {
		log.Println("Failed to add documents to the database:", err.Error())
		return fmt.Errorf("failed to add documents to the database: %w", err)
	}
	log.Println("Added documents to the database successfully")

	var collectionUUID string
	err = db.Conn.QueryRow(ctx, `SELECT uuid FROM documents WHERE name = $1 ORDER BY name limit 1`, cfg.Embedding.Model).Scan(&collectionUUID)
	if err != nil {
		log.Println("Failed to query collection UUID:", err.Error())
		return fmt.Errorf("failed to query collection UUID: %w", err)
	}

	batch := &pgx.Batch{}
	for i := range documents {
		log.Println("Inserting URL into database:", documents[i].Metadata["url"])
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
			DO UPDATE SET document_id = EXCLUDED.document_id`, cfg.Embedding.Name),
			documents[i].Metadata["url"], ids[i], collectionUUID)
	}
	log.Println("Inserting URLs into database...")
	// Execute the batch
	br := db.Conn.SendBatch(ctx, batch)

	if err := br.Close(); err != nil {
		return err
	}
	log.Printf("Inserted %d URLs into database successfully", len(documents))
	return nil
}
