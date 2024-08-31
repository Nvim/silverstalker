package api

import (
	"errors"
	"reflect"
	"slices"
)

/* Helpers to make game/team stats about participants */

type IndividualStats[T int | float64] struct {
	max T
	min T
	avg float64
}

type Stats[T int | float64] struct {
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
	ChampLevel                Stats[int]
	VisionScore               Stats[int]
	TimeSpentLiving           Stats[int]
	DamageDealtToChampions    Stats[int]
	LaneMinionsFirst10Minutes Stats[int]
}

// returns a slice of Stat where the player is minimum:
func getMins(match *MatchComputed) []Stats[int] {
	slice := make([]Stats[int], 0)

	v := reflect.ValueOf(*match)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		// field := t.Field(i)
		value := v.Field(i)
		if value.Kind() == reflect.Struct {
			stats := value.Interface().(Stats[int])

			teamMin := reflect.ValueOf(stats).FieldByName("isTeamMin")
			gameMin := reflect.ValueOf(stats).FieldByName("isGameMin")
			if teamMin.Bool() || gameMin.Bool() {
				slice = append(slice, stats)
			}
		}
	}
	return slice
}

// returns a Slice of stat where the player is under average:
func getBadRatios(match *MatchComputed) []Stats[int] {
	slice := make([]Stats[int], 0)

	v := reflect.ValueOf(*match)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		// field := t.Field(i)
		value := v.Field(i)
		if value.Kind() == reflect.Struct {
			stats := value.Interface().(Stats[int])

			teamRatio := reflect.ValueOf(stats).FieldByName("TeamRatio")
			gameRatio := reflect.ValueOf(stats).FieldByName("GameRatio")
			if teamRatio.Float() < 1.0 || gameRatio.Float() < 1.0 {
				slice = append(slice, stats)
			}
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
	// teamMates := filterParticipantsByTeamId(&participants, player.TeamId)

	var champLevels [10]int
	var visionScores [10]int
	var timeSpentLiving [10]int
	var damageDealtToChampions [10]int
	var laneMinions10 [10]int
	teamChampLevels := make([]int, 5)
	teamVisionScores := make([]int, 5)
	teamTimeSpentLiving := make([]int, 5)
	teamDamageDealtToChampions := make([]int, 5)
	teamLaneMinions10 := make([]int, 5)
	for i, p := range participants {
		if p.TeamId == player.TeamId {
			teamChampLevels[i] = p.ChampLevel
			teamVisionScores[i] = p.VisionScore
			teamTimeSpentLiving[i] = p.LongestTimeSpentLiving
			teamDamageDealtToChampions[i] = p.TotalDamageDealtToChampions
			teamLaneMinions10[i] = p.Challenges.LaneMinionsFirst10Minutes
		}
		champLevels[i] = p.ChampLevel
		visionScores[i] = p.VisionScore
		timeSpentLiving[i] = p.LongestTimeSpentLiving
		damageDealtToChampions[i] = p.TotalDamageDealtToChampions
		laneMinions10[i] = p.Challenges.LaneMinionsFirst10Minutes
	}

	champLevelStat, err := computeStat(champLevels, teamChampLevels, player.ChampLevel, "Niveau")
	if err != nil {
		return nil, err
	}
	visionScoreStat, err := computeStat(visionScores, teamVisionScores, player.VisionScore, "Score de vision")
	if err != nil {
		return nil, err
	}
	timeSpentLivingStat, err := computeStat(timeSpentLiving, teamTimeSpentLiving, player.LongestTimeSpentLiving, "Plus longue durée sans mourir")
	if err != nil {
		return nil, err
	}
	damageDealtToChampionsStat, err := computeStat(damageDealtToChampions, teamDamageDealtToChampions, player.TotalDamageDealtToChampions, "Dégâts aux champions ennemis")
	if err != nil {
		return nil, err
	}
	laneMinions10Stat, err := computeStat(laneMinions10, teamLaneMinions10, player.Challenges.LaneMinionsFirst10Minutes, "Farm à 10 minutes")
	if err != nil {
		return nil, err
	}

	computed := MatchComputed{
		champLevelStat,
		visionScoreStat,
		timeSpentLivingStat,
		damageDealtToChampionsStat,
		laneMinions10Stat,
	}

	return &computed, nil
}

func computeStat(gameScores [10]int, teamScores []int, playerScore int, name string) (Stats[int], error) {
	var max int = 0
	var min int = 1000000
	avg := 0.0

	for _, p := range teamScores {
		if p > max {
			max = p
		}
		if p < min {
			min = p
		}
		avg += float64(p)
	}
	avg = avg / 5
	teamStats := IndividualStats[int]{max, min, avg}
	max = 0
	min = 100000
	avg = 0.0
	for _, p := range gameScores {
		if p > max {
			max = p
		}
		if p < min {
			min = p
		}
		avg += float64(p)
	}
	avg = avg / 10
	gameStats := IndividualStats[int]{max, min, avg}

	isTeamMin := playerScore == teamStats.min
	isGameMin := playerScore == gameStats.min
	teamRatio := float64(playerScore) / teamStats.avg
	gameRatio := float64(playerScore) / gameStats.avg
	return Stats[int]{name, teamStats, gameStats, playerScore, isTeamMin, isGameMin, teamRatio, gameRatio}, nil
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
