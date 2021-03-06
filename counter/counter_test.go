package counter_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/zhixian001/jwt-cracker/counter"
)

func TestToString(t *testing.T) {
	chars := "012"
	maxLength := 4

	cntr := counter.MakeCounter(chars, maxLength)

	cntr.Increase()
	cntr.Increase()
	cntr.Increase()

	cntrResult := cntr.ToString()

	if cntrResult != "00" {
		t.Errorf("Counter ToString Error. (expected: %s, actual: %s)\n", "00", cntrResult)
	}
}

func TestLoadString(t *testing.T) {
	chars := "012"
	maxLength := 4

	cntr := counter.MakeCounter(chars, maxLength)

	cntr.Increase()
	cntr.Increase()
	cntr.Increase()

	cntr.LoadString("210")
	cntr.LoadString("21")

	cntrResult := cntr.ToString()

	if cntrResult != "21" {
		t.Errorf("Counter LoadString Error. (expected: %s, actual: %s)\n", "21", cntrResult)
	}
}

func TestToBigInt(t *testing.T) {
	chars := "012"
	maxLength := 4

	cntr := counter.MakeCounter(chars, maxLength)

	cntr.LoadString("000")

	toBigIntResult := cntr.ToBigInt()

	if toBigIntResult.String() != "13" {
		t.Errorf("Counter ToBigInt Error. (expected: %s, actual: %s)\n", "13", toBigIntResult)
	}

	cntr.LoadString("201")

	toBigIntResult = cntr.ToBigInt()

	if toBigIntResult.String() != "32" {
		t.Errorf("Counter ToBigInt Error. (expected: %s, actual: %s)\n", "32", toBigIntResult)
	}

	cntr.LoadString("00")

	toBigIntResult = cntr.ToBigInt()

	if toBigIntResult.String() != "4" {
		t.Errorf("Counter ToBigInt Error. (expected: %s, actual: %s)\n", "4", toBigIntResult)
	}
}

func TestLoadIntSlow(t *testing.T) {
	chars := "012"
	maxLength := 4

	cntr := counter.MakeCounter(chars, maxLength)

	for i := int64(1); i <= 120; i++ {
		cntr.LoadBigIntSlow(big.NewInt(i))
		println("\t" + cntr.ToString())
	}

	cntr.LoadBigIntSlow(big.NewInt(13))
	toStringResult := cntr.ToString()

	if toStringResult != "000" {
		t.Errorf("Counter LoadBigIntSlow Error. (expected: %s, actual: %s)\n", "000", toStringResult)
	}

	cntr.LoadBigIntSlow(big.NewInt(32))
	toStringResult = cntr.ToString()

	if toStringResult != "201" {
		t.Errorf("Counter LoadBigIntSlow Error. (expected: %s, actual: %s)\n", "201", toStringResult)
	}

	cntr.LoadBigIntSlow(big.NewInt(4))
	toStringResult = cntr.ToString()

	if toStringResult != "00" {
		t.Errorf("Counter LoadBigIntSlow Error. (expected: %s, actual: %s)\n", "00", toStringResult)
	}
}

func TestGetNumberOfDigitsInBijectiveSystem(t *testing.T) {
	chars := "012"
	maxLength := 4

	cntr := counter.MakeCounter(chars, maxLength)

	for n := 1; n <= 120; n++ {
		testingNumber := big.NewInt(int64(n))
		digits := cntr.GetNumberOfDigitsInBijectiveSystem(testingNumber)

		fmt.Printf("number: %d\t-> digit: %d\n", n, digits)

		switch n {
		case 1:
			if digits != 1 {
				t.Errorf("%d th number's digit should be 1\n", n)
			}
		case 4:
			if digits != 2 {
				t.Errorf("%d th number's digit should be 2\n", n)
			}
		case 12:
			if digits != 2 {
				t.Errorf("%d th number's digit should be 2\n", n)
			}
		case 13:
			if digits != 3 {
				t.Errorf("%d th number's digit should be 3\n", n)
			}
		default:
			//
		}
	}
}

func TestLoadInt(t *testing.T) {
	chars := "012"
	maxLength := 4

	cntr := counter.MakeCounter(chars, maxLength)

	for i := int64(1); i <= 120; i++ {
		cntr.LoadBigInt(big.NewInt(i))
		println("\t" + cntr.ToString())
	}

	cntr.LoadBigInt(big.NewInt(13))
	toStringResult := cntr.ToString()

	if toStringResult != "000" {
		t.Errorf("Counter LoadBigInt Error. (expected: %s, actual: %s)\n", "000", toStringResult)
	}

	cntr.LoadBigInt(big.NewInt(32))
	toStringResult = cntr.ToString()

	if toStringResult != "201" {
		t.Errorf("Counter LoadBigInt Error. (expected: %s, actual: %s)\n", "201", toStringResult)
	}

	cntr.LoadBigInt(big.NewInt(4))
	toStringResult = cntr.ToString()

	if toStringResult != "00" {
		t.Errorf("Counter LoadBigInt Error. (expected: %s, actual: %s)\n", "00", toStringResult)
	}
}
