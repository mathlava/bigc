// Some comments are quoted and modified from the following:
//     https://golang.org/src/math/big/rat.go
//       Copyright 2010 The Go Authors. All rights reserved.
//     https://golang.org/src/math/big/ratconv.go
//       Copyright 2015 The Go Authors. All rights reserved.
// Use of each source code is governed by a BSD-style license.

/*
Go (golang) library implementing rational complex numbers.
BigCは実部、虚部が共に有理数であるような複素数を扱うためのパッケージです。
リポジトリは https://github.com/mathlava/bigc です。
*/
package bigc

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/big"
	"strconv"
	"strings"
)

// A BigC object represents a rational complex number.
// BigCは複素数を表します
type BigC struct {
	re *big.Rat
	im *big.Rat
}

// NewBigC creates a new BigC with real-part r and imaginary-part i.
func NewBigC(r *big.Rat, i *big.Rat) *BigC {
	return &BigC{
		re: r,
		im: i,
	}
}

func (z *BigC) adjust(x *BigC) {
	*z.re = *x.re
	*z.im = *x.im
}

// Abs sets z to square of |x| (the absolute value of x) and returns z.
func (x *BigC) AbsSq() *big.Rat {
	res := new(big.Rat).Set(x.re)
	res.Mul(res, res)
	img := new(big.Rat).Set(x.im)
	return res.Add(res, img.Mul(img, img))
}

// Add sets z to the sum x+y and returns z.
func (z *BigC) Add(x *BigC, y *BigC) *BigC {
	z.adjust(x)
	z.re.Add(z.re, y.re)
	z.im.Add(z.im, y.im)
	return z
}

// Conj sets z to the conjugate complex number of x and returns z.
func (z *BigC) Conj(x *BigC) *BigC {
	z.adjust(x)
	z.im.Neg(z.im)
	return z
}

// Equal reports whether whether x equals z.
func (z *BigC) Equal(x *BigC) bool {
	return z.re.Cmp(x.re) == 0 && z.im.Cmp(x.im) == 0
}

// Imag returns the imaginary-part of x.
// The result is a reference to x's imaginary-part; it
// may change if a new value is assigned to x, and vice versa.
func (x *BigC) Imag() *big.Rat {
	return x.im
}

// Inv sets z to 1/x and returns z.
// If x == 0, Inv panics.
func (z *BigC) Inv(x *BigC) *BigC {
	z.adjust(x)
	denom := z.AbsSq()
	z.re.Quo(z.re, denom)
	z.im.Quo(z.im, denom)
	return z.Conj(z)
}

// IsPureImag reports whether whether x is a pure imaginary number.
func (x *BigC) IsPureImag() bool {
	return !x.IsReal() && x.re.Sign() == 0
}

// IsReal reports whether whether x is a real number.
func (x *BigC) IsReal() bool {
	return x.im.Sign() == 0
}

// Mul sets z to the product x*y and returns z.
func (z *BigC) Mul(x *BigC, y *BigC) *BigC {
	z.adjust(x)
	imag_temp := new(big.Rat).Set(z.re)

	z.re.Mul(z.re, y.re)
	real_temp := new(big.Rat).Set(z.im)
	real_temp.Mul(real_temp, y.im)
	z.re.Sub(z.re, real_temp)

	z.im.Mul(z.im, y.re)
	imag_temp.Mul(imag_temp, y.im)
	z.im.Add(z.im, imag_temp)
	return z
}

// Neg sets z to -x and returns z.
func (z *BigC) Neg(x *BigC) *BigC {
	z.adjust(x)
	z.re.Neg(z.re)
	z.im.Neg(z.im)
	return z
}

// Quo sets z to the quotient x/y and returns z.
// If y == 0, Quo panics.
func (z *BigC) Quo(x *BigC, y *BigC) *BigC {
	z.adjust(x)
	temp := new(BigC).Set(y)
	temp.Inv(temp)
	z.Mul(z, temp)
	return z
}

// Real returns the real-part of x.
// The result is a reference to x's real-part; it
// may change if a new value is assigned to x, and vice versa.
func (x *BigC) Real() *big.Rat {
	return x.re
}

// Set sets z to x (by making a copy of x) and returns z.
func (z *BigC) Set(x *BigC) *BigC {
	z = &BigC{
		re: x.re,
		im: x.im,
	}
	return z
}

