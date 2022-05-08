package wordle

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
	"github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
)

const maxGuess = 6
const wordSize = 5

type word [wordSize]byte
type wordHint [wordSize]uint

type Game struct {
	dict       *Dictionary
	guesses    [maxGuess]word
	hints      [maxGuess]wordHint
	guessIndex int  // the current attempt
	secret     word // the final answer
	secretHash *hashset.Set[byte]
	caps       bool // play the game in capitalized mode
	cheat      bool

	In  io.Reader // the stream game reads from, default to os.Stdin
	Out io.Writer // the stream game prints to, default to os.Stdout
}

func NewGame(secret string, file string) *Game {
	g := new(Game)
	g.secretHash = hashset.New(wordSize, generic.Equals[byte], generic.HashUint8)
	g.dict = NewDictionary(file)
	g.In = os.Stdin
	g.Out = color.Output

	if secret == "" {
		// randomly pick a secret word from dict
		var qualified []string
		for _, k := range g.dict.words {
			if len(k) == wordSize {
				qualified = append(qualified, k)
			}
		}
		r := rand.New(rand.NewSource(time.Now().Unix()))
		secret = qualified[r.Intn(len(qualified))]
	}

	if len(secret) != wordSize {
		panic("Secret must match word size")
	}
	if !g.dict.IsWord(secret) {
		panic("Secret must be a valid word")
	}

	for i := 0; i < wordSize; i++ {
		g.secret[i] = secret[i]
		g.secretHash.Put(secret[i])
	}

	return g
}

func (g *Game) Start() (win bool) {
	if g.cheat {
		color.White("The secret word is: %s", g.secret)
	}
	color.White("Start by typing a five letter word, then press ENTER.")
	r := bufio.NewScanner(g.In)
	for r.Scan() {
		if r.Err() == io.EOF {
			return
		}
		line := r.Text()
		var guess word
		// no matter what the user inputs, scanner the first five bytes only
		_, err := fmt.Sscanf(line, "%c%c%c%c%c", &guess[0], &guess[1], &guess[2], &guess[3], &guess[4])
		if err == io.EOF {
			color.White("you have not input enough characters, try again")
			continue
		}
		if err != nil {
			color.Red(err.Error())
			return false
		}

		// check guess is legal world
		if s := strings.ToLower(string(guess[:])); !g.dict.IsWord(s) {
			color.Red("%s: not a word, try again", s)
			continue
		}

		result := g.readWord(guess)
		g.printHints()
		if g.shouldStop(result) {
			color.Green("You win! The secret word was %s.\n", g.secret)
			return true
		}

		if g.guessIndex == maxGuess {
			// has exhausted all the chances
			color.Red("You lose by using up all the chances! The secret word was %s.\n", g.secret)
			return
		}
	}
	return
}

func (g *Game) SetCheat(cheat bool) {
	g.cheat = cheat
}

// readWord reads in a word from the user, as a new guess
func (g *Game) readWord(w word) wordHint {
	defer func() {
		g.guessIndex++
	}()
	result := g.validate(w)
	g.guesses[g.guessIndex] = w
	g.hints[g.guessIndex] = result
	return result
}

func (g *Game) shouldStop(hint wordHint) bool {
	for i := 0; i < wordSize; i++ {
		if hint[i] < 2 {
			return false
		}
	}
	return true
}

// validate compares the guess word with the secret answer,
// output the hints on each digit
// 0-incorrect
// 1-hit, but wrong position
// 2-bingo
func (g *Game) validate(guess word) (hint wordHint) {
	mismatchSet := hashset.New(0, generic.Equals[byte], generic.HashUint8)
	// 第一遍标绿
	for i := 0; i < wordSize; i++ {
		if guess[i] == g.secret[i] {
			hint[i] = 2
		} else {
			mismatchSet.Put(g.secret[i])
		}
	}

	// 对于非绿的字，判断是否在第一遍不匹配的字符集中
	for i := 0; i < wordSize; i++ {
		if hint[i] != 2 {
			if mismatchSet.Has(guess[i]) {
				hint[i] = 1
				mismatchSet.Remove(guess[i])
			}
		}
	}
	return
}

var (
	white  = color.New(color.FgWhite)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
)

// printHint clears the screen and print all historical guesses (colored by hints)
func (g *Game) printHint(guess word, hint wordHint) {
	for i := 0; i < wordSize; i++ {
		var printer *color.Color
		switch hint[i] {
		case 1:
			printer = yellow
		case 2:
			printer = green
		default:
			printer = white
		}
		// strings.ToUpper on a single byte
		printer.Fprintf(g.Out, "%c", unicode.ToUpper(rune(guess[i])))
	}
	fmt.Println()
}

func clearHistory(out io.Writer) {
	fmt.Fprintf(out, "\033[H\033[2J")
}

func (g *Game) printHints() {
	clearHistory(g.Out)
	for i := 0; i < g.guessIndex; i++ {
		g.printHint(g.guesses[i], g.hints[i])
	}
}
