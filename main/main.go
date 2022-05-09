package main

import (
	"flag"
	"fmt"

	"github.com/yangrq1018/wordle"
)

func main() {
	var (
		dict   = flag.String("dict", "large.txt", "path to dictionary file")
		secret = flag.String("secret", "", "provide a secret explicitly, if left empty, will be generated randomly")
		cheat  = flag.Bool("cheat", false, "print the secret word before playing the game")
		loop   = flag.Bool("loop", false, "keep playing")
	)
	flag.Parse()
	fmt.Printf("load dictionary %s\n", *dict)
	g := wordle.NewGame(*secret, *dict)
	g.SetCheat(*cheat)
	if *loop {
		g.Loop()
	}
	g.Start()
}
