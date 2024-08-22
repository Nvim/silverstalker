package api

/* Helpers to fetch data about matches */

import (
	"encoding/json"
)

type MatchMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchID"`
	Participants []string `json:"participants"`
}

type Challenge struct {
	DamagePerMinute           float64 `json:"damagePerMinute"`
	GameLength                float64 `json:"gameLength"`
	GoldPerMinute             float64 `json:"goldPerMinute"`
	Kda                       float64 `json:"kda"`
	LaneMinionsFirst10Minutes int     `json:"laneMinionsFirst10Minutes"`
	SoloKills                 int     `json:"soloKills"`
	TeamDamagePercentage      float64 `json:"teamDamagePercentage"`
}

type Participant struct {
	Lane                        string    `json:"lane"`
	RiotIDTagline               string    `json:"riotIdTagline"`
	SummonerName                string    `json:"summonerName"`
	Puuid                       string    `json:"puuid"`
	RiotIDGameName              string    `json:"riotIdGameName"`
	ChampionName                string    `json:"championName"`
	IndividualPosition          string    `json:"individualPosition"`
	Challenges                  Challenge `json:"challenges"`
	Summoner1Casts              int       `json:"summoner1Casts"`
	ObjectivesStolen            int       `json:"objectivesStolen"`
	TeamId                      int       `json:"teamId"`
	BaronKills                  int       `json:"baronKills"`
	DragonKills                 int       `json:"dragonKills"`
	ChampLevel                  int       `json:"champLevel"`
	DamageDealtToBuildings      int       `json:"damageDealtToBuildings"`
	DamageDealtToObjectives     int       `json:"damageDealtToObjectives"`
	DamageDealtToTurrets        int       `json:"damageDealtToTurrets"`
	DoubleKills                 int       `json:"doubleKills"`
	KillingSprees               int       `json:"killingSprees"`
	ItemsPurchased              int       `json:"itemsPurchased"`
	GoldEarned                  int       `json:"goldEarned"`
	GoldSpent                   int       `json:"goldSpent"`
	LongestTimeSpentLiving      int       `json:"longestTimeSpentLiving"`
	ParticipantID               int       `json:"participantId"`
	SightWardsBoughtInGame      int       `json:"sightWardsBoughtInGame"`
	VisionWardsBoughtInGame     int       `json:"visionWardsBoughtInGame"`
	WardsPlaced                 int       `json:"wardsPlaced"`
	ChampionID                  int       `json:"championId"`
	Summoner1ID                 int       `json:"summoner1Id"`
	Summoner2Casts              int       `json:"summoner2Casts"`
	Summoner2ID                 int       `json:"summoner2Id"`
	TotalDamageDealtToChampions int       `json:"totalDamageDealtToChampions"`
	TotalMinionsKilled          int       `json:"totalMinionsKilled"`
	TotalTimeSpentDead          int       `json:"totalTimeSpentDead"`
	TurretKills                 int       `json:"turretKills"`
	VisionScore                 int       `json:"visionScore"`
	Win                         bool      `json:"win"`
}

type MatchInfo struct {
	EndOfGameResult  string        `json:"endOfGameResult"`
	GameType         string        `json:"gameType"`
	GameName         string        `json:"gameName"`
	Participants     []Participant `json:"participants"`
	GameID           int64         `json:"gameId"`
	QueueID          int           `json:"queueId"`
	GameCreation     int64         `json:"gameCreation"`
	GameDuration     int           `json:"gameDuration"`
	GameEndTimestamp int64         `json:"gameEndTimestamp"`
}

type Match struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     MatchInfo     `json:"info"`
}

func GetMatchInfo(id string) (*Match, error) {
	url := "https://europe.api.riotgames.com/lol/match/v5/matches/" + id

	res, err := GetRiotApi(url)
	if err != nil {
		return nil, err
	}

	var match Match
	if err := json.Unmarshal(res, &match); err != nil {
		return nil, err
	}

	// fmt.Println(PrettyPrint(match))

	return &match, nil
}
