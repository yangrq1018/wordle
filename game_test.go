package wordle

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dictFilePath = "dictionary/large.txt"

func TestDictionary(t *testing.T) {
	d := NewDictionary(dictFilePath)
	assert.True(t, d.IsWord("cat"))
	assert.True(t, d.IsWord("tidal"))
	assert.False(t, d.IsWord("bilibili"))
}

func TestWinGame(t *testing.T) {
	g := NewGame("midst", dictFilePath)
	mockIn := bytes.NewBuffer(nil)
	mockOut := bytes.NewBuffer(nil)
	g.In = mockIn
	g.Out = mockOut

	_, err := fmt.Fprintf(mockIn,
		`daily
sword
unite
tidal
midst
`)
	assert.NoError(t, err)
	assert.True(t, g.Start())
}

func TestLoseGame(t *testing.T) {
	g := NewGame("midst", dictFilePath)
	mockIn := bytes.NewBuffer(nil)
	mockOut := bytes.NewBuffer(nil)
	g.In = mockIn
	g.Out = mockOut

	_, err := fmt.Fprintf(mockIn,
		`daily
sword
unite
tidal
plane
chair
`)
	assert.NoError(t, err)
	assert.False(t, g.Start())
}
