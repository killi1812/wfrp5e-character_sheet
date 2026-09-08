package cerror_test

import (
	"errors"
	"testing"

	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"
	"github.com/stretchr/testify/assert"
)

func TestCerrorDefinitions(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrBadDateFormat", cerror.ErrBadDateFormat, "bad date format, should be 2006-01-02"},
		{"ErrBadDateTimeFormat", cerror.ErrBadDateTimeFormat, "bad date and time format, should be 2006-01-02 15:04:05"},
		{"ErrBadTimeFormat", cerror.ErrBadTimeFormat, "bad time format, should be 15:04:05"},
		{"ErrBadUuid", cerror.ErrBadUuid, "failed to parse uuid"},
		{"ErrUnknownRole", cerror.ErrUnknownRole, "unknown role"},
		{"ErrInvalidCredentials", cerror.ErrInvalidCredentials, "invalid email or password"},
		{"ErrInvalidTokenFormat", cerror.ErrInvalidTokenFormat, "invalid token format"},
		{"ErrUserIsNil", cerror.ErrUserIsNil, "user is nil"},
		{"ErrBadRole", cerror.ErrBadRole, "role is not allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.Equal(t, tt.expected, tt.err.Error())
			assert.True(t, errors.Is(tt.err, tt.err))
		})
	}
}
