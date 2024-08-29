package main

import (
	"fmt"
	"log"
	"os"
	"time"

	bot "github.com/Nvim/silverstalker/Bot"

	api "github.com/Nvim/silverstalker/Api"
)

func main() {
	// err := godotenv.Load(".envrc")
	// if err != nil {
	// 	log.Fatal("Couldn't load .env: ", err)
	// 	return
	// }
	/* Lucas Init: */
	api.Lucas = new(api.PlayerInfo)
	lucas := api.Lucas
	lucas.GameName = "lucxsstbn"
	lucas.TagLine = "EUW"
	err := lucas.GetIDs()
	if err != nil {
		log.Fatal("Error getting player's IDs: " + err.Error())
		return
	}
	fmt.Println(api.PrettyPrint(lucas))

	/* Bot Init: */
	bot.BotToken = os.Getenv("BOT_TOKEN")
	err = bot.Init()
	if err != nil {
		log.Fatal("Error creating bot: " + err.Error())
		return
	}

	/* Get Lucas latest game: */
	matchIDs, err := lucas.GetLatestMatches()
	if err != nil {
		log.Fatal("Error getting player matches: " + err.Error())
		return
	}
	latestId := matchIDs[0]
	latestMatch, err := api.GetMatchInfo(latestId)
	if err != nil {
		log.Fatal("Error getting match info: " + err.Error())
		return
	}

	/* Get timestamps: */
	endTime := latestMatch.Info.GameEndTimestamp
	eventTimestamp := time.Unix(0, endTime*int64(time.Millisecond))
	// sleepDuration := time.Until(eventTimestamp.Add(20 * time.Minute))

	log.Println("Latest game ID: ", latestId)
	log.Println("Latest game ended at: ", eventTimestamp.UTC())

	// TODO: this calls api.GetMatchInfo a 2nd timekj
	stats, err := api.GetMatchStats(latestId)
	if err != nil {
		log.Fatal("Error gettting match stats: " + err.Error())
	}
	log.Println("Stats: " + stats)
	return

	// if sleepDuration > 0 {
	// 	log.Println("Sleeping until: ", (time.Now().Add(sleepDuration)).UTC())
	// 	time.Sleep(sleepDuration)
	// }
	// log.Println("Woke up from sleep")

	// Periodic fetching every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {

		/* Fetch Data: */
		log.Println("Fetching data...")
		matchIDs, err = lucas.GetLatestMatches()
		if err != nil {
			log.Fatal("Error getting player matches: " + err.Error())
			continue
		}
		if matchIDs[0] == latestId {
			// no new data
			log.Println("No new match data, latest match ID still is ", latestId)
			continue
		}
		latestId = matchIDs[0]
		latestMatch, err := api.GetMatchInfo(latestId)
		if err != nil {
			log.Fatal("Error getting match info: " + err.Error())
			continue
		}

		/* Get timestamps for newest game: */
		endTime := latestMatch.Info.GameEndTimestamp
		eventTimestamp := time.Unix(0, endTime*int64(time.Millisecond))

		log.Println("New latest game ID: ", latestId)
		log.Println("New latest game end: ", eventTimestamp.UTC())

		/* Send message: */
		_ = bot.SendMessage(fmt.Sprintf("Nouvelle game: %s", latestId))

		/* Go to sleep: */
		sleepDuration := time.Until(eventTimestamp.Add(20 * time.Minute))
		if sleepDuration > 0 {
			log.Println("Sleeping until: ", (time.Now().Add(sleepDuration)).UTC())
			time.Sleep(sleepDuration)
		}
		log.Println("Woke up from sleep")

		/* Reset Ticket */
		ticker.Reset(10 * time.Minute)
	}

	// err = bot.Listen()
	// if err != nil {
	// 	return
	// }
}
