package wordle

import (
	"bufio"
	"os"

	"github.com/zyedidia/generic/trie"
)

type Dictionary struct {
	trie  *trie.Trie[struct{}]
	words []string
}

func NewDictionary(file string) *Dictionary {
	d := new(Dictionary)
	d.trie = trie.New[struct{}]()

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		d.trie.Put(line, struct{}{})
		d.words = append(d.words, line)
	}
	return d
}

func (d *Dictionary) IsWord(word string) bool {
	return d.trie.Contains(word)
}
