package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
)

var (
	ApiToken   string = os.Getenv("API_TOKEN")
	ErrJson           = errors.New("can't unmarshal JSON")
	Lucas      *PlayerInfo
	fieldNames = map[string]string{
		"ChampLevel":                  "Niveau",
		"VisionScore":                 "Score de vision",
		"LongestTimeSpentLiving":      "Plus longue durée passée en vie",
		"TotalDamageDealtToChampions": "Dégâts aux champions ennemis",
		"LaneMinionsFirst10Minutes":   "Farm à 10 minutes",
		"GoldEarned":                  "Gold Obtenu",
		"WardsPlaced":                 "Wards placées",
		"SoloKills":                   "Solo Kills",
		"DamagePerMinute":             "Dégâts par minute",
		"Kda":                         "KDA",
		"GoldPerMinute":               "Gold Par minute",
	}
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

	return s, nil
}

func GetMatchMetaString(match *Match) (string, error) {
	playerIdx := slices.IndexFunc(match.Info.Participants, func(p Participant) bool {
		return p.Puuid == Lucas.PUUID
	})
	if playerIdx == -1 {
		return "", errors.New("couldn't find player's index")
	}
	player := match.Info.Participants[playerIdx]

	s := "🚨Nouvelle game! 🚨\n"
	if player.Win {
		s += "- Victoire🎉 (on va quand même te trash mon con)\n"
	} else {
		s += "- Défaite\n"
	}

	s += fmt.Sprintf("- Champ: %s (%s)\n", player.ChampionName, player.IndividualPosition)
	s += fmt.Sprintf("- %d/%d/%d (KDA: %.2f)\n", player.Kills, player.Deaths, player.Assists, player.Challenges.Kda)

	return s, nil
}

func GetMatchStatsString(match *Match) (string, error) {
	computed, err := ComputeStats(match, Lucas.PUUID)
	if err != nil {
		return "Error getting stats of game " + match.Metadata.MatchID, err
	}

	minSlice := getMins(computed)
	str := "Pires stats de la game: 🫵\n"
	for _, stat := range minSlice {
		str += fmt.Sprintf("* %s:: %.2f (Moyenne de l'équipe: %.2f, Moyenne de la game: %.2f)\n", fieldNames[stat.name], stat.playerStat, stat.teamStats.avg, stat.gameStats.avg)
	}
	if len(minSlice) < 4 {
		badSlice := getBadRatios(computed)
		for _, stat := range badSlice {
			if !SliceContains(minSlice, stat) {
				str += fmt.Sprintf("- %s: %.2f (Moyenne de l'équipe: %.2f, Moyenne de la game: %.2f)\n", fieldNames[stat.name], stat.playerStat, stat.teamStats.avg, stat.gameStats.avg)
			}
		}
	}

	return str, nil
}

// Message that will be sent by the bot:
func GetMatchDescString(match *Match) (string, error) {
	meta, err := GetMatchMetaString(match)
	if err != nil {
		return "Error gettting match stats: " + err.Error(), err
	}
	stats, err := GetMatchStatsString(match)
	if err != nil {
		return "Error gettting match stats: " + err.Error(), err
	}

	return fmt.Sprintf("%s%s", meta, stats), nil
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

func SliceContains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
