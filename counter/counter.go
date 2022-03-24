package counter

// #cgo CFLAGS: -Wall
// #cgo pkg-config: mpfr gmp
// #include "mpfr_calculation.h"
import "C"

import (
	"errors"
	"math"
	"math/big"
	"strings"
)

type counter struct {
	charsSplit []string
	counter    []int
	base       int
	width      int
	overflow   bool
}

// Couter of number of cases
func MakeCounter(chars string, maxLength int) *counter {
	charsSplit := strings.Split(chars, "")

	counter := &counter{
		charsSplit: charsSplit,
		counter:    make([]int, maxLength),
		base:       len(charsSplit),
		width:      maxLength,
		overflow:   false,
	}

	counter.counter[maxLength-1] = 1

	return counter
}

// Increase counter value by 1
func (cntr *counter) Increase() {
	lastDigitNewVal := cntr.counter[cntr.width-1] + 1

	// Addition logic
	carry := false

	if lastDigitNewVal > cntr.base {
		lastDigitNewVal = 1
		carry = true
	}

	cntr.counter[cntr.width-1] = lastDigitNewVal

	// update with carry
	for i := cntr.width - 2; i >= 0; i-- {
		if carry {
			newVal := cntr.counter[i] + 1

			if newVal > cntr.base {
				// carry continues
				newVal = 1
			} else {
				// carry ends
				carry = false
			}

			cntr.counter[i] = newVal
		} else {
			break
		}
	}

	// check carry enables new digit. overflow can be checked by evaluating (cntr.enabledWidth > cntr.width)
	// Overflow occur times = cntr.enabledWidth - cntr.width
	if carry {
		cntr.overflow = true
	}

}

// Get corresponding string combination representation of current counter value.
func (cntr *counter) ToString() string {
	if cntr.overflow {
		panic(errors.New("counter overflow"))
	}

	result := ""

	for _, charIndex := range cntr.counter {
		if charIndex != 0 {
			result += cntr.charsSplit[charIndex-1]
		}
	}

	return result
}

// Reset Counter
func (cntr *counter) reset() {
	cntr.overflow = false
	for i := range cntr.counter {
		cntr.counter[i] = 0
	}

	cntr.counter[cntr.width-1] = 1
}

// Get index of char inside 'charsSplit' array.
func (cntr *counter) charIndexOf(char string) int {
	searchIdx := -1

	for i := range cntr.charsSplit {
		if cntr.charsSplit[i] == char {
			searchIdx = i
			break
		}
	}

	return searchIdx
}

// Set counter value from string combination.
func (cntr *counter) LoadString(str string) {
	// TODO: check string
	inputStrSplitted := strings.Split(str, "")

	if len(inputStrSplitted) > cntr.width {
		panic(errors.New("counter string load error"))
	}

	cntr.reset()

	strIdx := 0
	for counterIdx := cntr.width - len(inputStrSplitted); counterIdx < cntr.width; counterIdx++ {
		cntr.counter[counterIdx] = cntr.charIndexOf(inputStrSplitted[strIdx]) + 1

		strIdx++
	}
}

// Convert counter value to bigint number
func (cntr *counter) ToBigInt() *big.Int {
	ret := big.NewInt(0)

	base := big.NewInt(int64(cntr.base))

	for i, v := range cntr.counter {
		digitExp := big.NewInt(int64(cntr.width - i - 1))

		digitExp.Exp(base, digitExp, nil)

		digitExp.Mul(digitExp, big.NewInt(int64(v)))

		ret.Add(ret, digitExp)
	}

	return ret
}

