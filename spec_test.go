package cron

import (
	"testing"
	"regexp"
	"strconv"
)

type expectation = map[bool][]string
type expectations = map[*regexp.Regexp]expectation

var samples = expectations{
	Every: {
		true:  {"*"},
		false: {"**", "2", "a"},
	},
	SingleOrDoubleDigit: {
		true:  {"1", "2", "10", "12"},
		false: {"a1", "123", "a", "."},
	},
	Step: {
		true:  {"/2", "/28"},
		false: {"/a", "2", "/123"},
	},
	EveryStep: {
		true:  {"*/2", "*/28"},
		false: {"*/a", "2", "/123"},
	},
	Alias: {
		true:  {"abb", "man"},
		false: {"m", "mann", "123"},
	},
	List: {
		true:  {"1,2,3,4", "12,1,42,2"},
		false: {"a,1,2,3", "1,2,3,", "123,2"},
	},
	RangeStep: {
		true:  {"12-31/23", "1-2/2", "1-31/32"},
		false: {"*/2", "1,2,4", "123"},
	},
	Range: {
		true:  {"1-2", "21-2"},
		false: {"123-2", "5-123", "123-123"},
	},
	Name: {
		true: {"@daily", "@annually", "@midnight", "@hourly"},
		false: {"@every", "monday"},
	},
}

func matchSamples(t *testing.T, r *regexp.Regexp) {
	match := func(exp string) bool {
		return r.MatchString(exp)
	}
	for expectation := range samples[r] {
		for _, sample := range samples[r][expectation] {
			if !expectation == match(sample) {
				t.Error("Expected", expectation, sample)
			}
		}
	}
}

func TestEveryRegexp(t *testing.T) {
	matchSamples(t, Every)
}

func TestSingleOrDoubleDigitRexexp(t *testing.T) {
	matchSamples(t, SingleOrDoubleDigit)
}

func TestStepRegexp(t *testing.T) {
	matchSamples(t, Step)
}

func TestEveryStepRegexp(t *testing.T) {
	matchSamples(t, EveryStep)
}

func TestAliasRegexp(t *testing.T) {
	matchSamples(t, Alias)
}

func TestRangeRegexp(t *testing.T) {
	matchSamples(t, Range)
}

func TestRangeStepRegexp(t *testing.T) {
	matchSamples(t, RangeStep)
}

func TestListRegexp(t *testing.T) {
	matchSamples(t, List)
}

func TestNameRegexp(t *testing.T) {
	matchSamples(t, Range)
}

func TestEveryStepSubmatch(t *testing.T) {
	expectedStep := uint64(32)
	sub := EveryStep.FindStringSubmatch("*/32")
	if len(sub) != 2 {
		t.Fatal("expected length of submatches", 2, "got", len(sub))
	}
	step, err := strconv.ParseUint(sub[1], 10, 8)
	if err != nil {
		t.Error(err)
	}
	if step != expectedStep {
		t.Error("expected", expectedStep, "got", step)
	}
}

func TestRangeStepSubmatch(t *testing.T) {
	expectedFrom := uint8(1)
	expectedTo := uint8(31)
	expectedStep := uint8(32)

	s := "1-31/32"

	sub := RangeStep.FindStringSubmatch(s)
	if len(sub) != 4 {
		t.Fatal("expected length of submatch", 4, "got", len(sub), s)
	}

	from, err := strconv.ParseUint(sub[1], 10, 8)
	if err != nil {
		t.Error(err)
	} else if uint8(from) != expectedFrom {
		t.Error("expected", expectedFrom, "got", from)
	}

	to, err := strconv.ParseUint(sub[2], 10, 8)
	if err != nil {
		t.Error(err)
	} else if uint8(to) != expectedTo {
		t.Error("expected", expectedTo, "got", to)
	}

	step, err := strconv.ParseUint(sub[3], 10, 8)
	if err != nil {
		t.Error(err)
	} else if uint8(step) != expectedStep {
		t.Error("expected", expectedStep, "got", step)
	}
}

func TestValidNumber(t *testing.T) {
	if spec.minute.InRange(60) {
		t.Error("60 should not be a valid number for minute field")
	}
	if !spec.minute.InRange(0) {
		t.Error("0 should be a valid number")
	}
	if !spec.minute.InRange(59) {
		t.Error("60 not valid number for minute")
	}
}

func TestValidAlias(t *testing.T) {
	if _, err := spec.dow.Unalias("mond"); err == nil {
		t.Error("should error on invalid alias")
	}
	if _, err := spec.dow.Unalias("mon"); err != nil {
		t.Error("mon should be a valid alias")
	}
}
