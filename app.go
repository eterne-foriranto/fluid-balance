package main

import (
	"fmt"
	"time"
)

func getBegin() time.Time {
	begin, err := time.Parse(Layout, Begin)
	treatErr(err)
	return begin
}

func getDurationHours() float64 {
	return time.Now().Sub(getBegin()).Hours()
}

func getPFreq() string {
	freq := float64(len(getPItems())) / getDurationHours() * 24
	return fmt.Sprintf("%v per day", freq)
}

func getPPeriod() string {
	period := getDurationHours() / float64(len(getPItems()))
	return fmt.Sprintf("every %v hours", period)
}

func getRate() string {
	total := 0
	for _, event := range getDrinkEvents() {
		total += event["volume"]
	}
	rate := float64(total) / getDurationHours() * 24
	return fmt.Sprintf("%v ml per day", rate)
}
