package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"regexp"
	"strconv"
	"strings"
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
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// parseMarkdownLinks converts markdown links [text](url) to plain text.
// Since maroto v2 doesn't support inline hyperlinks (each text component is
// rendered as a separate block which causes overlaps and spacing issues),
// we strip the markdown syntax and just display the link text inline.
// This provides the best reading experience for both standard and easy-to-read PDFs.
func parseMarkdownLinks(content, baseURL string, p props.Text) []core.Component {
	// Regex to match markdown links: [text](url)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)

	// Replace all markdown links with just the link text
	plainContent := linkRegex.ReplaceAllString(content, "$1")

	return []core.Component{text.New(plainContent, p)}
}

// Converts the given table to a PDF string representation.
// Uses maroto to create a PDF document with the table data.
// The PDF is returned as a string, which can be saved to a file or sent to cloud storage.
func GenerateEasyReadablePDF(table *models.Table, ho bool, lang models.Language, baseURL string) ([]byte, error) {
	m := getMaroto(ho)

	m.AddRows(getRows(*table, true, lang, baseURL)...)

	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}

func GenerateFullPDF(plan *models.Plan, ho bool, lang models.Language, baseURL string) ([]byte, error) {
	m := getMaroto(ho)
	m.AddAutoRow(col.New().Add(text.New(plan.Title, props.Text{Size: 18, Style: fontstyle.Bold, Align: align.Center, VerticalPadding: 2})))
	m.AddAutoRow(col.New().Add(text.New(plan.Description, props.Text{Size: 10, Top: 10, Bottom: 10, VerticalPadding: 2})))
	m.AddRows(getRows((*plan).Table, false, lang, baseURL)...)

	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}

// Converts the given plan to a PDF byte slice representation.
//
// Uses maroto to create a PDF document with the plan data.
// The PDF is returned as a byte slice, which can be saved to a file or sent to cloud storage.
func PlanToPDF(plan *models.Plan, ho, lf bool, lang models.Language, baseURL string) ([]byte, error) {
	if !lf {
		return GenerateFullPDF(plan, ho, lang, baseURL)
	}
	return GenerateEasyReadablePDF(&plan.Table, ho, lang, baseURL)

}

