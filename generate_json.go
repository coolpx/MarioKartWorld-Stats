package main

import (
	"os"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"fmt"
)

var fields = []string { "name", "speed_ground", "speed_water", "speed_air", "speed_antigravity", "acceleration", "weight", "handling_ground", "handling_water", "handling_air", "handling_antigravity", "traction", "mini_turbo", "trailing"}

// trailing column types (tires and gliders have no trailing column)
const (
	trailingNone = ""
	trailingType = "type"
	trailingVehicleSize = "vehicle_size"
)


type MultipartElement struct {
	Ground float64 `json:"ground"`
	Water float64 `json:"water"`
	Air float64 `json:"air"`
	AntiGravity float64 `json:"antigravity"`
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
	element.Ground = parseFloat64(fields[0])
	element.Water = parseFloat64(fields[1])
	element.Air = parseFloat64(fields[2])
	element.AntiGravity = parseFloat64(fields[3])
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

func generateFile(input string, output string, trailing string) error {
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
		outrecord.Speed = parseSpeed(record[fieldIndex("speed_ground"):])
		outrecord.Acceleration = parseFloat64(record[fieldIndex("acceleration")])
		outrecord.Weight = parseFloat64(record[fieldIndex("weight")])
		outrecord.Handling = parseHandling(record[fieldIndex("handling_ground"):])
		outrecord.Traction = parseFloat64(record[fieldIndex("traction")])
		outrecord.MiniTurbo = parseFloat64(record[fieldIndex("mini_turbo")])
		outrecord.Invincibility = parseFloat64(record[fieldIndex("invincibility")])
		switch (trailing) {
		case trailingType:
			outrecord.Type = record[fieldIndex("trailing")]
		case trailingVehicleSize:
			outrecord.VehicleSize = record[fieldIndex("trailing")]
		}
		
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

func readFile(name string, trailing string) {
	input := fmt.Sprintf("%s.csv", name)
	output := fmt.Sprintf("json/%s.json", name)
	err := generateFile(input, output, trailing)
	if (err != nil) {
		fmt.Println(err)
	}
}

func main() {
	readFile("characters", trailingVehicleSize)
	readFile("bodies", trailingNone)
	readFile("tires", trailingNone)
	readFile("gliders", trailingNone)
}
