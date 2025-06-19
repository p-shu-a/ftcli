package shared

import "testing"

func TestSuggestNewFileName(t *testing.T) {
	tests := []struct {
		input string
		ctr   int
		want  string
	}{
		{"test.txt", 1, "test.1.txt"},
		{"file.tar.gz", 3, "file.tar.3.gz"},
	}
	for _, tt := range tests {
		got := SuggestNewFileName(tt.input, tt.ctr)
		if got != tt.want {
			t.Errorf("SuggestNewFileName(%q,%d) = %q, want %q", tt.input, tt.ctr, got, tt.want)
		}
	}
}
