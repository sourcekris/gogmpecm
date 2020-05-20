# gogmpecm

golang bindings for the gmp-ecm library.

The Elliptic Curve Method for Integer Factorization (ECM)

The project provides limited but useful bindings for the ECM factorization functions from the
GMP-ECM library. GMP-ECM is a highly optimized implementation of Lenstra's elliptic curve 
factorization method. 

See http://ecm.gforge.inria.fr/ for more about GMP-ECM.

## Features

### Types
 * `type Mpz struct` Mpz type is a arbitrary precision integer from the GMP library.
 * `type Params struct` Params stores the ECM parameters.
 * `type B1 struct` B1 encodes an optimal value for B1 per bit length.
### Constructors
 * `NewMpz(x int64) *Mpz` NewMpz allocates and returns a new Fmpz set to x.
 * `NewParams() *Params` NewParams allocates and returns a new ECM Parameters struct.
### Assignments
 * `(z *Mpz) SetInt64(x int64) *Mpz` SetInt64 sets z to x and returns z.
 * `(z *Mpz) SetString(s string, base int) (*Mpz, bool)` SetString sets z to the value of s, 
   interpreted in the given base.
### Comparisons
 * `(z *Mpz) Cmp(y *Mpz) int` Cmp compares Mpz z and y and returns: -1 if x <  0 , 0 if x == 0 ,
   +1 if x >  0
 * `(z *Mpz) Sign() int` Sign returns: -1 if x <  0 , 0 if x == 0 , +1 if x >  0
### Helpers
 * `(z *Mpz) BitLen() int` BitLen returns the length of the absolute value of z in bits.
 * `(z *Mpz) OptimalB1() B1` OptimalB1 returns the best B1 struct given the bit length of z.
 * `(z *Mpz) OptimalB1Uint64() B1` OptimalB1Uint64 returns the best value of B1 given the bit length
   of z.
### Factorization
 * `(p *Params) Factor(n *Mpz) (*Mpz, error)` Factor will attempt to factor n given the ecm 
   parameters p and returns the factor if successful.
 * `(p *Params) FactorGivenB1(n *Mpz, b1 uint64) (*Mpz, error)` FactorGivenB1 will attempt to factor
   n given the ecm parameters p and returns the factor.

## Dependencies

 * Golang
 * GMP and GMP-ECM libraries must be installed. 
  * On modern Linux distributions this is achieved using the package manager. For example on Debian,
    and Ubuntu on Linux and Windows 10 WSL 2.0:

    ```
    $ sudo apt install libecm-dev
    ```

## Installation

 * Use go get to install the package
   ```
   $ go get github.com/sourcekris/gogmpecm
   ```

## Author

 * Kris Hunt (@ctfkris)