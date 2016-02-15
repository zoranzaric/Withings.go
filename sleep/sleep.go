package sleep

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

type Sleep struct {
	From    time.Time
	Till    time.Time
	Light   time.Duration
	Deep    time.Duration
	Rem     time.Duration
	Awake   time.Duration
	Wakeups int
}

func parseSleep(line []string) (Sleep, error) {
	sleep := Sleep{}
	if line[0] != "" {
		t, err := util.ParseTime(line[0])
		if err != nil {
			return sleep, err
		}
		sleep.From = t
	}

	if line[1] != "" {
		t, err := util.ParseTime(line[1])
		if err != nil {
			return sleep, err
		}
		sleep.Till = t
	}

	if line[2] != "" {
		d, err := time.ParseDuration(line[2] + "s")
		if err != nil {
			return sleep, err
		}
		sleep.Light = d
	}

	if line[3] != "" {
		d, err := time.ParseDuration(line[3] + "s")
		if err != nil {
			return sleep, err
		}
		sleep.Deep = d
	}

	if line[4] != "" {
		d, err := time.ParseDuration(line[4] + "s")
		if err != nil {
			return sleep, err
		}
		sleep.Rem = d
	}

	if line[5] != "" {
		d, err := time.ParseDuration(line[5] + "s")
		if err != nil {
			return sleep, err
		}
		sleep.Awake = d
	}

	if line[6] != "" {
		w, err := strconv.Atoi(line[6])
		if err != nil {
			return sleep, err
		}
		sleep.Wakeups = w
	}

	return sleep, nil
}

func Parse(filePath string) chan Sleep {
	c := make(chan Sleep)
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

			sleep, err := parseSleep(record)

			c <- sleep
		}
		close(c)
	}()

	return c
}

func (s Sleep) ToInfluxDBInsertString() string {
	return fmt.Sprintf("INSERT sleep start=%d,end=%d,light=%.0f,deep=%.0f,rem=%.0f,awake=%.0f,wakeups=%d %d", s.From.Unix(), s.Till.Unix(), s.Light.Seconds(), s.Deep.Seconds(), s.Rem.Seconds(), s.Awake.Seconds(), s.Wakeups, s.Till.UnixNano())
}
