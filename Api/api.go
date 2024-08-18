package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
)

var ApiToken string = ""

var JsonErr = errors.New("Can't unmarshal JSON!")

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

// Performs a GET request on the given URL
func GetRiotApi(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Println("Couldn't create request: ", err)
		return nil, err
	}
	req.Header.Add("X-Riot-Token", ApiToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("Error occured during request: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (p *PlayerInfo) getIDS() error {
	if p.GameName == "" || p.TagLine == "" {
		err := errors.New("Couldn't retrieve player's IDs without GameName and TagLine !")
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
		err = errors.New("Couldn't get info about player: empty SummonerID")
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
		err = errors.New("Player doesn't have any ranked games")
		return nil, err
	}
	rankedStats = &statsArray[rankedStatsIndex]
	return
}

func getLucasWinRate() (string, error) {
	lucas := new(PlayerInfo)
	lucas.GameName = "lucxsstbn"
	lucas.TagLine = "EUW"
	err := lucas.getIDS()
	if err != nil {
		return "Error getting player's IDs: " + err.Error(), err
	}
	fmt.Println(PrettyPrint(lucas))
	rankedStats, err := lucas.getRankedStats()
	if err != nil {
		return "Error getting player info: " + err.Error(), err
	}

	totalGames := rankedStats.Wins + rankedStats.Losses
	ratio := float64(rankedStats.Wins) / float64(totalGames) * 100

	var s string
	s = "Statistiques de " + lucas.GameName + " cette saison:\n"
	s += "Rang: " + rankedStats.Tier + " " + rankedStats.Rank + "\n"
	s += "Games: " + strconv.Itoa(totalGames) + "\n"
	s += "Victoires: " + strconv.Itoa(rankedStats.Wins) + " \n"
	s += "Défaites: " + strconv.Itoa(rankedStats.Losses) + " \n"
	s += "Ratio: " + strconv.FormatFloat(ratio, 'f', 4, 64) + "% \n"

	return s, nil
}

func Api() (string, error) {
	message, err := getLucasWinRate()
	if err != nil {
		return "An error happened, couldn't get player's stats", err
	}
	return message, nil
}
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
