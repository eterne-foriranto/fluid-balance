package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"time"
)

const Begin = "2023-06-26T11:05:55.869404707+03:00"
const Layout = time.RFC3339Nano

func main() {
	runBot()
}

func treatErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getConfigValue(sectionName string, key string) string {
	cnf, err := config.NewConfig("ini", "config.ini")

	if err != nil {
		panic(err)
	}

	section, err := cnf.GetSection(sectionName)

	if err != nil {
		panic(err)
	}
	return section[key]
}

type EventType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Event struct {
	ID     int       `json:"id"`
	Type   EventType `json:"type"`
	Time   time.Time `json:"time"`
	Drink  string    `json:"drink"`
	Volume int       `json:"volume"`
}
