package counter

import (
	"errors"
	"math/big"
	"strings"
)

type counter struct {
	alphabetSplit []string
	counter       []int
	base          int
	width         int
	overflow      bool
}

// Couter of number of cases
func MakeCounter(alphabet string, maxLength int) *counter {
	alphaSplit := strings.Split(alphabet, "")

	counter := &counter{
		alphabetSplit: alphaSplit,
		counter:       make([]int, maxLength),
		base:          len(alphaSplit),
		width:         maxLength,
		overflow:      false,
	}

	counter.counter[maxLength-1] = 1

	return counter
}

// TODO: Is overflow check required in this system?
// func (cntr *counter) checkOverflow() bool {
// return cntr.enabledWidth > cntr.width
// }

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

	for _, alphabetIndex := range cntr.counter {
		if alphabetIndex != 0 {
			result += cntr.alphabetSplit[alphabetIndex-1]
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

// Get index of alphabet inside 'alphabetSplit' array.
func (cntr *counter) alphabetIndexOf(alpha string) int {
	searchIdx := -1

	for i := range cntr.alphabetSplit {
		if cntr.alphabetSplit[i] == alpha {
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
		cntr.counter[counterIdx] = cntr.alphabetIndexOf(inputStrSplitted[strIdx]) + 1

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
// TODO: Implement this!!! critical logic
func (cntr *counter) LoadBigInt(inp *big.Int) {
	cntr.reset()

}

// Load from bigint number (iterative slow version - very slow. for debugging)
// decimal system to counter(bijective numerial system)
func (cntr *counter) LoadBigIntSlow(inp *big.Int) {
	cntr.reset()

	// inpCopied := big.NewInt(0).Set(inp)

	// base := big.NewInt(int64(cntr.base))
	// zero := big.NewInt(0)

	// for i := range cntr.counter {
	// 	calc := big.NewInt(0)

	// 	currentDigitMutiplier := calc.Exp(base, big.NewInt(int64(i)), nil)
	// 	currentDigitMaxValue := calc.Mul(currentDigitMutiplier, base)

	// 	compareResult := inpCopied.Cmp(currentDigitMaxValue)

	// 	if compareResult <= 0 {
	// 		// 이번 자리에서 끝
	// 		quotient := int(calc.Div(inpCopied, currentDigitMutiplier).Int64())

	// 		cntr.counter[cntr.width-1-i] = quotient

	// 		break
	// 	} else {
	// 		// 다음 자리수 계산
	// 	}

	// }

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

	// TOFIX: Add Algorithm
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
