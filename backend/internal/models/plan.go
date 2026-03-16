package models

import (
	"encoding/json"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/invopop/jsonschema"
	"github.com/tmc/langchaingo/schema"
)

type Language string

const (
	LanguageEN Language = "en"
	LanguageDE Language = "de"
)

type Planable interface {
	// Map represenation of the object with at least the plan_id without the table
	Map() map[string]any
	// Plan returns the basic Plan structure
	Plan() *Plan
}

type DonatedPlan struct {
	UserID string `db:"user_id" json:"user_id"`
	PlanID string `db:"plan_id" json:"plan_id"`
	// CreatedAt is the time the plan was donated as a datetime string
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	// Table is the table associated with the plan
	Table Table `db:"plan_table" json:"table"`
	// AllowSharing indicates if the plan can be used in the RAG system
	AllowSharing bool `db:"allow_sharing" json:"allow_sharing"`
}

func (d *DonatedPlan) Map() map[string]any {
	m := map[string]any{
		"user_id":     d.UserID,
		"plan_id":     d.PlanID,
		"created_at":  d.CreatedAt,
		"title":       d.Title,
		"description": d.Description,
	}

	return m
}

func (d *DonatedPlan) Plan() *Plan {
	return &Plan{
		PlanID:      d.PlanID,
		Title:       d.Title,
		Description: d.Description,
		Table:       d.Table,
	}
}

type ScrapedPlan struct {
	PlanID string `db:"plan_id"`
	URL    string `db:"url"`
	// CreatedAt is the time the plan was scraped as a datetime string
	CreatedAt   time.Time `db:"created_at"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Table       Table     `db:"plan_table"`
}

func (s *ScrapedPlan) Map() map[string]any {
	m := map[string]any{
		"plan_id":     s.PlanID,
		"url":         s.URL,
		"title":       s.Title,
		"created_at":  s.CreatedAt,
		"description": s.Description,
	}

	return m
}

func (s *ScrapedPlan) Plan() *Plan {
	return &Plan{
		PlanID:      s.PlanID,
		Title:       s.Title,
		Description: s.Description,
		Table:       s.Table,
	}
}

type GeneratedPlan struct {
	Title       string `db:"title" example:"Advanced Freestyle Training" jsonschema_description:"Title of the training plan"`
	Description string `db:"description" example:"A comprehensive training plan for improving freestyle technique" jsonschema_description:"Description or comments about the training plan"`
	Table       Table  `db:"table" jsonschema_description:"Structured table containing the training plan details"`
}

func (gp *GeneratedPlan) Map() map[string]any {
	m := map[string]any{
		"title":       gp.Title,
		"description": gp.Description,
		"plan_table":  gp.Table,
	}

	return m
}

func (gp *GeneratedPlan) Plan() *Plan {
	return &Plan{
		PlanID:      uuid.New().String(),
		Title:       gp.Title,
		Description: gp.Description,
		Table:       gp.Table,
	}
}

func GeneratedPlanSchema() (map[string]any, error) {
	schema := jsonschema.Reflect(&GeneratedPlan{})

	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON schema: %w", err)
	}
	var result map[string]any
	return result, json.Unmarshal(jsonSchema, &result)
}

// BASIC PLAN STRUCTURES

// Plan represents a swim training plan
// @Description A swim training plan with title, description, and structured table
type Plan struct {
	PlanID      string `db:"plan_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Table       Table  `db:"plan_table"`
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

// Table represents a training plan table with multiple rows
// @Description A structured training plan table containing exercise rows
type Table []Row

// Row represents a single exercise entry in a training plan
// @Description A single exercise entry with amount, distance, breaks, content, intensity and total volume. Supports nested Children for compound sets like 8 x (800 + 200).
type Row struct {
	Amount     int      `json:"Amount" example:"4" jsonschema_description:"Amount of repetitions"`
	Multiplier string   `json:"Multiplier" example:"x" jsonschema_description:"Multiplier for the distance (e.g. 'x' or 'times')"`
	Distance   int      `json:"Distance" example:"100" jsonschema_description:"Distance in meters. For parent rows with Children, this is auto-calculated as sum of children distances"`
	Break      string   `json:"Break" example:"20" jsonschema_description:"Break time typically in seconds. This needs to be a string, as other times are possible"`
	Content    string   `json:"Content" example:"Freestyle swim" jsonschema_description:"Content or description of the row"`
	Intensity  string   `json:"Intensity" example:"Easy" jsonschema_description:"Intensity level of the activity"`
	Sum        int      `json:"Sum" example:"400" jsonschema_description:"Total volume or sum for the row"`
	Children   []Row    `json:"Children,omitempty" jsonschema_description:"Nested exercise rows for compound sets (e.g., 8 x (800 + 200)). Parent Distance is auto-calculated from children"`
	Equipment  []string `json:"Equipment,omitempty" jsonschema:"description=Equipment needed for this row,enum=Flossen,enum=Kickboard,enum=Handpaddles,enum=Pull buoy,enum=Schnorchel" jsonschema_description:"Equipment needed for this specific row" example:"[\"Flossen\"]"`
}

func (r Row) String() string {
	equipmentStr := ""
	if len(r.Equipment) > 0 {
		equipmentStr = strings.Join(r.Equipment, ", ")
	}

	if len(r.Children) > 0 {
		childrenStr := ""
		for i, child := range r.Children {
			if i > 0 {
				childrenStr += " + "
			}
			childrenStr += fmt.Sprintf("%dm", child.Distance)
		}
		return fmt.Sprintf("| %d | %s | %d | %s | %dx(%s) | %s | %d | %s |", r.Amount, r.Multiplier, r.Distance, r.Break, r.Amount, childrenStr, r.Intensity, r.Sum, equipmentStr)
	}
	return fmt.Sprintf("| %d | %s | %d | %s | %s | %s | %d | %s |", r.Amount, r.Multiplier, r.Distance, r.Break, r.Content, r.Intensity, r.Sum, equipmentStr)
}

func (t *Table) String() string {
	tstr := "| Anzahl |  | Strecke(m) | Pause(s) | Inhalt | Intensität | Umfang | Ausrüstung |\n"
	tstr += "|---|---|---|---|---|---|---|---|\n"
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

// Recalculates the sum for each row and updates Distance for parent rows
// This is useful if the table has been modified and we need to update the sums
func (t *Table) UpdateSum() {
	total := 0
	for i, row := range *t {
		if strings.Contains(row.Content, "Gesamt") || strings.Contains(row.Content, "Total") {
			(*t)[i].Sum = total
		} else if len(row.Children) > 0 {
			// Nested row: recalculate children sums, then parent
			childrenSum := 0
			childrenDistance := 0
			for j, child := range row.Children {
				(*t)[i].Children[j].Sum = child.Amount * child.Distance
				childrenSum += (*t)[i].Children[j].Sum
				childrenDistance += child.Distance
			}
			// Update parent Distance to be sum of children distances
			(*t)[i].Distance = childrenDistance
			// Parent Sum = Amount × sum of children sums
			(*t)[i].Sum = row.Amount * childrenSum
			total += (*t)[i].Sum
		} else {
			// Flat row calculation (backward compatible)
			(*t)[i].Sum = row.Amount * row.Distance
			total += (*t)[i].Sum
		}
	}
}

// Returns the Header of the table
func (t *Table) Header(lang Language) []string {
	switch lang {
	case LanguageDE:
		return []string{"Anzahl", "", "Strecke(m)", "Pause(s)", "Inhalt", "Intensität", "Umfang"}
	default: // LanguageEN and any other unsupported languages
		return []string{"Amount", "", "Distance(m)", "Break(s)", "Content", "Intensity", "Volume"}
	}
}

// Returns the bottom row of the table
func (t *Table) Footer(lang Language) []string {
	sum := strconv.Itoa((*t)[len(*t)-1].Sum) + " m"
	switch lang {
	case LanguageDE:
		return []string{"KI-GENERIERT MIT SWIM-GEN.COM", "", "", "", "Gesamt", "", sum}
	default: // LanguageEN and any other unsupported languages
		return []string{"AI-GENERATED WITH SWIM-GEN.COM", "", "", "", "Total meters", "", sum}
	}
}

// Returns the json encoded table as a string
func (t *Table) JSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to marshal table to JSON: %w", err)
	}
	return string(bytes), nil
}

