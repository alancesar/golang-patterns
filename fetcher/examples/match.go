package main

import (
	"fmt"
	"golang-patterns/fetcher"
	"golang-patterns/internal/sleep"
	"log"
	"sync"
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

func championshipGetter() (string, error) {
	sleep.Random()
	return "Champions League", nil
}

func teamGetter(id int) (Team, error) {
	sleep.Random()
	if id == 1 {
		return Team{
			Name: "AFC Ajax",
			City: "Amsterdam",
		}, nil
	}

	return Team{
		Name: "Real Madrid CF",
		City: "Madrid",
	}, nil
}

func stadiumGetter() (Stadium, error) {
	sleep.Random()
	return Stadium{
		Name:     "Wembley Stadium",
		City:     "London",
		Capacity: 90000,
	}, nil
}

func main() {
	var match Match

	err := fetcher.New().
		With(func() (interface{}, error) {
			return stadiumGetter()
		}, func(value interface{}) {
			match.Stadium = value.(Stadium)
		}).
		With(func() (interface{}, error) {
			return championshipGetter()
		}, func(value interface{}) {
			match.Championship = value.(string)
		}).
		With(func() (interface{}, error) {
			return teamGetter(1)
		}, func(value interface{}) {
			match.Home = value.(Team)
		}).
		With(func() (interface{}, error) {
			return teamGetter(2)
		}, func(value interface{}) {
			match.Away = value.(Team)
		}).
		Fetch(&match)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(&match)
}
