package taqweem

import (
	"fmt"
	"log"
	"time"
	_ "time/tzdata"

	"github.com/hablullah/go-hijri"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type KSATime time.Time

func CurrentKSATime() KSATime {
	return KSATime(time.Now().UTC().Add(3 * time.Hour))
}

func KSATimeOf(t time.Time) KSATime {
	return KSATime(t.UTC().Add(3 * time.Hour))
}

func (t KSATime) Date() string {
	currdate := (time.Time(t)).Format("2006-01-02")
	return currdate
}

func (t KSATime) WeekDate() string {
	today := time.Time(t)
	year, month, day := today.AddDate(0, 0, -int(today.Weekday())).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, today.Location()).
		Format("2006-01-02")
}

func (t KSATime) MonthDate() string {
	now := time.Time(t)
	year, month, _ := now.Date()
	currentMonthStart := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	formattedDate := currentMonthStart.Format("2006-01-02")
	return formattedDate
}

func (t KSATime) CurrentDayStart() time.Time {
	return t.ksaDay(0)
}

func (t KSATime) NextDayStart() time.Time {
	return t.ksaDay(1)
}

func (t KSATime) ksaDay(inc int) time.Time {
	s := time.Time(t)
	year, month, day := s.Date()
	nd := time.Date(year, month, day+inc, 0, 0, 0, 0, s.Location())
	return nd
}

func (t KSATime) CurrentDayStartInUTC() time.Time {
	return t.ksaDayInUTC(0)
}

func (t KSATime) NextDayStartInUTC() time.Time {
	return t.ksaDayInUTC(1)
}

func (t KSATime) ksaDayInUTC(inc int) time.Time {
	s := time.Time(t).UTC()
	year, month, day := s.Date()
	nd := time.Date(year, month, day+inc, 0, 0, 0, 0, s.Location()).
		Add(-3 * time.Hour)
	return nd
}

func (t KSATime) DayContains(
	timeToCheck time.Time,
) bool {
	return timeToCheck.UTC().After(t.CurrentDayStartInUTC()) &&
		timeToCheck.UTC().Before(t.NextDayStartInUTC())
}

type HjriDate struct {
	year    uint
	month   uint
	day     uint
	weekday time.Weekday
	t       time.Time
}

func HijriOf(t time.Time) HjriDate {
	loc, _ := time.LoadLocation("Asia/Riyadh")
	v, _ := hijri.CreateUmmAlQuraDate(t.UTC().Add(3 * time.Hour))
	return HjriDate{
		year:    uint(v.Year),
		month:   uint(v.Month),
		day:     uint(v.Day),
		weekday: v.Weekday,
		t:       t.UTC().In(loc),
	}
}

func NowHijri() HjriDate {
	loc, _ := time.LoadLocation("Asia/Riyadh")
	t := time.Now().UTC()
	v, _ := hijri.CreateUmmAlQuraDate(t.Add(3 * time.Hour))
	return HjriDate{
		year:    uint(v.Year),
		month:   uint(v.Month),
		day:     uint(v.Day),
		weekday: v.Weekday,
		t:       t.In(loc),
	}
}

func (tn HjriDate) Gregorian() time.Time {
	return tn.t
}

func (t HjriDate) AddDate(years, months, days int) HjriDate {
	tm := t.t.AddDate(years, months, days)
	utm, _ := hijri.CreateUmmAlQuraDate(tm)
	return HjriDate{
		year:    uint(utm.Year),
		month:   uint(utm.Month),
		day:     uint(utm.Day),
		weekday: utm.Weekday,
		t:       tm,
	}
}

func (tn HjriDate) String() string {
	return fmt.Sprintf(
		"%04d-%02d-%02d",
		tn.year,
		tn.month,
		tn.day,
	)
}

func (tn HjriDate) Formatted() string {
	p := message.NewPrinter(
		language.Arabic,
	)
	log.Println(tn.month)
	return p.Sprintf(
		"%02d:%02d:%02d %d %s %d هـ",
		tn.t.Hour(),
		tn.t.Minute(),
		tn.t.Second(),
		tn.day,
		hjiriMonthsAR[tn.month],
		number.Decimal(tn.year, number.NoSeparator()),
	)
}

func (h HjriDate) WeekStartingDay() HjriDate {
	weekStartG := h.t.AddDate(0, 0, -int(h.t.Weekday()))
	wsg, _ := hijri.CreateUmmAlQuraDate(weekStartG)
	return HjriDate{
		year:    uint(wsg.Year),
		month:   uint(wsg.Month),
		day:     uint(wsg.Day),
		weekday: wsg.Weekday,
		t: time.Date(
			weekStartG.Year(),
			weekStartG.Month(),
			weekStartG.Day(),
			0,
			0,
			0,
			0,
			weekStartG.Location(),
		),
	}
}

