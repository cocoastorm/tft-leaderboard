package web

import (
	"context"
	"log"
	"time"

	"ktn-x.com/tft-leaderboard/app"
	"ktn-x.com/tft-leaderboard/data"
	"ktn-x.com/tft-leaderboard/tft"
)

type RankSync struct {
	jabber *app.Jabber
	pollDuration time.Duration
	success int
	fails int
}

func NewPoll(riot *tft.RiotClient, store *data.Store, t time.Duration) *RankSync {
	return &RankSync{
		jabber: &app.Jabber{
			Riot: riot,
			Store: store,
		},
		pollDuration: t,
	}
}

func (r *RankSync) Do(ctx context.Context) {
	// from: https://github.com/AlbinoDrought/np-scanner/blob/8444bda03bb7433e4f4b38e608bd28e2354a74c7/internal/web/web.go#L104
	// written this way so the timer fires immediately the first time
	// later on, it gets reset with the proper desired duration
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <- timer.C:
			err := r.jabber.UpdateRanks()
			if err != nil {
				r.fails += 1
				log.Printf("failed updating ranks [x%d]: %s", r.fails, err)
				timer.Reset(1 * time.Minute)
			} else {
				r.success += 1
				timer.Reset(r.pollDuration)
			}
		}
	}
}
