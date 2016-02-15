package main

import (
	"flag"
	"fmt"

	"github.com/zoranzaric/withings.go/sleep"
	"github.com/zoranzaric/withings.go/weight"
)

func main() {
	weightPath := flag.String("weight", "", "path to weight file")
	sleepPath := flag.String("sleep", "", "path to sleep file")
	flag.Parse()

	if *weightPath != "" {
		for w := range weight.Parse(*weightPath) {
			fmt.Printf("%v\n", w.ToInfluxDBInsertString())
		}
	}

	if *sleepPath != "" {
		for s := range sleep.Parse(*sleepPath) {
			fmt.Printf("%v\n", s.ToInfluxDBInsertString())
		}
	}
}
