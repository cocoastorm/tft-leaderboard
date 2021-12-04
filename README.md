# tft leaderboard

A dumb simple leaderboard showing tft rank.

> "Race to Iron"

## Getting Started

tft-leaderboard requires an api key to Riot Games' TFT API.
A development key should be adequate to get started, you can get one from: https://developer.riotgames.com/

### Golang Web Server

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

#### Configurable Environment Variables and/or Flags
Global Config:
* `TFT_DB` - path for bbolt kv database
* `TFT_API_KEY` - Riot Games API Key that has permission to use tft api routes

Web Server Config:
* `TFT_ADDRESS` - server address to listen on (pls. include the port number)
* `TFT_WRITE_TIMEOUT` - server write timeout
* `TFT_READ_TIMEOUT` - server read timeout
* `TFT_APP_PATH` - path to directory containing built frontend files to serve
* `TFT_APP_INDEXPATH` - path and/or filename of index file to serve
* `TFT_POLL` - duration of how long before the server polls riot games' api for tft ranks
* `TFT_GOAL` - human friendly string of what rank is the race to (eg. "DIAMONDS" - race to diamonds)

### Next.js Web App

Giving Yarn v2/3 a try. To get started make sure you have Node.js >= v16 installed.

1. enable corepack & install version 'berry' of yarnpkg
```bash
# if you have Node.js installed via Homebrew on macos
brew install corepack

# enable corepack
corepack enable

# install yarn v3
yarn version set berry
```

2. install node dependencies
```bash
# install dependencies
yarn install

# run server (must be built beforehand)
go build

# run development web app
yarn run dev
```

## To Do
- [] compress bbolt db file on exit?