// Load counter value from bigint number
// decimal system to counter(bijective numerial system)
func (cntr *counter) LoadBigInt(inp *big.Int) {
	// Cases (under base 3 bijective)
	// 1 - digit
	// 	dec: [1  ,  3]
	// 2 - digit
	//  dec: [4  , 12]
	// 3 - digit
	// 	dec: [13 , 39]
	// 4 - digit
	// 	dec: [40 , 120]
	// ...
	// n - digit
	//  dec: [, ]

	cntr.reset()

	// find number of digit candidates
	numberOfDigits := cntr.GetNumberOfDigitsInBijectiveSystem(inp)

	base := big.NewInt(int64(cntr.base))
	bigIntOne := big.NewInt(1)
	bigIntZero := big.NewInt(0)

	// fill all digits
	subtractor := new(big.Int).Exp(base, big.NewInt(int64(numberOfDigits)), nil)
	subtractor.Sub(subtractor, bigIntOne)
	subtractHelper := new(big.Int).Sub(base, bigIntOne)
	subtractor.Div(subtractor, subtractHelper)

	// Copy from input
	remaining := new(big.Int).Set(inp)

	// Subtract subtractor from remaining
	remaining.Sub(remaining, subtractor)

	// Set other digit numbers
	// common base n system (with 0)
	if numberOfDigits > cntr.width {
		panic(errors.New("digit is wider than counter width"))
	}

	digitDifference := cntr.width - numberOfDigits

	// fill counter
	for i := range cntr.counter {
		if i+digitDifference == cntr.width {
			break
		}

		currentDigitMultiplier := new(big.Int).Exp(base, big.NewInt(int64(cntr.width-i-digitDifference-1)), nil)

		additionalDigit := 0

		for additionalDigitCandidate := cntr.base - 1; additionalDigitCandidate > 0; additionalDigitCandidate-- {
			subtractor = subtractor.Mul(currentDigitMultiplier, big.NewInt(int64(additionalDigitCandidate)))

			subtractHelper = subtractHelper.Sub(remaining, subtractor)

			// Compare
			if subtractHelper.Cmp(bigIntZero) >= 0 {
				remaining.Set(subtractHelper)

				additionalDigit = additionalDigitCandidate
				break
			}
		}

		cntr.counter[i+digitDifference] = 1 + additionalDigit
	}
}

// Calculate how many digits are needed to express input decimal integer in bijective format
func (cntr *counter) GetNumberOfDigitsInBijectiveSystem(number *big.Int) int {
	base := big.NewInt(int64(cntr.base))

	baseLn := math.Log(float64(cntr.base))

	bigIntOne := big.NewInt(1)

	// left
	leftOver := new(big.Int).Sub(base, bigIntOne)
	leftOver.Mul(leftOver, number)
	leftOver.Add(leftOver, base)

	leftOverLog := float64(C.log_e(C.CString(leftOver.String())))

	leftValue := math.Ceil((leftOverLog / baseLn) - 1)

	// right
	rightOver := new(big.Int).Sub(base, bigIntOne)
	rightOver.Mul(rightOver, number)
	rightOver.Add(rightOver, bigIntOne)

	rightOverLog := float64(C.log_e(C.CString(rightOver.String())))
	rightValue := math.Floor(rightOverLog / baseLn)

	// Result
	return int(math.Min(leftValue, rightValue))
}

// Load from bigint number (iterative slow version - very slow. for debugging)
// decimal system to counter(bijective numerial system)
func (cntr *counter) LoadBigIntSlow(inp *big.Int) {
	cntr.reset()

	maxInt := big.NewInt(0)

	base := big.NewInt(int64(cntr.base))

	for i := range cntr.counter {
		toAdd := big.NewInt(int64(i + 1))
		toAdd.Exp(base, toAdd, nil)
		maxInt.Add(maxInt, toAdd)
	}

	if maxInt.Cmp(inp) < 0 {
		panic(errors.New("counter int load error (out of range)"))
	}

	for {
		currentInt := cntr.ToBigInt()
		if currentInt.Cmp(inp) == 0 {
			break
		}
		cntr.Increase()
	}
}

// Set counter value to max value.
func (cntr *counter) ToMaxValue() {
	cntr.reset()
	for i := range cntr.counter {
		cntr.counter[i] = cntr.base
	}
}

// Set counter value to min value.
func (cntr *counter) ToMinValue() {
	cntr.reset()
}
