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

const (
	maxGuess     = 6
	wordSize     = 5
	alphabetSize = 26
)

type word [wordSize]byte
type wordHint [wordSize]uint

func (w word) String() string {
	return string(w[:])
}

type Game struct {
	dict          *Dictionary
	guesses       [maxGuess]word
	hints         [maxGuess]wordHint
	guessIndex    int  // the current attempt
	secret        word // the final answer
	cheat         bool
	loop          bool
	gameAttempted int

	In  io.Reader // the stream game reads from, default to os.Stdin
	Out io.Writer // the stream game prints to, default to os.Stdout
}

func NewGame(secret string, file string) *Game {
	g := new(Game)
	g.dict = NewDictionary(file)
	g.In = os.Stdin
	g.Out = color.Output

	if secret == "" {
		g.genSecret()
	} else {
		g.setSecret(secret)
	}
	return g
}

func (g *Game) genSecret() {
	// randomly pick a word from dictionary
	var qualified []string
	for _, k := range g.dict.words {
		if len(k) == wordSize {
			qualified = append(qualified, k)
		}
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	g.setSecret(qualified[r.Intn(len(qualified))])
}

func (g *Game) setSecret(secret string) {
	if len(secret) != wordSize {
		panic("Secret must match word size")
	}
	if !g.dict.IsWord(secret) {
		panic("Secret must be a valid word")
	}
	for i := 0; i < wordSize; i++ {
		g.secret[i] = secret[i]
	}
}

func (g *Game) Start() (win bool) {
	r := bufio.NewScanner(g.In)

restart:
	if g.loop && g.gameAttempted > 0 {
		// prompt for user ENTER
		fmt.Fprintf(g.Out, "Press [ENTER] to restart or [Ctrl+C] to exit\n")
		bufio.NewReader(g.In).ReadLine()
		cls(g.Out)
	}
	if g.cheat {
		color.White("The secret word is %s", g.secret)
	}
	color.White("Start by typing a five letter word, then press [ENTER]")

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
		if s := strings.ToLower(guess.String()); !g.dict.IsWord(s) {
			color.Red("%s: not a word, try again", s)
			continue
		}

		result := g.readWord(guess)
		g.screen()
		if g.shouldStop(result) {
			g.win()
			if g.loop {
				g.reset()
				g.genSecret()
				goto restart
			}
			return true
		}

		if g.guessIndex == maxGuess {
			// has exhausted all the chances
			g.lose()
			if g.loop {
				g.reset()
				g.genSecret()
				goto restart
			}
			return
		}
	}
	return
}

// note reset does not reset the secret word
func (g *Game) reset() {
	g.guessIndex = 0
}

func (g *Game) win() {
	g.gameAttempted++
	color.Green("You win! The secret word is %s.\n", g.secret)
	color.White(getMeaningOfWord(g.secret))
}

func (g *Game) lose() {
	g.gameAttempted++
	color.Red("You lose by using up all the chances! The secret word is %s.\n", g.secret)
	color.White(getMeaningOfWord(g.secret))
}

func (g *Game) SetCheat(cheat bool) {
	g.cheat = cheat
}

func (g *Game) Loop() {
	g.loop = true
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

// validate compares the guess word with the secret answer, output the hints on each position
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

// printHint clears the screen and print all historical guesses (colored by hints)
func (g *Game) printHint(guess word, hint wordHint) {
	for i := 0; i < wordSize; i++ {
		// strings.ToUpper on a single byte
		printer(hint[i]).Fprintf(g.Out, "%c", unicode.ToUpper(rune(guess[i])))
	}
}

func (g *Game) knownLetters() ([26]byte, [26]*color.Color) {
	var (
		locRight  = hashset.New(0, generic.Equals[byte], generic.HashUint8)
		inWord    = hashset.New(0, generic.Equals[byte], generic.HashUint8)
		notInWord = hashset.New(0, generic.Equals[byte], generic.HashUint8)
	)
	for i := 0; i < g.guessIndex; i++ {
		for j := range g.hints[i] {
			switch g.hints[i][j] {
			case 2:
				locRight.Put(g.guesses[i][j])
			case 1:
				inWord.Put(g.guesses[i][j])
			case 0:
				notInWord.Put(g.guesses[i][j])
			}
		}
	}

	var (
		alphabet [alphabetSize]byte
		hints    [alphabetSize]*color.Color
	)

	for i := byte('a'); i <= byte('z'); i++ {
		alphabet[i-'a'] = i + 'A' - 'a'
		if locRight.Has(i) {
			hints[i-'a'] = green
		} else if inWord.Has(i) {
			hints[i-'a'] = yellow
		} else if notInWord.Has(i) {
			hints[i-'a'] = nil // invisible
		} else {
			hints[i-'a'] = white
		}
	}
	return alphabet, hints
}

func (g *Game) screen() {
	cls(g.Out)
	for i := 0; i < maxGuess; i++ {
		if i < g.guessIndex {
			g.printHint(g.guesses[i], g.hints[i])
		} else {
			fmt.Fprint(g.Out, strings.Repeat("_", wordSize))
		}
		fmt.Println()
	}

	letters, hints := g.knownLetters()
	fmt.Println("Knowledge:")
	for i := range letters {
		switch hints[i] {
		case nil:
			fmt.Fprint(g.Out, " ")
		default:
			hints[i].Fprintf(g.Out, "%c", letters[i])
		}

	}
	fmt.Println()
}
