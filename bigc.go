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

type BigC struct {
	Real *big.Rat
	Imag *big.Rat
}

func (a *BigC) Add(b *BigC) *BigC {
	a.Real.Add(a.Real, b.Real)
	a.Imag.Add(a.Imag, b.Imag)
	return a
}

func (a *BigC) Sub(b *BigC) *BigC {
	a.Real.Sub(a.Real, b.Real)
	a.Imag.Sub(a.Imag, b.Imag)
	return a
}

func (a *BigC) Mul(b *BigC) *BigC {
	imag_temp := new(big.Rat).Set(a.Real)

	a.Real.Mul(a.Real, b.Real)
	real_temp := new(big.Rat).Set(a.Imag)
	real_temp.Mul(real_temp, b.Imag)
	a.Real.Sub(a.Real, real_temp)

	a.Imag.Mul(a.Imag, b.Real)
	imag_temp.Mul(imag_temp, b.Imag)
	a.Imag.Add(a.Imag, imag_temp)
	return a
}

func (a *BigC) Inv() *BigC {
	temp := new(big.Rat).Set(a.Real)
	temp.Mul(temp, temp)
	temp2 := new(big.Rat).Set(a.Imag)
	temp2.Mul(temp2, temp2)
	temp.Add(temp, temp2)

	a.Real.Quo(a.Real, temp)
	a.Imag.Quo(a.Imag, temp)
	a.Imag.Neg(a.Imag)

	return a
}

func (a *BigC) Set(b *BigC) *BigC {
	a = &BigC{
		Real: b.Real,
		Imag: b.Imag,
	}
	return a
}

func (a *BigC) Quo(b *BigC) *BigC {
	temp := new(BigC).Set(b)
	temp.Inv()
	a.Mul(temp)
	return a
}

func (a *BigC) Neg() *BigC {
	a.Real.Neg(a.Real)
	a.Imag.Neg(a.Imag)
	return a
}

func Real(a *BigC) *BigC {
	return &BigC{
		Real: new(big.Rat).Set(a.Real),
		Imag: big.NewRat(0, 1),
	}
}

func Imag(a *BigC) *BigC {
	return &BigC{
		Real: big.NewRat(0, 1),
		Imag: new(big.Rat).Set(a.Imag),
	}
}

func (a *BigC) ToString() string {
	if a.Real.Sign() == 0 && a.Imag.Sign() == 0 {
		return "0"
	}
	i_sign_char := ""
	i_str := ""
	if a.Imag.Sign() == 1 {
		i_sign_char = "+"
	}
	denom := fmt.Sprintf("/%s", a.Imag.Denom().String())
	if a.Imag.Denom().Cmp(big.NewInt(1)) == 0 {
		denom = ""
	}
	if a.Imag.Num().Cmp(big.NewInt(1)) == 0 {
		i_str = fmt.Sprintf("i%s", denom)
	} else if a.Imag.Num().Cmp(big.NewInt(-1)) == 0 {
		i_str = fmt.Sprintf("-i%s", denom)
	} else {
		i_str = fmt.Sprintf("%si%s", a.Imag.Num().String(), denom)
	}
	if a.Real.Sign() == 0 {
		return i_str
	}
	return fmt.Sprintf("%s%s%s", a.Real.RatString(), i_sign_char, i_str)
}

func NewBigCFromString(expr string) (*BigC, error) {
	no_w := strings.Join(strings.Fields(strings.TrimSpace(expr)), "")
	ast, err := parser.ParseExpr(no_w)
	if err != nil {
		return nil, err
	}
	a, err := parseExpr(ast)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func parseExpr(ex interface{}) (*BigC, error) {
	switch node := ex.(type) {
	case *ast.UnaryExpr:
		x, err := parseExpr(node.X)
		if err != nil {
			return nil, err
		}
		switch node.Op {
		case token.ADD:
			return x, nil
		case token.SUB:
			x.Neg()
			return x, nil
		}
		return nil, errors.New("unexpected operator.")
	case *ast.BinaryExpr:
		x, e1 := parseExpr(node.X)
		if e1 != nil {
			return nil, e1
		}
		y, e2 := parseExpr(node.Y)
		if e2 != nil {
			return nil, e2
		}
		switch node.Op {
		case token.ADD:
			x.Add(y)
			return x, nil
		case token.SUB:
			x.Sub(y)
			return x, nil
		case token.MUL:
			x.Mul(y)
			return x, nil
		case token.QUO:
			x.Quo(y)
			return x, nil
		}
		return nil, errors.New("unexpected operator.")
	case *ast.BasicLit:
		switch node.Kind {
		case token.INT:
			num, err := strconv.ParseInt(node.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			return &BigC{
				Real: big.NewRat(num, 1),
				Imag: big.NewRat(0, 1),
			}, nil
		case token.FLOAT:
			num, err := strconv.ParseFloat(node.Value, 64)
			if err != nil {
				return nil, err
			}
			var res big.Rat
			res.SetFloat64(num)
			return &BigC{
				Real: &res,
				Imag: big.NewRat(0, 1),
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
				Real: big.NewRat(0, 1),
				Imag: &res,
			}, nil
		}
	case *ast.Ident:
		if node.Name == "i" {
			return &BigC{
				Real: big.NewRat(0, 1),
				Imag: big.NewRat(1, 1),
			}, nil
		}
		return nil, errors.New("parse error.")
	case *ast.ParenExpr:
		i, err := parseExpr(node.X)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return nil, errors.New("parse error.")
}
