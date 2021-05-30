package db

import (
	"fmt"
	"strconv"
	"tot/auth"

	"github.com/go-redis/redis"
	"github.com/muesli/termenv"
)

type User struct {
	Name  string
	Score float64
}

type Leaderboard []User

var rdb = auth.Redis()

func (u *User) Add() {
	rdb.ZAdd("leaderboard", redis.Z{Score: u.Score, Member: u.Name})
}

func Get() Leaderboard {
	sl := rdb.ZRevRangeWithScores("leaderboard", 0, 15)
	leaderboard := Leaderboard{}
	for _, member := range sl.Val() {
		u := User{Name: member.Member.(string), Score: member.Score}
		leaderboard = append(leaderboard, u)
	}
	return leaderboard
}

func (u *User) Update() {
	rdb.ZRem("leaderboard", u.Name)
	u.Add()
}

func (l Leaderboard) Print() {
	termenv.AltScreen()
	fmt.Println("Leaderboard:")
	for rank, user := range l {
		rank += 1
		if rank == 1 {
			fmt.Print(termenv.String(strconv.Itoa(rank)).Foreground(termenv.ANSIBrightCyan).Bold())
		} else {
			fmt.Print(termenv.String(strconv.Itoa(rank)).Bold())
		}
		fmt.Print(" ", user.Name+"   ", user.Score, "\n")
	}
}
