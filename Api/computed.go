package api

/* Helpers to make game/team stats about participants */

type IndividualStats[T int | float64] struct {
	max T
	min T
	avg float64
}

type Stats[T int | float64] struct {
	teamStats IndividualStats[T]
	gameStats IndividualStats[T]
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

func getChampLevelStats(participants *[]Participant, teamMates *[]Participant, teamId int) (Stats[int], error) {
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

	return Stats[int]{teamStats, gameStats}, nil
}

func getVisionScoreStats(participants *[]Participant, teamMates *[]Participant, teamId int) (Stats[int], error) {
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

	return Stats[int]{teamStats, gameStats}, nil
}

func getTimeSpentLivingStats(participants *[]Participant, teamMates *[]Participant, teamId int) (Stats[int], error) {
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

	return Stats[int]{teamStats, gameStats}, nil
}

func getDamageDealtToChampionsStats(participants *[]Participant, teamMates *[]Participant, teamId int) (Stats[int], error) {
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

	return Stats[int]{teamStats, gameStats}, nil
}
