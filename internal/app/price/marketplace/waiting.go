package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func getWaitTime(now time.Time) time.Duration {
	curMinutes := now.Minute()
	callPeriod := 30
	var callMinutes int

	if curMinutes < callPeriod {
		callMinutes = callPeriod - curMinutes
	} else {
		callMinutes = callPeriod - (curMinutes - callPeriod)
	}

	callTime := getCallTimeFromMinutes(now, callMinutes)
	variation := 10

	return custTime.GetWaitDurWithRandomComp(now, callTime, variation)
}

func getCallTimeFromMinutes(now time.Time, minutes int) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minutes, 0, 0, now.Location())
}