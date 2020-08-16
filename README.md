# bigc

BigCは実部、虚部が有理数であるような複素数を扱うためのパッケージです。

## Usage

#### type BigC

    type BigC struct {
    }


BigCは複素数を表します

#### func  NewBigC

    func NewBigC(r *big.Rat, i *big.Rat) *BigC

NewBigC creates a new BigC with real-part r and imaginary-part i.

#### func  ParseString

    func ParseString(expr string) (*BigC, error)

ParseString returns a new BigC instance of the result of the expression expr.
Arithmetic operations and parentheses are supported.

#### func (*BigC) AbsSq

    func (x *BigC) AbsSq() *big.Rat

Abs sets z to square of |x| (the absolute value of x) and returns z.

#### func (*BigC) Add

    func (z *BigC) Add(x *BigC, y *BigC) *BigC

Add sets z to the sum x+y and returns z.

#### func (*BigC) Conj

    func (z *BigC) Conj(x *BigC) *BigC

Conj sets z to the conjugate complex number of x and returns z.

#### func (*BigC) Equal

    func (z *BigC) Equal(x *BigC) bool

Equal reports whether x equals z.

#### func (*BigC) FloatString

    func (x *BigC) FloatString(prec int) string

FloatString returns a string representation of x in decimal form with prec
digits of precision after the radix point. The last digit is rounded to nearest,
with halves rounded away from zero.

#### func (*BigC) Imag

    func (x *BigC) Imag() *big.Rat

Imag returns the imaginary-part of x. The result is a reference to x's
imaginary-part; it may change if a new value is assigned to x, and vice versa.

#### func (*BigC) Inv

    func (z *BigC) Inv(x *BigC) *BigC

Inv sets z to 1/x and returns z. If x == 0, Inv panics.

#### func (*BigC) IsPureImag

    func (x *BigC) IsPureImag() bool

IsPureImag reports whether x is a pure imaginary number.

#### func (*BigC) IsReal

    func (x *BigC) IsReal() bool

IsReal reports whether x is a real number.

#### func (*BigC) Mul

    func (z *BigC) Mul(x *BigC, y *BigC) *BigC

Mul sets z to the product x*y and returns z.

#### func (*BigC) Neg

    func (z *BigC) Neg(x *BigC) *BigC

Neg sets z to -x and returns z.

#### func (*BigC) Quo

    func (z *BigC) Quo(x *BigC, y *BigC) *BigC

Quo sets z to the quotient x/y and returns z. If y == 0, Quo panics.

#### func (*BigC) Real

    func (x *BigC) Real() *big.Rat

Real returns the real-part of x. The result is a reference to x's real-part; it
may change if a new value is assigned to x, and vice versa.

#### func (*BigC) Set

    func (z *BigC) Set(x *BigC) *BigC

Set sets z to x (by making a copy of x) and returns z.

#### func (*BigC) String

    func (x *BigC) String() string

String returns a string exact representation of x.

#### func (*BigC) Sub

    func (z *BigC) Sub(x *BigC, y *BigC) *BigC

Sub sets z to the difference x-y and returns z.
