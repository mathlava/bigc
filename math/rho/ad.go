package rho

import (
	"math/big"

	"github.com/mathlava/bigc"
)

var (
	val0   = big.NewInt(0)
	val1_2 = big.NewRat(1, 2)
)

// numに素微分を行い, numを返します.
func ArithmeticDerivative(num *bigc.BigC) *bigc.BigC {
	if num.IsReal() {
		ad_rat(num.Real())
		return num
	}
	if num.IsPureImag() {
		ad_rat(num.Imag())
		return num
	}
	cache := new(bigc.BigC).Set(num).AbsSq()
	num.Imag().Quo(num.Imag(), cache).Mul(num.Imag(), val1_2)
	num.Real().Quo(num.Real(), cache).Mul(num.Real(), val1_2)
	ad_rat(cache)
	num.Real().Mul(num.Real(), cache)
	num.Imag().Mul(num.Imag(), cache)
	return num
}

func ad_rat(num *big.Rat) *big.Rat {
	sign := num.Sign()
	num.Abs(num)
	ba_ := ad_int(new(big.Int).Set(num.Denom()))
	ba_.Mul(ba_, num.Num())
	ad_int(num.Num()).Mul(num.Num(), num.Denom()).Sub(num.Num(), ba_)
	num.Denom().Mul(num.Denom(), num.Denom())
	if sign == -1 {
		num.Neg(num)
	}
	return num
}

func ad_int(num *big.Int) *big.Int {
	sign := num.Sign()
	num.Abs(num)
	if num.Cmp(val1) == 0 {
		return num.Set(val0)
	}
	if sign == 0 {
		panic("undefined.")
	}
	cache := new(big.Int).Set(num)
	num.Set(val0)
	add := new(big.Int)
	for _, p := range Primes(new(big.Int).Set(cache)) {
		add.Set(cache)
		add.Div(add, p)
		num.Add(num, add)
	}
	if sign == -1 {
		num.Neg(num)
	}
	return num
}
