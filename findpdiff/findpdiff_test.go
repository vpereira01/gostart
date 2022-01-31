package findpdiff

import (
	"math/big"
	"testing"
)

type findPrimesFactorsDiffTest struct {
	input          *big.Int
	expectedOutput tryFindPrimesDiffResult
}

type TryToFindPrimeFactorsTest struct {
	input          *big.Int
	expectedOutput TryToFindPrimeFactorsResult
}

func getBigInt(bigIntAsString string) *big.Int {
	result, _ := new(big.Int).SetString(bigIntAsString, 10)
	return result
}

var TryToFindPrimeFactorsTests = []TryToFindPrimeFactorsTest{
	{getBigInt("13322573880505234223"), TryToFindPrimeFactorsResult{true, getBigInt("3637339949"), getBigInt("3662724427")}},
	{getBigInt("13065662271585108503"), TryToFindPrimeFactorsResult{true, getBigInt("3407226799"), getBigInt("3834691097")}},
	{getBigInt("16726821280062444821"), TryToFindPrimeFactorsResult{true, getBigInt("4082519047"), getBigInt("4097181443")}},
	{getBigInt("13940064823464665839"), TryToFindPrimeFactorsResult{true, getBigInt("3249705307"), getBigInt("4289639677")}},
	{getBigInt("14042965475976382991"), TryToFindPrimeFactorsResult{true, getBigInt("3433522931"), getBigInt("4089958261")}},
	{getBigInt("16648983298230638723"), TryToFindPrimeFactorsResult{true, getBigInt("4040346197"), getBigInt("4120682359")}},
	{getBigInt("13797531713940771079"), TryToFindPrimeFactorsResult{true, getBigInt("3244421819"), getBigInt("4252693541")}},
	{getBigInt("13782080400846584441"), TryToFindPrimeFactorsResult{true, getBigInt("3695810159"), getBigInt("3729109399")}},
}

var findPrimesFactorsDiffTests = []findPrimesFactorsDiffTest{
	{getBigInt("13322573880505234223"), tryFindPrimesDiffResult{true, getBigInt("25384478")}},   // 4*n % 9 = 0
	{getBigInt("13065662271585108503"), tryFindPrimesDiffResult{true, getBigInt("427464298")}},  // 4*n % 9 = 1
	{getBigInt("13371152238833765231"), tryFindPrimesDiffResult{true, getBigInt("667716974")}},  // 4*n % 9 = 2
	{getBigInt("16648983298230638723"), tryFindPrimesDiffResult{true, getBigInt("80336162")}},   // 4*n % 9 = 3
	{getBigInt("14330847798237362867"), tryFindPrimesDiffResult{true, getBigInt("613730734")}},  // 4*n % 9 = 4
	{getBigInt("13797531713940771079"), tryFindPrimesDiffResult{true, getBigInt("1008271722")}}, // 4*n % 9 = 5
	{getBigInt("13940064823464665839"), tryFindPrimesDiffResult{true, getBigInt("1039934370")}}, // 4*n % 9 = 6
	{getBigInt("14042965475976382991"), tryFindPrimesDiffResult{true, getBigInt("656435330")}},  // 4*n % 9 = 6
	{getBigInt("13257863350502626553"), tryFindPrimesDiffResult{true, getBigInt("517219168")}},  // 4*n % 9 = 6
	{getBigInt("13489600177963835063"), tryFindPrimesDiffResult{true, getBigInt("287625418")}},  // 4*n % 9 = 7
	{getBigInt("13782080400846584441"), tryFindPrimesDiffResult{true, getBigInt("33299240")}},   // 4*n % 9 = 8
}

func TestTryToFindPrimeFactors(t *testing.T) {
	for _, test := range TryToFindPrimeFactorsTests {
		output := TryToFindPrimeFactors(test.input)
		if output.foundFactors != test.expectedOutput.foundFactors {
			t.Errorf("Output %v not equal to expected %v", output.foundFactors, test.expectedOutput.foundFactors)
		}

		if output.foundFactors && output.smallFactor.Cmp(test.expectedOutput.smallFactor) != 0 {
			t.Errorf("Output %q not equal to expected %q", output.smallFactor, test.expectedOutput.smallFactor)
		}

		if output.foundFactors && output.bigFactor.Cmp(test.expectedOutput.bigFactor) != 0 {
			t.Errorf("Output %q not equal to expected %q", output.bigFactor, test.expectedOutput.bigFactor)
		}
	}
}

func TestFindPrimesFactorsDiff(t *testing.T) {
	for _, test := range findPrimesFactorsDiffTests {
		output := findPrimesFactorsDiff(test.input)
		if output.wasFound != test.expectedOutput.wasFound {
			t.Errorf("Output %v not equal to expected %v", output.wasFound, test.expectedOutput.wasFound)
		}

		if output.wasFound && output.diffFound.Cmp(test.expectedOutput.diffFound) != 0 {
			t.Errorf("Output %q not equal to expected %q", output.diffFound, test.expectedOutput.diffFound)
		}
	}
}
