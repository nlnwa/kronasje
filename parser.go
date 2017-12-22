package cron

import (
	"fmt"
	"strings"
	"strconv"
	"math"
)

type Parser interface {
	Parse(expression string) *Schedule
}

type parser struct{}

var cronParser = new(parser)

func Parse(expression string) (Schedule, error) {
	return cronParser.Parse(expression)
}

// Parses a cron expression to a struct which represents the cron expression string as bitfields.
// Since the maximum number of bits needed is 60 (for minutes) an uint64 can be (and is) used to represent each bitfield.
//
// A bit in position 1 means (according to the the spec): 0th minute, 0 hour, 1st, January, Sunday.
// A bit in position 3 means (according to the spec): 2nd minute, 2nd hour, 3 dom, March, Tuesday.
// A bit in position ...
func (p *parser) Parse(expression string) (Schedule, error) {
	fields := strings.Fields(expression)
	nrOfFields := len(fields)

	if nrOfFields != 1 && nrOfFields != 5 {
		return nil, fmt.Errorf("number of fields expected to be either 1 or 5, got %d", nrOfFields)
	}

	if nrOfFields == 1 {
		return p.parseNamedExpression(fields[0])
	} else {
		minute, err := parseField(fields[0], spec.minute)
		hour, err := parseField(fields[1], spec.hour)
		dom, err := parseField(fields[2], spec.dom)
		month, err := parseField(fields[3], spec.month)
		dow, err := parseField(fields[4], spec.dow)

		// return last error
		if err != nil {
			return nil, err
		}

		return &BitCron{minute, hour, dom, month, dow}, nil
	}
}

// Parse named expressions like @yearly, @daily, etc..
func (p *parser) parseNamedExpression(value string) (Schedule, error) {
	if expression, ok := names[value]; ok {
		return p.Parse(expression)
	} else {
		return nil, fmt.Errorf("no such named cron expression")
	}
}

// offset increments the value by 1 if the spec minimum for the field is 0,
// additionally if the field equals dow (day of week) then 7 (sunday's alternative value) is wrapped to 0.
func offset(fieldSpec *FieldSpec, value uint8) uint8 {
	sundayAsSeven := fieldSpec == spec.dow && value == fieldSpec.Max
	if fieldSpec.Min == 0 && !sundayAsSeven {
		return value + 1
	}
	return value
}

// bits return a bit field where all the bits less significant than or equal to the value bit is set to 1's.
func bits(value uint8) uint64 {
	return math.MaxUint64 >> (64 - value)
}

// bit returns the bit field where the bit number of the value is set to 1.
func bit(value uint8) uint64 {
	return 1 << (value - 1)
}

// field converts the value to a bit field
func field(value uint8) uint64 {
	return bit(value)
}

// rangeField returns a bit field where the bits higher than or equal to from and lower than or equal to to are set to 1's.
func rangeField(from uint8, to uint8) uint64 {
	if from == to {
		return field(from)
	}
	return bits(to) ^ bits(from-1)
}

// rangeField returns a bit field where the bits including from are set to 1's at every step lower than or equal to to.
func stepField(from uint8, to uint8, step uint8) uint64 {
	value := uint64(0)
	for i := from; i <= to; i = i + step {
		value |= bit(i)
	}
	return value
}

func listField(values []uint8) uint64 {
	var value uint64
	for _, val := range values {
		value |= bit(val)
	}
	return value
}

func parseEvery(fieldSpec *FieldSpec) (uint64, error) {
	return rangeField(offset(fieldSpec, fieldSpec.Min), offset(fieldSpec, fieldSpec.Max)), nil
}

func parseSingleOrDoubleDigit(value string, fieldSpec *FieldSpec) (uint64, error) {
	num, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	if !fieldSpec.InRange(uint8(num)) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", num, fieldSpec.Min, fieldSpec.Max)
	} else {
		return field(offset(fieldSpec, uint8(num))), nil
	}
}

func parseEveryStep(value string, fieldSpec *FieldSpec) (uint64, error) {
	sub := EveryStep.FindStringSubmatch(value)
	step, err := strconv.ParseUint(sub[1], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	return stepField(offset(fieldSpec, fieldSpec.Min), offset(fieldSpec, fieldSpec.Max), uint8(step)), nil
}

func parseRange(value string, fieldSpec *FieldSpec) (uint64, error) {
	sub := Range.FindStringSubmatch(value)
	from, err := strconv.ParseUint(sub[1], 10, 8)
	to, err := strconv.ParseUint(sub[2], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	if !fieldSpec.InRange(uint8(from)) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", from, fieldSpec.Min, fieldSpec.Max)
	}
	if !fieldSpec.InRange(uint8(to)) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", to, fieldSpec.Min, fieldSpec.Max)
	}
	return rangeField(offset(fieldSpec, uint8(from)), offset(fieldSpec, uint8(to))), nil
}

func parseRangeStep(value string, fieldSpec *FieldSpec) (uint64, error) {
	sub := RangeStep.FindStringSubmatch(value)
	from, err := strconv.ParseUint(sub[1], 10, 8)
	to, err := strconv.ParseUint(sub[2], 10, 8)
	step, err := strconv.ParseUint(sub[3], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	if !fieldSpec.InRange(uint8(from)) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", from, fieldSpec.Min, fieldSpec.Max)
	}
	if !fieldSpec.InRange(uint8(to)) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", to, fieldSpec.Min, fieldSpec.Max)
	}
	return stepField(offset(fieldSpec, uint8(from)), offset(fieldSpec, uint8(to)), uint8(step)), nil
}

func parseAlias(alias string, fieldSpec *FieldSpec) (uint64, error) {
	number, err := fieldSpec.Unalias(alias)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	return field(offset(fieldSpec, number)), nil
}

// NB: Doesn't allow ranges in the list
func parseList(value string, fieldSpec *FieldSpec) (uint64, error) {
	strValues := strings.Split(value, ",")
	values := make([]uint8, len(strValues))
	for i := range strValues {
		val, err := strconv.ParseUint(strings.TrimSpace(strValues[i]), 10, 8)
		if err != nil {
			return 0, fmt.Errorf("%v", err)
		}
		if !fieldSpec.InRange(uint8(val)) {
			return 0, fmt.Errorf("expected %d in the range %d-%d", val, fieldSpec.Min, fieldSpec.Max)
		}
		values[i] = offset(fieldSpec, uint8(val))
	}
	return listField(values), nil
}

// parseField parses any field of a BitCron expression
func parseField(value string, fieldSpec *FieldSpec) (uint64, error) {
	switch {
	case Every.MatchString(value):
		return parseEvery(fieldSpec)

	case SingleOrDoubleDigit.MatchString(value):
		return parseSingleOrDoubleDigit(value, fieldSpec)

	case EveryStep.MatchString(value):
		return parseEveryStep(value, fieldSpec)

	case Range.MatchString(value):
		return parseRange(value, fieldSpec)

	case RangeStep.MatchString(value):
		return parseRangeStep(value, fieldSpec)

	case Alias.MatchString(value):
		return parseAlias(value, fieldSpec)

	case List.MatchString(value):
		return parseList(value, fieldSpec)

	default:
		return 0, fmt.Errorf("field %v does not match any pattern", value)
	}
}
