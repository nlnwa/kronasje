package kronasje

import (
	"testing"
	"math"
)

/*
func equals(t *testing.T, actual uint64, expected uint64, msg ...string) {
	if actual != expected {
		if len(msg) > 0 {
			t.Error(msg)
		} else {
			t.Error("Expected", expected, "got", actual)
		}
	}
}
func ok(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
*/

func TestParseSingleOrDoubleDigit(t *testing.T) {
	actual, err := parseSingleOrDoubleDigit("1", spec.minute)
	expected := uint64(2)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseSingleOrDoubleDigit("2", spec.minute)
	expected = uint64(4)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseSingleOrDoubleDigit("3", spec.minute)
	expected = uint64(8)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseSingleOrDoubleDigit("3", spec.dom)
	expected = 4
	ok(t, err)
	equals(t, expected, actual)

}

func TestParseEvery(t *testing.T) {
	actual, err := parseEvery(spec.dow)
	ok(t, err)
	expected := bits(spec.dow.Max)
	equals(t, expected, actual)

	actual, err = parseEvery(spec.dom)
	ok(t, err)
	expected = bits(spec.dom.Max)
	equals(t, expected, actual)

	actual, err = parseEvery(spec.minute)
	ok(t, err)
	expected = bits(offset(spec.minute, spec.minute.Max))
	equals(t, expected, actual)

	actual, err = parseEvery(spec.hour)
	ok(t, err)
	expected = bits(offset(spec.hour, spec.hour.Max))
	equals(t, expected, actual)
}

func TestParseEveryStep(t *testing.T) {
	actual, err := parseEveryStep("*/1", spec.minute)
	expected := bits(60)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseEveryStep("*/2", spec.minute)
	expected = uint64(0)
	for i := uint8(1); i <= 60; i += 2 {
		expected |= bit(i)
	}
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseEveryStep("*/3", spec.month)
	expected = uint64(0)
	for i := uint8(1); i <= spec.month.Max; i += 3 {
		expected |= bit(i)
	}
	ok(t, err)
	equals(t, expected, actual)
}

func TestParseRangeStep(t *testing.T) {
	actual, err := parseRangeStep("2-5/2", spec.minute)
	expected := bit(3) + bit(5)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRangeStep("2-5/3", spec.minute)
	expected = bit(3) + bit(6) // 100100
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRangeStep("2-5/3", spec.dom)
	expected = bit(2) + bit(5) // 10010
	ok(t, err)
	equals(t, expected, actual)
}

func TestParseList(t *testing.T) {
	actual, err := parseList("1, 2, 3", spec.dom)
	expected := bit(3) + bit(2) + bit(1)
	ok(t, err)
	equals(t, expected, actual)
}

// Test the stepField function.
//
// func stepField(from uint8, to uint8, step uint8) (uint64, error)
func TestStepField(t *testing.T) {
	actual := stepField(1, 2, 1)
	expected := bit(2) + bit(1)
	equals(t, expected, actual)

	// 11100
	actual = stepField(3, 5, 1)
	expected = bit(5) + bit(4) + bit(3)

	equals(t, expected, actual)

	actual = stepField(3, 5, 2)
	expected = bit(5) + bit(3)
	equals(t, expected, actual)

	actual = stepField(3, 5, 3)
	expected = bit(3)
	equals(t, expected, actual)
}

func TestRangeFieldFunction(t *testing.T) {
	// It should return just bit value if from equals to
	actual := rangeField(2, 2)
	expected := bit(2)
	equals(t, expected, actual)

	// e.g. 2-3 should produce binary 5 (110)
	actual = rangeField(2, 3)
	expected = uint64(6)
	equals(t, expected, actual)

	// e.g. 3-6 should produce binary 55 (111100)
	actual = rangeField(3, 6)
	expected = bits(6) ^ (bit(3) - 1)
	// 111111 ^ 111
	equals(t, expected, actual)
	//1000000 -> 111111 ^ 11
	expected = (bit(7) - 1) ^ bits(3-1)
	equals(t, expected, actual)

}

func TestBitsFunction(t *testing.T) {
	// 1111111
	actual := bits(7)
	expected := uint64(127)
	equals(t, expected, actual)
	// 111
	actual = bits(3)
	expected = uint64(7)
	equals(t, expected, actual)
}

func TestBitFunction(t *testing.T) {
	actual := bit(2)
	expected := uint64(2)
	equals(t, expected, actual)
	// binary 8 (1000)
	actual = bit(4)
	expected = uint64(8)
	equals(t, expected, actual)

	actual = uint64(math.Pow(2, 63.0))
	expected = uint64(1) << 63
	equals(t, expected, actual)

	// bits(n)+1 should equal bit(n+1)  and
	actual = bit(7)
	expected = bits(6) + 1
	equals(t, expected, actual)
	// vica versa
	actual = bit(7) - 1
	expected = bits(6)
	equals(t, expected, actual)
}
