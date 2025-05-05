package enums

type FormatTimeEnum struct{}

var FormatTime FormatTimeEnum

func (FormatTimeEnum) Default() string {
	return "2006-01-02 15:04:05.000 -0700 -07"
}

func (FormatTimeEnum) DataHour() string {
	return "2006-01-02T15:04:05"
}
