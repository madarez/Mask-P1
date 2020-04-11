package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/psilva261/timsort/v2"
)

// GPS data structure holding lat and long in radians, i.e., phi and lambda
type GPSSentence [2]float64

// stringify GPS data
func (g *GPSSentence) String() string {

	var lat, long string
	var d float64

	if d = GetDegrees(g[0]); d < 0 {
		lat = fmt.Sprintf("%8.6fS", -d)
	} else {
		lat = fmt.Sprintf("%8.6fN", d)
	}

	if d = GetDegrees(g[1]); d < 0 {
		long = fmt.Sprintf("%8.6fW", -d)
	} else {
		long = fmt.Sprintf("%8.6fE", d)
	}

	return fmt.Sprintf("[%v %v]", lat, long)
}

// A mapping of time to GPS data
type LocHistory map[int]*GPSSentence

// order a location history by time
func (trh *LocHistory) Sort() []int {
	keys := make([]int, 0, len(*trh))
	for k := range *trh {
		keys = append(keys, k)
	}
	timsort.Ints(keys, func(a, b int) bool { return a < b })
	return keys
}

// A mapping of individuals' id to their location history
type IndivTrack map[string]*LocHistory

// construct an individual track collection
func NewIndivTrack() *IndivTrack {
	var it IndivTrack = make(map[string]*LocHistory)
	return &it
}

// get the location of an individual at a timestamp
func (it *IndivTrack) Get(id string, t int) *GPSSentence {
	return (*(*it)[id])[t]
}

// delete a timestamp key from an individual's location history
func (it *IndivTrack) DeleteTime(id string, t int) {
	delete(*(*it)[id], t)
}

// delete an ID from individual track collections
func (it *IndivTrack) DeleteId(id string) {
	delete(*it, id)
}

// push a row from csv into individual track collection
func (it *IndivTrack) Push(record []string) {

	var err error

	time_nano, err := strconv.ParseInt(record[3], 10, 64)
	if err != nil {
		log.Panicf("Time %v is not reconstructable: %v", record[3], err)
	}

	var time_sec int

	// if the timestamp is outdated skip
	if time_nano = (time_nano - checkpoint) / int64(time.Second); math.MinInt32 < time_nano {
		time_sec = int(time_nano)
	} else {
		return
	}

	id := record[0]

	lat, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		log.Panicf("Lat %v is not reconstructable: %v", record[1], err)
	}

	long, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		log.Panicf("Long %v is not reconstructable: %v", record[2], err)
	}

	lh, ok := (*it)[id]
	if !ok { // if location history have not been created yet, initialize it
		var tmp LocHistory = make(map[int](*GPSSentence), 1)
		(*it)[id], lh = &tmp, &tmp
	}
	// store data in radiants
	(*lh)[time_sec] = &GPSSentence{GetRadiant(lat), GetRadiant(long)}
}

// order the individual track collection by time
func (it *IndivTrack) Sort() map[string][]int {
	ids := make(map[string][]int, len(*it))
	for id := range *it {
		ids[id] = (*it)[id].Sort()
	}
	return ids
}
