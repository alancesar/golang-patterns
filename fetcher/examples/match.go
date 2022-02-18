package main

import (
	"fmt"
	"golang-patterns/fetcher"
	"math/rand"
	"sync"
	"time"
)

type (
	Match struct {
		sync.Mutex
		Home         Team
		Away         Team
		Championship string
		Stadium      Stadium
	}

	Team struct {
		Name string
		City string
	}

	Stadium struct {
		Name     string
		City     string
		Capacity int
	}
)

func randomSleep() {
	ms := rand.Intn(10) * 100
	time.Sleep(time.Millisecond * time.Duration(ms))
}

func championshipGetter() string {
	randomSleep()
	return "Champions League"
}

func teamGetter(id int) Team {
	randomSleep()
	if id == 1 {
		return Team{
			Name: "AFC Ajax",
			City: "Amsterdam",
		}
	}

	return Team{
		Name: "Real Madrid CF",
		City: "Madrid",
	}
}

func stadiumGetter() Stadium {
	randomSleep()
	return Stadium{
		Name:     "Wembley Stadium",
		City:     "London",
		Capacity: 90000,
	}
}

func main() {
	var match Match

	fetcher.New().
		With(func() interface{} {
			return stadiumGetter()
		}, func(value interface{}) {
			match.Stadium = value.(Stadium)
		}).
		With(func() interface{} {
			return championshipGetter()
		}, func(value interface{}) {
			match.Championship = value.(string)
		}).
		With(func() interface{} {
			return teamGetter(1)
		}, func(value interface{}) {
			match.Home = value.(Team)
		}).
		With(func() interface{} {
			return teamGetter(2)
		}, func(value interface{}) {
			match.Away = value.(Team)
		}).Fetch(&match)

	fmt.Println(&match)
}
