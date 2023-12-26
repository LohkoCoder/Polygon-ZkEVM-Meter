package gasUsedL12

import "sync"

var wg *sync.WaitGroup

func Fetch2Months() {
	startHeight := 17945670
	endHeight := startHeight + 4000*100

	for startHeight <= endHeight {
		WatchBatches(int64(startHeight), int64(startHeight+100))
		startHeight += 101
	}

}
