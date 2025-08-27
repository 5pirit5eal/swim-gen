package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/breakline"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Converts the given table to a PDF string representation.
// Uses maroto to create a PDF document with the table data.
// The PDF is returned as a string, which can be saved to a file or sent to cloud storage.
func TableToPDF(table models.Table) ([]byte, error) {
	m := getMaroto()

	m.AddRows(getRows(table)...)

	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}

func PlanToPDF(plan *models.Plan) ([]byte, error) {
	m := getMaroto()

	m.RegisterHeader(text.NewAutoRow(plan.Title, props.Text{Size: 18, Style: fontstyle.Bold, Align: align.Center}))

	m.AddAutoRow(col.New().Add(text.New(plan.Description, props.Text{Size: 10, Top: 10, Bottom: 10})))

	m.AddRows(getRows((*plan).Table)...)

	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}

// Uploads the given pdf to cloud storage and returns the URI of the uploaded file.
func UploadPDF(ctx context.Context, serviceAccount, bucketName, objectName string, pdfData []byte) (string, error) {
	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	buf := bytes.NewBuffer(pdfData)
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Creates an object writer.
	wc := bucket.Object(objectName).NewWriter(ctx)
	wc.ContentType = "application/pdf"

	// Uploads the PDF data.
	if _, err := io.Copy(wc, buf); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	// Close and flush the contents
	err = wc.Close()
	if err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}

	// Generate signed url (V4 Signing)
	opts := storage.SignedURLOptions{
		GoogleAccessID: serviceAccount,
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
	}

	url, err := bucket.SignedURL(objectName, &opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %w", err)
	}

	return url, nil
}

func GenerateFilename() string {
	return path.Join("anonymous", uuid.NewString()+".pdf")
}

func getMaroto() core.Maroto {
	cfg := config.NewBuilder().
		WithMaxGridSize(25).
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)
	return m
}

func getRows(table models.Table) []core.Row {
	if len(table) < 2 {
		return make([]core.Row, 0)
	}
	// A row consists of 7 columns based on models.Row
	headerRow := row.New(8)
	headerProps := props.Text{Style: fontstyle.Bold, Align: align.Center, Top: 2, Bottom: 2}
	darkGray := &props.Color{Red: 200, Green: 200, Blue: 200}
	lightGray := &props.Color{Red: 240, Green: 240, Blue: 240}

	for i, title := range table.Header() {
		if i == 1 {
			headerRow.Add(text.NewCol(1, title, headerProps))
		} else if i == 4 {
			headerRow.Add(text.NewCol(9, title, headerProps))
		} else {
			headerRow.Add(text.NewCol(3, title, headerProps))
		}
	}
	headerRow.WithStyle(&props.Cell{BackgroundColor: darkGray})

	rows := []core.Row{headerRow}
	p := props.Text{
		Align:  align.Center,
		Top:    2,
		Bottom: 2,
		// Left:   2,
		// Right:  2,
	}
	for i, content := range table {
		row := row.New()
		if i < len(table)-1 {
			row.Add(col.New(3).Add(text.New(strconv.Itoa(content.Amount), p)).WithStyle(&props.Cell{}))
			row.Add(col.New(1).Add(text.New(content.Multiplier, p)))
			row.Add(col.New(3).Add(text.New(strconv.Itoa(content.Distance), p)))
			row.Add(col.New(3).Add(text.New(content.Break, p)))
			row.Add(col.New(9).Add(text.New(content.Content, props.Text{Align: align.Center, Top: 2, Bottom: 2, Left: 2, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy})))
			row.Add(col.New(3).Add(text.New(content.Intensity, p)))
			row.Add(col.New(3).Add(text.New(strconv.Itoa(content.Sum), p)))

			if (i+1)%2 == 0 {
				row.WithStyle(&props.Cell{BackgroundColor: lightGray})
			}

		} else {
			row.Add(
				col.New(16),
				text.NewCol(6, "Gesamtumfang", headerProps),
				text.NewCol(3, strconv.Itoa(content.Sum), headerProps),
			).WithStyle(&props.Cell{BackgroundColor: darkGray})
		}
		rows = append(rows, row)
	}
	return rows
}
