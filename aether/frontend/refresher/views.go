// Frontend > Refresher > Views
// This package contains the special locations generated by the refresher loop, such as the home view, and the popular. These are called views, and they don't actually compile or generate new data in each way, they're just collections of the existing compiled data in a different order. These are compiled and updated, not rendered on the fly.

package refresher

import (
	// "aether-core/aether/frontend/beapiconsumer"
	// "aether-core/aether/frontend/clapiconsumer"
	"aether-core/aether/frontend/festructs"
	// "aether-core/aether/io/api"
	// pbstructs "aether-core/aether/protos/mimapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	// // "github.com/davecgh/go-spew/spew"
	// // "fmt"
	// "encoding/json"
	"strings"
	// "sync"
	"github.com/asdine/storm/q"
	"time"
)

// GenerateHomeView gets the top 10 most popular items in the communities you subscribe to, and sort them by rank.
func GenerateHomeView() {
	logging.Logf(1, "Home view generator is running")
	start := time.Now()
	// get subscribed boards fingerprints
	sbs := globals.FrontendConfig.ContentRelations.GetAllSubbedBoards()
	// Get the underlying compiled boards
	subbedBoardFps := []string{}
	for k, _ := range sbs {
		if !sbs[k].Notify {
			continue
		}
		subbedBoardFps = append(subbedBoardFps, sbs[k].Fingerprint)
	}
	boardCarriers := *getBoardsByFpList(subbedBoardFps)
	var thrs festructs.CThreadBatch
	for k, _ := range boardCarriers {
		// thrlen := min(len(boardCarriers[k].Threads), 10)
		// boardThreads := boardCarriers[k].Threads[0:thrlen]
		boardThreads := *(boardCarriers[k].GetTopThreadsForView(10))
		for j, _ := range boardThreads {
			boardThreads[j].ViewMeta_BoardName = boardCarriers[k].Boards[0].Name
		}
		thrs = append(thrs, boardThreads...)
	}
	thrs.SortByScore()
	globals.KvInstance.Save(&festructs.HomeViewCarrier{
		Id:      1,
		Threads: thrs,
	})
	elapsed := time.Since(start)
	logging.Logf(1, "Home view generator took %v seconds.", elapsed.Seconds())
}

/*subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(resp.Board.Fingerprint)

we also need to care about notify - it controls what gets into the home view.

*/

// GeneratePopularView gets the top 10 most popular items in each of the whitelisted communities and sorts them by rank.
func GeneratePopularView() {
	logging.Logf(1, "Popular view generator is running")
	start := time.Now()
	boardCarriers := []festructs.BoardCarrier{}
	// // check if sfwlist is disabled
	// if globals.FrontendConfig.ContentRelations.SFWList.GetSFWListDisabled() {
	// 	// sfwlist disabled
	// 	err := globals.KvInstance.All(&boardCarriers)
	// 	if err != nil {
	// 		logging.Logf(1, "Getting boards while SFWList disabled errored out. Error: %v", err)
	// 		return
	// 	}
	// } else {
	// 	// sfwlist enabled
	// 	boardCarriers = *getBoardsByFpList(globals.FrontendConfig.ContentRelations.SFWList.Boards)
	// 	logging.Logf(1, "base board carriers length: %v", len(boardCarriers))
	// 	logging.Logf(1, "sfwlist length: %v", len(globals.FrontendConfig.ContentRelations.SFWList.Boards))

	// }
	boardCarriers = *getBoardsByFpList(globals.FrontendConfig.ContentRelations.SFWList.Boards)
	/*
		^ This is a little weird - if this runs before the sfwlist is pulled in, it will result in an empty popular list. But if we make it so that in the case it's empty it generates the popular list from all communities, we might have NSFW threads surfacing up for people who haven't opted in for that.

		Here, I'm opting to show nothing instead of showing potentially risky data. It's a compromise.
	*/
	logging.Logf(1, "base board carriers length: %v", len(boardCarriers))
	logging.Logf(1, "sfwlist length: %v", len(globals.FrontendConfig.ContentRelations.SFWList.Boards))
	var thrs festructs.CThreadBatch
	for k, _ := range boardCarriers {
		// thrlen := min(len(boardCarriers[k].Threads), 10)
		// boardThreads := boardCarriers[k].Threads[0:thrlen]
		// thrs = append(thrs, boardThreads...)
		boardThreads := *(boardCarriers[k].GetTopThreadsForView(10))
		for j, _ := range boardThreads {
			boardThreads[j].ViewMeta_BoardName = boardCarriers[k].Boards[0].Name
		}
		thrs = append(thrs, boardThreads...)
	}
	thrs.SortByScore()
	existingPopularView := festructs.PopularViewCarrier{}
	err := globals.KvInstance.One("Id", 1, &existingPopularView)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		logging.Logf(1, "Popular view fetch in new popular view creation encountered an error. Error: %v", err)
		return
	}
	if len(thrs) > 0 || len(existingPopularView.Threads) == 0 {
		globals.KvInstance.Save(&festructs.PopularViewCarrier{
			Id:      1,
			Threads: thrs,
		})
	} else {
		logging.Logf(1, "Popular view produced zero threads and thus bailed on updating. This is something that should be looked at.") // TODO FUTURE
	}

	elapsed := time.Since(start)
	logging.Logf(1, "Popular items count: %v", len(thrs))
	logging.Logf(1, "Popular view generator took %v seconds.", elapsed.Seconds())
}

func getBoardsByFpList(boardFingerprints []string) *[]festructs.BoardCarrier {
	query := globals.KvInstance.Select(q.In("Fingerprint", boardFingerprints))
	var bcs []festructs.BoardCarrier
	query.Find(&bcs)
	return &bcs
}

/*----------  Internal util functions  ----------*/

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
