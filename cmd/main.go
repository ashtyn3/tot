package main

import (
	"os"
	"strings"
	"tot/query"
	"tot/ui"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

func search(w *wow.Wow) string {
	str := query.RandomRepo()
	w.Stop()
	return str
}
func main() {
	w := wow.New(os.Stdout, spin.Get(spin.Dots), " Getting typing material")
	w.Start()
	text := search(w)
	if len(text) > 250 {
		text = text[:170]
	}
	ui.Run(strings.Replace(text, "\"", "", -1))
}
