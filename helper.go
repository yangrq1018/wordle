package wordle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/fatih/color"
)

// Terminal support

func cls(out io.Writer) {
	fmt.Fprintf(out, "\033[H\033[2J")
}

var (
	white  = color.New(color.FgWhite)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
)

func printer(i uint) (printer *color.Color) {
	switch i {
	case 1:
		printer = yellow
	case 2:
		printer = green
	default:
		printer = white
	}
	return printer
}

const wordAPI = "https://api.dictionaryapi.dev/api/v2/entries/en/"

type meaningOfWord struct {
	Word      string      `json:"word"`
	Phonetics []phonetics `json:"phonetics"`
	Meanings  []meanings  `json:"meanings"`
}
type phonetics struct {
	Text  string `json:"text"`
	Audio string `json:"audio,omitempty"`
}

type meanings struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Definitions  []struct {
		Definition string        `json:"definition"`
		Example    string        `json:"example"`
		Synonyms   []interface{} `json:"synonyms"`
		Antonyms   []interface{} `json:"antonyms"`
	} `json:"definitions"`
}

// Word meaing API
func getMeaningOfWord(w word) string {
	res, err := http.Get(wordAPI + w.String())
	if err != nil {
		return err.Error()
	}

	var means []meaningOfWord
	_ = json.NewDecoder(res.Body).Decode(&means)
	res.Body.Close()

	if len(means) == 0 {
		return "No meaning found"
	}
	return buildExplanation(w, means[0])
}

func buildExplanation(w word, mean meaningOfWord) string {
	tmpl := template.New("word_meaning_template")
	tmpl, _ = tmpl.Parse(`Meaning of {{ .Word }}
{{- range .Mean.Meanings}}
- As {{ .PartOfSpeech }}, {{ (index .Definitions 0).Definition }}
{{- if (index .Definitions 0).Example }}
  Example: {{ (index .Definitions 0).Example }}
{{- end }}
{{- end }}
	`)

	buf := bytes.NewBuffer(nil)
	err := tmpl.Execute(buf, map[string]interface{}{
		"Word": w.String(),
		"Mean": mean,
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}
