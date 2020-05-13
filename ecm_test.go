package ecm

import "testing"

func TestCmp(t *testing.T) {
	a := NewMpz(1024)
	b := NewMpz(1024)
	c := NewMpz(2048)

	if a.Cmp(b) != 0 {
		t.Error("Cmp(): expected a.Cmp(b) == 0")
	}

	if a.Cmp(c) > 0 {
		t.Error("Cmnp(): expected a.Cmp(c) < 0")
	}

	if c.Cmp(a) < 0 {
		t.Error("Cmp(): expected c.Cmp(a) > 0")
	}
}

func TestSign(t *testing.T) {
	a := NewMpz(1024)
	b := NewMpz(-1024)
	c := NewMpz(0)

	if c.Sign() != 0 {
		t.Error("Sign(): expected Sign of 0 to be 0")
	}

	if a.Sign() <= 0 {
		t.Error("Sign(): expected Sign of 1024 to be positive")
	}

	if b.Sign() >= 0 {
		t.Error("Sign(): expected Sign of -1024 to be negative")
	}
}

func TestBitLen(t *testing.T) {
	for _, tc := range []struct {
		name string
		n    *Mpz
		want int
	}{
		{
			name: "bitlen of 0 is 0",
			n:    NewMpz(0),
			want: 0,
		},
		{
			name: "bitlen of 1024 is 11",
			n:    NewMpz(1024),
			want: 11,
		},
	} {
		got := tc.n.BitLen()
		if got != tc.want {
			t.Errorf("BitLen(): %s want / got mismatch: %d / %d", tc.name, tc.want, got)
		}
	}
}
