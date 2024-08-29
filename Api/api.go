package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

var (
	ApiToken string = os.Getenv("API_TOKEN")
	ErrJson         = errors.New("can't unmarshal JSON")
	Lucas    *PlayerInfo
)

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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errmsg := "Unexpected status code: " + strconv.Itoa(resp.StatusCode)
		return nil, errors.New(errmsg)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetLucasStats() (string, error) {
	lucas := Lucas

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
	s += "matchIDs des dernieres games: \n"

	// for i, m := range matchIDs {
	// 	s += ("- " + m + ": (")
	// 	match, err := getMatchInfo(matchIDs[i])
	// 	if err != nil {
	// 		return "Error: " + err.Error(), err
	// 	}
	//
	// 	index := slices.IndexFunc(match.Info.Participants, func(p Participant) bool {
	// 		return p.Puuid == lucas.PUUID
	// 	})
	//
	// 	if index == -1 {
	// 		return "Error \n", errors.New("coulnd't find Lucas in participants list")
	// 	}
	// 	l := match.Info.Participants[index]
	//
	// 	stats, err := getChampLevelStats(&match.Info.Participants, l.TeamId)
	// 	if err != nil {
	// 		return "Error: " + err.Error(), err
	// 	}
	//
	// 	s += "Game avg level: " + strconv.FormatFloat(stats.gameStats.avg, 'f', 4, 64) + ", "
	// 	s += "Team avg level: " + strconv.FormatFloat(stats.teamStats.avg, 'f', 4, 64) + ", "
	// 	s += "Lucas level: " + strconv.Itoa(l.ChampLevel) + ")\n"
	// }

	return s, nil
}

func GetMatchStats(matchId string) (string, error) {
	computed, err := ComputeStats(matchId, Lucas.PUUID)
	if err != nil {
		return "Error getting stats of game " + matchId, err
	}

	str := "\n"

	s := reflect.ValueOf(computed).Elem()
	typeOfS := s.Type()
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		// str += fmt.Sprintf("%s: %#v", s.Type().Field(i).Name, field.Interface())
		str += fmt.Sprintf("%s = %v\n\n", typeOfS.Field(i).Name, field.Interface())
	}

	return str, nil
}

func Api() (string, error) {
	message, err := GetLucasStats()
	if err != nil {
		return "An error happened, couldn't get player's stats", err
	}
	return message, nil
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