// Validate recursively validates the table structure
func (t *Table) Validate() error {
	return t.validateRowDepth(0)
}

func (t *Table) validateRowDepth(depth int) error {
	if depth > 5 {
		return fmt.Errorf("maximum nesting depth (5) exceeded")
	}

	for i, row := range *t {
		if row.Amount < 0 {
			return fmt.Errorf("row %d has negative amount: %d", i, row.Amount)
		}
		if row.Distance < 0 {
			return fmt.Errorf("row %d has negative distance: %d", i, row.Distance)
		}

		if len(row.Children) > 0 {
			if row.Amount == 0 {
				return fmt.Errorf("row %d has children but Amount = 0", i)
			}
			childrenTable := Table(row.Children)
			if err := childrenTable.validateRowDepth(depth + 1); err != nil {
				return fmt.Errorf("row %d children: %w", i, err)
			}
		}
	}
	return nil
}

// FlattenTable converts a nested table to a flat representation for display
func (t *Table) FlattenTable(indent string) []string {
	lines := []string{}
	for _, row := range *t {
		if len(row.Children) > 0 {
			childrenStr := ""
			for i, child := range row.Children {
				if i > 0 {
					childrenStr += " + "
				}
				childrenStr += fmt.Sprintf("%dm", child.Distance)
			}
			lines = append(lines, fmt.Sprintf("%s%d x (%s) - %s (Sum: %dm)", indent, row.Amount, childrenStr, row.Content, row.Sum))
			childrenTable := Table(row.Children)
			lines = append(lines, childrenTable.FlattenTable(indent+"  ")...)
		} else {
			lines = append(lines, fmt.Sprintf("%s%d x %dm - %s (Sum: %dm)", indent, row.Amount, row.Distance, row.Content, row.Sum))
		}
	}
	return lines
}

// GetTotalVolume calculates total volume from row sums
// Assumes UpdateSum() has been called to set correct Sum values
func (t *Table) GetTotalVolume() int {
	total := 0
	for _, row := range *t {
		if !strings.Contains(row.Content, "Gesamt") && !strings.Contains(row.Content, "Total") {
			total += row.Sum
		}
	}
	return total
}

type Document struct {
	Plan Planable
	Meta *Metadata
}

func (doc Document) ToLangChainDoc() (schema.Document, error) {
	genericPlan := doc.Plan.Plan()
	// Create a map of the plan
	planMap := doc.Plan.Map()

	if _, found := planMap["plan_id"]; !found {
		return schema.Document{}, fmt.Errorf("plan_id not found in plan map")
	}

	// Add the metadata to the map
	maps.Copy(planMap, StructToMap(doc.Meta))

	// Add the description to the plan descriptions
	genericPlan.Description += "\n" + doc.Meta.Reasoning
	// Create a document from the plan
	return schema.Document{
		PageContent: genericPlan.String(),
		Metadata:    planMap,
	}, nil
}
