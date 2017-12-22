package cron

import (
	"testing"
	"time"
)

func TestNext(t *testing.T) {

	now := time.Date(2017, 12, 31, 23, 59, 0, 0, time.UTC)

	// every next is now
	schedule, err := Parse("* * * * *")
	ok(t, err)
	expected := now
	equals(t, expected, schedule.Next(now))

	// 2 minutes from now
	schedule, err = Parse("2 * * * *")
	ok(t, err)
	expected = time.Date(2018, 1, 1, 0, 2, 0, 0, time.UTC)
	equals(t, expected, schedule.Next(now))

	// 2 minutes from now
	schedule, err = Parse("2 * * * *")
	ok(t, err)
	expected = time.Date(2018, 1, 1, 0, 2, 0, 0, time.UTC)
	equals(t, expected, schedule.Next(now))

	// first day of month
	schedule, err = Parse("* * 1 * *")
	ok(t, err)
	expected = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	equals(t, expected, schedule.Next(now))

	// 3 hours from now
	schedule, err = Parse("* 2 * * *")
	ok(t, err)
	expected = time.Date(2018, 1, 1, 2, 0, 0, 0, time.UTC)
	equals(t, expected, schedule.Next(now))

	schedule, err = Parse("* * 2 1 *")
	ok(t, err)
	expected = time.Date(2018, 1, 2, 0, 0, 0, 0, time.UTC)
	equals(t, expected, schedule.Next(now))

	// two months from now
	schedule, err = Parse("* * * 2 *")
	ok(t, err)
	expected = time.Date(2018, 2, 0, 0, 0, 0, 0, time.UTC)
	equals(t, expected, schedule.Next(now))

	schedule, err = Parse("* * * * 1")
	ok(t, err)
	expected = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	equals(t, expected, schedule.Next(now))

	schedule, err = Parse("* * * * */3")
	ok(t, err)
	expected = now
	equals(t, expected, schedule.Next(now))
}
