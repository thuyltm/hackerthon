package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Complete the maxSubsetSum function below.

func maxSubsetSum(arr []int32) int32 {
	start := time.Now()
	result := divideAndConquer(arr, 0, len(arr)-1)
	elapsed := time.Since(start)
	fmt.Printf("Simulated operation took %d seconds\n", int(elapsed/time.Second))
	log.Println(result.maxValue)
	//log.Println(result.maxList)
	//printResult(result)
	//result2 := maxSubsetSumWithoutChannel(arr, 0, len(arr)-1)
	//log.Println(result2.maxValue)
	//log.Println(result.maxList)
	//printResult(result2)
	return result.maxValue
}

func printResult(result Result) {
	log.Println("maxValue", result.maxValue)
	log.Println("maxList", result.maxList)
	log.Println("mapMaxValue", result.mapMaxValue)
	log.Println("mapMaxList", result.mapNextListMax)
	log.Println("mapNonAdjacentList", result.mapNonAdjacentList)
	log.Println("tailList", result.tailList)
	log.Println("adjacentTailList", result.adjacentTailList)
	log.Println("sortedKeyList", result.sortedKeyList)
}

type Result struct {
	mapNextListMax     map[int][]int
	mapMaxValue        map[int]int32
	maxList            []int
	maxValue           int32
	mapNonAdjacentList map[int][]int
	tailList           []int
	sortedKeyList      []int
	adjacentTailList   []int
}

func divideAndConquer(arr []int32, start int, end int) Result {
	if end-start <= 100 {
		return maxSubsetSumWithoutChannel(arr, start, end)
	}
	middle := start + (end-start)/2
	result1 := divideAndConquer(arr, start, middle)
	result2 := divideAndConquer(arr, middle, end)
	if len(result1.maxList) == 0 && len(result2.maxList) == 0 {
		return Result{}
	}
	if len(result1.maxList) == 0 && len(result2.maxList) > 0 {
		return result2
	}
	if len(result1.maxList) > 0 && len(result2.maxList) == 0 {
		return result1
	}
	result1 = updateWhenMerge2Result(result1, result2)
	result1 = reUpdateMaxValueForTail(result1, arr)
	result1 = reUpdateMaxValue4AdjacentTail(arr, result1)
	result1 = reUpdateHeadList(result1, arr)
	//log.Println("==============Start Special Case result2.maxValue > result1.maxCase ===================")
	if result2.maxValue > result1.maxValue {
		result1.maxValue = result2.maxValue
		result1.maxList = result2.maxList
	}
	//log.Println("==============End Special Corner Case===================")
	//log.Println(result1)
	//tempResult := maxSubsetSumWithoutChannel(arr, start, end)
	//debug(tempResult, result1, result2, arr, start, end, middle)
	if result1.sortedKeyList[len(result1.sortedKeyList)-1] == result2.sortedKeyList[0] {
		result1.sortedKeyList = append(result1.sortedKeyList, result2.sortedKeyList[1:]...)
	} else {
		result1.sortedKeyList = append(result1.sortedKeyList, result2.sortedKeyList...)
	}
	result1.tailList = result2.tailList
	result1.adjacentTailList = result2.adjacentTailList
	return result1
}

func debug(tempResult Result, result1 Result, result2 Result, arr []int32, start int, end int, middle int) {
	if tempResult.maxValue != result1.maxValue && end-start < 40 {
		log.Println("==============debug===================")
		log.Println(arr[start:end], "size=", end-start)
		log.Println(arr[start:middle])
		log.Println(arr[middle-1 : end])
		log.Println(result1.maxList, "value=", result1.maxValue)
		log.Println(tempResult.maxList, "value=", tempResult.maxValue)
		printResult(result1)
		printResult(result2)
		log.Println("==============end debug===================")
	}
}

func updateWhenMerge2Result(result1 Result, result2 Result) Result {
	tailList1 := result1.tailList
	sortedKeyList2 := result2.sortedKeyList
	maxList2 := result2.maxList
	loopSortedKeyList2 := []int{maxList2[0]}
	if len(maxList2) >= 2 {
		comparableKey2 := maxList2[1]
		i := 0
		for {
			if i == len(sortedKeyList2) {
				break
			}
			key := sortedKeyList2[i]
			if key == maxList2[0] {
				i++
				continue
			}
			if key > comparableKey2 {
				break
			}
			loopSortedKeyList2 = append(loopSortedKeyList2, key)
			i++
		}
	} else {
		//log.Println("================Start Special Case len(maxList)=1=========================")
		loopSortedKeyList2 = make([]int, len(sortedKeyList2))
		copy(loopSortedKeyList2, sortedKeyList2)
		//log.Println("================End Special Case=========================================")
	}
	for _, key := range sortedKeyList2 {
		result1.mapNonAdjacentList[key] = result2.mapNonAdjacentList[key]
		result1.mapMaxValue[key] = result2.mapMaxValue[key]
		result1.mapNextListMax[key] = result2.mapNextListMax[key]
	}
	//Result1 must connect to the result2 from head to tail. FORCE!!!!!!!!!!!!!!
	for _, key := range result1.sortedKeyList {
		for _, key2 := range loopSortedKeyList2 {
			if key+1 < key2 && !slices.Contains(tailList1, key2) {
				result1.mapNonAdjacentList[key] = append(result1.mapNonAdjacentList[key], key2)
			}
		}
	}
	return result1
}

