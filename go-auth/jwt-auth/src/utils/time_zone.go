package utils

import "time"

func TimeNowBrazil() time.Time {
	brazilLocation, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		InternalServerError("Error while loading Brasília time zone: " + err.Error())
	}
	return time.Now().In(brazilLocation)
}

func TimeBrazil(valueTime time.Time) time.Time {
	brazilLocation, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		InternalServerError("Error while loading Brasília time zone: " + err.Error())
	}
	return valueTime.In(brazilLocation)
}
