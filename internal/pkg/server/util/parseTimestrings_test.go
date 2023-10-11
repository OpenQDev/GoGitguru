package util

import (
	"testing"
	"time"
)

func TestParseTimes(t *testing.T) {
	tests := []struct {
		name        string
		timeStrings []string
		want        []time.Time
		wantErr     bool
	}{
		{
			name:        "Valid time strings",
			timeStrings: []string{"2022-01-01T00:00:00Z", "2022-12-31T23:59:59Z"},
			want:        []time.Time{time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC)},
			wantErr:     false,
		},
		{
			name:        "Invalid time string",
			timeStrings: []string{"invalid"},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimes(tt.timeStrings...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equal(got, tt.want) {
				t.Errorf("ParseTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equal(a, b []time.Time) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !v.Equal(b[i]) {
			return false
		}
	}
	return true
}
