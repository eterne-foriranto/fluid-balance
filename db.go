package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const BaseUrl = "http://127.0.0.1:9088/api/v1/db/ldb/namespaces/"

type Db struct {
	Name string `json:"name"`
}

type IdEntity struct {
	Id int `json:"id"`
}

type Aggregation struct {
	Value  int      `json:"value"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

type AbstractResponse struct {
	Aggregations []Aggregation `json:"aggregations"`
	Namespaces   []string      `json:"namespaces"`
	CacheEnabled bool          `json:"cache_enabled"`
	TotalItems   int           `json:"total_items"`
}

type EventResponse struct {
	AbstractResponse
	Items []map[string]string
}

type DrinkResponse struct {
	AbstractResponse
	Items []map[string]int
}

type Response struct {
	Items        []IdEntity    `json:"items"`
	Aggregations []Aggregation `json:"aggregations"`
	Namespaces   []string      `json:"namespaces"`
	CacheEnabled bool          `json:"cache_enabled"`
	TotalItems   int           `json:"total_items"`
}

type ResponseDrink struct {
	Items        []map[string]string `json:"items"`
	Aggregations []Aggregation       `json:"aggregations"`
	Namespaces   []string            `json:"namespaces"`
	CacheEnabled bool                `json:"cache_enabled"`
	TotalItems   int                 `json:"total_items"`
}

type Drink struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func getLastId(ns string) int {
	url := BaseUrl + ns + "/items?limit=1&sort_field=id&sort_order=desc&fields=id&format=json"
	resp, err := http.Get(url)
	treatErr(err)
	body := resp.Body
	defer body.Close()
	entireBody, err := io.ReadAll(body)
	var result Response
	err = json.Unmarshal(entireBody, &result)
	treatErr(err)
	return result.Items[0].Id
}

func postDrink(drink string) {
	url := BaseUrl + "drinks/items"
	id := getLastId("drinks") + 1
	drinkItem := Drink{id, drink}
	drinkJson, err := json.Marshal(drinkItem)
	treatErr(err)
	drinkReader := strings.NewReader(string(drinkJson))
	_, err = http.Post(url, "application/json", drinkReader)
}

func getDrinks() []string {
	url := BaseUrl + "drinks/items?format=json"
	resp, err := http.Get(url)
	treatErr(err)
	body := resp.Body
	defer body.Close()
	entireBody, err := io.ReadAll(body)
	var result ResponseDrink
	err = json.Unmarshal(entireBody, &result)
	treatErr(err)
	drinks := make([]string, 0)
	for _, item := range result.Items {
		drinks = append(drinks, item["name"])
	}
	return drinks
}

func putDrink(newDrink string) {
	oldDrinks := getDrinks()
	for _, oldDrink := range oldDrinks {
		if newDrink == oldDrink {
			return
		}
	}
	postDrink(newDrink)
}

func postEvent(event *Event) {
	eventJson, err := json.Marshal(event)
	treatErr(err)
	eventReader := strings.NewReader(string(eventJson))
	url := BaseUrl + "ns/items"
	_, err = http.Post(url, "application/json", eventReader)
}

func httpGet(ns string, qs map[string]string) []byte {
	url := BaseUrl + ns + "/items?"
	params := make([]string, 0)
	for key, value := range qs {
		params = append(params, fmt.Sprintf("%v=%v", key, value))
	}
	url += strings.Join(params, "&")
	resp, err := http.Get(url)
	treatErr(err)
	body := resp.Body
	defer body.Close()
	entireBody, err := io.ReadAll(body)
	treatErr(err)
	return entireBody
}

func getEvents(ns string, filter string, fields string) []byte {
	qs := map[string]string{
		"filter": filter,
		"fields": fields,
	}
	rawResponse := httpGet(ns, qs)
	return rawResponse
}

func getPItems() []map[string]string {
	qs := map[string]string{
		"filter": "type.name='p'",
		"fields": "time",
	}
	rawResponse := httpGet("ns", qs)
	var result EventResponse
	err := json.Unmarshal(rawResponse, &result)
	treatErr(err)
	return result.Items
}

func getDrinkEvents() []map[string]int {
	bytes := getEvents("ns", "type.name='drink'", "volume")
	var result DrinkResponse
	err := json.Unmarshal(bytes, &result)
	treatErr(err)
	return result.Items
}
