package dto

type VerseResponse struct {
	Verse *Verse `json:"verse"`
}

type VerseCSV struct {
	Translation string `csv:"ayah_en"`
}

type Verse struct {
	Id       int    `json:"id"`
	VerseKey string `json:"verse_key"`
}

type Translation struct {
	Text string `json:"text"`
}

type AudioFiles struct {
	Audios []Audio `json:"audio_files"`
}

type Audio struct {
	Url string `json:"url"`
}
