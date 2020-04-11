package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tidwall/rbang"
)

var allTracks = NewIndivTrack()

var checkpoint int64 = 1586186202000000000 //time.Now().UnixNano()

func main() {
	<-feeder(allTracks)
	s := allTracks.Sort()
	var tr rbang.RTree // data struct to hold the geographical data
	var ts []int
	var i int
	var oldSuffix, newSuffix string // query time suffix designated to ids
	var oldLoc, omin, omax, newLoc, nmin, nmax *GPSSentence
	for t := 0; t < 5; t++ {
		oldSuffix, newSuffix = strconv.Itoa(t-1), strconv.Itoa(t)
		// Prune old times from s and allTracks and take the current location as centre
		for id := range s {
			ts = s[id]
			oldLoc = allTracks.Get(id, ts[0])
			switch cap(ts) {
			case 0:
				log.Panicf("Time arrays belonging to ID %v is empty!", id)
			case 1:
				newLoc = allTracks.Get(id, ts[0])
			default:
				for i = 0; i < cap(ts)-1; i++ {
					if ts[i+1] <= t {
						allTracks.DeleteTime(id, ts[i])
						// (*s)[id] = ts[i+1:]
					} else {
						break
					}
				}
				ts, s[id] = ts[i:], ts[i:]
				switch cap(s[id]) {
				case 0:
					allTracks.DeleteId(id)
					delete(s, id)
				case 1:
					newLoc = allTracks.Get(id, ts[0])
				default:
					t1 := ts[0]
					g1 := allTracks.Get(id, t1)
					if t1 >= t {
						newLoc = g1
					} else {
						t2 := ts[1]
						g2 := allTracks.Get(id, t2)
						newLoc = InterpolatePath(g1, g2, t1, t2, t)
					}
				}
			}
			// check the intersections of id at current location with location at current timestamp seen so far
			tr.Search(*newLoc, *newLoc,
				func(min, max [2]float64, value interface{}) bool {
					if otherId := value.(string); strings.HasSuffix(otherId, newSuffix) {
						fmt.Printf("[%v]: id %v was close to id %v\n",
							t, id, strings.TrimSuffix(otherId, newSuffix))
					}
					return true
				},
			)
			// updating the new location into the RTree
			omin, omax = ProximateBox(oldLoc)
			nmin, nmax = ProximateBox(newLoc)
			tr.Replace(*omin, *omax, id+oldSuffix,
				*nmin, *nmax, id+newSuffix)
		}
	}
}
