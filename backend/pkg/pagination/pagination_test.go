package pagination

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		page, per   string
		wantPage    int
		wantPerPage int
	}{
		{"defaults on empty", "", "", DefaultPage, DefaultPerPage},
		{"valid values", "3", "50", 3, 50},
		{"non-numeric falls back", "abc", "xyz", DefaultPage, DefaultPerPage},
		{"zero and negative fall back", "0", "-5", DefaultPage, DefaultPerPage},
		{"per_page capped at max", "1", "9999", 1, MaxPerPage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.page, tt.per)
			if got.Page != tt.wantPage || got.PerPage != tt.wantPerPage {
				t.Errorf("Parse(%q, %q) = %+v, want page=%d per_page=%d",
					tt.page, tt.per, got, tt.wantPage, tt.wantPerPage)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	p := Params{Page: 3, PerPage: 20}
	if got := p.Offset(); got != 40 {
		t.Errorf("Offset() = %d, want 40", got)
	}
}

func TestNewMeta(t *testing.T) {
	meta := NewMeta(Params{Page: 2, PerPage: 20}, 45)
	if meta.TotalPages != 3 {
		t.Errorf("TotalPages = %d, want 3", meta.TotalPages)
	}
	empty := NewMeta(Params{Page: 1, PerPage: 20}, 0)
	if empty.TotalPages != 1 {
		t.Errorf("TotalPages for empty result = %d, want 1", empty.TotalPages)
	}
}
