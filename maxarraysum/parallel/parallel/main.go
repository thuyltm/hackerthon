package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	maxSubsetSumUnderXChannel     chan Input              = make(chan Input)
	combineChannel                chan SortedInput        = make(chan SortedInput)
	maxSubsetSumAboveXChannel     chan CombineSortedInput = make(chan CombineSortedInput)
	sortCombineSortedInputChannel chan SortedInput        = make(chan SortedInput)
	finishChannel                 chan int32              = make(chan int32)
)

func maxSubsetSumUnder10000(arr []int32,
	maxSubsetSumUnderXChannel <-chan Input) {
	data := <-maxSubsetSumUnderXChannel
	result := divideAndConquer(arr, data.start, data.end)
	combineChannel <- SortedInput{
		batch:  data.batch,
		result: result,
	}
}

func combineMaxSubsetSumAbove10000(arr []int32,
	maxSubsetSumAboveXChannel <-chan CombineSortedInput) {
	data := <-maxSubsetSumAboveXChannel
	result := combine(arr, data.result[0], data.result[1])
	//log.Println("Finish combineMaxSubsetSumAbove10000 ", data.batch, data.result[0].start, data.result[1].end)
	sortCombineSortedInputChannel <- SortedInput{
		batch:  data.batch,
		result: result,
	}
}

func loopSortCombineSortedInput(arr []int32, batchNumber int) {
	log.Printf("Wait sortCombine %d", batchNumber)
	result := make(map[int]Result)
	for range batchNumber {
		data := <-sortCombineSortedInputChannel
		log.Println("Receive sorted data above 10000 ", data.batch, data.result.start, data.result.end, data.result.maxValue)
		result[data.batch] = data.result
	}
	if len(result) == 2 {
		log.Println("Start Combine Last 2 Batch")
		finalResult := combine(arr, result[0], result[1])
		log.Println("Finish Combine Last 2 Batch")
		finishChannel <- finalResult.maxValue
	} else {
		loopCombineChannel(arr, result)
	}
}

func loopCombineChannel(arr []int32, result map[int]Result) {
	totalBatch := len(result)
	log.Println("totalBatch", totalBatch)
	j := totalBatch - 1
	k := j - 1
	i := totalBatch/2 - 1
	if totalBatch%2 > 0 {
		i = i + 1
	}
	for {
		if j <= 0 {
			break
		}
		go combineMaxSubsetSumAbove10000(arr, maxSubsetSumAboveXChannel)
		maxSubsetSumAboveXChannel <- CombineSortedInput{
			result: []Result{result[k], result[j]},
			batch:  i,
		}
		j = k - 1
		k = j - 1
		i--
	}

	if totalBatch%2 > 0 {
		go loopSortCombineSortedInput(arr, totalBatch/2+1)
		sortCombineSortedInputChannel <- SortedInput{
			batch:  i,
			result: result[0],
		}
	} else {
		go loopSortCombineSortedInput(arr, totalBatch/2)
	}
}

const BATCH_NUMBER = 10000

func divide(arr []int32) int32 {
	batchNumber := len(arr) / BATCH_NUMBER
	log.Println(batchNumber)
	for i := range batchNumber {
		startIndex := BATCH_NUMBER * i
		go maxSubsetSumUnder10000(arr, maxSubsetSumUnderXChannel)
		maxSubsetSumUnderXChannel <- Input{
			start: startIndex,
			end:   startIndex + BATCH_NUMBER,
			batch: i,
		}
	}
	remainNumber := len(arr) % BATCH_NUMBER
	go maxSubsetSumUnder10000(arr, maxSubsetSumUnderXChannel)
	maxSubsetSumUnderXChannel <- Input{
		start: BATCH_NUMBER * batchNumber,
		end:   BATCH_NUMBER*batchNumber + remainNumber - 1,
		batch: batchNumber,
	}
	result := make(map[int]Result)
	for range batchNumber + 1 {
		data := <-combineChannel
		result[data.batch] = data.result
	}
	go loopCombineChannel(arr, result)
	maxValue := <-finishChannel
	log.Printf("Final Result %d", maxValue)
	return maxValue
}

func main() {
	file, err := os.Open("../input09.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1024*1024)

	os.Setenv("OUTPUT_PATH", "result.txt")
	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	n := int32(nTemp)
	log.Printf("n: %d\n", n)

	arrTemp := strings.Split(readLine(reader), " ")

	var arr []int32

	for i := 0; i < int(n); i++ {
		arrItemTemp, err := strconv.ParseInt(arrTemp[i], 10, 64)
		checkError(err)
		arrItem := int32(arrItemTemp)
		arr = append(arr, arrItem)
	}
	start := time.Now()
	res := divide(arr)
	elapsed := time.Since(start)
	fmt.Printf("Simulated operation took %d seconds\n", int(elapsed/time.Second))
	writer := bufio.NewWriterSize(stdout, 1024*1024)
	fmt.Fprintf(writer, "%d\n", res)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
