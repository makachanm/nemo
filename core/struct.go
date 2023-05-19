package build

type TimeStamp struct {
	Year  int
	Month int
	Day   int
	Hour  int
	Min   int
}

func MakeTimeStamp(year int, month int, day int, hour int, min int) TimeStamp {
	return TimeStamp{
		Year:  year,
		Month: month,
		Day:   day,
		Hour:  hour,
		Min:   min,
	}
}

func (t *TimeStamp) StampSize() int {
	return t.Year + t.Month + t.Day + t.Hour + t.Min
}

func (t *TimeStamp) isBiggerStamp(src TimeStamp, cmp TimeStamp) bool {
	if src.Year > cmp.Year {
		return true
	} else if src.Month > cmp.Month {
		return true
	} else if src.Month == cmp.Month {
		if src.Day > cmp.Day {
			return true
		} else if src.Hour > cmp.Hour {
			return true
		} else if src.Min > cmp.Min {
			return true
		}
	}

	return false
}
