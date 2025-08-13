package dto_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MDx3R/ef-test/internal/transport/http/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMonthYear_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.MonthYear
		expected string
	}{
		{
			name:     "valid date",
			input:    dto.MonthYear("08-2025"),
			expected: `"08-2025"`,
		},
		{
			name:     "empty string",
			input:    dto.MonthYear("null"),
			expected: `"null"`,
		},
		{
			name:     "empty string",
			input:    dto.MonthYear(""),
			expected: `"null"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestMonthYear_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "valid date",
			jsonData: `"08-2025"`,
			expected: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "null value",
			jsonData: `null`,
			expected: time.Time{},
			wantErr:  false,
		},
		{
			name:     "empty string",
			jsonData: `""`,
			expected: time.Time{},
			wantErr:  false,
		},
		{
			name:     "invalid format",
			jsonData: `"2025-08-01"`,
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var my dto.MonthYear
			err := json.Unmarshal([]byte(tt.jsonData), &my)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, my.ToTime().Equal(tt.expected), "expected %v, got %v", tt.expected, my.ToTime())
		})
	}
}

func TestMonthYear_IntegrationWithStruct(t *testing.T) {
	type TestStruct struct {
		ID   uuid.UUID     `json:"id"`
		Date dto.MonthYear `json:"date"`
	}

	expectedDate := dto.MonthYear("12-2025")
	obj := TestStruct{
		ID:   uuid.MustParse("45b09176-3055-4be6-9571-1cfb8cece98e"),
		Date: expectedDate,
	}

	// Marshal
	data, err := json.Marshal(obj)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"12-2025"`)

	// Unmarshal
	var parsed TestStruct
	err = json.Unmarshal(data, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, expectedDate, parsed.Date)
}
