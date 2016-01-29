package ghprj

import "testing"

var epsilon = 0.001

func TestFuzz(t *testing.T) {
	tests := []struct {
		choice string
		query  string
		score  float64
	}{
		{"angular-tooltips", "api", 0.0},
		{"r101-api", "api", 0.375},
		{"api", "api", 1.0},
		{"lol", "", 1.0},
		{"r101-api", "rapi", 0.285},
	}

	for _, tt := range tests {
		s := score(tt.choice, tt.query)
		if s < tt.score-epsilon || s > tt.score+epsilon {
			t.Errorf("expected score to be %f for %s and %s, got %f", tt.score, tt.choice, tt.query, s)
		}
	}
}
