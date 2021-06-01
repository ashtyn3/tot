package ui

import (
	"fmt"
	"math"
	"strings"
	"time"
	"tot/utils"
	"unicode"

	"github.com/ahmetb/go-cursor"
	"github.com/muesli/termenv"
	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync() // cosmestic purpose
}

var index = 0
var printable = []string{}
var done = false
var usrInput = []string{""}
var state = 0
var back = false
var timeString = ""

func Timer(secs int) {
	for range time.Tick(1 * time.Second) {

		if done == true {
			reset()
			fmt.Println(cursor.MoveTo(1, 1)+"WPM: ", utils.WPM(index, secs))
			mins := int(secs / 60)
			seconds := secs % 60
			fmt.Printf("Time: %d:%02d\n", mins, seconds)
			fmt.Println("characters typed: ", index)
			break
		}

		if secs == 0 {
			done = true
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
			timeString = termenv.String(fmt.Sprintf("(%d:%02d)", mins, seconds)).Foreground(termenv.ANSIBrightYellow).String()
		} else if seconds <= 5 && mins == 0 {
			timeString = termenv.String(fmt.Sprintf("(%d:%02d)", mins, seconds)).Foreground(termenv.ANSIBrightRed).String()
		} else {
			timeString = fmt.Sprintf("(%d:%02d)", mins, seconds)
		}
	}
}

var lastCycle = []string{}

func update(text string) {
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
	fmt.Print(cursor.MoveTo(3, 0))
	termenv.ClearLine()
	fmt.Print("\r" + strings.Join(printable, ""))
}

var keypressed = false

func Run(text string) {
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()
	length := len(text)
	for range text {
		usrInput = append(usrInput, "")
	}
	reset()
	fmt.Println("Press ESC button to quit")
	go Timer(120)
	go func() {
		if done == false {
			for range time.Tick(10 * time.Millisecond) {
				fmt.Print(cursor.MoveTo(4, 0) + cursor.ClearEntireLine() + "\r" + timeString)
				if keypressed {
					update(text)
					keypressed = false
					back = false
				}
			}
		} else {
			fmt.Print(cursor.MoveTo(4, 0) + cursor.ClearEntireLine())
		}
	}()
	if strings.HasSuffix(text, ".") == false || strings.HasSuffix(text, "?") == false || strings.HasSuffix(text, "!") {
		fmt.Println(text + "...")
	} else {
		fmt.Println(text)
	}
keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			keypressed = true
			//termenv.MoveCursor(0, 1)
			//termenv.ClearLine()
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
				if index != length && ev.Ch <= unicode.MaxASCII {
					index++
					usrInput[index] = string(ev.Ch)
				}
			}

			if index == length && state == 0 {
				done = true
				termenv.ClearScreen()
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
