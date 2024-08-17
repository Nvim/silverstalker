package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
)

var ApiToken string = ""

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
func GetRiotApi(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Couldn't create request: ", err)
	}
	req.Header.Add("X-Riot-Token", ApiToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error occured during request: ", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body
}

func (p *PlayerInfo) getIDS() {
	if p.GameName == "" || p.TagLine == "" {
		fmt.Println("Error: can't retrieve player's IDs without GameName and TagLine !")
		return
	}
	puidUrl := "https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/" + p.GameName + "/" + p.TagLine

	var puidResponse AccountJSON
	if err := json.Unmarshal(GetRiotApi(puidUrl), &puidResponse); err != nil {
		fmt.Println("Can't unmarshal json!")
	}

	var summonerResponse SummonerJSON
	summonerUrl := "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/" + puidResponse.Puuid
	if err := json.Unmarshal(GetRiotApi(summonerUrl), &summonerResponse); err != nil {
		fmt.Println("Can't unmarshal json!")
	}
	p.PUUID = puidResponse.Puuid
	p.AccountID = summonerResponse.AccountID
	p.SummonerID = summonerResponse.ID
}

func (p *PlayerInfo) getRankedStats() (rankedStats LeagueStats) {
	if p.SummonerID == "" {
		fmt.Println("Error: can't get info about player: empty SummonerID")
		return
	}

	statsUrl := "https://euw1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + p.SummonerID

	var statsArray []LeagueStats
	if err := json.Unmarshal(GetRiotApi(statsUrl), &statsArray); err != nil {
		fmt.Println("Error parsing Leagues json!")
	}

	rankedStatsIndex := slices.IndexFunc(statsArray, func(s LeagueStats) bool {
		return s.QueueType == "RANKED_SOLO_5x5"
	})
	rankedStats = statsArray[rankedStatsIndex]
	return
}

func getLucasWinRate() string {
	lucas := new(PlayerInfo)
	lucas.GameName = "lucxsstbn"
	lucas.TagLine = "EUW"
	lucas.getIDS()
	fmt.Println(PrettyPrint(lucas))
	rankedStats := lucas.getRankedStats()
	// if rankedStats != nil {
	// 	fmt.Println("Couldn't get player's stats!")
	// 	return "Error"
	// }

	totalGames := rankedStats.Wins + rankedStats.Losses
	ratio := float64(rankedStats.Wins) / float64(totalGames) * 100

	var s string
	s = "Statistiques de " + lucas.GameName + " cette saison:\n"
	s += "Rang: " + rankedStats.Tier + " " + rankedStats.Rank + "\n"
	s += "Games: " + strconv.Itoa(totalGames) + "\n"
	s += "Victoires: " + strconv.Itoa(rankedStats.Wins) + " \n"
	s += "Défaites: " + strconv.Itoa(rankedStats.Losses) + " \n"
	s += "Ratio: " + strconv.FormatFloat(ratio, 'f', 4, 64) + "% \n"

	return s
}

func Api() string {
	// lucas := new(PlayerInfo)
	// lucas.GameName = "lucxsstbn"
	// lucas.TagLine = "EUW"
	// lucas.getIDS()
	// fmt.Println(PrettyPrint(lucas))
	return getLucasWinRate()
}
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
