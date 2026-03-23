package models

import "strings"

// SanitizeString removes invalid UTF-8 sequences, null bytes, and known ad-injection
// artifacts from a string to prevent PostgreSQL encoding errors and data pollution.
func SanitizeString(s string) string {
	// Drop invalid UTF-8 byte sequences (e.g. 0x80–0xFF garbage from scraped content)
	s = strings.ToValidUTF8(s, "")
	// Remove null bytes (valid UTF-8 but rejected by PostgreSQL text columns).
	// Use string([]byte{0x00}) to construct the null byte at runtime so that editors
	// and formatters cannot silently strip the literal  from the source file.
	s = strings.ReplaceAll(s, string([]byte{0x00}), "")
	// Remove Google AdSense push artifact injected by some scraped pages
	s = strings.ReplaceAll(s, "(adsbygoogle = window.adsbygoogle || []).push({})", "")
	return s
}

// SanitizeRows recursively sanitizes all string fields of each row in the table, including
// nested subrows, to ensure all text is clean and free of invalid characters or ad artifacts.
// All string fields are sanitized because they all appear in PageContent via Row.String().
func SanitizeRows(table *Table) {
	for i := range *table {
		row := &(*table)[i]
		row.Content = SanitizeString(row.Content)
		row.Multiplier = SanitizeString(row.Multiplier)
		row.Break = SanitizeString(row.Break)
		row.Intensity = SanitizeString(row.Intensity)
		if len(row.SubRows) > 0 {
			subTable := Table(row.SubRows)
			SanitizeRows(&subTable)
			row.SubRows = []Row(subTable)
		}
	}
}
