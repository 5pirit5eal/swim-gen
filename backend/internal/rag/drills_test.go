package rag

import (
	"strings"
	"testing"
)

func TestBuildDrillDataQuery(t *testing.T) {
	whereClause := "cmetadata->>'language' = $1"

	t.Run("alphabetical order by default with LIMIT/OFFSET pushed into SQL", func(t *testing.T) {
		query := buildDrillDataQuery(whereClause, 2, false)

		if !strings.Contains(query, "ORDER BY cmetadata->>'title'") {
			t.Errorf("expected alphabetical ORDER BY, got: %s", query)
		}
		if !strings.Contains(query, "LIMIT $2 OFFSET $3") {
			t.Errorf("expected LIMIT/OFFSET placeholders in SQL so Postgres returns at most limit rows, got: %s", query)
		}
	})

	t.Run("random order when requested, LIMIT still pushed into SQL", func(t *testing.T) {
		query := buildDrillDataQuery(whereClause, 2, true)

		if !strings.Contains(query, "ORDER BY RANDOM()") {
			t.Errorf("expected ORDER BY RANDOM(), got: %s", query)
		}
		if !strings.Contains(query, "LIMIT $2 OFFSET $3") {
			t.Errorf("expected LIMIT/OFFSET placeholders in SQL so only the sampled rows leave the DB, got: %s", query)
		}
		if strings.Contains(query, "cmetadata->>'title'") {
			t.Errorf("random query must not order alphabetically, got: %s", query)
		}
	})
}
