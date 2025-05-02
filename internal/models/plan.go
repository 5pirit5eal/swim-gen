package models

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/tmc/langchaingo/schema"
)

type Mappable interface {
	// Map represenation of the object with at least the plan_id
	Map() map[string]any
	Plan() *Plan
}

type DonatedPlan struct {
	UserID string `db:"user_id"`
	PlanID string `db:"plan_id"`
	// CreatedAt is the time the plan was donated as a datetime string
	CreatedAt   string `db:"created_at"`
	Title       string `db:"title"`
	Description string `db:"description"`
	// Table is the table associated with the plan
	Table Table `db:"plan_table"`
}

func (d *DonatedPlan) Map() map[string]any {
	m := map[string]any{
		"plan_id":    d.PlanID,
		"created_at": d.CreatedAt,
		"title":      d.Title,
	}

	return m
}

func (d *DonatedPlan) Plan() *Plan {
	return &Plan{
		Title:       d.Title,
		Description: d.Description,
		Table:       d.Table,
	}
}

type ScrapedPlan struct {
	PlanID string `db:"plan_id"`
	URL    string `db:"url"`
	// CreatedAt is the time the plan was scraped as a datetime string
	CreatedAt   string `db:"created_at"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Table       Table  `db:"plan_table"`
}

func (s *ScrapedPlan) Map() map[string]any {
	m := map[string]any{
		"plan_id": s.PlanID,
		"url":     s.URL,
		"title":   s.Title,
	}

	return m
}

func (s *ScrapedPlan) Plan() *Plan {
	return &Plan{
		Title:       s.Title,
		Description: s.Description,
		Table:       s.Table,
	}
}

type Plan struct {
	PlanID      string `db:"plan_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Table       Table  `db:"table"`
}

func (p *Plan) Map() map[string]any {
	m := map[string]any{
		"plan_id":     p.PlanID,
		"title":       p.Title,
		"description": p.Description,
	}

	return m
}
func (p *Plan) Plan() *Plan {
	return p
}

func (p *Plan) String() string {
	return fmt.Sprintf("%s:\n %s\n %s", p.Title, p.Description, p.Table.String())
}

func PlanToDoc(plan Mappable, metadata *Metadata) schema.Document {
	genericPlan := plan.Plan()
	// Create a map of the plan
	planMap := plan.Map()

	// Add the metadata to the map
	maps.Copy(planMap, StructToMap(metadata))

	// Add the description to the plan descriptions
	genericPlan.Description += "\n" + metadata.Reasoning
	// Create a document from the plan
	return schema.Document{
		PageContent: genericPlan.String(),
		Metadata:    planMap,
	}
}

type Table []Row

type Row struct {
	Amount     int    `json:"Amount" jsonschema_description:"Amount of repetitions"`
	Multiplier string `json:"Multiplier" jsonschema_description:"Multiplier for the distance (e.g. 'x' or 'times')"`
	Distance   int    `json:"Distance" jsonschema_description:"Distance in meters"`
	Break      string `json:"Break" jsonschema_description:"Break time typically in seconds. This needs to be a string, as other times are possible"`
	Content    string `json:"Content" jsonschema_description:"Content or description of the row"`
	Intensity  string `json:"Intensity" jsonschema_description:"Intensity level of the activity"`
	Sum        int    `json:"Sum" jsonschema_description:"Total volume or sum for the row"`
}

func (r Row) String() string {
	return fmt.Sprintf("| %d | %s | %d | %s | %s | %s | %d |", r.Amount, r.Multiplier, r.Distance, r.Break, r.Content, r.Intensity, r.Sum)
}

func (t *Table) String() string {
	tstr := "| Anzahl |  | Strecke(m) | Pause(s) | Inhalt | Intensität | Umfang |\n"
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

// Recalculates the sum for each row
// This is useful if the table has been modified and we need to update the sums
func (t *Table) UpdateSum() {
	total := 0
	for i, row := range *t {
		if row.Content == "Gesamt" {
			(*t)[i].Sum = total
		} else {
			(*t)[i].Sum = row.Amount * row.Distance
			total += (*t)[i].Sum
		}
	}
}

// Returns the Header of the table
//
// | Anzahl |  | Strecke(m) | Pause(s) | Inhalt | Intensität | Umfang |
func (t *Table) Header() []string {
	return []string{"Anzahl", "", "Strecke(m)", "Pause(s)", "Inhalt", "Intensität", "Umfang"}
}

// Returns the json encoded table as a string
func (t *Table) JSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to marshal table to JSON: %w", err)
	}
	return string(bytes), nil
}
