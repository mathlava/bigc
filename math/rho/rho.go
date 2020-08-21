package rho

import (
	"math/big"
)

// numsの最大公約数を返します. ※副作用があります
func Gcd(nums ...*big.Int) *big.Int {
	switch len(nums) {
	case 2:
		if nums[0].Cmp(nums[1]) == -1 {
			n := nums[1]
			nums[1] = nums[0]
			nums[0] = n
		}
		if nums[1].Sign() == 0 {
			return nums[0]
		}
		return Gcd(nums[1], nums[0].Mod(nums[0], nums[1]))
	case 1:
		return nums[0]
	default:
		g := Gcd(nums[0], nums[1])
		nums = nums[1:]
		nums[0] = g
		return Gcd(nums...)
	}
}

var val1 = big.NewInt(1)

// ポラード・ロー法を試行し, 約数の一つの求値を試みます.
func Rho(N *big.Int) (d *big.Int, success bool) {
	d = new(big.Int)
	f := func(a *big.Int) *big.Int {
		a.Mul(a, a)
		a.Add(a, val1)
		a.Mod(a, N)
		return a
	}
	x := big.NewInt(57)
	x.Mod(x, N)
	y := new(big.Int).Set(x)
	f(y)
	count := new(big.Int).Set(N)
	count.Sqrt(count)
	count.Sqrt(count)
	for {
		count.Sub(count, val1)
		if count.Sign() == -1 {
			break
		}
		d.Set(f(f(y)))
		d.Sub(d, f(x))
		d.Abs(d)
		d.Mod(d, N)
		d = Gcd(d, new(big.Int).Set(N))
		if d.Cmp(val1) != 0 {
			if d.Sign() != 0 && d.Cmp(N) != 0 {
				return d, true
			}
			break
		}
	}
	return nil, false
}

const test_count = 200

// numsを素因数分解します. ※副作用があります
func Primes(nums ...*big.Int) []*big.Int {
	if len(nums) == 1 {
		if nums[0].Cmp(val1) == 0 || nums[0].Sign() == 0 {
			return nums
		} else if nums[0].Sign() == -1 {
			nums[0].Neg(nums[0])
			nums = Primes(nums...)
			return append(nums, big.NewInt(-1))
		}
	}
	count := len(nums) - 1
	for i := 0; ; i++ {
		if !nums[i].ProbablyPrime(test_count) || nums[0].Cmp(val1) == 0 {
			primes := nums[:i]
			nums = nums[i:]
			nums = append(nums, primes...)
			break
		}
		if i == count {
			return nums
		}
	}
	n := new(big.Int)
	var ok bool
	n, ok = Rho(nums[0])
	if !ok {
		n = new(big.Int).Set(val1)
		for {
			n.Add(n, val1)
			cache := new(big.Int).Set(nums[0])
			if cache.Mod(cache, n).Sign() == 0 {
				break
			}
		}
	}
	nums = append(nums, n)
	nums[0].Div(nums[0], n)
	return Primes(nums...)
}
