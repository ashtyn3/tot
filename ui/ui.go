package ui

import (
	"fmt"
	"math"
	"strings"
	"time"
	"tot/utils"

	"github.com/muesli/termenv"
	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync() // cosmestic purpose
}

var index = 0
var printable = []string{}
var done = false

func Timer(secs int) {
	for range time.Tick(1 * time.Second) {
		if done == true {
			termenv.RestoreCursorPosition()
			reset()
			fmt.Println("WPM: ", utils.WPM(index, secs))
			mins := int(secs / 60)
			seconds := secs % 60
			fmt.Printf("Time: %d:%02d\n", mins, seconds)
			fmt.Println("characters typed: ", index)
			break
		}

		if secs == 0 {
			done = true
			termenv.RestoreCursorPosition()
			reset()
			fmt.Println("WPM: ", utils.WPM(index, secs))
			mins := int(secs / 60)
			seconds := secs % 60
			fmt.Printf("Time: %d:%02d\n", mins, seconds)
			fmt.Println("characters typed: ", index)
			break
		}
		secs -= 1
		mins := int(math.Trunc(float64(secs / 60)))
		seconds := secs % 60
		termenv.CursorForward(1)
		if seconds <= 10 && seconds > 5 && mins == 0 {
			fmt.Print(termenv.String(fmt.Sprintf("(%d:%02d)", mins, seconds)).Foreground(termenv.ANSIBrightYellow))
		} else if seconds <= 5 && mins == 0 {
			fmt.Print(termenv.String(fmt.Sprintf("(%d:%02d)", mins, seconds)).Foreground(termenv.ANSIBrightRed))
		} else {
			fmt.Printf("(%d:%02d)", mins, seconds)
		}
	}
}

func Run(text string) {
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()
	length := len(text)
	usrInput := []string{""}
	for range text {
		usrInput = append(usrInput, "")
	}
	reset()
	fmt.Println("Press ESC button to quit")
	if strings.HasSuffix(text, ".") == false || strings.HasSuffix(text, "?") == false || strings.HasSuffix(text, "!") {
		fmt.Println(text + "...")
	} else {
		fmt.Println(text)
	}
	termenv.SaveCursorPosition()
	go Timer(120)
	state := 0

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			//termenv.MoveCursor(0, 1)
			//termenv.ClearLine()
			back := false
			switch ev.Key {
			case term.KeyF1:
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeySpace:
				if index != length {
					index++
					usrInput[index] = " "
				}
			case term.KeyBackspace | term.KeyBackspace2:
				if index-1 != length {
					if index != -1 {
						back = true
						index--
						usrInput[index+1] = ""
						printable = printable[:len(printable)-1]
					}
				}
			case term.KeyEnter:
				if index == length && state == 0 {
					reset()
					break keyPressListenerLoop
				}
			default:
				// we only want to read a single character or one key pressed event
				if index != length {
					index++
					usrInput[index] = string(ev.Ch)
				}
			}
			if strings.Join(usrInput, "") == text[0:index] {
				state = 0
				termenv.Reset()
				if back == false {
					printable = append(printable, termenv.String(usrInput[index]).Foreground(termenv.ANSIBrightGreen).Bold().String())
				}
				termenv.Reset()
			} else {
				state = 1
				termenv.Reset()
				if back == false {
					printable = append(printable, termenv.String(usrInput[index]).Foreground(termenv.ANSIBrightRed).Bold().String())
				}
				termenv.Reset()
			}
			termenv.ClearLine()
			termenv.RestoreCursorPosition()
			fmt.Print(strings.Join(printable, ""))

			if index == length {
				done = true
				termenv.ClearLine()
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
