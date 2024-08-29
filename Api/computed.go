package api

import (
	"errors"
	"slices"
)

/* Helpers to make game/team stats about participants */

type IndividualStats[T int | float64] struct {
	max T
	min T
	avg float64
}

type Stats[T int | float64] struct {
	teamStats  IndividualStats[T]
	gameStats  IndividualStats[T]
	playerStat T
}

type MatchComputed struct {
	ChampLevel             Stats[int]
	VisionScore            Stats[int]
	TimeSpentLiving        Stats[int]
	DamageDealtToChampions Stats[int]
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

	champLevel, err := getChampLevelStats(&participants, team)
	if err != nil {
		return nil, err
	}
	visionScore, err := getVisionScoreStats(&participants, team)
	if err != nil {
		return nil, err
	}
	timeSpentLiving, err := getTimeSpentLivingStats(&participants, team)
	if err != nil {
		return nil, err
	}

	damageDealtToChampions, err := getDamageDealtToChampionsStats(&participants, team)
	if err != nil {
		return nil, err
	}

	champLevel.playerStat = player.ChampLevel
	visionScore.playerStat = player.VisionScore
	timeSpentLiving.playerStat = player.LongestTimeSpentLiving
	damageDealtToChampions.playerStat = player.TotalDamageDealtToChampions

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

func getChampLevelStats(participants *[]Participant, teamMates *[]Participant) (Stats[int], error) {
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

	return Stats[int]{teamStats, gameStats, 0}, nil
}

func getVisionScoreStats(participants *[]Participant, teamMates *[]Participant) (Stats[int], error) {
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

	return Stats[int]{teamStats, gameStats, 0}, nil
}

func getTimeSpentLivingStats(participants *[]Participant, teamMates *[]Participant) (Stats[int], error) {
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

	return Stats[int]{teamStats, gameStats, 0}, nil
}

func getDamageDealtToChampionsStats(participants *[]Participant, teamMates *[]Participant) (Stats[int], error) {
	max := 0
	min := 100000
	avg := 0.0

	for _, p := range *teamMates {
		if p.TotalDamageDealtToChampions > max {
			max = p.TotalDamageDealtToChampions
		}
		if p.LongestTimeSpentLiving < min {
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
		if p.LongestTimeSpentLiving > max {
			max = p.TotalDamageDealtToChampions
		}
		if p.LongestTimeSpentLiving < min {
			min = p.TotalDamageDealtToChampions
		}
		avg += float64(p.TotalDamageDealtToChampions)
	}
	avg = avg / 10
	gameStats := IndividualStats[int]{max, min, avg}

	return Stats[int]{teamStats, gameStats, 0}, nil
}
