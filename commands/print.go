// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"math"
	"time"
)

var statusCodeDist map[int]int = make(map[int]int)

func (b *Boom) Print() {
	total := b.end.Sub(b.start)
	var avgTotal int64
	var fastest, slowest time.Duration

	for {
		select {
		case r := <-b.results:
			statusCodeDist[r.statusCode]++

			avgTotal += r.duration.Nanoseconds()
			if fastest.Nanoseconds() == 0 || r.duration.Nanoseconds() < fastest.Nanoseconds() {
				fastest = r.duration
			}
			if r.duration.Nanoseconds() > slowest.Nanoseconds() {
				slowest = r.duration
			}
		default:
			rps := float64(b.N) / total.Seconds()
			fmt.Printf("\nSummary:\n")
			fmt.Printf("  total:\t%v secs\n", total.Seconds())
			fmt.Printf("  slowest:\t%v secs\n", slowest.Seconds())
			fmt.Printf("  fastest:\t%v secs\n", fastest.Seconds())
			fmt.Printf("  average:\t%v secs\n", float64(avgTotal)/float64(b.N)*math.Pow(10, 9)) // TODO: in seconds
			fmt.Printf("  requests/sec:\t%v\n", rps)
			fmt.Printf("  speed index:\t%v\n", speedIndex(rps))
			b.printStatusCodes()
			return
		}
	}
}

func (b *Boom) printStatusCodes() {
	fmt.Printf("\nStatus code distribution:\n")
	for code, num := range statusCodeDist {
		fmt.Printf("  [%d]\t%d responses\n", code, num)
	}
}

func speedIndex(rps float64) string {
	if rps > 500 {
		return "Whoa, pretty neat"
	} else if rps > 100 {
		return "Pretty good"
	} else if rps > 50 {
		return "Meh"
	} else {
		return "Hahahaha"
	}
}