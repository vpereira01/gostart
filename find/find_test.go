package find

import (
	"math/big"
	"testing"
)

type findPrimesSubTest struct {
	input, expectedOutput *big.Int
}

func getBigInt(bigIntAsString string) *big.Int {
	result, _ := new(big.Int).SetString(bigIntAsString, 10)
	return result
}

var findPrimesSubTests = []findPrimesSubTest{
	{getBigInt("13322573880505234223"), getBigInt("25384478")},
	{getBigInt("13065662271585108503"), getBigInt("427464298")},
	{getBigInt("13940064823464665839"), getBigInt("1039934370")},
	{getBigInt("14042965475976382991"), getBigInt("656435330")},
	{getBigInt("16648983298230638723"), getBigInt("80336162")},
	{getBigInt("13797531713940771079"), getBigInt("1008271722")},
	{getBigInt("13782080400846584441"), getBigInt("33299240")},
	{getBigInt("13489600177963835063"), getBigInt("287625418")},
	{getBigInt("13257863350502626553"), getBigInt("517219168")},
}

func TestFindPrimesSub(t *testing.T) {
	for _, test := range findPrimesSubTests {
		if output := FindPrimesSub(test.input); output.Cmp(test.expectedOutput) != 0 {
			t.Errorf("Output %q not equal to expected %q", output, test.expectedOutput)
		}
	}
}
