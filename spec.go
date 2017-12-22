package cron

import (
	"regexp"
	"fmt"
)

// This spec tries to adhere to the 4th Berkely Distribution of the crontab
// manual (man 5 crontab) dated 19 April 2010.

const (
	start               = `^`
	end                 = `$`
	every               = `\*`
	singleOrDoubleDigit = `([\d]{1,2})`
	alias               = `([[:alpha:]]{3})`
	step                = `/` + singleOrDoubleDigit
	numberRange         = singleOrDoubleDigit + `-` + singleOrDoubleDigit
	list                = singleOrDoubleDigit + `(?:,\s*` + singleOrDoubleDigit + `)*`
	name                = `@[[:alpha:]]+`
)

var (
	Every               = regexp.MustCompile(start + every + end)
	Step                = regexp.MustCompile(start + step + end)
	EveryStep           = regexp.MustCompile(start + every + step + end)
	SingleOrDoubleDigit = regexp.MustCompile(start + singleOrDoubleDigit + end)
	Alias               = regexp.MustCompile(start + alias + end)
	Range               = regexp.MustCompile(start + numberRange + end)
	List                = regexp.MustCompile(start + list + end)
	RangeStep           = regexp.MustCompile(start + numberRange + step + end)
	Name                = regexp.MustCompile(start + name + end)
)

// Days and months can be specified with named aliases such as "mon", "jan", etc.
type Aliases map[string]uint8

// Every field has a minimum and maximum value and possibly aliases.
type FieldSpec struct {
	Min     uint8
	Max     uint8
	Aliases Aliases
}

// Unalias returns the value aliased as alias. Error returned if the field has no such alias or no aliases.
func (f *FieldSpec) Unalias(alias string) (uint8, error) {
	if f.Aliases == nil {
		return 0, fmt.Errorf("field has no aliases")
	}
	if number, ok := f.Aliases[alias]; !ok {
		return 0, fmt.Errorf(`"%v" is not a valid alias`, alias)
	} else {
		return number, nil
	}
}

// InRange returns a boolean indicating if the given number lies in the range of the minimum and maximum value of the field spec.
func (f *FieldSpec) InRange(number uint8) bool {
	if number < f.Min || number > f.Max {
		return false
	} else {
		return true
	}
}

func (f *FieldSpec) String() string {
	return fmt.Sprintf("min %v, max %v, aliases %+v", f.Min, f.Max, f.Aliases)
}

type fields struct {
	minute *FieldSpec
	hour   *FieldSpec
	dom    *FieldSpec
	month  *FieldSpec
	dow    *FieldSpec
}

var spec = &fields{
	minute: &FieldSpec{0, 59, nil},
	hour:   &FieldSpec{0, 23, nil},
	dom:    &FieldSpec{1, 31, nil},
	month: &FieldSpec{1, 12,
		Aliases{
			"jan": 1,
			"feb": 2,
			"mar": 3,
			"apr": 4,
			"may": 5,
			"jun": 6,
			"jul": 7,
			"aug": 8,
			"sep": 9,
			"okt": 10,
			"nov": 11,
			"des": 12,
		},
	},
	dow: &FieldSpec{0, 7,
		Aliases{
			"sun": 0,
			"mon": 1,
			"tue": 2,
			"wed": 3,
			"thu": 4,
			"fri": 5,
			"sat": 6,
			// "sun": 7,
		},
	},
}

// Common cron expressions can be specified using names
var names = map[string]string{
	"@yearly":   "0 0 1 1 *", // 1st day in the 1st month at midnight
	"@annually": "@yearly",
	"@monthly":  "0 0 1 * *", // 1st day of every month at midnight
	"@weekly":   "0 0 * * 0", // Every sunday at midnight
	"@daily":    "0 0 * * *", // Every day at noon
	"@midnight": "@daily",
	"@hourly":   "0 * * * *", // Every hour
}
