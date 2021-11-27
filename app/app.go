package app

import (
	"fmt"
	"log"

	"ktn-x.com/tft-leaderboard/data"
	"ktn-x.com/tft-leaderboard/tft"
)

type Jabber struct {
	Riot *tft.RiotClient
	Store *data.Store
}

func (app *Jabber) Import(pool ParticipantPool) ([]*data.Contestant, error) {
	var (
		participants = pool.GetParticipants()
		contestants = make([]*data.Contestant, 0, len(participants))
		errorBag = NewErrorBag(len(participants))
	)

	for _, name := range participants {
		summoner, err := app.Riot.Summoner(name)
		if err != nil {
			errorBag.Add(fmt.Errorf("failed fetching riot summoner %s: %s", name, err))
			continue
		}

		contestants = append(contestants, &data.Contestant{
			Summoner: summoner,
		})
	}

	if err := app.Store.UpdateContestants(contestants); err != nil {
		errorBag.Add(err)
	}

	return contestants, errorBag.Error("failed importing contestants")
}

func (app *Jabber) UpdateRanks() error {
	contestants, err := app.Store.ListContestants()
	if err != nil {
		return err
	}

	items := make([]*tft.TftPair, 0, len(contestants))
	problems := NewErrorBag(len(contestants))

	log.Println("updating ranks")

	for _, contestant := range contestants {
		summoner := contestant.Summoner
		rank, err := app.Riot.TftRanked(summoner.Id)

		if err != nil {
			problems.Add(err)
			continue
		}

		pair := &tft.TftPair{
			Summoner: summoner,
			Rank: rank,
		}

		items = append(items, pair)

		if rank != nil {
			log.Println(data.PrintResult(pair))
		}
	}

	err = app.Store.UpdateContestantRanks(items)
	if err != nil {
		return err
	}

	return problems.Error("failed fetching tft league ranks")
}
