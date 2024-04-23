package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowTime, err := time.Parse("20060102", now)
	checkErr(w, err, "parsing string now to date")

	dateTime, err := time.Parse("20060102", date)
	checkErr(w, err, "parsing string date to date")

	err = errorIfEmpty(repeat)
	checkErr(w, err, "repeat not empty")

	nextDate, err := CalculateNextDate(dateTime, nowTime, repeat)
	checkErr(w, err, "bad calculating")

	log.Println("[Info] FOR now =", now, "date =", date, "repeat =", repeat)
	log.Println("[Info] nextDate =", nextDate.Format("20060102"), "err =", err)

	w.Write([]byte(nextDate.Format("20060102")))
}

func CalculateNextDate(date, now time.Time, repeat string) (time.Time, error) {
	nextDate := date
	repeat_array := strings.Split(repeat, " ")
	switch repeat_array[0] {
	case "y":
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if nextDate.After(now) {
				break
			}
		}
	case "d":
		if len(repeat_array) != 2 {
			return nextDate, errors.New("wrong repeat size")
		}
		days, err := strconv.Atoi(repeat_array[1])
		if err != nil {
			return nextDate, err
		}
		if days > 400 {
			return nextDate, errors.New("more than 400 days")
		}
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				break
			}
		}
	default:
		return nextDate, errors.New("wrong repeat format")
	}
	return nextDate, nil
}

func errorIfEmpty(s string) error {
	if s == "" {
		return errors.New("empty repeat")
	}
	return nil
}

func checkErr(w http.ResponseWriter, err error, s string) {
	if err != nil {
		log.Println("[Error] Failed: " + s)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("[Info] Success: " + s)
}
