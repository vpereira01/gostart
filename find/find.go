package find

import (
	"fmt"
	"math"
	"math/big"
)

// Find 2 prime factors using the primes' subtraction as a faster auxiliary calculation
// base equations (maxima notation):
// solve([p1*p2 = numberN, p2-p1=subEst], [p1, p2]); (only positive solutions)
//   p1(subEst) :=  (sqrt(subEst^2  + 4*numberN) - subEst)/2;
//   p2(subEst) :=  (sqrt(subEst^2  + 4*numberN) + subEst)/2;
func FindPrimesSub(numberN *big.Int) *big.Int {
	// In order to avoid big float precision errors lets assume a default precision double the number n
	//   this is required because big floats initialize a precision and if a mul() requires a bigger precision
	//   to store its results the precision is not increased but the result is rounded to meet the precision.
	bigFloatsDefaultPrecision := uint(numberN.BitLen() * 2)

	// Calculate 4*numberN since it's reused a lot
	fourTimesNumberN := new(big.Int)
	fourTimesNumberN.SetInt64(4)
	fourTimesNumberN.Mul(fourTimesNumberN, numberN)

	// We want to find perfect squares made by 4*numberN + subEst^2
	// so we first start by doing the startNum = math.ceil(math.sqrt(4*numberN))
	//   since there too big errors with calculating this we use math.floor instead
	//   also because the sqrt(subEst^2  + 4 * numberN) must be even, because
	//     a) subEst is the result of biggerOddNumber-smallOddNumber=evenNumber
	//     b) 4*numberN is even because evenNumber*oddNumber=evenNumber
	//     c) the sqrt(evenNumber)=evenNumber
	startNumBigFloat := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	fourTimesNumberNBigFloat := new(big.Float).SetInt(fourTimesNumberN)
	startNumBigFloat.Sqrt(fourTimesNumberNBigFloat)
	roundDownToPair(startNumBigFloat)
	startNum := new(big.Int)
	startNumBigFloat.Int(startNum)

	twoBigInt := new(big.Int).SetInt64(2)

	return tryFindPrimesSub(fourTimesNumberN, startNum, twoBigInt, bigFloatsDefaultPrecision)
}

// Tries to find primes subtraction for a given startNum and step
func tryFindPrimesSub(fourTimesNumberN *big.Int, initialStartNum *big.Int, step *big.Int, bigFloatsPrecision uint) *big.Int {
	value1Mod9BigInt := new(big.Int)
	nineBigInt := new(big.Int).SetInt64(9)

	value1Mod100BigInt := new(big.Int)
	hundredBigInt := new(big.Int).SetInt64(100)

	countSteps := uint(0)
	countPerfectSquareEstimated := uint(0)

	startNum := new(big.Int).Set(initialStartNum)
	value1 := new(big.Int)
	subEstBigFloat := new(big.Float).SetPrec(bigFloatsPrecision)

	for {
		value1.Mul(startNum, startNum)       // startNum**2
		value1.Sub(value1, fourTimesNumberN) // startNum**2 - 4*numberN

		// due to starting value calculations and rounding, initial values can be wrong
		if value1.Sign() == -1 {
			startNum.Add(startNum, step)
			continue
		}

		// Check if value1 is a perfect square using shortcuts because sqrt() calculation is costly
		// reference https://mathworld.wolfram.com/SquareNumber.html
		value1Mod100BigInt.Mod(value1, hundredBigInt)
		value1Mod100 := value1Mod100BigInt.Int64()
		value1Mod9BigInt.Mod(value1, nineBigInt)
		value1Mod9 := value1Mod9BigInt.Int64()

		if (value1Mod9 == 1 || value1Mod9 == 4 || value1Mod9 == 7 || value1Mod9 == 0) &&
			(value1Mod100 == 0 || value1Mod100 == 1 || value1Mod100 == 4 || value1Mod100 == 9 ||
				value1Mod100 == 16 || value1Mod100 == 21 || value1Mod100 == 24 || value1Mod100 == 25 ||
				value1Mod100 == 29 || value1Mod100 == 36 || value1Mod100 == 41 || value1Mod100 == 44 ||
				value1Mod100 == 49 || value1Mod100 == 56 || value1Mod100 == 61 || value1Mod100 == 64 ||
				value1Mod100 == 69 || value1Mod100 == 76 || value1Mod100 == 81 || value1Mod100 == 84 ||
				value1Mod100 == 89 || value1Mod100 == 96) {

			countPerfectSquareEstimated++

			subEstBigFloat.SetInt(value1)
			subEstBigFloat.Sqrt(subEstBigFloat) // sqrt(startNum**2 - 4*numberN)

			if subEstBigFloat.IsInt() {
				subEst := new(big.Int)
				subEstBigFloat.Int(subEst)
				fmt.Printf("-----> subEst: %10v countSteps: %10v countPerfectSquareEstimated: %10v\n", subEst, countSteps, countPerfectSquareEstimated)
				return subEst
			}
		}

		startNum.Add(startNum, step)
		countSteps++

		if countSteps >= math.MaxUint32 {
			panic("failed to find result")
		}
	}
}

// Sets bigFloat to its floor value and then subtract one if odd
// warning: changes precision
func roundDownToPair(bigFloat *big.Float) {
	// if already rounded, it's an int, nothing to do
	if bigFloat.IsInt() {
		return
	}

	tempBigInt := new(big.Int)

	bigFloat.Int(tempBigInt)
	oneBigInt := new(big.Int).SetInt64(1)
	twoBigInt := new(big.Int).SetInt64(2)

	modBigInt := new(big.Int)
	modBigInt.Mod(tempBigInt, twoBigInt)
	if modBigInt.Int64() == 1 {
		tempBigInt.Sub(tempBigInt, oneBigInt)
	}
	bigFloat.SetInt(tempBigInt)
}
