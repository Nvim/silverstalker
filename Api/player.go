package api

/* Helpers to fetch data about players */

import (
	"encoding/json"
	"errors"
	"slices"
)

type PlayerInfo struct {
	GameName   string
	TagLine    string
	PUUID      string
	SummonerID string
	AccountID  string
}

type AccountJSON struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type SummonerJSON struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

// Stats of a player in a league
type LeagueStats struct {
	LeagueID     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerID   string `json:"summonerId"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
}

type MatchID []string

func (p *PlayerInfo) GetIDs() error {
	if p.GameName == "" || p.TagLine == "" {
		err := errors.New("couldn't retrieve player's IDs without GameName and TagLine")
		return err
	}
	puidUrl := "https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/" + p.GameName + "/" + p.TagLine

	var puidResponse AccountJSON
	res, err := GetRiotApi(puidUrl)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, &puidResponse); err != nil {
		return err
	}

	var summonerResponse SummonerJSON
	summonerUrl := "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/" + puidResponse.Puuid
	res, err = GetRiotApi(summonerUrl)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, &summonerResponse); err != nil {
		return err
	}
	p.PUUID = puidResponse.Puuid
	p.AccountID = summonerResponse.AccountID
	p.SummonerID = summonerResponse.ID

	return nil // No errors :)
}

func (p *PlayerInfo) getRankedStats() (rankedStats *LeagueStats, err error) {
	if p.SummonerID == "" {
		err = errors.New("couldn't get info about player: empty SummonerID")
		return nil, err
	}

	statsUrl := "https://euw1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + p.SummonerID

	var statsArray []LeagueStats
	res, err := GetRiotApi(statsUrl)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(res, &statsArray); err != nil {
		return nil, err
	}

	rankedStatsIndex := slices.IndexFunc(statsArray, func(s LeagueStats) bool {
		return s.QueueType == "RANKED_SOLO_5x5"
	})
	if rankedStatsIndex == -1 {
		err = errors.New("player doesn't have any ranked games")
		return nil, err
	}
	rankedStats = &statsArray[rankedStatsIndex]
	return
}

func (p *PlayerInfo) GetLatestMatches() (MatchID, error) {
	if p.PUUID == "" {
		err := errors.New("couldn't get player's matches: empty PUUID")
		return nil, err
	}

	matchesURL := "https://europe.api.riotgames.com/lol/match/v5/matches/by-puuid/" + p.PUUID + "/ids?queue=420&start=0&count=20"

	res, err := GetRiotApi(matchesURL)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Raw response body: ", string(res))
	var matchIDs MatchID
	if err := json.Unmarshal(res, &matchIDs); err != nil {
		return nil, err
	}

	return matchIDs, nil
}
