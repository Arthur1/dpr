package deploypackage

import (
	"testing"
	"time"

	"github.com/Arthur1/dpr/internal/tagdb"
	"github.com/stretchr/testify/assert"
)

func TestGetLastUpdatedAtFromTagRows(t *testing.T) {
	t.Run("returns latest UpdatedAt", func(t *testing.T) {
		tagRows := []*tagdb.TagRow{
			{UpdatedAt: timegen(t, "2024-01-01T00:00:00+09:00")},
			{UpdatedAt: timegen(t, "2024-01-30T09:00:00+09:00")},
			{UpdatedAt: timegen(t, "2023-12-01T12:00:00+09:00")},
		}
		actual := getLastUpdatedAtFromTagRows(tagRows)
		expected := timegen(t, "2024-01-30T09:00:00+09:00")
		assert.Equal(t, expected, actual)
	})

	t.Run("throws panic if input is empty", func(t *testing.T) {
		assert.Panics(t, func() {
			getLastUpdatedAtFromTagRows([]*tagdb.TagRow{})
		})
	})
}

func timegen(t *testing.T, str string) time.Time {
	ti, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t.Error(err)
	}
	return ti
}
