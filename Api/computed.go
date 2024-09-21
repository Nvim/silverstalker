package api

import (
	"errors"
	"reflect"
	"slices"
)

var fields = map[string]bool{
	"ChampLevel":                  false,
	"VisionScore":                 false,
	"LongestTimeSpentLiving":      false,
	"TotalDamageDealtToChampions": false,
	"GoldEarned":                  false,
	"WardsPlaced":                 false,
	"SoloKills":                   true,
	"DamagePerMinute":             true,
	"LaneMinionsFirst10Minutes":   true,
	"Kda":                         true,
	"GoldPerMinute":               true,
}

/* Helpers to make game/team stats about participants */

type IndividualStats struct {
	max float64
	min float64
	avg float64
}

type Stats struct {
	name       string
	teamStats  IndividualStats
	gameStats  IndividualStats
	playerStat float64
	isTeamMin  bool    // is player the worst of his team
	isGameMin  bool    // is player the worst in the game
	TeamRatio  float64 // ratio between team avg and players score
	GameRatio  float64 // ratio between game avg and players score
}

type MatchComputed struct {
	stats map[string]Stats // TODO: generic type instead of hard-coded int
}

// returns a slice of Stat where the player is minimum:
func getMins(match *MatchComputed) []Stats {
	slice := make([]Stats, 0)

	for _, stat := range match.stats {
		if stat.isTeamMin || stat.isGameMin {
			slice = append(slice, stat)
		}
	}
	return slice
}

// returns a Slice of stat where the player is under average:
func getBadRatios(match *MatchComputed) []Stats {
	slice := make([]Stats, 0)

	for _, stat := range match.stats {
		if stat.TeamRatio < 1.0 || stat.GameRatio < 1.0 {
			slice = append(slice, stat)
		}
	}
	return slice
}

func ComputeStats(match *Match, puiid string) (*MatchComputed, error) {
	playerIdx := slices.IndexFunc(match.Info.Participants, func(p Participant) bool {
		return p.Puuid == Lucas.PUUID
	})

	if playerIdx == -1 {
		return nil, errors.New("couldn't find player's PUID in participants list")
	}

	player := match.Info.Participants[playerIdx]
	participants := match.Info.Participants
	statsMap := make(map[string]Stats)

	for field, isChallengeField := range fields {
		var teamStat [5]float64
		var gameStat [10]float64
		teamIdx := 0
		for i, p := range participants {
			var score float64
			if isChallengeField {
				score = getChallengeFieldInt(&p, field)
			} else {
				score = getFieldInt(&p, field)
			}
			if p.TeamId == player.TeamId {
				teamStat[teamIdx] = score
				teamIdx++
			}
			gameStat[i] = score
		}
		var playerScore float64
		if isChallengeField {
			playerScore = getChallengeFieldInt(&player, field)
		} else {
			playerScore = getFieldInt(&player, field)
		}
		stat, err := computeStat(gameStat, teamStat, playerScore, field)
		if err != nil {
			return nil, err
		}
		statsMap[field] = stat
	}

	computed := MatchComputed{
		statsMap,
	}

	return &computed, nil
}

func getFieldInt(participant *Participant, field string) float64 {
	r := reflect.ValueOf(participant)
	f := reflect.Indirect(r).FieldByName(field)
	if f.CanFloat() {
		return f.Float()
	}
	return float64(f.Int())
}

func getChallengeFieldInt(participant *Participant, field string) float64 {
	r := reflect.ValueOf(participant.Challenges)
	f := reflect.Indirect(r).FieldByName(field)
	if f.CanFloat() {
		return f.Float()
	}
	return float64(f.Int())
}

func computeStat(gameScores [10]float64, teamScores [5]float64, playerScore float64, name string) (Stats, error) {
	var max, min float64
	var avg float64
	max, min, avg = getMaxMinAvg(teamScores[:])
	teamStats := IndividualStats{max, min, avg}
	max, min, avg = getMaxMinAvg(gameScores[:])
	gameStats := IndividualStats{max, min, avg}

	isTeamMin := playerScore == teamStats.min
	isGameMin := playerScore == gameStats.min
	teamRatio := float64(playerScore) / teamStats.avg
	gameRatio := float64(playerScore) / gameStats.avg
	return Stats{name, teamStats, gameStats, playerScore, isTeamMin, isGameMin, teamRatio, gameRatio}, nil
}

func getMaxMinAvg(scores []float64) (max float64, min float64, avg float64) {
	max = 0.0
	min = 1000000.0
	avg = 0.0

	for _, val := range scores {
		if val > max {
			max = val
		}
		if val < min {
			min = val
		}
		avg += float64(val)
	}
	avg = avg / float64(len(scores))

	return max, min, avg
}

// func filterParticipantsByTeamId(participants *[]Participant, teamId int) *[]Participant {
// 	var filteredParticipants []Participant
// 	for _, participant := range *participants {
// 		if participant.TeamId == teamId {
// 			filteredParticipants = append(filteredParticipants, participant)
// 		}
// 	}
// 	return &filteredParticipants
// }
