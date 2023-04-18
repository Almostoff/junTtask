package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

func main() {
	URL := "http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=11/11/2020"

	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	valCurs := &ValCurs{}

	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}
	err = decoder.Decode(&valCurs)
	if err != nil {
		fmt.Println("Error decoding XML:", err)
		return
	}

	// Ищю максимальное и минимальное значение курса валюты, а также среднее значение курса рубля
	var maxValute Valute
	var minValute Valute
	var totalValue float64
	for i, valute := range valCurs.Valutes {
		value := strings.ReplaceAll(valute.Value, ",", ".")
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Error parsing value for %s: %v\n", valute.CharCode, err)
			continue
		}
		if i == 0 {
			maxValute = valute
			minValute = valute
		} else {
			maxValue := strings.ReplaceAll(maxValute.Value, ",", ".")
			max, _ := strconv.ParseFloat(maxValue, 64)
			if val > max {
				maxValute = valute
			}
			minValue := strings.ReplaceAll(minValute.Value, ",", ".")
			min, _ := strconv.ParseFloat(minValue, 64)
			if val < min {
				minValute = valute
			}
		}
		totalValue += val
	}

	averageValue := totalValue / float64(len(valCurs.Valutes))

	fmt.Printf("Max value: %s (%s) on %s\n", maxValute.Value, maxValute.Name, valCurs.Date)
	fmt.Printf("Min value: %s (%s) on %s\n", minValute.Value, minValute.Name, valCurs.Date)
	fmt.Printf("Average value: %.4f\n", averageValue)
}
