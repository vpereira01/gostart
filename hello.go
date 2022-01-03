// Hello Word in Go by Vivek Gite
package main

// Import OS and fmt packages
import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"
)

// Let us start
func main() {

	//3072 is the number of bits for RSA
	bitSize := 512

	//Generate RSA keys
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}
	fmt.Printf("privateKey.D: %v\n", privateKey.D)
	fmt.Printf("privateKey.E: %v\n", privateKey.E)
	fmt.Printf("privateKey.N: %v\n", privateKey.N)
	fmt.Printf("privateKey.Primes: %v\n", privateKey.Primes)

	nAsBigFloat := new(big.Float)
	nAsBigFloat.SetInt(privateKey.N)
	firstPrimeAsBigFloat := new(big.Float)
	firstPrimeAsBigFloat.SetInt(privateKey.Primes[0])
	divResultAsBigFloat := new(big.Float)
	divResultAsBigFloat.Quo(firstPrimeAsBigFloat, nAsBigFloat)
	fmt.Printf("divResultAsBigFloat: %v\n", divResultAsBigFloat)
	// divCenas, _ := divResultAsBigFloat.Float64()
	// divResultAsBigFloat.SetFloat64(divCenas)
	// fmt.Printf("divResultAsBigFloat: %v\n", divResultAsBigFloat)

	divResultAsBigFloat.Mul(divResultAsBigFloat, nAsBigFloat)
	divResult := new(big.Int)
	divResultAsBigFloat.Int(divResult)
	fmt.Printf("divResult: %v\n", divResult)
	diff := new(big.Int)
	diff.Sub(divResult, privateKey.Primes[0])
	fmt.Printf("diff     : %v\n", diff)
	return

	twoBigInt := new(big.Int)
	twoBigInt.SetInt64(2)
	var probablyPrimes []*big.Int
	oddNumber := new(big.Int)
	oddNumber.SetString("10263280077814176196883978050069", 10)

	for ; oddNumber.BitLen() <= 250; oddNumber.Add(oddNumber, twoBigInt) {
		if oddNumber.ProbablyPrime(20) {
			newProbablyPrime := new(big.Int)
			newProbablyPrime.Set(oddNumber)
			probablyPrimes = append(probablyPrimes, newProbablyPrime)
		}

		if len(probablyPrimes)%50 == 0 {
			fmt.Printf("current oddNumber=%v, found %v primes, probablyPrimes: %v\n", oddNumber, len(probablyPrimes), probablyPrimes)
		}
	}

	// easy rsa number 10263280077814176196883978050069
	a, _ := new(big.Int).SetString("10263280077814176196883978050069", 10)
	b, _ := new(big.Int).SetString("1103758723386229", 10)
	var c = new(big.Int)
	c = c.Mul(a, b)
	fmt.Printf("c: %v\n", c)
	expectedProduct, _ := new(big.Int).SetString("1522605027922533360535618378132637429718068114961380688657908494580122963258952897654000350692006139", 10)
	leftFactor, isLeftFactorSet := new(big.Int).SetString("37975227936943673922808872755445627854565536638199", 10)
	rightFactor, isRightFactorSet := new(big.Int).SetString("40094690950920881030683735292761468389214899724061", 10)

	if !isLeftFactorSet || !isRightFactorSet {
		fmt.Printf("conversion to big.int failed: isLeftFactorSet=%v isRightFactorSet=%v\n", isLeftFactorSet, isRightFactorSet)
		return
	}
	var product = new(big.Int)
	product.Mul(leftFactor, rightFactor)
	if product.Cmp(expectedProduct) != 0 {
		fmt.Printf("product and expected product differ\n")
		fmt.Printf("product: %v\n", product)
		fmt.Printf("expectedProduct: %v\n", expectedProduct)
	}
	fmt.Printf("result: %v\n", product)

	fmt.Printf("expectedProduct.IsInt64(): %v\n", expectedProduct.IsInt64())
	fmt.Printf("expectedProduct.Int64(): %v\n", expectedProduct.Int64())
}