func (h HjriDate) MonthStartingDay() HjriDate {
	wsg := hijri.UmmAlQuraDate{
		Day:     1,
		Month:   int64(h.month),
		Year:    int64(h.year),
		Weekday: h.weekday,
	}

	return HjriDate{
		year:    uint(wsg.Year),
		month:   uint(wsg.Month),
		day:     uint(wsg.Day),
		weekday: wsg.Weekday,
		t:       wsg.ToGregorian().Add(-3 * time.Hour).In(h.t.Location()),
	}
}

func (t HjriDate) CurrentDayStart() HjriDate {
	return t.dayTime(0)
}

func (t HjriDate) NextMonthStart() HjriDate {
	year := t.year

	nextMonth := (t.month + 1) % 12
	if nextMonth == 1 {
		year++
	}

	nd := hijri.UmmAlQuraDate{
		Day:   1,
		Month: int64(nextMonth),
		Year:  int64(year),
	}

	gnd := nd.ToGregorian().Add(-3 * time.Hour).In(t.t.Location())

	return HjriDate{
		year:    uint(nd.Year),
		month:   uint(nd.Month),
		day:     uint(nd.Day),
		weekday: gnd.Weekday(),
		t:       gnd,
	}
}

func (t HjriDate) NextDayStart() HjriDate {
	return t.dayTime(1)
}

func (h HjriDate) dayTime(inc int) HjriDate {
	year, month, day := h.t.Date()
	nd := time.Date(year, month, day+inc, 0, 0, 0, 0, h.t.Location())

	return HjriDate{
		year:    uint(h.year),
		month:   uint(h.month),
		day:     uint(h.day),
		weekday: h.weekday,
		t:       nd,
	}
}

func (h HjriDate) DayContains(
	timeToCheck time.Time,
) bool {
	return timeToCheck.UTC().After(h.CurrentDayStart().t.UTC()) &&
		timeToCheck.UTC().Before(h.NextDayStart().t.UTC())
}

var hjiriMonthsAR = map[uint]string{
	1:  "محرم",
	2:  "صفر",
	3:  "ربيع الأول",
	4:  "ربيع الثاني",
	5:  "جمادى الأولى",
	6:  "جمادى الثانية",
	7:  "رجب",
	8:  "شعبان",
	9:  "رمضان",
	10: "شوال",
	11: "ذو القعدة",
	12: "ذو الحجة",
}

// converts time.Duration to the most significant unit with mintues as minimum unit.
func DurationToMSU(d time.Duration) string {
	const (
		minutesPerHour = 60
		hoursPerDay    = 24
	)

	// Convert the duration to minutes
	minutes := int(d.Minutes())

	// Calculate days, hours, and remaining minutes
	days := minutes / (minutesPerHour * hoursPerDay)
	hours := (minutes % (minutesPerHour * hoursPerDay)) / minutesPerHour
	remainingMinutes := minutes % minutesPerHour

	result := ""
	if days > 0 {
		return fmtDurationInAR(days, DAY)
	}
	if hours > 0 {
		return fmtDurationInAR(hours, HOUR)
	}

	if remainingMinutes > 0 {
		return fmtDurationInAR(remainingMinutes, MINUTE)
	}

	return result
}

type TimeUnit int

const (
	DAY TimeUnit = iota
	HOUR
	MINUTE
)

func fmtDurationInAR(value int, unit TimeUnit) string {
	p := message.NewPrinter(
		language.Arabic,
	)

	switch value {
	case 1:
		switch unit {
		case DAY:
			return p.Sprintf("%s", "يوم")
		case HOUR:
			return p.Sprintf("%s", "ساعة")
		case MINUTE:
			return p.Sprintf("%s", "دقيقة")
		}
	case 2:
		switch unit {
		case DAY:
			return p.Sprintf("%s", "يومان")
		case HOUR:
			return p.Sprintf("%s", "ساعات")
		case MINUTE:
			return p.Sprintf("%s", "دقيقتان")
		}
	case 3, 4, 5, 6, 7, 8, 9, 10:
		switch unit {
		case DAY:
			return p.Sprintf("%d %s", value, "أيام")
		case HOUR:
			return p.Sprintf("%d %s", value, "ساعات")
		case MINUTE:
			return p.Sprintf("%d %s", value, "دقائق")
		}
	default:
		switch unit {
		case DAY:
			return p.Sprintf("%d %s", value, "يوما")
		case HOUR:
			return p.Sprintf("%d %s", value, "ساعة")
		case MINUTE:
			return p.Sprintf("%d %s", value, "دقيقة")
		}
	}
	return ""
}
