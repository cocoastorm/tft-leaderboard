package tft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type RiotApi struct {
	Region string
	ApiKey string
	ApiUrl string
}

type RiotClient struct {
	cli *http.Client
	api *RiotApi
}

func NewRiotApiDefault(key string) RiotApi {
	return RiotApi{
		Region: "NA1",
		ApiKey: key,
		ApiUrl: "https://na1.api.riotgames.com",
	}
}

func NewRiotClient(api *RiotApi) *RiotClient {
	if err := api.validate(); err != nil {
		panic(err)
	}

	cli := &http.Client{
		Timeout: time.Second * 15,
	}

	return &RiotClient{
		cli,
		api,
	}
}

func NewRiot(key string) *RiotClient {
	api := NewRiotApiDefault(key)
	return NewRiotClient(&api)
}

func (c *RiotApi) validate() error {
	errMsg := "riot api config errors"
	errors := make([]string, 0, 3)

	if c.Region == "" {
		errors = append(errors, fmt.Sprintln("\tRegion should be set"))
	}

	if c.ApiUrl == "" {
		errors = append(errors, fmt.Sprintln("\tBase API URl should be set"))
	}

	if c.ApiKey == "" {
		errors = append(errors, fmt.Sprintln("\tAPI Key should be set"))
	}

	if len(errors) != 0 {
		return fmt.Errorf("%s:\n%s", errMsg, strings.Join(errors, "\n"))
	}

	return nil
}

func (r *RiotClient) build(method string, p string) (*http.Request, error) {
	to, err := url.Parse(r.api.ApiUrl)
	if err != nil {
		return nil, err
	}

	to.Path = path.Join(to.Path, p)

	req, err := http.NewRequest(method, to.String(), nil)
	if err != nil {
		return nil, err
	}

	// consume json
	req.Header.Add("Accept", "application/json")

	// api key auth
	if r.api.ApiKey == "" {
		return nil, fmt.Errorf("riot token is required")
	}
	req.Header.Add("X-Riot-Token", r.api.ApiKey)

	return req, nil
}

func (r *RiotClient) Summoner(name string) (*Summoner, error) {
	summonerUrl := fmt.Sprintf("/tft/summoner/v1/summoners/by-name/%s", name)
	req, err := r.build("GET", summonerUrl)
	if err != nil {
		return nil, err
	}

	resp, err := r.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	summoner := Summoner{}

	err = json.NewDecoder(resp.Body).Decode(&summoner)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s summoner: %s", name, err)
	}

	return &summoner, nil
}

func (r *RiotClient) TftRanks(summonerId string) ([]*TftLeague, error) {
	leagueUrl := fmt.Sprintf("/tft/league/v1/entries/by-summoner/%s", summonerId)
	req, err := r.build("GET", leagueUrl)
	if err != nil {
		return nil, err
	}

	resp, err := r.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var (
		errMsg = fmt.Sprintf("failed to fetch tft league rank for summoner [%s]:", summonerId)
		ranks []*TftLeague
	)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&ranks)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", errMsg, err)
	}

	return ranks, nil
}

func (r *RiotClient) TftRanked(summonerId string) (*TftLeague, error) {
	queueType := "RANKED_TFT"

	ranks, err := r.TftRanks(summonerId)
	if err != nil {
		return nil, err
	}

	for _, rank := range ranks {
		if rank.QueueType == queueType {
			return rank, nil
		}
	}

	// explicitly don't spit out error for nil tft rank
	fmt.Printf("summoner [%s] has no tft league rank of queue type: %s\n",  summonerId, queueType)
	return nil, nil
}
