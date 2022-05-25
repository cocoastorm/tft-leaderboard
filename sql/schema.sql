CREATE TABLE `users` IF NOT EXISTS (
  id UUID PRIMARY KEY,
  accountId VARCHAR (255) NOT NULL ,
  profileIconId VARCHAR (255) NOT NULL,
  revisionDate timestamp NOT NULL,
  name VARCHAR (255) NOT NULL,
  puuid VARCHAR (255) NOT NULL,
  summonerLevel integer NOT NULL,
);

CREATE TABLE `tft_ranks` IF NOT EXISTS (
  id UUID PRIMARY KEY,

  leagueId VARCHAR (255) NOT NULL,
  summonerId VARCHAR (255) NOT NULL,
  summonerName VARCHAR (255) NOT NULL,

  queueType VARCHAR (255) NOT NULL,

  ratedTier VARCHAR (255) NOT NULL,
  ratedRating VARCHAR (255) NOT NULL,
  tier VARCHAR (255) NOT NULL,
  rank VARCHAR (255) NOT NULL,
  leaguePoints integer NOT NULL,
  wins integer NOT NULL,
  losses integer NOT NULL,
  hotStreak boolean NOT NULL,
);
