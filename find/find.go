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
	roundDownToPairAndNineMultiple(startNumBigFloat)
	startNum := new(big.Int)
	startNumBigFloat.Int(startNum)

	// Let's try to speed up our guess my only generating startNum guesses that produce numbers expected to the perfect squares
	// We can do this because given fourTimesNumberN % 9 we can now if startNum guesses % 9 values
	// if fourTimesNumberN % 9 = 0, startNum % 9 must be 0,1,2,3,4,5,6,7,8
	// if fourTimesNumberN % 9 = 1, startNum % 9 must be 1,8
	// if fourTimesNumberN % 9 = 2, startNum % 9 must be 0,3,6
	// if fourTimesNumberN % 9 = 3, startNum % 9 must be 1,2,4,5,7,8
	// if fourTimesNumberN % 9 = 4, startNum % 9 must be 2,7
	// if fourTimesNumberN % 9 = 5, startNum % 9 must be 0,3,6
	// if fourTimesNumberN % 9 = 6, startNum % 9 must be 1,2,4,5,7,8
	// if fourTimesNumberN % 9 = 7, startNum % 9 must be 4,5
	// if fourTimesNumberN % 9 = 8, startNum % 9 must be 0,3,6
	// "Calculated" using https://www.wolframalpha.com/input/?i=solve+%7B+%28x%5E2-4*n%29+mod+9+%3D+0%2C+4*n+mod+9+%3D+8%2C+x%3Cn%2C+x%3E0%2C+n%3E0+%7D+for+x
	nineModsAndStartIncrements := make(map[uint][]uint)
	nineModsAndStartIncrements[0] = []uint{0, 10, 2, 12, 4, 14, 6, 16, 8}
	nineModsAndStartIncrements[1] = []uint{10, 8}
	nineModsAndStartIncrements[2] = []uint{0, 12, 6}
	nineModsAndStartIncrements[3] = []uint{10, 2, 4, 14, 16, 8}
	nineModsAndStartIncrements[4] = []uint{2, 16}
	nineModsAndStartIncrements[5] = []uint{0, 12, 6}
	nineModsAndStartIncrements[6] = []uint{10, 2, 4, 14, 16, 8}
	nineModsAndStartIncrements[7] = []uint{4, 14}
	nineModsAndStartIncrements[8] = []uint{0, 12, 6}

	nineBigInt := new(big.Int).SetInt64(9)
	fourTimesNumberNMod9BigInt := new(big.Int)
	fourTimesNumberNMod9BigInt.Mod(fourTimesNumberN, nineBigInt)
	fourTimesNumberNMod9 := uint(fourTimesNumberNMod9BigInt.Uint64())
	eighteenBigInt := new(big.Int).SetInt64(18)

	results := make(chan *big.Int)
	stopFlag := false

	for _, startIncrement := range nineModsAndStartIncrements[fourTimesNumberNMod9] {
		startNumAdjusted := new(big.Int)
		startNumAdjusted.SetUint64(uint64(startIncrement))
		startNumAdjusted.Add(startNumAdjusted, startNum)
		go tryFindPrimesSub(fourTimesNumberN, startNumAdjusted, eighteenBigInt, bigFloatsDefaultPrecision, results, &stopFlag)
	}

	result := <-results
	stopFlag = true
	return result
}

// Tries to find primes subtraction for a given startNum and step
func tryFindPrimesSub(fourTimesNumberN *big.Int, initialStartNum *big.Int, step *big.Int, bigFloatsPrecision uint, results chan *big.Int, stopFlag *bool) {
	fmt.Printf("fourTimesNumberN: %v initialStartNum: %v step: %v\n", fourTimesNumberN, initialStartNum, step)

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

		if value1Mod100 == 0 || value1Mod100 == 1 || value1Mod100 == 4 || value1Mod100 == 9 ||
			value1Mod100 == 16 || value1Mod100 == 21 || value1Mod100 == 24 || value1Mod100 == 25 ||
			value1Mod100 == 29 || value1Mod100 == 36 || value1Mod100 == 41 || value1Mod100 == 44 ||
			value1Mod100 == 49 || value1Mod100 == 56 || value1Mod100 == 61 || value1Mod100 == 64 ||
			value1Mod100 == 69 || value1Mod100 == 76 || value1Mod100 == 81 || value1Mod100 == 84 ||
			value1Mod100 == 89 || value1Mod100 == 96 {

			countPerfectSquareEstimated++

			subEstBigFloat.SetInt(value1)
			subEstBigFloat.Sqrt(subEstBigFloat) // sqrt(startNum**2 - 4*numberN)

			if subEstBigFloat.IsInt() {
				subEst := new(big.Int)
				subEstBigFloat.Int(subEst)
				fmt.Printf("-----> subEst: %10v countSteps: %10v countPerfectSquareEstimated: %10v\n", subEst, countSteps, countPerfectSquareEstimated)
				results <- subEst
			}
		}

		startNum.Add(startNum, step)
		countSteps++

		if countSteps >= math.MaxUint32 {
			fmt.Printf("stoping due math.MaxUint32\n")
			break
		}

		if *stopFlag {
			break
		}
	}
}

// warning: changes precision
func roundDownToPairAndNineMultiple(bigFloat *big.Float) {
	// if already rounded, it's an int, nothing to do
	if bigFloat.IsInt() {
		return
	}

	// convert to big.Int
	tempBigInt := new(big.Int)
	bigFloat.Int(tempBigInt)

	// subtract until multiple of nine
	nineBigInt := new(big.Int).SetInt64(9)
	modBigInt := new(big.Int)
	modBigInt.Mod(tempBigInt, nineBigInt)
	tempBigInt.Sub(tempBigInt, modBigInt)

	// if not even, reduce more 9 units
	twoBigInt := new(big.Int).SetInt64(2)
	modBigInt.Mod(tempBigInt, twoBigInt)
	if modBigInt.Int64() == 1 {
		tempBigInt.Sub(tempBigInt, nineBigInt)
	}

	bigFloat.SetInt(tempBigInt)
}
