package util

import (
"time"
)

type TimeDif struct {
	Hours   int
	Minutes int
	Seconds int
	Days    int
	Months  int
	Years   int
}


func CurrDate() string {
	t := time.Now()
	returnVal := t.Format("20060102")
	return returnVal
}


func CurrDateTime() string {
	t := time.Now()
	returnVal := t.Format("20060102150405")
	return returnVal
}

func CurrDateTimeFmt() string {
	t := time.Now()
	returnVal := t.Format("2006-01-02 15:04:05")
	return returnVal
}


func (self *TimeDif) adjust(daysInMonth int) {
	if self.Seconds > 59 {
		self.Minutes++
		self.Seconds -= 60
	}
	if self.Minutes > 59 {
		self.Hours++
		self.Minutes -= 60
	}
	if self.Hours > 23 {
		self.Days++
		self.Hours -= 24
	}
	//daysInMonth := daysIn(time.Month(int(to.Month())-1), to.Year())
	if self.Days >= daysInMonth {
		self.Days -= daysInMonth
		self.Months++
	}

	if self.Months > 11 {
		self.Months = self.Months - 12
		self.Years++
	}
}

func Difference(from, to time.Time) TimeDif {
	//
	diff := diffBetweenToYearRound(from, to)

	diffEnd := diffBetweenFromYearRound(to)

	if diffEnd.Years >= 0 {
		diff.Years += diffEnd.Years
		diff.Months += diffEnd.Months
		diff.Days += diffEnd.Days

		diff.adjust(daysIn(time.Month(int(to.Month())-1), to.Year()))
	}
	/////////////////////Time Calculation/////////////////////////////
	//It's calculated with a trick, assuming from date is 1 day before the to date
	tempT := from.AddDate(diff.Years, diff.Months, diff.Days-1) //creates 1 day prior to to.Date with time from.Date
	dur := int(to.Sub(tempT).Seconds())
	if dur >= 24*3600 { //to remove one day hence this day already calculated above
		dur -= 24 * 3600
	}
	diff.Hours = dur / 3600                  //seconds to Hour
	dur -= diff.Hours * 3600                 // removing Calculated hour from duration
	diff.Minutes = dur / 60                  // seconds to Minute
	diff.Seconds = dur - (diff.Minutes * 60) //reminding seconds
	////Fix date with time//////////////////
	diff.adjust(daysIn(time.Month(int(to.Month())-1), to.Year()))

	return diff
}

func diffBetweenToYearRound(from, to time.Time) TimeDif {
	toRoundYear := time.Date(to.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	return roudDefference(from, toRoundYear)
}

func diffBetweenFromYearRound(to time.Time) TimeDif {
	fromDiffAdjestTime := time.Date(to.Year(), 1, 2, 0, 0, 0, 0, time.UTC)
	return roudDefference(fromDiffAdjestTime, to)

}

func roudDefference(from, to time.Time) TimeDif {
	diff := TimeDif{0, 0, 0, 0, 0, 0}

	//toRoundYear := time.Date(to.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	diff.Years = to.Year() - from.Year()
	diff.Months = int(to.Month()) - int(from.Month())
	diff.Days = to.Day() - from.Day()
	///////////////////Date Calculation/////////////////////////////////
	if diff.Months < 0 || (diff.Months == 0 && diff.Days < 0) {
		diff.Years--
		if diff.Days < 0 {
			diff.Days *= -1
		}

		if diff.Months == 0 && diff.Days != 0 {
			diff.Months = 11
		} else {
			diff.Months = 11 - (diff.Months * -1)
		}
		diff.Days = daysIn(time.Month(int(to.Month())-1), to.Year()) - diff.Days
	} else if diff.Days < 0 {
		diff.Days *= -1
		diff.Days = daysIn(time.Month(int(to.Month())-1), to.Year()) - diff.Days
		diff.Months--
	}
	return diff
}

func IsLepYear(year int) bool {
	// Checks whether given year is leap year or not
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	// As only month and year given last day of month returns
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}