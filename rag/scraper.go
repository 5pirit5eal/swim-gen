package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type Plan struct {
	URL, Title, Description string
	Table                   Table
}

func (p *Plan) String() string {
	return fmt.Sprintf("%s:\n %s\n %s", p.Title, p.Description, p.Table.String())
}

func (p *Plan) Map() map[string]any {
	m := make(map[string]any)
	m["url"] = p.URL
	m["title"] = p.Title
	m["description"] = p.Description
	m["table"] = p.Table
	return m
}

type Table []Row

type Row struct {
	Amount     int
	Multiplier string
	Distance   int
	Break      string
	Content    string
	Intensity  string
	Sum        int
}

func (r Row) String() string {
	return fmt.Sprintf("| %d | %s | %d | %s | %s | %s | %d |", r.Amount, r.Multiplier, r.Distance, r.Break, r.Content, r.Intensity, r.Sum)
}

func (t *Table) String() string {
	tstr := "Anzahl |  | Strecke(m) | Pause(s) | Inhalt | Intensität | Umfang |\n"
	tstr += "|---|---|---|---|---|---|---|\n"
	for _, row := range *t {
		tstr += row.String() + "\n"
	}
	return tstr
}

// Adds a final row to the table with the total sum
func (t *Table) AddSum() {
	sum := 0
	for _, row := range *t {
		sum += row.Sum
	}
	*t = append(*t, Row{Content: "Gesamt", Sum: sum})
}

type URLMap struct {
	mux sync.Mutex
	m   map[string]Plan
}

func NewURLMap(alreadyVisited []string) *URLMap {
	m := make(map[string]Plan)
	for _, url := range alreadyVisited {
		m[url] = Plan{}
	}
	return &URLMap{
		mux: sync.Mutex{},
		m:   m,
	}
}

func (um *URLMap) Store(key string, value Plan) {
	um.mux.Lock()
	defer um.mux.Unlock()
	um.m[key] = value
}

func (um *URLMap) Load(key string) (Plan, bool) {
	um.mux.Lock()
	defer um.mux.Unlock()
	value, found := um.m[key]
	return value, found
}

func (um *URLMap) Len() int {
	return len(um.m)
}

type KeyValuePair struct {
	URL  string
	Plan Plan
}

func (um *URLMap) Range() <-chan KeyValuePair {
	ch := make(chan KeyValuePair)
	go func() {
		defer close(ch)
		for k, v := range um.m {
			um.mux.Lock()
			kvp := KeyValuePair{URL: k, Plan: v}
			um.mux.Unlock()
			ch <- kvp
		}
	}()
	return ch
}

var TABLE_HEADER = []string{
	"Anzahl",
	"Multiplikator",
	"Strecke(m)",
	"Pause(s)",
	"Inhalt",
	"Intensität",
	"Umfang",
}

func Scrape(alreadyVisited []string, seeds ...string) (*URLMap, error) {

	// Create a new scraper
	scraper := colly.NewCollector(
		colly.AllowedDomains("docswim.de"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.MaxDepth(1),
	)
	scraper.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})
	scraper.OnError(func(_ *colly.Response, err error) {
		// Handle errors during scraping
		log.Println("Error occurred while scraping: ", err)
	})
	scraper.OnResponse(func(r *colly.Response) {
		log.Println("Visited:", r.Request.URL.String())
	})

	// Create a map to track visited URLs and plans
	visitedURLs := NewURLMap(alreadyVisited)

	scraper.OnHTML("a[href]", func(e *colly.HTMLElement) {
		log.Println("Found link:", e.Attr("href"))
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
		if _, found := visitedURLs.Load(href); !found {
			// Mark the href as visited
			visitedURLs.Store(href, Plan{})
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

		table := make(Table, 0)
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
					} else {
						log.Println("No previous row to append content to")
					}
					return
				} else {
					log.Println("Skipping empty row")
					return
				}
			}

			row := Row{
				Amount:     amount,
				Multiplier: r.ChildText("td:nth-child(2)"),
				Distance:   distance,
				Break:      r.ChildText("td:nth-child(4)"),
				Content:    r.ChildText("td:nth-child(5)"),
				Intensity:  r.ChildText("td:nth-child(6)"),
				Sum:        sum,
			}

			// Append the row to the table
			log.Println("Added row to table:", row.String())
			table = append(table, row)

		})

		log.Println("Found description:", desc)
		log.Println("Found table with", len(table), "rows")

		if len(table) > 0 {
			table.AddSum()
		}

		url := e.Request.URL.String()

		visitedURLs.Store(url, Plan{
			URL:         url,
			Title:       title,
			Description: desc,
			Table:       table,
		})
	})

	scraper.OnScraped(func(r *colly.Response) {
		log.Println("Scraping finished for:", r.Request.URL)
	})

	// Visit each seed URL and scrape the data
	for _, url := range seeds {

		err := scraper.Visit(url)
		if err != nil {
			log.Println("Error visiting seed URL:", err)
			return nil, err
		}
	}
	scraper.Wait()
	return visitedURLs, nil
}

func improvePlan(ctx context.Context, model llms.Model, p Plan, c chan schema.Document, ec chan error) {
	ms, err := MetadataSchema()
	if err != nil {
		log.Println("Failed in retrieving Schema")
		ec <- err
		return
	}
	var metadata Metadata
	// Enhance scraped documents with gemini and create meaningful metadata
	query := fmt.Sprintf(scrapeTemplateStr, p.Title, p.Description, p.Table.String(), ms)
	answer, err := llms.GenerateFromSinglePrompt(ctx, model, query, llms.WithResponseMIMEType("application/json"))
	if err != nil {
		log.Println("Error when generating answer with LLM:", err)
		ec <- fmt.Errorf("LLM generation error: %w", err)
		return
	}
	log.Println("Successful answer from LLM: ", answer)

	planMap := p.Map()
	// Parse the answer as JSON
	err = json.Unmarshal([]byte(answer), &metadata)
	if err != nil {
		log.Println("Error parsing LLM response:", err)
		ec <- fmt.Errorf("JSON unmarshal error: %w", err)
		return
	}
	// Add the results to the map
	maps.Copy(planMap, StructToMap(metadata))
	// Create request body by converting the plans into documents
	c <- schema.Document{
		PageContent: p.String(),
		Metadata:    planMap,
	}
}
