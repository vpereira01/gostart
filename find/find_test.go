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
	{getBigInt("13322573880505234223"), getBigInt("25384478")},   // 4*n % 9 = 0
	{getBigInt("13065662271585108503"), getBigInt("427464298")},  // 4*n % 9 = 1
	{getBigInt("13371152238833765231"), getBigInt("667716974")},  // 4*n % 9 = 2
	{getBigInt("16648983298230638723"), getBigInt("80336162")},   // 4*n % 9 = 3
	{getBigInt("14330847798237362867"), getBigInt("613730734")},  // 4*n % 9 = 4
	{getBigInt("13797531713940771079"), getBigInt("1008271722")}, // 4*n % 9 = 5
	{getBigInt("13940064823464665839"), getBigInt("1039934370")}, // 4*n % 9 = 6
	{getBigInt("14042965475976382991"), getBigInt("656435330")},  // 4*n % 9 = 6
	{getBigInt("13257863350502626553"), getBigInt("517219168")},  // 4*n % 9 = 6
	{getBigInt("13489600177963835063"), getBigInt("287625418")},  // 4*n % 9 = 7
	{getBigInt("13782080400846584441"), getBigInt("33299240")},   // 4*n % 9 = 8
}

func TestFindPrimesSub(t *testing.T) {
	for _, test := range findPrimesSubTests {
		if output := FindPrimesSub(test.input); output.Cmp(test.expectedOutput) != 0 {
			t.Errorf("Output %q not equal to expected %q", output, test.expectedOutput)
		}
	}
}
