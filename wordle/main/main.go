package main

import (
	"flag"

	"github.com/yangrq1018/wordle"
)

func main() {
	var (
		dict   = flag.String("dict", "large.txt", "")
		secret = flag.String("secret", "", "provide a secret explicitly, if left empty, will be generated randomly")
		cheat  = flag.Bool("cheat", false, "")
	)
	flag.Parse()
	g := wordle.NewGame(*secret, *dict)
	g.SetCheat(*cheat)
	g.Start()
}
