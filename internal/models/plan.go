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
