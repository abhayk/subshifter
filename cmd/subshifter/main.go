package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	file := flag.String("file", "", "The path to the subtitle file.")
	offset := flag.Int("offset", 0, "The shift offset. To increment by half a second provide +500. To decrement -500")
	flag.Parse()

	// check for mandatory arguments
	if strings.Compare(*file, "") == 0 {
		fmt.Println("Usage: subshifter <file> <offset>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputFile, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.Create(*file + ".tmp")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-->") {
			line = updateTimesWithOffset(line, *offset)
		}
		_, err := fmt.Fprintln(outputFile, line)
		if err != nil {
			log.Fatal(err)
		}
	}

	inputFile.Close()
	outputFile.Close()

	// if the original file was abc.srt then back it up with the name abc.srt.orig
	err = os.Rename(*file, *file+".orig")
	if err != nil {
		log.Fatal(err)
	}
	// rename the tmp file to the original subtitle file
	err = os.Rename(*file+".tmp", *file)
	if err != nil {
		log.Fatal(err)
	}
}

func updateTimesWithOffset(line string, offset int) string {
	split := strings.Split(line, "-->")
	start := addOffset(parseTime(split[0]), offset)
	end := addOffset(parseTime(split[1]), offset)
	return timeToString(start) + " --> " + timeToString(end)
}

func timeToString(inputTime time.Time) string {
	return strings.ReplaceAll(inputTime.Format("15:04:05.000"), ".", ",")
}

func addOffset(inputTime time.Time, offsetMs int) time.Time {
	return inputTime.Add(time.Millisecond * time.Duration(offsetMs))
}

func parseTime(str string) time.Time {
	time, err := time.Parse("15:04:05", strings.ReplaceAll(strings.TrimSpace(str), ",", "."))
	if err != nil {
		log.Fatal(err)
	}
	return time
}
