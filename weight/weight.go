package weight

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/zoranzaric/withings.go/util"
)

type Weight struct {
	Time    time.Time
	Weight  float64
	Fat     float64
	Comment string
}

func parseWeight(line []string) (Weight, error) {
	weight := Weight{}

	if line[0] != "" {
		t, err := util.ParseTime(line[0])
		if err != nil {
			return weight, err
		}
		weight.Time = t
	}

	if line[1] != "" {
		w, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return weight, err
		}
		weight.Weight = w
	}

	if line[2] != "" {
		fat, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return weight, err
		}
		weight.Fat = fat
	}

	if line[4] != "" {
		weight.Comment = line[4]
	}

	return weight, nil
}
func Parse(filePath string) chan Weight {
	c := make(chan Weight)

	go func() {
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		reader := csv.NewReader(f)
		defer f.Close()
		reader.Comma = ','

		readFirstLine := false

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if !readFirstLine {
				readFirstLine = true
				continue
			}

			weight, err := parseWeight(record)
			if err != nil {
				log.Fatal(err)
			}

			c <- weight
		}
		close(c)
	}()

	return c
}

func (w Weight) ToInfluxDBInsertString() string {
	if w.Fat == 0 {
		return fmt.Sprintf("INSERT weight weight=%f %d", w.Weight, w.Time.UnixNano())
	} else {
		return fmt.Sprintf("INSERT weight weight=%f,fat=%f,fat_p=%.2f %d", w.Weight, w.Fat, w.Fat/w.Weight*100, w.Time.UnixNano())
	}
}
