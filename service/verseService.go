package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sully/vid-gen-v2/dto"
)

type VerseService struct {
	TempAudioPath    string
	TempImagepath    string
	IntermediatePath string
	OutputPath       string

	audioUrl       string
	imageUrl       string
	randomVerseUrl string
	csvVerses      []dto.VerseCSV
}

func NewVerseService(tempAudioPath, tempImagePath, outputPath string) (vs *VerseService) {
	var err error
	vs = new(VerseService)
	vs.OutputPath = outputPath
	vs.TempAudioPath = tempAudioPath
	vs.TempImagepath = tempImagePath
	vs.IntermediatePath = "temp/intermediate.mp4"
	vs.audioUrl = "https://api.quran.com/api/v4/recitations/3/by_ayah/"
	vs.imageUrl = "https://source.unsplash.com/random/200x300/?nature"
	vs.randomVerseUrl = "https://api.quran.com/api/v4/verses/random?language=ar&translations=131"
	vs.csvVerses, err = LoadVerses()
	if err != nil {
		log.Fatal("Error loading csv")
	}
	return
}

func (vs *VerseService) GetVerseAndDownloadFiles() (*Verse, error) {

	var verseResponse dto.VerseResponse
	var audioResponse dto.AudioFiles

	resp, err := http.Get(vs.randomVerseUrl)

	if err != nil {
		fmt.Println("Error getting random verse")
		return &Verse{}, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&verseResponse); err != nil {
		fmt.Println("Error dencoding verse")
		return &Verse{}, err
	}

	audioResp, err := http.Get(vs.audioUrl + verseResponse.Verse.VerseKey)

	if err != nil {
		fmt.Println("Error getting audio")
		return &Verse{}, err
	}

	if err := json.NewDecoder(audioResp.Body).Decode(&audioResponse); err != nil {
		fmt.Println("Error decoding audio")
		return &Verse{}, err
	}

	verse := new(Verse)
	verse.VerseKey = verseResponse.Verse.VerseKey
	verse.VerseNumber = verseResponse.Verse.Id
	verse.AudioUrl = "https://verses.quran.com/" + audioResponse.Audios[0].Url
	verse.Translation = vs.csvVerses[verseResponse.Verse.Id-1].Translation
	fmt.Println(verse.AudioUrl)

	vs.downloadTempFiles(verse.AudioUrl)

	if err := Convert(*verse, vs.OutputPath, vs.TempImagepath, vs.TempAudioPath, vs.IntermediatePath); err != nil {
		return &Verse{}, err
	}

	defer vs.deleteTempFiles()
	return verse, nil
}

func (vs *VerseService) downloadTempFiles(audioUrl string) error {
	if err := DownloadFile(vs.TempAudioPath, audioUrl); err != nil {
		return err
	}
	if err := DownloadFile(vs.TempImagepath, vs.imageUrl); err != nil {
		return err
	}

	return nil
}

func (vs *VerseService) deleteTempFiles() {
	os.Remove(vs.TempAudioPath)
	os.Remove(vs.TempImagepath)
	os.Remove(vs.IntermediatePath)
}

type Verse struct {
	VerseKey    string
	VerseNumber int
	AudioUrl    string
	Translation string
}
