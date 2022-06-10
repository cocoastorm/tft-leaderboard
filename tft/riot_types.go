package tft

// taken from https://developer.riotgames.com/apis
// 2021-11-07

type Summoner struct {
	AccountId     string `json:"accountId"`
	ProfileIconId int64  `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	Name          string `json:"name"`
	Id            string `json:"id"`
	Puuid         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
}

type MiniSeries struct {
	Losses   int64  `json:"losses"`
	Progress string `json:"progress"`
	Target   int64  `json:"target"`
	Wins     int64  `json:"wins"`
}

type TftLeague struct {
	LeagueId     string      `json:"leagueId"`
	SummonerId   string      `json:"summonerId"`
	SummonerName string      `json:"summonerName"`
	QueueType    string      `json:"queueType"`
	RatedTier    string      `json:"ratedTier"`
	RatedRating  int64       `json:"ratedRating"`
	Tier         string      `json:"tier"`
	Rank         string      `json:"rank"`
	LeaguePoints int64       `json:"leaguePoints"`
	Wins         int64       `json:"wins"`
	Losses       int64       `json:"losses"`
	HotStreak    bool        `json:"hotStreak"`
	Veteran      bool        `json:"veteran"`
	FreshBlood   bool        `json:"freshBlood"`
	Inactive     bool        `json:"inactive"`
	MiniSeries   *MiniSeries `json:"miniSeries"`
}

type TftPair struct {
	Summoner *Summoner  `json:"summoner"`
	Rank     *TftLeague `json:"rank"`
}

type TftRanks struct {
	Summoner *Summoner    `json:"summoner"`
	Ranks    []*TftLeague `json:"ranks"`
}
