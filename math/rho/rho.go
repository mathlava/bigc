package rho

import (
	"context"
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

// RhoAsyncの結果を表します。
type RhoResult struct {
	// 約数
	Divisor *big.Int
	// 処理が成功したかどうか
	Success bool
}

// ポラード・ロー法を試行し, 約数の一つの求値を試みます.
func RhoAsync(ctx context.Context, N *big.Int) chan RhoResult {
	res := make(chan RhoResult)
	go func() {
		d := new(big.Int)
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
	loop:
		for {
			select {
			case <-ctx.Done():
				return
			default:
				count.Sub(count, val1)
				if count.Sign() == -1 {
					break loop
				}
				d.Set(f(f(y)))
				d.Sub(d, f(x))
				d.Abs(d)
				d.Mod(d, N)
				d = Gcd(d, new(big.Int).Set(N))
				if d.Cmp(val1) != 0 {
					if d.Sign() != 0 && d.Cmp(N) != 0 {
						res <- RhoResult{
							Divisor: d,
							Success: true,
						}
						return
					}
					break loop
				}
			}
		}
		res <- RhoResult{
			Divisor: nil,
			Success: false,
		}
	}()
	return res
}

// ポラード・ロー法を試行し, 約数の一つの求値を試みます.
func Rho(N *big.Int) (d *big.Int, success bool) {
	res := <-RhoAsync(context.Background(), N)
	return res.Divisor, res.Success
}

const test_count = 200

// numsを素因数分解します. ※副作用があります
func Primes(num *big.Int) []*big.Int {
	ch := PrimesAsync(context.Background(), num)
	res := <-ch
	return res
}

// numsを素因数分解します. ※副作用があります
func PrimesAsync(ctx context.Context, num *big.Int) chan []*big.Int {
	res := make(chan []*big.Int)
	go func() {
		negative := false
		if num.CmpAbs(val1) == 0 || num.Sign() == 0 {
			res <- []*big.Int{num}
			return
		} else if num.Sign() == -1 {
			num.Neg(num)
			negative = true
		}
		nums := []*big.Int{num}

	loop:
		for {
			count := len(nums) - 1
		search_prime:
			for i := 0; ; i++ {
				if !nums[i].ProbablyPrime(test_count) || nums[0].Cmp(val1) == 0 {
					// 合成数が見つかったら
					// 合成数を先頭にもってくる
					primes := nums[:i]
					nums = nums[i:]
					nums = append(nums, primes...)
					break search_prime
				}
				if i == count {
					break loop
				}
			}
			var rho_res RhoResult
			childCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			select {
			case <-ctx.Done():
				return
			case rho_res = <-RhoAsync(childCtx, nums[0]):
				n := rho_res.Divisor
				if !rho_res.Success {
					n = new(big.Int).Set(val1)
				search_divisor:
					for {
						n.Add(n, val1)
						cache := new(big.Int).Set(nums[0])
						if cache.Mod(cache, n).Sign() == 0 {
							break search_divisor
						}
					}
				}
				nums = append(nums, n)
				nums[0].Div(nums[0], n)
			}
		}
		if negative {
			nums = append(nums, big.NewInt(-1))
		}
		res <- nums
	}()
	return res
}
