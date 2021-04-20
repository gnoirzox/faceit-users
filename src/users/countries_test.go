package users

import "testing"

// This is an integration test;
// it expects to have the Postgres DB to be running
func TestIsValidCountry(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"fra", true},
		{"FRA", true},
		{"RUS", true},
		{"RES", false},
		{"RS", false},
		{"rs", false},
		{"United Kingdom", false},
	}

	for _, test := range tests {
		var c Country
		c.IsoAlphaCode = test.input

		if got, _ := c.IsValidCountry(); got != test.want {
			t.Errorf("c.IsoAlphaCode = %q c.IsValidCountry() = %v", test.input, got)
		}
	}
}
