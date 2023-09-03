package util

import (
	"sync"
)

// Merge fans multiple error channels in to a single error channel
func Merge(errChans ...<-chan error) <-chan error {
	mergedChan := make(chan error)

	// Create a WaitGroup that waits for all error channels to close
	var wg sync.WaitGroup
	wg.Add(len(errChans))
	go func() {
		// When all error channels are closed, close the merged channel
		wg.Wait()
		close(mergedChan)
	}()

	// Wait for each error channel to close
	for i := range errChans {
		go func(errChan <-chan error) {
			for err := range errChan {
				if err != nil {
					// Fan the contents of each error channel into the merged channel
					mergedChan <- err
				}
			}
			// Tell the WaitGroup that one of the error channels is closed
			wg.Done()
		}(errChans[i])
	}

	return mergedChan
}
