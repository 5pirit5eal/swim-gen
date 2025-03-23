package rag

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/golobby/dotenv"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
	"github.com/tmc/langchaingo/schema"
)

type MockModel struct{}

func (m *MockModel) GenerateContent(ctx context.Context, prompt []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	metadata := Metadata{
		// Mock metadata fields
		Freistil:           true,
		Brust:              true,
		Delfin:             true,
		Ruecken:            true,
		Schwierigkeitsgrad: Anfaenger,
		Trainingstyp:       Techniktraining,
	}
	metadataJSON, _ := json.Marshal(metadata)
	return &llms.ContentResponse{Choices: []*llms.ContentChoice{&llms.ContentChoice{Content: string(metadataJSON)}}}, nil
}

func (m *MockModel) Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	return "", nil
}

func TestImprovePlan(t *testing.T) {
	ctx := context.Background()
	model := &MockModel{}
	plan := Plan{
		URL:         "http://example.com",
		Title:       "Test Plan",
		Description: "This is a test plan",
		Table: Table{
			{Amount: 1, Multiplier: "x", Distance: 100, Break: "10s", Content: "Swim", Intensity: "High", Sum: 100},
		},
	}

	c := make(chan schema.Document)
	ec := make(chan error)

	go improvePlan(ctx, model, plan, c, ec)

	select {
	case doc := <-c:
		assert.Equal(t, plan.String(), doc.PageContent)
		assert.Contains(t, doc.Metadata, "url")
		assert.Contains(t, doc.Metadata, "title")
		assert.Contains(t, doc.Metadata, "description")
		assert.Contains(t, doc.Metadata, "table")
	case err := <-ec:
		t.Fatalf("Error improving plan: %v", err)
	}
}

func TestImprovePlanWithLLMCall(t *testing.T) {
	config := Config{}
	file, err := os.Open("../.env")
	if err != nil {
		log.Fatal(err)
	}
	if err := dotenv.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	model, err := vertex.New(
		ctx, googleai.WithCloudProject(config.ProjectID),
		googleai.WithCloudLocation(config.Region),
		googleai.WithDefaultModel(config.Model),
		googleai.WithHarmThreshold(googleai.HarmBlockLowAndAbove),
	)
	if err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		URL:         "https://docswim.de/index.php/2017/07/10/trainingsplan-01-grundlagen-fundament-3-700m/",
		Title:       "Trainingsplan #01: Grundlagen-Fundament / 3.700m",
		Description: "Grundlagenausdauer: Immer aktuell\nDie Trainingseinheit der Woche dient der Entwicklung der Grundlagenausdauer mit vorwiegend extensiven Belastungen. Zur Sicherung eines Belastungswechsels und einer leichten Variabilität in der Belastung, eignen sich moderate und kurze intensive Anteile in Form von Tempowechselspielen oder leichten Temposteigerungen.\nDiese Einheit kann auch als regenerative Einheit bestens im Trainingskonzept eingefügt werden. Dann einfach die GA2-Anteile durch GA1 oder sogar ReKom ersetzen. Schon ist es die optimale Vor- oder Nachbereitungseinheit!",
		Table: Table{
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: " Inhalt", Intensity: "", Sum: 0},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 1, Multiplier: "x", Distance: 300, Break: "", Content: "Einschwimmen", Intensity: "", Sum: 300},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 2, Multiplier: "x", Distance: 600, Break: "60", Content: "Je: 50m GA1-2 + 150m GA1", Intensity: "GA1-2", Sum: 1200},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 2, Multiplier: "x", Distance: 400, Break: "45", Content: "Je: 200m GA1 + 200m GA1-2 leicht gesteig.", Intensity: "GA1-1", Sum: 800},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 2, Multiplier: "x", Distance: 200, Break: "45", Content: "Je: 100m GA1 + 50m GA1-2 + 50m GA2", Intensity: "GA1+GA2", Sum: 400},
			{Amount: 2, Multiplier: "x", Distance: 200, Break: "45", Content: "Je: 100m GA1 + 50m GA1-2 + 50m GA2", Intensity: "GA1+GA2", Sum: 400},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 2, Multiplier: "x", Distance: 100, Break: "30", Content: "Jeweils: 75m GA1 + 25m GA2", Intensity: "GA1+GA2", Sum: 200},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 1, Multiplier: "x", Distance: 200, Break: "", Content: "Ausschwimmen", Intensity: "", Sum: 200},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "", Intensity: "", Sum: 0},
			{Amount: 0, Multiplier: "", Distance: 0, Break: "", Content: "Gesamt", Intensity: "", Sum: 3500},
		},
	}
	c := make(chan schema.Document)
	ec := make(chan error)
	go improvePlan(ctx, model, plan, c, ec)
	select {
	case doc := <-c:
		assert.Equal(t, plan.String(), doc.PageContent)
		assert.Contains(t, doc.Metadata, "url")
		assert.Contains(t, doc.Metadata, "title")
		assert.Contains(t, doc.Metadata, "description")
		assert.Contains(t, doc.Metadata, "table")
		assert.Equal(t, doc.Metadata["freistil"], true)
		assert.Equal(t, doc.Metadata["brust"], false)
		assert.Equal(t, doc.Metadata["delfin"], false)
		assert.Equal(t, doc.Metadata["ruecken"], false)
		assert.Equal(t, doc.Metadata["lagen"], false)
		assert.Equal(t, doc.Metadata["schwierigkeitsgrad"], Fortgeschritten)
		assert.Equal(t, doc.Metadata["trainingstyp"], Grundlagen)
	case err := <-ec:
		t.Fatalf("Error improving plan: %v", err)
	}

}
func TestScrape(t *testing.T) {
	alreadyVisited := []string{}
	seeds := []string{
		"https://docswim.de/index.php/2017/07/10/trainingsplan-01-grundlagen-fundament-3-700m/",
		"https://docswim.de/index.php/2019/09/02/trainingsplan-99-kraulschwimmen-lernen-der-kraul-kurs-teil-2-2-1-700m/",
	}

	plans, err := Scrape(alreadyVisited, seeds...)
	assert.NoError(t, err)
	assert.NotNil(t, plans)
	assert.Greater(t, plans.Len(), 0)

	n := 0
	for kvp := range plans.Range() {
		if len(kvp.Plan.Table) != 0 {
			assert.NotEmpty(t, kvp.URL)
			assert.NotEmpty(t, kvp.Plan.Title)
			assert.NotEmpty(t, kvp.Plan.Description)
			n++
		}
	}
	assert.Equal(t, n, 2)
}
