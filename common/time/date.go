package time

import (
	"time"
	//"fmt"

	"github.com/pkg/errors"

	. "iugo.fleet/common/logger"
)

type Date int

const DATE_FORMAT = "2006-01-02 15:04:05"

const DATE_FORMAT_2 = "02/01/2006 15:04:05"

const DATE_FORMAT_3 = "02/01/2006"

const DateFormatIncorrect = "Date format is not correct"

const MilliSecondsInDay = 86400000

var TimeLoc *time.Location

func init() {
	var err error

	// FIXME replace location in the future to support international use cases
	TimeLoc, err = time.LoadLocation("Europe/Istanbul")
	if err != nil {
		LogError(err)
	}

	time.Local = TimeLoc
}

func (d Date) ToUnixMillis() (int64, error) {
	t, err := d.ToTime()

	return t.UnixNano()/1e6, err
}

func (d Date) ToUnix() (int64, error) {
	t, err := d.ToTime()

	return t.Unix(), err
}

func (d Date) AddDays(n int) (Date) {
	t, err := d.ToTime()
	if err != nil {
	    LogError(err)
	}

	t = t.Add(time.Duration(n) * 24 * time.Hour)

	return fromTime(t)
}

func (d Date) GetStartUnixMillis() (int64, error) {
	return d.ToUnixMillis()
}

func (d Date) GetEndUnixMillis() (int64, error) {
	m, err := d.ToUnixMillis()
	return m + MilliSecondsInDay - 1, err
}

func FromUnixMillis(t int64) (Date) {
	return fromTime(time.Unix(0, t * 1e6))
}

func FromUnix(t int64) (Date) {
	return fromTime(time.Unix(t, 0))
}

func GetMondayFromUnixMillis(t int64) (Date) {

	/* First lets get date from given milliseconds */
	d, _ := FromUnixMillis(t).ToTime()

	// iterate back to Monday
	for d.Weekday() != time.Monday {
		d = d.AddDate(0, 0, -1)
	}

	// Convert time to date and return
	return fromTime(d)
}

func GetStartOfMonthFromUnixMillis(t int64) (Date) {
	/* First lets get date from given milliseconds */
	d, _ := FromUnixMillis(t).ToTime()

	year  := d.Year()
	month := int(d.Month())
	m := (year*100) + (month * 1)


	return Date(m)
}

func GetStartOfMonthFromUnix(t int64) (Date) {
	/* First lets get date from given milliseconds */
	d, _ := FromUnix(t).ToTime()

	year  := d.Year()
	month := int(d.Month())
	m := (year*10000) + (month * 100) + 1


	return Date(m)
}

func GetStartOfNextMonthFromUnix(t int64) (Date) {
	/* First lets get date from given milliseconds */
	d, _ := FromUnix(t).ToTime()

	e := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)

	year  := e.Year()
	month := int(e.Month())
	m := (year*10000) + (month * 100) + 1



	return Date(m)
}

func GetStartOfYearFromUnix(t int64) (Date) {
	/* First lets get date from given milliseconds */
	d, _ := FromUnix(t).ToTime()

	year  := d.Year()
	m := (year*10000) + (1 * 100) + 1


	return Date(m)
}

func GetStartOfNextYearFromUnix(t int64) (Date) {
	/* First lets get date from given milliseconds */
	d, _ := FromUnix(t).ToTime()

	e := time.Date(d.Year(), 1, 1, 0, 0, 0, 0, time.UTC).AddDate(1, 0, 0)

	year  := e.Year()
	m := (year*10000) + (1 * 100) + 1



	return Date(m)
}

func GetMonthStart(t int64) (int64){
	/* First lets get date from given milliseconds */
	d, _ := FromUnixMillis(t).ToTime()

	year  := d.Year()
	month := int(d.Month())
	m := (year*10000) + (month * 100) + 01

	date := Date(m)

	ms, _  := date.ToUnixMillis()

	return ms
}

func DaysInMonth(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}


// private functions
func (d Date) ToTime() (time.Time, error) {
	year  := d / 10000
	month := (d - (year*10000)) / 100
	day   := d - ((year*10000) + (month * 100))

	if year < 1000 || month > 12 || month < 1 || day > 31 || day < 1 {
		return time.Now(), errors.New(DateFormatIncorrect)
	}

	//fmt.Println(": \t\t", year, month, day, t.UnixNano()/1e6)

	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, TimeLoc)

	return t, nil
}

func fromTime(t time.Time) (Date) {
	year  := t.Year()
	month := int(t.Month())
	day   := t.Day()
	d := (year*10000) + (month * 100) + day

	//fmt.Println(": \t\t", d)

	return Date(d)
}

func SecondsToTime(t int64)(time.Time){

	return time.Unix(0,t*int64(time.Millisecond))

}