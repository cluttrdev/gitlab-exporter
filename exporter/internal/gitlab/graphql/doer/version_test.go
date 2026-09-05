package doer

import (
	"testing"
)

func TestParseServerVersion(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantMajor int
		wantMinor int
		wantPatch int
		wantErr   bool
	}{
		{"empty", "", 0, 0, 0, true},
		{"garbage", "garbage", 0, 0, 0, true},
		{"major only", "17", 17, 0, 0, false},
		{"major.minor", "17.4", 17, 4, 0, false},
		{"major.minor.patch", "17.4.2", 17, 4, 2, false},
		{"ee suffix", "17.4.2-ee", 17, 4, 2, false},
		{"pre suffix", "16.10.0-pre", 16, 10, 0, false},
		{"rc suffix", "2.0.0-rc1", 2, 0, 0, false},
		{"major zero-minor", "17.0", 17, 0, 0, false},
		{"only suffix patch", "17.4.-pre", 17, 4, 0, false},
		{"negative not supported", "-1.2.3", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseServerVersion(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseServerVersion(%q) error=%v wantErr=%v", tt.in, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if v.Major != tt.wantMajor || v.Minor != tt.wantMinor || v.Patch != tt.wantPatch {
				t.Errorf("ParseServerVersion(%q) = {%d,%d,%d}, want {%d,%d,%d}",
					tt.in, v.Major, v.Minor, v.Patch, tt.wantMajor, tt.wantMinor, tt.wantPatch)
			}
		})
	}
}

func TestGateAllows(t *testing.T) {
	tests := []struct {
		name    string
		server  string
		minV    string
		maxV    string
		want    bool
		wantErr bool
	}{
		{"no bounds", "2.1.0", "", "", true, false},
		{"min satisfied, no max", "2.1.0", "2.0.0", "", true, false},
		{"min unsatisfied", "2.1.0", "2.2.0", "", false, false},
		{"both bounds, in range", "2.1.0", "2.0.0", "3.0.0", true, false},
		{"both bounds, below upper", "2.1.0", "2.0.0", "2.0.5", false, false},

		// inclusive boundaries
		{"equal min, inclusive", "2.0.0", "2.0.0", "", true, false},
		{"equal max, inclusive", "2.0.0", "", "2.0.0", true, false},
		{"just below max", "2.0.0", "", "2.0.1", true, false},

		// pre-release pinning: -rc1 is ignored for comparison
		{"rc1 against same base", "2.0.0-rc1", "2.0.0", "", true, false},
		{"rc1 against same base with maxV", "2.0.0-rc1", "", "2.0.0", true, false},

		// EE / pre suffixes ignored
		{"ee in range", "17.4.2-ee", "17.0", "18.0", true, false},
		{"pre below min", "15.10.0-pre", "16.0", "", false, false},

		// only-major / only-minor in bounds
		{"min specified as major only", "17.4.2", "17", "", true, false},
		{"min specified as major.minor only", "17.4.2", "17.0", "", true, false},
		{"max specified as major only, exceeded", "18.0.0", "", "17", false, false},

		// malformed bounds
		{"malformed min", "2.1.0", "not-a-version", "", false, true},
		{"malformed max", "2.1.0", "", "not-a-version", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGate(tt.server)
			if err != nil {
				t.Fatalf("NewGate(%q) error: %v", tt.server, err)
			}
			got, err := g.Allows(tt.minV, tt.maxV)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Allows(%q, %q) error=%v wantErr=%v", tt.minV, tt.maxV, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Allows(server=%q, min=%q, max=%q) = %v, want %v",
					tt.server, tt.minV, tt.maxV, got, tt.want)
			}
		})
	}
}

func TestNewGateMalformedServer(t *testing.T) {
	if _, err := NewGate(""); err == nil {
		t.Errorf("NewGate(\"\") should error")
	}
	if _, err := NewGate("garbage"); err == nil {
		t.Errorf("NewGate(\"garbage\") should error")
	}
}

func TestGateNilReceiver(t *testing.T) {
	var g *Gate
	if _, err := g.Allows("", ""); err == nil {
		t.Errorf("nil gate Allows should error")
	}
}
