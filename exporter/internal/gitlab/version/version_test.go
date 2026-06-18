package version

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    GitLabVersion
		wantErr bool
	}{
		{
			name:  "major only",
			input: "17",
			want:  GitLabVersion{Major: 17, Raw: "17"},
		},
		{
			name:  "major minor",
			input: "17.0",
			want:  GitLabVersion{Major: 17, Minor: 0, Raw: "17.0"},
		},
		{
			name:  "major minor patch",
			input: "17.4.2",
			want:  GitLabVersion{Major: 17, Minor: 4, Patch: 2, Raw: "17.4.2"},
		},
		{
			name:  "ee suffix",
			input: "17.4.2-ee",
			want:  GitLabVersion{Major: 17, Minor: 4, Patch: 2, Raw: "17.4.2-ee"},
		},
		{
			name:  "pre suffix",
			input: "16.10.0-pre",
			want:  GitLabVersion{Major: 16, Minor: 10, Patch: 0, Raw: "16.10.0-pre"},
		},
		{
			name:  "double-digit minor",
			input: "16.11",
			want:  GitLabVersion{Major: 16, Minor: 11, Raw: "16.11"},
		},
		{
			name:  "high major",
			input: "99.99",
			want:  GitLabVersion{Major: 99, Minor: 99, Raw: "99.99"},
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "non-numeric leading",
			input:   "garbage",
			wantErr: true,
		},
		{
			name:    "non-numeric minor",
			input:   "17.x",
			wantErr: true,
		},
		{
			name:  "patch with only suffix is tolerated",
			input: "17.4.-ee",
			want:  GitLabVersion{Major: 17, Minor: 4, Patch: 0, Raw: "17.4.-ee"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) = %v, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("Parse(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name string
		a    GitLabVersion
		b    GitLabVersion
		want int
	}{
		{
			name: "equal major minor",
			a:    GitLabVersion{Major: 17, Minor: 0},
			b:    GitLabVersion{Major: 17, Minor: 0},
			want: 0,
		},
		{
			name: "equal all components",
			a:    GitLabVersion{Major: 17, Minor: 4, Patch: 2},
			b:    GitLabVersion{Major: 17, Minor: 4, Patch: 2},
			want: 0,
		},
		{
			name: "a major greater",
			a:    GitLabVersion{Major: 18, Minor: 0},
			b:    GitLabVersion{Major: 17, Minor: 99},
			want: 1,
		},
		{
			name: "a major less",
			a:    GitLabVersion{Major: 16, Minor: 11},
			b:    GitLabVersion{Major: 17, Minor: 0},
			want: -1,
		},
		{
			name: "a minor greater",
			a:    GitLabVersion{Major: 17, Minor: 5},
			b:    GitLabVersion{Major: 17, Minor: 0},
			want: 1,
		},
		{
			name: "a minor less",
			a:    GitLabVersion{Major: 17, Minor: 0},
			b:    GitLabVersion{Major: 17, Minor: 1},
			want: -1,
		},
		{
			name: "a patch greater",
			a:    GitLabVersion{Major: 17, Minor: 0, Patch: 5},
			b:    GitLabVersion{Major: 17, Minor: 0, Patch: 0},
			want: 1,
		},
		{
			name: "a patch less",
			a:    GitLabVersion{Major: 17, Minor: 0, Patch: 0},
			b:    GitLabVersion{Major: 17, Minor: 0, Patch: 1},
			want: -1,
		},
		{
			name: "patch ignored when major differs",
			a:    GitLabVersion{Major: 18, Minor: 0, Patch: 0},
			b:    GitLabVersion{Major: 17, Minor: 99, Patch: 99},
			want: 1,
		},
		{
			name: "raw field is ignored",
			a:    GitLabVersion{Major: 17, Minor: 0, Raw: "alpha"},
			b:    GitLabVersion{Major: 17, Minor: 0, Raw: "beta"},
			want: 0,
		},
		{
			name: "comparison against parsed literal",
			a:    mustParse(t, "17.4.2-ee"),
			b:    mustParse(t, "17.0"),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.a, tt.b); got != tt.want {
				t.Fatalf("Compare(%+v, %+v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
			// Compare must be antisymmetric: swapping arguments negates
			// the result.
			if got, want := Compare(tt.b, tt.a), -tt.want; got != want {
				t.Fatalf("Compare(%+v, %+v) = %d, want %d (antisymmetry)", tt.b, tt.a, got, want)
			}
		})
	}
}

func mustParse(t *testing.T, s string) GitLabVersion {
	t.Helper()
	v, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q): %v", s, err)
	}
	return v
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name string
		v    GitLabVersion
		want bool
	}{
		{"zero value", GitLabVersion{}, true},
		{"only major set", GitLabVersion{Major: 17}, false},
		{"only minor set", GitLabVersion{Minor: 1}, false},
		{"only patch set", GitLabVersion{Patch: 1}, false},
		{"only raw set", GitLabVersion{Raw: "17.0"}, false},
		{"fully populated", GitLabVersion{Major: 17, Minor: 0, Raw: "17.0"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsZero(); got != tt.want {
				t.Fatalf("%+v.IsZero() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}
