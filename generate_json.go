package main

import (
	"os"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"fmt"
)

var fields = []string { "name", "speed_paved", "speed_offroad", "speed_water", "acceleration", "mini_turbo", "weight", "handling_paved", "handling_offroad", "handling_water"}

type MultipartElement struct {
	Paved float64 `json:"paved"`
	Offroad float64 `json:"offroad"`
	Water float64 `json:"water"`
}

type Speed struct {
	MultipartElement
}

type Handling struct {
	MultipartElement
}

type Record struct {
	Name string `json:"name"`
	Speed Speed `json:"speed"`
	Acceleration float64 `json:"acceleration"`
	Weight float64 `json:"weight"`
	Handling Handling `json:"handling"`
	Traction float64 `json:"traction"`
	MiniTurbo float64 `json:"mini_turbo"`
	Invincibility float64 `json:"invincibility"`
	Type string `json:"type,omitempty"`
	VehicleSize string `json:"vehicle_size,omitempty"`
}

func fieldIndex(field string) (int) {
	for i, v := range(fields) {
		if (v == field) {
			return i
		}
	}
	return 0
}

func parseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (element *MultipartElement) parseMultipartElement(fields []string) {
	element.Paved = parseFloat64(fields[0])
	element.Offroad = parseFloat64(fields[1])
	element.Water = parseFloat64(fields[2])
}

func parseSpeed(fields []string) Speed {
	var s Speed
	s.parseMultipartElement(fields)
	return s
}
func parseHandling(fields []string) Handling {
	var h Handling
	h.parseMultipartElement(fields)
	return h
}

func generateFile(input string, output string) error {
	infile, err := os.Open(input)
	if (err != nil) {
		return err
	}
	
	// Read in the CSV file
	reader := csv.NewReader(infile)
	records, err := reader.ReadAll()
	infile.Close()
	
	// Iterate through the records
	records = records[1:]
	outrecords := make([]Record, len(records))
	for idx, record := range(records) {
		var outrecord Record
		
		outrecord.Name = record[fieldIndex("name")]
		outrecord.Speed = parseSpeed(record[fieldIndex("speed_paved"):])
		outrecord.Acceleration = parseFloat64(record[fieldIndex("acceleration")])
		outrecord.Weight = parseFloat64(record[fieldIndex("weight")])
		outrecord.Handling = parseHandling(record[fieldIndex("handling_paved"):])
		outrecord.Traction = parseFloat64(record[fieldIndex("traction")])
		outrecord.MiniTurbo = parseFloat64(record[fieldIndex("mini_turbo")])
		outrecord.Invincibility = parseFloat64(record[fieldIndex("invincibility")])
		
		outrecords[idx] = outrecord
	}
	
	outfile, err := os.Create(output)
	if (err != nil) {
		return err
	}
	b, err := json.MarshalIndent(outrecords, "", "  ")
	outfile.Write(b)
	outfile.Close()
	
	return err
}

func readFile(name string) {
	input := fmt.Sprintf("%s.csv", name)
	output := fmt.Sprintf("json/%s.json", name)
	err := generateFile(input, output)
	if (err != nil) {
		fmt.Println(err)
	}
}

func main() {
	readFile("characters")
	readFile("vehicles")
}
