package format_test

import (
	"testing"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/util/format"
	"github.com/stretchr/testify/assert"
)

func TestFormatConstants(t *testing.T) {
	refTime := time.Date(2026, 9, 8, 14, 30, 45, 0, time.UTC)

	t.Run("DateFormat", func(t *testing.T) {
		formatted := refTime.Format(format.DateFormat)
		assert.Equal(t, "2026-09-08", formatted)

		parsed, err := time.Parse(format.DateFormat, formatted)
		assert.NoError(t, err)
		assert.Equal(t, 2026, parsed.Year())
		assert.Equal(t, time.Month(9), parsed.Month())
		assert.Equal(t, 8, parsed.Day())
	})

	t.Run("DateTimeFormat", func(t *testing.T) {
		formatted := refTime.Format(format.DateTimeFormat)
		assert.Equal(t, "2026-09-08 14:30:45", formatted)

		parsed, err := time.Parse(format.DateTimeFormat, formatted)
		assert.NoError(t, err)
		assert.Equal(t, 2026, parsed.Year())
		assert.Equal(t, 14, parsed.Hour())
		assert.Equal(t, 30, parsed.Minute())
		assert.Equal(t, 45, parsed.Second())
	})

	t.Run("TimeFormat", func(t *testing.T) {
		formatted := refTime.Format(format.TimeFormat)
		assert.Equal(t, "14:30:45", formatted)

		parsed, err := time.Parse(format.TimeFormat, formatted)
		assert.NoError(t, err)
		assert.Equal(t, 14, parsed.Hour())
		assert.Equal(t, 30, parsed.Minute())
		assert.Equal(t, 45, parsed.Second())
	})
}
