package main

import (
	"runtime"
	"fmt"
	"log"
	"time"
	"bytes"
)

type Assertions struct {
	Bytes   []byte
	Result   int
}

func maxOnesAfterRemoveItem(b []byte) int {
	var i, result int

	x := map[int]int{}

	i = 0;
	result = 0;
	for _, v := range b {
		if (v==1){
			x[i] ++;
		}else {
			i++;
		}

		if (i==0){
			result = x[i]
		}

		if (i>0 && (x[i] + x[i-1]) > result){
			result = x[i] + x[i-1];
		}
	}
	return result;
}

func main() {

	var Assertions = []Assertions{
		{Bytes: []byte{0}, Result: 0},
		{Bytes: []byte{1}, Result: 1},
		{Bytes: []byte{0, 0}, Result: 0},
		{Bytes: []byte{0, 1}, Result: 1},
		{Bytes: []byte{1, 0}, Result: 1},
		{Bytes: []byte{1, 1}, Result: 2},
		{Bytes: []byte{1, 1, 1}, Result: 3},
		{Bytes: []byte{1, 1, 0, 1, 1}, Result: 4},
		{Bytes: []byte{1, 1, 1, 0, 1, 0, 1, 1}, Result: 4},
		{Bytes: []byte{1, 1, 0, 1, 1, 0, 1, 1, 1}, Result: 5},
		{Bytes: []byte{1, 1, 0, 1, 1, 0, 1, 1, 1, 0}, Result: 5},
	}

	log.Printf("Memory usage before start")
	printMemUsage()

	assertAlg(Assertions);

	printMemUsage()
	log.Printf("---------")
	// Force GC to clear up, should see a memory drop
	runtime.GC()

	assertAlg2(Assertions);

	log.Printf("Memory after finish")
	printMemUsage()
}

func assertAlg(Assertions []Assertions){
	defer timeTrack(time.Now(), "function assertAlg()")
	for _, v := range Assertions{
		assert(v.Bytes, v.Result, true);
	}
}

func assertAlg2(Assertions []Assertions){
	defer timeTrack(time.Now(), "function assertAlg2()")
	for _, v := range Assertions{
		assertV2(v.Bytes, v.Result, true);
	}
}

func assert(bytes []byte, result int, showMessage bool) {
	afterRemoveItem := maxOnesAfterRemoveItem(bytes)
	parseResult(bytes, afterRemoveItem, result, showMessage);
}

func assertV2(bytes []byte, result int, showMessage bool) {
	afterRemoveItem := maxOnesAfterRemoveItemV2(bytes)
	parseResult(bytes, afterRemoveItem, result, showMessage);
}

func parseResult(bytes []byte, afterRemoveItem int,  result int, showMessage bool){
	if !showMessage{
		return;
	}

	if afterRemoveItem != result {
		log.Printf("Error: longest sequence for %v is %v, must be %v", bytes, afterRemoveItem, result)
	}else{
		log.Printf("longest sequence for %v is %v", bytes, result)
	}
}

func maxOnesAfterRemoveItemV2(b []byte) int {
	if (len(b) == 1){
		return int(b[0])
	}

	s := bytes.Split(b, []byte{0})

	pairsCount := len(s);
	if (pairsCount == 1){
		return len(s[0])
	}

	var result, i, prevPairLength int;

	result = 0;
	i = 0
	prevPairLength = 0;

	for i < pairsCount{
		pairLength := len(s[i])
		if (i==0){
			result = pairLength
		}else {
			if (prevPairLength+pairLength > result) {
				result = prevPairLength + pairLength;
			}
		}
		prevPairLength = pairLength;
		i++;
	}

	return result;
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v kB", bToKb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v kB", bToKb(m.TotalAlloc))
	fmt.Printf("\tSys = %v kB", bToKb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