// FloatString returns a string representation of x in decimal form with prec
// digits of precision after the radix point. The last digit is rounded to
// nearest, with halves rounded away from zero.
func (x *BigC) FloatString(prec int) string {
	if x.re.Sign() == 0 && x.im.Sign() == 0 {
		return "0"
	}
	if x.re.Sign() == 0 {
		return fmt.Sprintf("%si", x.im.FloatString(prec))
	}
	if x.im.Sign() == 0 {
		return fmt.Sprintf("%s", x.re.FloatString(prec))
	}
	if x.im.Sign() == 1 {
		return fmt.Sprintf("%s+%s", x.re.FloatString(prec), x.im.FloatString(prec))
	}
	return fmt.Sprintf("%s%s", x.re.FloatString(prec), x.im.FloatString(prec))
}

// String returns a string exact representation of x.
func (x *BigC) String() string {
	if x.re.Sign() == 0 && x.im.Sign() == 0 {
		return "0"
	}
	i_sign_char := ""
	i_str := ""
	if x.im.Sign() == 1 {
		i_sign_char = "+"
	}
	denom := fmt.Sprintf("/%s", x.im.Denom().String())
	if x.im.Denom().Cmp(big.NewInt(1)) == 0 {
		denom = ""
	}
	if x.im.Num().Cmp(big.NewInt(1)) == 0 {
		i_str = fmt.Sprintf("i%s", denom)
	} else if x.im.Num().Cmp(big.NewInt(-1)) == 0 {
		i_str = fmt.Sprintf("-i%s", denom)
	} else if x.im.Sign() == 0 {
		i_str = ""
	} else {
		i_str = fmt.Sprintf("%si%s", x.im.Num().String(), denom)
	}
	if x.re.Sign() == 0 {
		return i_str
	}
	return fmt.Sprintf("%s%s%s", x.re.RatString(), i_sign_char, i_str)
}

// Sub sets z to the difference x-y and returns z.
func (z *BigC) Sub(x *BigC, y *BigC) *BigC {
	z.adjust(x)
	z.re.Sub(z.re, y.re)
	z.im.Sub(z.im, y.im)
	return z
}

// ParseString returns a new BigC instance of the result of the expression expr.
// Arithmetic operations and parentheses are supported.
func ParseString(expr string) (*BigC, error) {
	no_w := strings.Join(strings.Fields(strings.TrimSpace(expr)), "")
	ast, err := parser.ParseExpr(no_w)
	if err != nil {
		return nil, err
	}
	a, err := walk(ast)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func walk(ex interface{}) (*BigC, error) {
	switch node := ex.(type) {
	case *ast.BinaryExpr:
		x, e1 := walk(node.X)
		if e1 != nil {
			return nil, e1
		}
		y, e2 := walk(node.Y)
		if e2 != nil {
			return nil, e2
		}
		switch node.Op {
		case token.ADD:
			x.Add(x, y)
			return x, nil
		case token.SUB:
			x.Sub(x, y)
			return x, nil
		case token.MUL:
			x.Mul(x, y)
			return x, nil
		case token.QUO:
			x.Quo(x, y)
			return x, nil
		}
		return nil, errors.New("unexpected operator.")
	case *ast.UnaryExpr:
		x, err := walk(node.X)
		if err != nil {
			return nil, err
		}
		switch node.Op {
		case token.ADD:
			return x, nil
		case token.SUB:
			x.Neg(x)
			return x, nil
		}
		return nil, errors.New("unexpected operator.")
	case *ast.ParenExpr:
		i, err := walk(node.X)
		if err != nil {
			return nil, err
		}
		return i, nil
	case *ast.BasicLit:
		switch node.Kind {
		case token.INT:
			num, err := strconv.ParseInt(node.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			return &BigC{
				re: big.NewRat(num, 1),
				im: big.NewRat(0, 1),
			}, nil
		case token.FLOAT:
			num, err := strconv.ParseFloat(node.Value, 64)
			if err != nil {
				return nil, err
			}
			var res big.Rat
			res.SetFloat64(num)
			return &BigC{
				re: &res,
				im: big.NewRat(0, 1),
			}, nil
		case token.IMAG:
			if node.Value[len(node.Value)-1] != 'i' {
				return nil, errors.New("unknown error.")
			}
			num, err := strconv.ParseFloat(node.Value[:len(node.Value)-1], 64)
			if err != nil {
				return nil, err
			}
			var res big.Rat
			res.SetFloat64(num)
			return &BigC{
				re: big.NewRat(0, 1),
				im: &res,
			}, nil
		}
	case *ast.Ident:
		if node.Name == "i" {
			return &BigC{
				re: big.NewRat(0, 1),
				im: big.NewRat(1, 1),
			}, nil
		}
		return nil, errors.New("unexpected identier.")
	}
	return nil, errors.New("parse error.")
}
