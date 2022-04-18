package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"

func cleanStr(str string) string {
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ToLower(str)
	return str
}

func WordCount(text string) map[string]int {
	freq := make(map[string]int)
	words := strings.Fields(text)
	wordsLen := len(words)
	workers := 100
	var wg sync.WaitGroup
	ch := make(chan map[string]int, workers+1)
	BATCH_SIZE := wordsLen / workers

	// work on the text in batches of BATCH_SIZE
	for i, j := 0, BATCH_SIZE; i < wordsLen; i, j = j, (j + BATCH_SIZE) {

		// this will occur at the end
		if wordsLen < j {
			j = wordsLen
		}
		wg.Add(1)

		go func(x, y int) {
			subFreq := make(map[string]int)

			for _, word := range words[x:y] {
				subFreq[cleanStr(word)]++
			}

			ch <- subFreq
			wg.Done()
		}(i, j)
	}

	wg.Wait()
	close(ch)

	for m := range ch {
		for w, count := range m {
			freq[w] += count
		}
	}

	return freq

}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	data, err := os.ReadFile(DataFile)
	if err != nil {
		panic(err)
	}

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
