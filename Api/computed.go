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
	"LaneMinionsFirst10Minutes":   true,
}

type IntOrFloat interface {
	int | float64
}

/* Helpers to make game/team stats about participants */

type IndividualStats[T IntOrFloat] struct {
	max T
	min T
	avg float64
}

type Stats[T IntOrFloat] struct {
	name       string
	teamStats  IndividualStats[T]
	gameStats  IndividualStats[T]
	playerStat T
	isTeamMin  bool    // is player the worst of his team
	isGameMin  bool    // is player the worst in the game
	TeamRatio  float64 // ratio between team avg and players score
	GameRatio  float64 // ratio between game avg and players score
}

type MatchComputed struct {
	stats map[string]interface{}
}

// returns a slice of Stat where the player is minimum:
func getMins[T IntOrFloat](match *MatchComputed) []Stats[T] {
	slice := make([]Stats[T], 0)

	for _, stat := range match.stats {
		s := stat.(Stats[T])
		if s.isTeamMin || s.isGameMin {
			slice = append(slice, s)
		}
	}
	return slice
}

// returns a Slice of stat where the player is under average:
func getBadRatios[T IntOrFloat](match *MatchComputed[T]) []Stats[T] {
	slice := make([]Stats[T], 0)

	for _, stat := range match.stats {
		if stat.TeamRatio < 1.0 || stat.GameRatio < 1.0 {
			slice = append(slice, stat)
		}
	}
	return slice
}

func ComputeStats[T IntOrFloat](match *Match, puiid string) (*MatchComputed[T], error) {
	playerIdx := slices.IndexFunc(match.Info.Participants, func(p Participant) bool {
		return p.Puuid == Lucas.PUUID
	})

	if playerIdx == -1 {
		return nil, errors.New("couldn't find player's PUID in participants list")
	}

	player := match.Info.Participants[playerIdx]
	participants := match.Info.Participants
	statsMap := make(map[string]Stats[T])

	for field, isChallengeField := range fields {
		var teamStat [5]T
		var gameStat [10]T
		teamIdx := 0
		for i, p := range participants {
			var score T
			if isChallengeField {
				score = getChallengeFieldInt[T](&p, field)
			} else {
				score = getFieldInt[T](&p, field)
			}
			if p.TeamId == player.TeamId {
				teamStat[teamIdx] = T(score)
				teamIdx++
			}
			gameStat[i] = T(score)
		}
		var playerScore T
		if isChallengeField {
			playerScore = getChallengeFieldInt[T](&player, field)
		} else {
			playerScore = getFieldInt[T](&player, field)
		}
		stat, err := computeStat(gameStat, teamStat, playerScore, field)
		if err != nil {
			return nil, err
		}
		statsMap[field] = stat
	}

	computed := MatchComputed[T]{
		statsMap,
	}

	return &computed, nil
}

func getFieldInt[T IntOrFloat](participant *Participant, field string) T {
	r := reflect.ValueOf(participant)
	f := reflect.Indirect(r).FieldByName(field)
	if f.CanFloat() {
		return T(f.Float())
	}
	return T(f.Int())
	// switch v := f.(type) {
	//  case int:
	//   return int(f.Int())
	//   case float64:
	//   return float64(f.Float())
	// }
}

func getChallengeFieldInt[T IntOrFloat](participant *Participant, field string) T {
	r := reflect.ValueOf(participant.Challenges)
	f := reflect.Indirect(r).FieldByName(field)
	if f.CanFloat() {
		return T(f.Float())
	}
	return T(f.Int())
}

func computeStat[T IntOrFloat](gameScores [10]T, teamScores [5]T, playerScore T, name string) (Stats[T], error) {
	var max, min T
	var avg float64
	max, min, avg = getMaxMinAvg(teamScores[:])
	teamStats := IndividualStats[T]{max, min, avg}
	max, min, avg = getMaxMinAvg(gameScores[:])
	gameStats := IndividualStats[T]{max, min, avg}

	isTeamMin := playerScore == teamStats.min
	isGameMin := playerScore == gameStats.min
	teamRatio := float64(playerScore) / teamStats.avg
	gameRatio := float64(playerScore) / gameStats.avg
	return Stats[T]{name, teamStats, gameStats, playerScore, isTeamMin, isGameMin, teamRatio, gameRatio}, nil
}

func getMaxMinAvg[T IntOrFloat](scores []T) (max T, min T, avg float64) {
	max = 0
	min = 1000000
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
