package release

import "testing"

func TestFilename(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"Trilogie, chapitre 1 : Automatique", "Trilogie, chapitre 1 - Automatique"},
	}
	for _, test := range tests {
		if out := Filename(test.in); out != test.out {
			t.Errorf("want: %s, have: %s", test.out, out)
		}
	}
}
