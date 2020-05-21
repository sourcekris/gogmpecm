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

func TestOptimalB1(t *testing.T) {
	for _, tc := range []struct {
		name string
		n    *Mpz
		want uint64
	}{
		{
			name: "optimal B1 for 0 is 1358",
			n:    NewMpz(0),
			want: 1358,
		},
		{
			name: "optimal B1 for 9128309128301 is 4537",
			n:    NewMpz(9128309128301),
			want: 4537,
		},
	} {
		got := tc.n.OptimalB1()
		if got.B1 != tc.want {
			t.Errorf("OptimalB1(): %s want / got mismatch: %d / %d", tc.name, tc.want, got.B1)
		}
	}
}

func TestFactor(t *testing.T) {
	for _, tc := range []struct {
		name string
		n    string
		want *Mpz
	}{
		{
			name: "factors easily",
			n:    "187072209578355573530071658587684226515959365500927",
			want: NewMpz(2349023),
		},
	} {
		p := NewParams()
		n, _ := new(Mpz).SetString(tc.n, 10)
		factor, err := p.Factor(n)
		if err != nil {
			t.Errorf("Factor() %s failed: expected no error got error: %v", tc.name, err)
		}

		if factor.Cmp(tc.want) != 0 {
			t.Errorf("Factor(): %s want / got mismatch: %v / %v", tc.name, tc.want, factor)
		}
	}
}
