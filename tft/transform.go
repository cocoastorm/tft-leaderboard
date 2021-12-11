package tft

func (summoner *Summoner) transform() *WSummoner {
	return &WSummoner{
		Id: summoner.Id,
		Name: summoner.Name,
		Level: summoner.SummonerLevel,
		ProfileIconId: summoner.ProfileIconId,
	}
}

func (league *TftLeague) transform() *WTftLeague {
	return &WTftLeague{
		QueueType: league.QueueType,
		LeaguePoints: league.LeaguePoints,
		Tier: league.Tier,
		Rank: league.Rank,
		Wins: league.Wins,
		Losses: league.Losses,
		Total: league.Wins + league.Losses,
	}
}

func (pair *TftPair) Transform() *WTftPair {
	w := &WTftPair{}

	if pair.Summoner != nil {
		w.Summoner = pair.Summoner.transform()
	}

	if pair.Rank != nil {
		w.Rank = pair.Rank.transform()
	}

	return w
}
