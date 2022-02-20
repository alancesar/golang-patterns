package main

import (
	"fmt"
	"golang-patterns/internal/sleep"
	"golang-patterns/newfetcher"
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
		return newfetcher.New(match, stadiumGetter, func(match *Match, stadium Stadium) error {
			match.Stadium = stadium
			return nil
		})
	}

	championshipFetcher := func() error {
		return newfetcher.New(match, championshipGetter, func(match *Match, championship string) error {
			match.Championship = championship
			return nil
		})
	}

	homeTeamFetcher := func(id int) error {
		return newfetcher.NewWithParam(match, id, teamGetter, func(match *Match, home Team) error {
			match.Home = home
			return nil
		})
	}

	awayTeamFetcher := func(id int) error {
		return newfetcher.NewWithParam(match, id, teamGetter, func(match *Match, away Team) error {
			match.Away = away
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
