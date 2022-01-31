package findpsum

import (
	"math/big"
	"testing"
)

type tryToFindPrimeFactorsTest struct {
	inputTargetNumberN               *big.Int
	inputPrimeFactorsSumEstimatation *big.Int
	expectedOutput                   TryToFindPrimeFactorsResult
}

func getBigInt(bigIntAsString string) *big.Int {
	result, _ := new(big.Int).SetString(bigIntAsString, 10)
	return result
}

var tryToFindPrimeFactorsTests = []tryToFindPrimeFactorsTest{
	{getBigInt("13322573880505234223"), getBigInt("7300114376"), TryToFindPrimeFactorsResult{true, getBigInt("3637339949"), getBigInt("3662724427")}},
	{getBigInt("13322573880505234223"), getBigInt("7300214376"), TryToFindPrimeFactorsResult{true, getBigInt("3637339949"), getBigInt("3662724427")}},
	{getBigInt("13322573880505234223"), getBigInt("7300014376"), TryToFindPrimeFactorsResult{true, getBigInt("3637339949"), getBigInt("3662724427")}},
	{getBigInt("13322573880505234223"), getBigInt("7300914376"), TryToFindPrimeFactorsResult{true, getBigInt("3637339949"), getBigInt("3662724427")}},
	{getBigInt("17951149905535786241"), getBigInt("8479111442"), TryToFindPrimeFactorsResult{true, getBigInt("4198438661"), getBigInt("4275672781")}},
}

func TestTryToFindPrimeFactors(t *testing.T) {
	for _, test := range tryToFindPrimeFactorsTests {
		output := TryToFindPrimeFactors(test.inputTargetNumberN, test.inputPrimeFactorsSumEstimatation)
		if output.foundFactors != test.expectedOutput.foundFactors {
			t.Errorf("Output %v not equal to expected %v", output.foundFactors, test.expectedOutput.foundFactors)
		}
		if output.foundFactors && output.smallFactor.Cmp(test.expectedOutput.smallFactor) != 0 {
			t.Errorf("Output %v not equal to expected %v", output.smallFactor, test.expectedOutput.smallFactor)
		}
		if output.foundFactors && output.bigFactor.Cmp(test.expectedOutput.bigFactor) != 0 {
			t.Errorf("Output %v not equal to expected %v", output.bigFactor, test.expectedOutput.bigFactor)
		}
	}
}
