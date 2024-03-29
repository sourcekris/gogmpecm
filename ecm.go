// Package ecm wraps the GMP-ECM library.
package ecm

/*
#cgo darwin CPPFLAGS: -I/opt/homebrew/include/
#cgo darwin LDFLAGS: -L/opt/homebrew/lib/
#cgo LDFLAGS: -lgmp -lecm -lm
#include <gmp.h>
#include <ecm.h>
#include <stdlib.h>

// Macros

void ecm_set_sigma(ecm_params q, mpz_t sig) {
	q->sigma[0] = *sig;
}

mpz_t *ecm_get_sigma(ecm_params q) {
	return &q->sigma;
}

*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Mpz type is a arbitrary precision integer from the GMP library.
type Mpz struct {
	i    C.mpz_t
	init bool
}

// Params stores the ECM parameters.
type Params struct {
	i    C.ecm_params
	init bool
}

// B1 encodes an optimal value for B1 per bit length.
type B1 struct {
	Bits   int
	B1     uint64
	Curves int
}

var (
	// Zero is a Mpz with value 0.
	Zero = NewMpz(0)
	// One is a Mpz with value 1.
	One = NewMpz(1)
)

// OptimalB1s encodes the optimal B1 params from
// https://members.loria.fr/PZimmermann/records/ecm/params.html
var OptimalB1s = []B1{
	{30, 1358, 2},
	{35, 1270, 5},
	{40, 1629, 10},
	{45, 4537, 10},
	{50, 12322, 9},
	{55, 12820, 18},
	{60, 21905, 21},
	{65, 24433, 41},
	{70, 32918, 66},
	{75, 64703, 71},
	{80, 76620, 119},
	{85, 155247, 123},
	{90, 183849, 219},
	{95, 245335, 321},
	{100, 445657, 339},
	{105, 643986, 468},
	{110, 1305195, 439},
	{115, 1305195, 818},
	{120, 3071166, 649},
	{125, 3784867, 949},
	{130, 4572523, 1507},
	{135, 7982718, 1497},
	{140, 9267681, 2399},
	{145, 22025673, 1826},
	{150, 22025673, 3159},
	{155, 26345943, 4532},
	{160, 35158748, 6076},
	{165, 46919468, 8177},
	{170, 47862548, 14038},
	{175, 153319098, 7166},
	{180, 153319098, 12017},
	{185, 188949210, 16238},
	{190, 410593604, 13174},
	{195, 496041799, 17798},
	{200, 491130495, 29584},
	{205, 1067244762, 23626},
	{210, 1056677983, 38609},
	{215, 1328416470, 49784},
	{220, 1315263832, 81950},
	{225, 2858117139, 63461},
}

// mpzFinalize releases the memory allocated to the Mpz.
func mpzFinalize(z *Mpz) {
	if z.init {
		runtime.SetFinalizer(z, nil)
		C.mpz_clear(&z.i[0])
		z.init = false
	}
}

// mpzDoinit initializes an Mpz type.
func (z *Mpz) mpzDoinit() {
	if z.init {
		return
	}
	z.init = true
	C.mpz_init(&z.i[0])
	runtime.SetFinalizer(z, mpzFinalize)
}

// NewMpz allocates and returns a new Fmpz set to x.
func NewMpz(x int64) *Mpz {
	return new(Mpz).SetInt64(x)
}

// string returns z in the base given
func (z *Mpz) string(base int) string {
	if z == nil {
		return "<nil>"
	}
	z.mpzDoinit()
	p := C.mpz_get_str(nil, C.int(base), &z.i[0])
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s
}

// String returns a string value of z in base 10.
func (z *Mpz) String() string {
	return z.string(10)
}

// NewParams allocates and returns a new ECM Parameters struct.
func NewParams() *Params {
	p := new(Params)
	C.ecm_init(&p.i[0])
	return p
}

// SetSigma sets the value of the ecm_params->sigma to sigma.
func (p *Params) SetSigma(sigma *Mpz) {
	C.ecm_set_sigma(&p.i[0], &sigma.i[0])
}

// GetSigma gets the value of the ecm_params->sigma.
func (p *Params) GetSigma() *Mpz {
	z := new(Mpz)
	z.mpzDoinit()
	z.i = *C.ecm_get_sigma(&p.i[0])
	return z
}

// SetInt64 sets z to x and returns z.
func (z *Mpz) SetInt64(x int64) *Mpz {
	z.mpzDoinit()
	y := C.long(x)
	C.mpz_set_si(&z.i[0], y)
	return z
}

// SetString sets z to the value of s, interpreted in the given base,
// and returns z and a boolean indicating success. If SetString fails,
// the value of z is undefined but the returned value is nil.
//
// The base argument must be 0 or a value from 2 through MaxBase. If the base
// is 0, the string prefix determines the actual conversion base. A prefix of
// ``0x'' or ``0X'' selects base 16; the ``0'' prefix selects base 8, and a
// ``0b'' or ``0B'' prefix selects base 2. Otherwise the selected base is 10.
//
func (z *Mpz) SetString(s string, base int) (*Mpz, bool) {
	z.mpzDoinit()
	if base != 0 && (base < 2 || base > 36) {
		return nil, false
	}
	// Skip leading + as mpz_set_str doesn't understand them
	if len(s) > 1 && s[0] == '+' {
		s = s[1:]
	}
	// mpz_set_str incorrectly parses "0x" and "0b" as valid
	if base == 0 && len(s) == 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X' || s[1] == 'b' || s[1] == 'B') {
		return nil, false
	}
	p := C.CString(s)
	defer C.free(unsafe.Pointer(p))
	if C.mpz_set_str(&z.i[0], p, C.int(base)) < 0 {
		return nil, false
	}
	return z, true // err == io.EOF => scan consumed all of s
}

// Cmp compares Mpz z and y and returns:
//   -1 if z <  y
//    0 if z == y
//   +1 if z >  y
func (z *Mpz) Cmp(y *Mpz) int {
	z.mpzDoinit()
	y.mpzDoinit()
	r := int(C.mpz_cmp(&z.i[0], &y.i[0]))
	if r < 0 {
		r = -1
	} else if r > 0 {
		r = 1
	}
	return r
}

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
//
func (z *Mpz) Sign() int {
	z.mpzDoinit()
	return z.Cmp(Zero)
}

// BitLen returns the length of the absolute value of z in bits.
// The bit length of 0 is 0.
func (z *Mpz) BitLen() int {
	z.mpzDoinit()
	if z.Sign() == 0 {
		return 0
	}
	return int(C.mpz_sizeinbase(&z.i[0], 2))
}

// OptimalB1 returns the best B1 struct given the bit length of z.
func (z *Mpz) OptimalB1() B1 {
	l := z.BitLen() / 2
	for _, b := range OptimalB1s {
		if b.Bits > l {
			return b
		}
	}

	return OptimalB1s[len(OptimalB1s)-1]
}

// OptimalB1Uint64 returns the best value of B1 given the bit length of z.
func (z *Mpz) OptimalB1Uint64() uint64 {
	l := z.BitLen() / 2
	for _, b := range OptimalB1s {
		if b.Bits > l {
			return b.B1
		}
	}

	return OptimalB1s[len(OptimalB1s)-1].B1
}

// Factor will attempt to factor n given the ecm parameters p and returns the factor if successful
// or an error. It will attempt to find its own optimal value for B1.
func (p *Params) Factor(n *Mpz) (*Mpz, error) {
	fac := NewMpz(0)
	res := int(C.ecm_factor(&fac.i[0], &n.i[0], C.double(n.OptimalB1Uint64()), &p.i[0]))
	if res > 0 {
		return fac, nil
	}
	return nil, fmt.Errorf("ecm_factor failed")
}

// FactorGivenB1 will attempt to factor n given the ecm parameters p and returns the factor if
// successful or an error. It will use the users value for B1.
func (p *Params) FactorGivenB1(n *Mpz, b1 uint64) (*Mpz, error) {
	fac := NewMpz(0)
	res := int(C.ecm_factor(&fac.i[0], &n.i[0], C.double(b1), &p.i[0]))
	if res > 0 {
		return fac, nil
	}
	return nil, fmt.Errorf("ecm_factor failed")
}
