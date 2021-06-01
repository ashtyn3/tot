package ui

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tot/auth"
	"tot/db"
	"tot/query"

	"github.com/ahmetb/go-cursor"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/go-redis/redis"
	"github.com/muesli/termenv"
	term "github.com/nsf/termbox-go"
)

func search(w *wow.Wow) string {
	str := query.RandomRepo()
	w.Stop()
	return str
}

var rdb = auth.Redis()

func StartScreen() {
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()
	reset()
	fmt.Println("Press ESC button to quit")
	home, _ := os.UserHomeDir()
	f, Err := os.ReadFile(home + "/.tot")
	if os.IsNotExist(Err) {
		cursor.MoveTo(0, 0)
		term.Close()
		termenv.AltScreen()
		termenv.ClearScreen()
		fmt.Println("Please type in user name")
		var name string
		fmt.Scanln(&name)
		os.WriteFile(home+"/.tot", []byte(name), 0644)
		rdb.ZAdd("leaderboard", redis.Z{Score: 0, Member: name})
		StartScreen()
	}
	score := rdb.ZScore("leaderboard", string(f))
	currentScore := score.Val()
	cUser := db.User{Name: string(f), Score: currentScore}

	rank := rdb.ZRevRank("leaderboard", cUser.Name)
	iRank := strconv.Itoa(int(rank.Val() + 1))
	fmt.Println("Current ranking: ", termenv.String(iRank).Foreground(termenv.ANSIBrightYellow).Bold())
keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			//termenv.MoveCursor(0, 1)
			//termenv.ClearLine()
			switch ev.Key {
			case term.KeyF1:
			case term.KeyEsc:
				termenv.ExitAltScreen()
				break keyPressListenerLoop
			default:

				if strings.ToLower(string(ev.Ch)) == "g" {
					cursor.MoveTo(0, 0)
					term.Close()
					termenv.AltScreen()
					w := wow.New(os.Stdout, spin.Get(spin.Dots), " Getting typing material")
					w.Start()
					text := search(w)
					text = strings.Replace(text, "\"", "", -1)
					reg := regexp.MustCompile(`[^\x00-\x7F]`)
					text = string(reg.ReplaceAll([]byte(text), []byte(" ")))
					if len(text) > 250 {
						text = text[:170]
					}
					Run(text)
					break keyPressListenerLoop
				}
				if strings.ToLower(string(ev.Ch)) == "l" {
					cursor.MoveTo(0, 0)
					fmt.Println()
					db.Get().Print()
				}
				if strings.ToLower(string(ev.Ch)) == "h" {
					cursor.MoveTo(0, 0)
					fmt.Println()
					fmt.Println(termenv.String("h").Bold().String() + ": Shows help text")
					fmt.Println(termenv.String("g").Bold().String() + ": Opens new typeracing game")
					fmt.Println(termenv.String("l").Bold().String() + ": Shows leaderboard")
				}
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
