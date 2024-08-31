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
	ChampLevel             Stats[int]
	VisionScore            Stats[int]
	TimeSpentLiving        Stats[int]
	DamageDealtToChampions Stats[int]
}

// returns a slice of Stat where the player is minimum:
func (match MatchComputed) getMins() []Stats[int] {
	slice := make([]Stats[int], 0)

	v := reflect.ValueOf(match)
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
func (match MatchComputed) getBadRatios() []Stats[int] {
	slice := make([]Stats[int], 0)

	v := reflect.ValueOf(match)
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

func ComputeStats(matchId string, puiid string) (*MatchComputed, error) {
	match, err := GetMatchInfo(matchId)
	if err != nil {
		return nil, err
	}

	playerIdx := slices.IndexFunc(match.Info.Participants, func(p Participant) bool {
		return p.Puuid == Lucas.PUUID
	})

	if playerIdx == -1 {
		return nil, errors.New("Couldn't find player's PUID in participants list")
	}

	player := match.Info.Participants[playerIdx]
	participants := match.Info.Participants
	team := filterParticipantsByTeamId(&participants, player.TeamId)

	champLevel, err := getChampLevelStats(&participants, team, player.ChampLevel)
	if err != nil {
		return nil, err
	}
	visionScore, err := getVisionScoreStats(&participants, team, player.VisionScore)
	if err != nil {
		return nil, err
	}
	timeSpentLiving, err := getTimeSpentLivingStats(&participants, team, player.LongestTimeSpentLiving)
	if err != nil {
		return nil, err
	}

	damageDealtToChampions, err := getDamageDealtToChampionsStats(&participants, team, player.TotalDamageDealtToChampions)
	if err != nil {
		return nil, err
	}

	computed := MatchComputed{
		champLevel,
		visionScore,
		timeSpentLiving,
		damageDealtToChampions,
	}

	return &computed, nil
}

func filterParticipantsByTeamId(participants *[]Participant, teamId int) *[]Participant {
	var filteredParticipants []Participant
	for _, participant := range *participants {
		if participant.TeamId == teamId {
			filteredParticipants = append(filteredParticipants, participant)
		}
	}
	return &filteredParticipants
}

func getChampLevelStats(participants *[]Participant, teamMates *[]Participant, playerStat int) (Stats[int], error) {
	max := 1
	min := 18
	avg := 0.0

	for _, p := range *teamMates {
		if p.ChampLevel > max {
			max = p.ChampLevel
		}
		if p.ChampLevel < min {
			min = p.ChampLevel
		}
		avg += float64(p.ChampLevel)
	}
	avg = avg / 5
	teamStats := IndividualStats[int]{max, min, avg}

	max = 1
	min = 18
	avg = 0.0
	for _, p := range *participants {
		if p.ChampLevel > max {
			max = p.ChampLevel
		}
		if p.ChampLevel < min {
			min = p.ChampLevel
		}
		avg += float64(p.ChampLevel)
	}
	avg = avg / 10
	gameStats := IndividualStats[int]{max, min, avg}

	isTeamMin := playerStat == teamStats.min
	isGameMin := playerStat == gameStats.min
	teamRatio := float64(playerStat) / teamStats.avg
	gameRatio := float64(playerStat) / gameStats.avg
	return Stats[int]{"Niveau", teamStats, gameStats, playerStat, isTeamMin, isGameMin, teamRatio, gameRatio}, nil
}

func getVisionScoreStats(participants *[]Participant, teamMates *[]Participant, playerStat int) (Stats[int], error) {
	max := 0
	min := 100000
	avg := 0.0

	for _, p := range *teamMates {
		if p.VisionScore > max {
			max = p.VisionScore
		}
		if p.VisionScore < min {
			min = p.VisionScore
		}
		avg += float64(p.VisionScore)
	}
	avg = avg / 5
	teamStats := IndividualStats[int]{max, min, avg}
	max = 0
	min = 100000
	avg = 0.0
	for _, p := range *participants {
		if p.VisionScore > max {
			max = p.VisionScore
		}
		if p.VisionScore < min {
			min = p.VisionScore
		}
		avg += float64(p.VisionScore)
	}
	avg = avg / 10
	gameStats := IndividualStats[int]{max, min, avg}

	isTeamMin := playerStat == teamStats.min
	isGameMin := playerStat == gameStats.min
	teamRatio := float64(playerStat) / teamStats.avg
	gameRatio := float64(playerStat) / gameStats.avg
	return Stats[int]{"Score de vision", teamStats, gameStats, playerStat, isTeamMin, isGameMin, teamRatio, gameRatio}, nil
}

func getTimeSpentLivingStats(participants *[]Participant, teamMates *[]Participant, playerStat int) (Stats[int], error) {
	max := 0
	min := 100000
	avg := 0.0

	for _, p := range *teamMates {
		if p.LongestTimeSpentLiving > max {
			max = p.LongestTimeSpentLiving
		}
		if p.LongestTimeSpentLiving < min {
			min = p.LongestTimeSpentLiving
		}
		avg += float64(p.LongestTimeSpentLiving)
	}
	avg = avg / 5
	teamStats := IndividualStats[int]{max, min, avg}
	max = 0
	min = 100000
	avg = 0.0
	for _, p := range *participants {
		if p.LongestTimeSpentLiving > max {
			max = p.LongestTimeSpentLiving
		}
		if p.LongestTimeSpentLiving < min {
			min = p.LongestTimeSpentLiving
		}
		avg += float64(p.LongestTimeSpentLiving)
	}
	avg = avg / 10
	gameStats := IndividualStats[int]{max, min, avg}

	isTeamMin := playerStat == teamStats.min
	isGameMin := playerStat == gameStats.min
	teamRatio := float64(playerStat) / teamStats.avg
	gameRatio := float64(playerStat) / gameStats.avg
	return Stats[int]{"Plus longue durée passée en vie", teamStats, gameStats, playerStat, isTeamMin, isGameMin, teamRatio, gameRatio}, nil
}

func getDamageDealtToChampionsStats(participants *[]Participant, teamMates *[]Participant, playerStat int) (Stats[int], error) {
	max := 0
	min := 100000
	avg := 0.0

	for _, p := range *teamMates {
		if p.TotalDamageDealtToChampions > max {
			max = p.TotalDamageDealtToChampions
		}
		if p.TotalDamageDealtToChampions < min {
			min = p.TotalDamageDealtToChampions
		}
		avg += float64(p.TotalDamageDealtToChampions)
	}
	avg = avg / 5
	teamStats := IndividualStats[int]{max, min, avg}
	max = 0
	min = 100000
	avg = 0.0
	for _, p := range *participants {
		if p.TotalDamageDealtToChampions > max {
			max = p.TotalDamageDealtToChampions
		}
		if p.TotalDamageDealtToChampions < min {
			min = p.TotalDamageDealtToChampions
		}
		avg += float64(p.TotalDamageDealtToChampions)
	}
	avg = avg / 10
	gameStats := IndividualStats[int]{max, min, avg}

	isTeamMin := playerStat == teamStats.min
	isGameMin := playerStat == gameStats.min
	teamRatio := float64(playerStat) / teamStats.avg
	gameRatio := float64(playerStat) / gameStats.avg
	return Stats[int]{"Dégâts aux champions ennemis", teamStats, gameStats, playerStat, isTeamMin, isGameMin, teamRatio, gameRatio}, nil
}
