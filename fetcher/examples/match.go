package main

import (
	"fmt"
	"golang-patterns/fetcher"
	"golang-patterns/internal/sleep"
	"golang.org/x/sync/errgroup"
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
	match := &Match{}
	stadiumFetcher := func() error {
		return fetcher.New(match, func() (interface{}, error) {
			return stadiumGetter()
		}, func(match sync.Locker, stadium interface{}) error {
			match.(*Match).Stadium = stadium.(Stadium)
			return nil
		})
	}

	championshipFetcher := func() error {
		return fetcher.New(match, func() (interface{}, error) {
			return championshipGetter()
		}, func(match sync.Locker, championship interface{}) error {
			match.(*Match).Championship = championship.(string)
			return nil
		})
	}

	homeTeamFetcher := func(id interface{}) error {
		return fetcher.NewWithParam(match, id, func(id interface{}) (interface{}, error) {
			return teamGetter(id.(int))
		}, func(match sync.Locker, home interface{}) error {
			match.(*Match).Home = home.(Team)
			return nil
		})
	}

	awayTeamFetcher := func(id interface{}) error {
		return fetcher.NewWithParam(match, id, func(id interface{}) (interface{}, error) {
			return teamGetter(id.(int))
		}, func(match sync.Locker, away interface{}) error {
			match.(*Match).Away = away.(Team)
			return nil
		})
	}

	var group errgroup.Group
	group.Go(stadiumFetcher)
	group.Go(championshipFetcher)
	group.Go(func() error {
		return homeTeamFetcher(1)
	})
	group.Go(func() error {
		return awayTeamFetcher(2)
	})
	if err := group.Wait(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(match)
}
