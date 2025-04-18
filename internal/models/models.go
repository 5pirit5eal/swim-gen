package models

import "fmt"

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
	Amount     int    `json:"Amount"`     // Amount of repetitions
	Multiplier string `json:"Multiplier"` // Multiplier for the distance (e.g., "x", "times")
	Distance   int    `json:"Distance"`   // Distance in meters
	Break      string `json:"Break"`      // Break time in seconds
	Content    string `json:"Content"`    // Content or description of the row
	Intensity  string `json:"Intensity"`  // Intensity level of the activity
	Sum        int    `json:"Sum"`        // Total volume or sum for the row
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

type Document struct {
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
type AddRequest struct {
	Documents []Document `json:"documents"`
}

type QueryRequest struct {
	Content string `json:"content"`
	// This has to be a map[string]any to support langchaingo filter syntax
	// pgvector currently only supports key=value pairs though
	Filter map[string]any `json:"filter,omitempty"`
}

type RAGResponse struct {
	Description string `json:"description"`
	Plan        string `json:"plan,omitempty"`
	Table       Table  `json:"table"`
}

type ChooseResponse struct {
	Idx         int    `json:"index"`
	Description string `json:"description"`
}
