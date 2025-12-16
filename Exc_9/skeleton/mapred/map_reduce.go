package mapred

import (
	"regexp"
	"strings"
	"sync"
)

type MapReduce struct {
}

// Run executes the MapReduce algorithm on the input text.
// It acts as the "Master" node, orchestrating Mappers and Reducers.
func (mr MapReduce) Run(input []string) map[string]int {
	// 1. MAP PHASE
	// We use a channel to collect results from all mappers concurrently.
	// Since we don't know exactly how many words we'll get, we treat it as a stream.
	mappedStream := make(chan []KeyValue)
	var mapWg sync.WaitGroup

	// Launch a goroutine for every line/block of text in the input
	for _, textBlock := range input {
		mapWg.Add(1)
		go func(text string) {
			defer mapWg.Done()
			// Call the user-defined mapper
			res := mr.wordCountMapper(text)
			mappedStream <- res
		}(textBlock)
	}

	// We need a separate goroutine to close the channel once all mappers are done.
	// If we don't do this, the range loop below will wait forever (deadlock).
	go func() {
		mapWg.Wait()
		close(mappedStream)
	}()

	// 2. SHUFFLE (GROUP) PHASE
	// Here we organize the data. We take the flattened list of KeyValues
	// and group them by Key (Word) -> List of Values (Counts).
	intermediate := make(map[string][]int)

	for keyValues := range mappedStream {
		for _, kv := range keyValues {
			intermediate[kv.Key] = append(intermediate[kv.Key], kv.Value)
		}
	}

	// 3. REDUCE PHASE
	// Now we process each word and its counts concurrently.
	reducedStream := make(chan KeyValue)
	var reduceWg sync.WaitGroup

	for key, values := range intermediate {
		reduceWg.Add(1)
		go func(k string, v []int) {
			defer reduceWg.Done()
			// Call the user-defined reducer
			reduced := mr.wordCountReducer(k, v)
			reducedStream <- reduced
		}(key, values)
	}

	// Again, close the results channel when all reducers finish
	go func() {
		reduceWg.Wait()
		close(reducedStream)
	}()

	// 4. COLLECT RESULTS
	// Convert the stream of reduced KeyValues into the final map required by the signature
	finalResult := make(map[string]int)
	for kv := range reducedStream {
		finalResult[kv.Key] = kv.Value
	}

	return finalResult
}

// wordCountMapper takes a raw string, cleans it, and splits it into words.
// It emits a KeyValue pair (word, 1) for every occurrence.
func (mr MapReduce) wordCountMapper(text string) []KeyValue {
	// We want to process case-insensitive
	text = strings.ToLower(text)

	// RegEx to replace anything that IS NOT a letter (a-z) with a space.
	// This handles cases like "exc9" -> "exc " effectively stripping the number.
	// It also removes punctuation like ".", "!", etc.
	reg := regexp.MustCompile("[^a-z]+")
	cleanText := reg.ReplaceAllString(text, " ")

	// Split by whitespace to get individual words
	words := strings.Fields(cleanText)

	var results []KeyValue
	for _, w := range words {
		results = append(results, KeyValue{
			Key:   w,
			Value: 1, // Every word counts as 1 initially
		})
	}
	return results
}

// wordCountReducer sums up the occurrences of a specific word.
func (mr MapReduce) wordCountReducer(key string, values []int) KeyValue {
	count := 0
	for _, v := range values {
		count += v
	}

	return KeyValue{
		Key:   key,
		Value: count,
	}
}
