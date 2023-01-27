package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	os_exec "os/exec"
	"strings"
)

func Convert(verse Verse, outputPath, tempImagePath, tempAudioPath, intermediatePath string) error {

	writeToFile("translation.txt", ReplaceSubstring(ReplaceSubstring(newLineAfterThreeWords(verse.Translation), "˹", " "), "˺", " "))

	template := "ffmpeg -y -r 25 -i %s -i %s %s"
	// template2 := "ffmpeg -y -i %s -vf \"drawtext=fontfile=/font/Amejo.ttf:text='Quran\n" + strings.Join(strings.Split(verse.VerseKey, ":"), "-") + "':fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=(w-text_w)/2:y=(h-text_h)/2\" -codec:a copy %s"
	var test string
	if len(verse.Translation) > 1 {
		test = "ffmpeg -y -i %s -vf \"drawtext=textfile=translation.txt:x=(w-text_w)/2:y=(h-text_h)/2: fontsize=15:fontcolor=yellow@0.9: box=1: boxcolor=black@0.6\" -c:a copy output.mp4"
	} else {
		test = "ffmpeg -y -i %s -vf \"drawtext=textfile=translation.txt:x=(w-text_w)/2:y=(h-text_h)/2: fontsize=5:fontcolor=yellow@0.9: box=1: boxcolor=black@0.6\" -c:a copy output.mp4"
	}

	if err := exec(template, tempImagePath, tempAudioPath, intermediatePath); err != nil {
		fmt.Println("Error converting")
		return err
	}
	if err := exec(test, intermediatePath); err != nil {
		fmt.Println("Error converting")
		return err
	}

	return nil
}

func writeToFile(filePathstring, text string) {
	f, err := os.Create(filePathstring)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(text)

	if err2 != nil {
		log.Fatal(err)
	}
}

func splitBy(s string, n int) []string {
	var ss []string
	for i := 1; i < len(s); i++ {
		if i%n == 0 {
			ss = append(ss, s[:i])
			s = s[i:]
			i = 1
		}
	}
	ss = append(ss, s)
	return ss
}

func newLineAfterThreeWords(s string) string {
	words := strings.Fields(s)
	var result string
	for i, word := range words {
		result += word
		if (i+1)%3 == 0 && i != len(words)-1 {
			result += "\n"
		} else if i != len(words)-1 {
			result += " "
		}
	}
	return result
}

func ReplaceSubstring(s string, oldString string, newString string) string {
	return strings.Replace(s, oldString, newString, -1)
}

func exec(template string, params ...interface{}) error {
	cmd := fmt.Sprintf(template, params...)
	fmt.Println("Running command : " + cmd)

	interpreter := ""
	interpreterArgs := ""

	interpreter = "sh"
	interpreterArgs = "-c"
	c, err := os_exec.Command(interpreter, interpreterArgs, cmd).CombinedOutput()

	fmt.Printf("ffmpeg out:\n%s\n", string(c))

	return err
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