func reUpdateMaxValueForTail(result1 Result, arr []int32) Result {
	tailList := result1.tailList
	var maxValue int32
	var maxList []int
	for _, key := range tailList {
		nonAdjList := result1.mapNonAdjacentList[key]
		if len(nonAdjList) == 0 {
			continue
		}
		var maxValueKey int32
		var maxAdjKey int
		for _, nonAdjKey := range nonAdjList {
			newValue := arr[key] + result1.mapMaxValue[nonAdjKey]
			if newValue > maxValueKey {
				maxValueKey = newValue
				maxAdjKey = nonAdjKey
			}
		}
		result1.mapMaxValue[key] = maxValueKey
		result1.mapNextListMax[key] = append([]int{key}, result1.mapNextListMax[maxAdjKey]...)
		newMaxValue := result1.mapMaxValue[key]
		if newMaxValue > maxValue {
			maxValue = newMaxValue
			maxList = result1.mapNextListMax[key]
		}
	}
	result1.maxValue = maxValue
	result1.maxList = maxList
	return result1
}

func reUpdateMaxValue4AdjacentTail(arr []int32, result1 Result) Result {
	adjacentTailList := result1.adjacentTailList
	maxValueTail := result1.maxValue
	maxListTail := result1.maxList
	for _, key := range adjacentTailList {
		nonAdjacentKeyList := result1.mapNonAdjacentList[key]
		var newMaxValue int32
		var maxKey int
		for _, nonAdjacentKey := range nonAdjacentKeyList {
			currentValue := arr[key] + result1.mapMaxValue[nonAdjacentKey]
			if currentValue > newMaxValue {
				newMaxValue = currentValue
				maxKey = nonAdjacentKey
			}
		}
		result1.mapMaxValue[key] = newMaxValue
		result1.mapNextListMax[key] = append([]int{key}, result1.mapNextListMax[maxKey]...)
		if newMaxValue > maxValueTail {
			maxValueTail = newMaxValue
			maxListTail = result1.mapNextListMax[key]
		}
	}
	result1.maxValue = maxValueTail
	result1.maxList = maxListTail
	return result1
}

func reUpdateHeadList(result1 Result, arr []int32) Result {
	sortedKeyList := result1.sortedKeyList
	tailKeyList := result1.tailList
	adjacentTailList := result1.adjacentTailList
	maxValueTail := result1.maxValue
	maxListTail := result1.maxList
	maxValue := maxValueTail
	maxList := maxListTail
	for i := len(sortedKeyList) - len(tailKeyList) - len(adjacentTailList) - 1; i >= 0; i-- {
		key := sortedKeyList[i]
		nonAdjList := result1.mapNonAdjacentList[key]
		var maxValueKey int32
		var maxKey int
		for _, nonAdjKey := range nonAdjList {
			currentValue := arr[key] + result1.mapMaxValue[nonAdjKey]
			if currentValue > maxValueKey {
				maxValueKey = currentValue
				maxKey = nonAdjKey
			}
			if len(maxListTail) >= 1 && nonAdjKey == maxListTail[0] {
				break
			}
		}
		if len(maxListTail) >= 1 &&
			maxValueKey < arr[key]+maxValueTail && key+1 < maxListTail[0] {
			maxKey = maxListTail[0]
			maxValueKey = arr[key] + maxValueTail
			//NonAdjacentList must include more the maxKey. FORCE!!!!!!!!!!!!
			result1.mapNonAdjacentList[key] = append([]int{maxKey}, result1.mapNonAdjacentList[key]...)
		}
		result1.mapMaxValue[key] = maxValueKey
		result1.mapNextListMax[key] = append([]int{key}, result1.mapNextListMax[maxKey]...)
		if maxValueKey > maxValue {
			maxValue = maxValueKey
			maxList = result1.mapNextListMax[key]
		}
	}
	result1.maxValue = maxValue
	result1.maxList = maxList
	return result1
}

func maxSubsetSumWithoutChannel(arr []int32, left int, right int) Result {
	sortedKeyList := []int{}
	mapNonAdjacentList, sortedKeyList := createMapNonAdjacentList(arr, left, right, sortedKeyList)
	return calcAndReturnListMaxSum(mapNonAdjacentList, arr, sortedKeyList)
}