// Uploads the given pdf to cloud storage and returns the URI of the uploaded file.
func UploadPDF(ctx context.Context, serviceAccount, bucketName, objectName string, pdfData []byte) (string, error) {
	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer func() { _ = client.Close() }()

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

func GenerateStoragePath(username, planID, title string) string {
	sanitizedTitle := sanitizeFilename(title)
	if sanitizedTitle == "" {
		sanitizedTitle = "training-plan"
	}
	filename := sanitizedTitle + ".pdf"

	if username != "" {
		return path.Join(username, filename)
	}

	if planID != "" {
		return path.Join(planID, filename)
	}

	return GenerateFilename()
}

func sanitizeFilename(name string) string {
	name = strings.ToLower(name)

	// Transliterate German characters
	replacements := []struct {
		original    string
		replacement string
	}{
		{"ä", "ae"},
		{"ö", "oe"},
		{"ü", "ue"},
		{"ß", "ss"},
	}

	for _, r := range replacements {
		name = strings.ReplaceAll(name, r.original, r.replacement)
	}

	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// Remove all characters except alphanumeric, underscores, and hyphens
	reg := regexp.MustCompile(`[^a-z0-9_\-]`)
	name = reg.ReplaceAllString(name, "")

	// Collapse multiple underscores
	reg = regexp.MustCompile(`_{2,}`)
	name = reg.ReplaceAllString(name, "_")

	return name
}

func getMaroto(ho bool) core.Maroto {
	cfg := config.NewBuilder().
		WithMaxGridSize(25).
		WithLeftMargin(10).
		WithTopMargin(15).
		WithBottomMargin(15).
		WithRightMargin(10)

	if ho {
		cfg = cfg.WithTopMargin(5).WithBottomMargin(5).WithOrientation(orientation.Horizontal)
	}

	m := maroto.New(cfg.Build())
	return m
}

// Convert table rows to maroto rows
// lf indicates if large font should be used
// baseURL is prepended to relative URLs in markdown links
func getRows(table models.Table, lf bool, lang models.Language, baseURL string) []core.Row {
	if len(table) < 2 {
		return make([]core.Row, 0)
	}
	// A row consists of 7 columns based on models.Row
	headerRow := row.New()
	headerProps := props.Text{Style: fontstyle.Bold, Align: align.Center, Top: 2, Bottom: 2, VerticalPadding: 1}
	if lf {
		headerProps.Top = 3
		headerProps.Bottom = 3
		headerProps.Size = 12
		headerProps.VerticalPadding = 1.5
	}
	darkGray := &props.Color{Red: 200, Green: 200, Blue: 200}

	for i, title := range table.Header(lang) {
		switch i {
		case 1:
			headerRow.Add(text.NewCol(1, title, headerProps))
		case 4:
			headerRow.Add(text.NewCol(9, title, headerProps))
		default:
			headerRow.Add(text.NewCol(3, title, headerProps))
		}
	}
	headerRow.WithStyle(&props.Cell{BackgroundColor: darkGray})

	rows := []core.Row{headerRow}
	p := props.Text{Align: align.Center, Top: 2, Bottom: 2, VerticalPadding: 1}
	if lf {
		p.Size = 16
		p.Top = 3
		p.Bottom = 3
		p.VerticalPadding = 1.5
	}
	rowIndex := 0
	for i, content := range table {
		// Skip the last row if it's a footer/total row
		if i == len(table)-1 {
			sloganProps := props.Text{Size: headerProps.Size, Align: align.Left, Top: p.Top, Bottom: p.Bottom, Left: 2, Style: fontstyle.BoldItalic, VerticalPadding: headerProps.VerticalPadding}
			footer := table.Footer(lang)
			footerRow := row.New()
			footerRow.Add(
				text.NewCol(7, footer[0], sloganProps),
				col.New(12),
				text.NewCol(3, footer[4], headerProps),
				text.NewCol(3, footer[6], headerProps),
			).WithStyle(&props.Cell{BackgroundColor: darkGray})
			rows = append(rows, footerRow)
			break
		}

		// Add main row
		mainRow := createRow(content, p, rowIndex%2 == 1)
		rows = append(rows, mainRow)
		rowIndex++

		// Add subrows recursively
		if len(content.SubRows) > 0 {
			subRows := createSubRows(content.SubRows, p, lang, baseURL, rowIndex)
			rows = append(rows, subRows...)
			// Update rowIndex based on number of subrows added
			rowIndex += len(subRows)
		}
	}
	return rows
}

// createRow creates a standard row for the PDF
func createRow(content models.Row, p props.Text, alternateStyle bool) core.Row {
	lightGray := &props.Color{Red: 240, Green: 240, Blue: 240}
	row := row.New()
	row.Add(col.New(3).Add(text.New(strconv.Itoa(content.Amount), p)))
	row.Add(col.New(1).Add(text.New(content.Multiplier, p)))
	row.Add(col.New(3).Add(text.New(strconv.Itoa(content.Distance), p)))
	row.Add(col.New(3).Add(text.New(content.Break, p)))

	// Parse markdown links in content and render as text/link components
	contentCol := col.New(9)
	segments := parseMarkdownLinks(content.Content, "", p)
	contentCol.Add(segments...)
	row.Add(contentCol)

	row.Add(col.New(3).Add(text.New(content.Intensity, p)))
	row.Add(col.New(3).Add(text.New(strconv.Itoa(content.Sum), p)))

	if alternateStyle {
		row.WithStyle(&props.Cell{BackgroundColor: lightGray})
	}

	return row
}

// aggregateSubRowContent recursively aggregates nested subrow content into a string
// Format: "content (distance1 Child1 + distance2 Child2)" or "(distance1 Child1 + distance2 Child2)" if no content
// Omits break times and intensities for aggregated subrows
func aggregateSubRowContent(subRows []models.Row, baseURL string, p props.Text) string {
	if len(subRows) == 0 {
		return ""
	}

	parts := make([]string, len(subRows))
	for i, subRow := range subRows {
		// Build content for this subrow
		content := subRow.Content
		if len(subRow.SubRows) > 0 {
			// Recursively aggregate nested subrows
			nestedContent := aggregateSubRowContent(subRow.SubRows, baseURL, p)
			if nestedContent != "" {
				content = fmt.Sprintf("%s (%s)", content, nestedContent)
			}
		}

		// Format: distance + content (if content exists)
		if content != "" {
			parts[i] = fmt.Sprintf("%dm %s", subRow.Distance, content)
		} else {
			parts[i] = fmt.Sprintf("%dm", subRow.Distance)
		}
	}

	return strings.Join(parts, " + ")
}

// createSubRows creates visual rows for subrows with indentation
// Nested subrows are aggregated into the parent subrow's content instead of creating new rows
func createSubRows(subRows []models.Row, p props.Text, lang models.Language, baseURL string, startRowIndex int) []core.Row {
	rows := make([]core.Row, 0)
	veryLightGray := &props.Color{Red: 248, Green: 248, Blue: 248}

	for i, subRow := range subRows {
		row := row.New()
		// Empty amount field for subrows
		row.Add(col.New(3).Add(text.New("", p)))
		// Empty multiplier field for subrows
		row.Add(col.New(1).Add(text.New("", p)))
		// Distance value in distance column
		row.Add(col.New(3).Add(text.New(strconv.Itoa(subRow.Distance), p)))
		// Empty break field for aggregated subrows
		row.Add(col.New(3).Add(text.New("", p)))

		// Content with indent indicator and aggregated nested subrow content
		contentCol := col.New(9)
		indentText := subRow.Content
		if len(subRow.SubRows) > 0 {
			aggregated := aggregateSubRowContent(subRow.SubRows, baseURL, p)
			if aggregated != "" {
				indentText = fmt.Sprintf("%s (%s)", indentText, aggregated)
			}
		}
		segments := parseMarkdownLinks(indentText, baseURL, p)
		contentCol.Add(segments...)
		row.Add(contentCol)

		// Empty intensity field for aggregated subrows
		row.Add(col.New(3).Add(text.New("", p)))
		// Empty sum field for subrows
		row.Add(col.New(3).Add(text.New("", p)))

		// Alternate background colors for subrows
		if (startRowIndex+i)%2 == 1 {
			row.WithStyle(&props.Cell{BackgroundColor: veryLightGray})
		}

		rows = append(rows, row)
	}

	return rows
}
