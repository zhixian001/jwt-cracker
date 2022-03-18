package counter_test

import (
	"math/big"
	"testing"

	"github.com/zhixian001/jwt-cracker/counter"
)

func TestToString(t *testing.T) {
	alphabet := "012"
	maxLength := 4

	cntr := counter.MakeCounter(alphabet, maxLength)

	cntr.Increase()
	cntr.Increase()
	cntr.Increase()

	cntrResult := cntr.ToString()

	if cntrResult != "00" {
		t.Errorf("Counter ToString Error. (expected: %s, actual: %s)\n", "00", cntrResult)
	}
}

func TestLoadString(t *testing.T) {
	alphabet := "012"
	maxLength := 4

	cntr := counter.MakeCounter(alphabet, maxLength)

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
	alphabet := "012"
	maxLength := 4

	cntr := counter.MakeCounter(alphabet, maxLength)

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

func TestLoadInt(t *testing.T) {
	alphabet := "012"
	maxLength := 4

	cntr := counter.MakeCounter(alphabet, maxLength)

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

// [0, 1, 2]

// base: 3

//     0  1  ---- 1
//     1  1  ---- 2
//     2  1  ---- 3
//   0 0  2  ---- 4
//   0 1  2  ---- 5
//   0 2  2  ---- 6
//   1 0  2  ---- 7
//   1 1  2  ---- 8
//   1 2  2  ---- 9
//   2 0  2  ---- 10
//   2 1  2  ---- 11
//   2 2  2  ---- 12
// 0 0 0  3  ---- 13

// 000  13
// 001
// 002
// 010
// 011
// 012
// 020
// 021
// 022
// 100

// 2 0 1  3  ---- 32

// 9 3 1

//  27+3+2
