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

	fetcher.New(func(locker sync.Locker, event fetcher.Event) {
		switch event.Name {
		case "STADIUM":
			match.Stadium = event.Data.(Stadium)
		case "CHAMPIONSHIP":
			match.Championship = event.Data.(string)
		case "HOME_TEAM":
			match.Home = event.Data.(Team)
		case "AWAY_TEAM":
			match.Away = event.Data.(Team)
		}
	}).AddProducer("STADIUM", func() interface{} {
		return stadiumGetter()
	}).AddProducer("CHAMPIONSHIP", func() interface{} {
		return championshipGetter()
	}).AddProducer("HOME_TEAM", func() interface{} {
		return teamGetter(1)
	}).AddProducer("AWAY_TEAM", func() interface{} {
		return teamGetter(2)
	}).Fetch(&match)

	fmt.Println(&match)
}
