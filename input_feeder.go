package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

// reading data from a csv source
func feeder(it *IndivTrack) <-chan bool {
	csvFile, _ := os.Open("sample.csv")
	r := csv.NewReader(bufio.NewReader(csvFile))
	r.FieldsPerRecord = 4 // expecting four fields of id, geo lat and long, and timestamp
	r.ReuseRecord = true
	finished := make(chan bool)

	go func() {
		var err error
		var record []string
		for {
			// Read each record from csv
			record, err = r.Read()
			if err == io.EOF {
				break
			}
			it.Push(record)
		}
		close(finished)
	}()
	return finished
}