func createMapNonAdjacentList(arr []int32, left int, right int,
	sortedKeyList []int) (map[int][]int, []int) {
	mapNonAdjacentList := make(map[int][]int)
	prefix := arr[left]
	if prefix > 0 {
		sortedKeyList = append(sortedKeyList, left)
	}
	if left > right-2 { //stop condition
		return mapNonAdjacentList, sortedKeyList
	}
	if prefix <= 0 && right-left > 2 {
		return createMapNonAdjacentList(arr, left+1, right, sortedKeyList)
	}
	if right-left == 2 {
		if arr[left+1] > 0 {
			sortedKeyList = append(sortedKeyList, left+1)
		}
		if arr[left+2] > 0 {
			sortedKeyList = append(sortedKeyList, left+2)
		}
		if arr[left] > 0 && arr[right] > 0 {
			mapNonAdjacentList[left] = append(mapNonAdjacentList[left], right)
		}
		if arr[left] > 0 && arr[right] < 0 {
			mapNonAdjacentList[left] = []int{}
		}
		if arr[left+1] > 0 {
			mapNonAdjacentList[left+1] = []int{}
		}
		if arr[left+2] > 0 {
			mapNonAdjacentList[left+2] = []int{}
		}
		return mapNonAdjacentList, sortedKeyList
	}
	mapNonAdjacentList, sortedKeyList = createMapNonAdjacentList(arr, left+1, right, sortedKeyList)
	for i := left + 2; i <= right; i++ {
		if arr[i] > 0 {
			mapNonAdjacentList[int(left)] = append(mapNonAdjacentList[int(left)], int(i))
		}
	}
	if len(mapNonAdjacentList[int(left)]) == 0 && arr[left] > 0 {
		mapNonAdjacentList[int(left)] = []int{}
	}
	return mapNonAdjacentList, sortedKeyList
}

func calcAndReturnListMaxSum(mapNonAdjacentList map[int][]int, arr []int32, sortedKeyList []int) Result {
	maxList := []int{}
	tailList := []int{}
	var maxValue int32
	calcResult := Result{
		mapNextListMax: make(map[int][]int),
		mapMaxValue:    make(map[int]int32),
	}
	adjacentTailList := []int{}
	for i := len(sortedKeyList) - 1; i >= 0; i-- {
		key := sortedKeyList[i]
		if len(mapNonAdjacentList[key]) == 0 {
			tailList = append(tailList, key)
		}
		if len(mapNonAdjacentList[key]) == 1 || len(mapNonAdjacentList[key]) == 2 {
			adjacentTailList = append(adjacentTailList, key)
		}
		value, tempMaxList := getSubsetListByKey(mapNonAdjacentList, key, arr, calcResult)
		if value > maxValue {
			maxValue = value
			maxList = tempMaxList
		}
	}
	calcResult.maxValue = maxValue
	calcResult.maxList = maxList
	calcResult.mapNonAdjacentList = mapNonAdjacentList
	calcResult.tailList = tailList
	calcResult.sortedKeyList = sortedKeyList
	calcResult.adjacentTailList = adjacentTailList
	return calcResult
}

func getSubsetListByKey(mapNonAdjacentList map[int][]int, key int,
	arr []int32, calcResult Result) (int32, []int) {
	value, exists := calcResult.mapMaxValue[key]
	if exists {
		return value, calcResult.mapNextListMax[key]
	}
	adjancentList := mapNonAdjacentList[key]
	if len(adjancentList) == 0 && arr[key] > 0 {
		calcResult.mapMaxValue[key] = arr[key]
		calcResult.mapNextListMax[key] = []int{key}
		return arr[key], []int{key}
	}
	var tempMaxList []int
	var temp int32
	var result int32
	var addedValue int32
	for _, value := range mapNonAdjacentList[key] {
		addedValue, tempMaxList = getSubsetListByKey(mapNonAdjacentList, value, arr, calcResult)
		if tempMaxList[0] != key {
			addedValue = arr[key] + addedValue
			tempMaxList2 := append([]int{key}, tempMaxList...)
			tempMaxList = tempMaxList2
		}
		temp = addedValue
		if temp > result && temp > 0 {
			result = temp
			calcResult.mapNextListMax[key] = tempMaxList
			calcResult.mapMaxValue[key] = result
		}
	}
	return calcResult.mapMaxValue[key], calcResult.mapNextListMax[key]
}

func main() {
	file, err := os.Open("input09.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1024*1024)

	os.Setenv("OUTPUT_PATH", "result.txt")
	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	n := int32(nTemp)

	arrTemp := strings.Split(readLine(reader), " ")

	var arr []int32

	for i := 0; i < int(n); i++ {
		arrItemTemp, err := strconv.ParseInt(arrTemp[i], 10, 64)
		checkError(err)
		arrItem := int32(arrItemTemp)
		arr = append(arr, arrItem)
	}

	res := maxSubsetSum(arr)

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
