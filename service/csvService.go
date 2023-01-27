package service

import (
	"fmt"
	"os"
	"sully/vid-gen-v2/dto"

	"github.com/gocarina/gocsv"
)

func LoadVerses() (csvVerses []dto.VerseCSV, e error) {
	const fileName string = "quran-dataset.csv"
	in, err := os.Open(fileName)
	if err != nil {
		return csvVerses, fmt.Errorf("error opening file: %s", err.Error())
	}
	defer in.Close()

	if err := gocsv.UnmarshalFile(in, &csvVerses); err != nil {
		return csvVerses, fmt.Errorf("error unmarshalling csv file: %s", err.Error())
	}

	return
}
