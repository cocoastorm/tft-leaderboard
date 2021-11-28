# tft leaderboard

A dumb simple leaderboard showing tft rank.

> "Race to Iron"

## Getting Started

tft-leaderboard requires an api key to Riot Games' TFT API.
A development key should be adequate to get started, you can get one from: https://developer.riotgames.com/

As a Golang application, we have a few go modules that need to be installed.
With Golang installed and the respective `GOPATH` set up:

```bash
# install go dependencies
go get
# build application
go build
```

After building, to get started tft-leaderboard needs to know what tft players to add to the board.
Add them with `tft-leaderboard import --i <file>`, where the input file is a text file with a summoner name on each line:

```bash
echo -n 'Scarra' > names.txt

./tft-leaderboard import --i names.txt
```

Finally, you can run the web server with:

`./tft-leaderboard serve --api-key="<RIOT API KEY>"`

This should launch a web server listening on the default port `8080`. A sorted list of players and their rank is available at `/api/leaderboard`.
A frontend web page will be available soon :tm:

## To Do
- [] front web page
- [] compress bbolt db file on exit?
