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
const COMPARE = false
const DEBUG = false

func maxSubsetSum(arr []int32) int32 {
	start := time.Now()
	result := divideAndConquer(arr, 0, len(arr)-1)
	elapsed := time.Since(start)
	fmt.Printf("Simulated operation took %d seconds\n", int(elapsed/time.Second))
	log.Println(result.maxValue)
	//printResult(result)
	/* result2 := maxSubsetSumWithoutChannel(arr, 0, len(arr)-1)
	log.Println(result2.maxValue) */
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

type UpdateResult struct {
	maxValueTail int32
	maxListTail  []int
	maxTailKey   int
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
	result1 = updateWhenMerge2Result(arr, result1, result2)
	result1 = reUpdateMaxValue4All(arr, result1, result2)
	//log.Println("==============Start Special Case result2.maxValue > result1.maxCase ===================")
	if result2.maxValue > result1.maxValue {
		result1.maxValue = result2.maxValue
		result1.maxList = result2.maxList
	}
	//log.Println("==============End Special Corner Case===================")
	//log.Println(result1)
	if COMPARE {
		tempResult := maxSubsetSumWithoutChannel(arr, start, end)
		debug(tempResult, result1, result2, arr, start, end, middle)
	}
	if result1.sortedKeyList[len(result1.sortedKeyList)-1] == result2.sortedKeyList[0] {
		result1.sortedKeyList = append(result1.sortedKeyList, result2.sortedKeyList[1:]...)
	} else {
		result1.sortedKeyList = append(result1.sortedKeyList, result2.sortedKeyList...)
	}
	updateTailList := []int{}
	updateAdjacentTailList := []int{}
	for i := len(result1.sortedKeyList) - 1; i >= 0; i-- {
		currentKey := result1.sortedKeyList[i]
		nonAdjTailList := result1.mapNonAdjacentList[currentKey]
		if len(nonAdjTailList) == 0 {
			updateTailList = append(updateTailList, currentKey)
			continue
		}
		if len(nonAdjTailList) > 0 && slices.Contains(updateTailList, nonAdjTailList[0]) {
			updateAdjacentTailList = append(updateAdjacentTailList, currentKey)
			continue
		}
		break
	}
	result1.tailList = updateTailList
	result1.adjacentTailList = updateAdjacentTailList
	return result1
}

func debug(tempResult Result, result1 Result, result2 Result, arr []int32, start int, end int, middle int) {
	if tempResult.maxValue != result1.maxValue {
		log.Println("==============debug===================")
		/* log.Println(arr[start:end], "size=", end-start)
		log.Println(arr[start:middle])
		log.Println(arr[middle-1 : end]) */
		log.Println(result1.maxList, "value=", result1.maxValue)
		log.Println(tempResult.maxList, "value=", tempResult.maxValue)
		log.Println(result2.maxList, "value=", result2.maxValue)
		log.Println(result2.sortedKeyList)
		//printResult(result1)
		for i := 0; i < len(result1.mapNextListMax[7577]); i++ {
			log.Print(result1.mapMaxValue[result1.mapNextListMax[7577][i]],
				tempResult.mapMaxValue[tempResult.mapNextListMax[7577][i]])
		}
		/* log.Println(result1.mapMaxValue[7579], tempResult.mapMaxValue[7579])
		log.Println(result1.mapMaxValue[7580], tempResult.mapMaxValue[7580])
		log.Println(arr[7579], arr[7580])
		log.Println(result1.mapMaxValue[7581])
		log.Println(result1.mapMaxValue[7582]) */

		log.Println(result1.adjacentTailList)
		log.Println(result1.tailList)
		for i := 0; i < len(result1.maxList); i++ {
			if result1.maxList[i] != tempResult.maxList[i] {
				log.Printf("Diff result1 %d, temp %d", result1.maxList[i], tempResult.maxList[i])
				log.Println(result1.maxList[i], result1.mapNextListMax[result1.maxList[i]],
					result1.mapMaxValue[result1.maxList[i]])
				log.Println(tempResult.maxList[i], result1.mapNextListMax[tempResult.maxList[i]],
					result1.mapMaxValue[tempResult.maxList[i]])
				log.Println(tempResult.maxList[i], tempResult.mapNextListMax[tempResult.maxList[i]],
					result1.mapMaxValue[tempResult.maxList[i]])
				if i-1 >= 0 {
					log.Printf("Parent %d", result1.maxList[i-1])
					log.Println(result1.mapNonAdjacentList[result1.maxList[i-1]])
				}
				break
			}
		}
		//printResult(result1)
		//printResult(result2)
		log.Println("==============end debug===================")
		log.Fatal("Error")
	}
}

func updateWhenMerge2Result(arr []int32, result1 Result, result2 Result) Result {
	tailList1 := result1.tailList
	sortedKeyList2 := result2.sortedKeyList
	maxList2 := result2.maxList
	connectSortedKeyList2 := createConnectSortedKeyList2(sortedKeyList2)
	for i := len(connectSortedKeyList2) - 1; i >= 0; i-- {
		key := connectSortedKeyList2[i]
		if slices.Contains(maxList2, key) || slices.Contains(result2.tailList, key) ||
			slices.Contains(result2.adjacentTailList, key) {
			continue
		}
		result2 = updateOtherHead(arr, result2, key, 0)
		if DEBUG {
			if !checkUnique(result2.mapNextListMax[key]) {
				log.Fatal("Error Unique in updateOtherHead of updateWhenMerge2Result")
			}
			if !checkInvalid(result2.mapNextListMax[key]) {
				log.Fatal("Error Invalid in updateOtherHead of updateWhenMerge2Result")
			}
		}
	}
	for _, key := range sortedKeyList2 {
		if !slices.Contains(result1.sortedKeyList, key) {
			result1.mapNonAdjacentList[key] = result2.mapNonAdjacentList[key]
			result1.mapMaxValue[key] = result2.mapMaxValue[key]
			result1.mapNextListMax[key] = result2.mapNextListMax[key]
		}
	}

	for _, key := range result1.tailList {
		for _, key2 := range connectSortedKeyList2 {
			if key+1 < key2 && !slices.Contains(tailList1, key2) &&
				!slices.Contains(result1.mapNonAdjacentList[key], key2) {
				result1.mapNonAdjacentList[key] = append(result1.mapNonAdjacentList[key], key2)
			}
		}
		if DEBUG {
			if len(result1.mapNonAdjacentList[key]) == 0 && len(result2.sortedKeyList) >= 4 {
				log.Fatal("Error mapNonAdjacentList of the tail key is empty ", key)
			}
		}
	}
	for _, key := range result1.adjacentTailList {
		for _, key2 := range connectSortedKeyList2 {
			if key+1 < key2 && !slices.Contains(tailList1, key2) &&
				!slices.Contains(result1.mapNonAdjacentList[key], key2) {
				result1.mapNonAdjacentList[key] = append(result1.mapNonAdjacentList[key], key2)
			}
		}
	}
	return result1
}

func createConnectSortedKeyList2(sortedKeyList2 []int) []int {
	if len(sortedKeyList2) < 4 {
		return sortedKeyList2
	}
	return sortedKeyList2[0:4]
}

func reUpdateMaxValue4All(arr []int32, result1 Result, result2 Result) Result {
	result1, updateResult := reUpdateMaxValueForTailOfResult1(arr, result1)
	sortKeyList := result1.sortedKeyList
	if len(sortKeyList) == len(result1.tailList) {
		result1.maxValue = updateResult.maxValueTail
		result1.maxList = updateResult.maxListTail
		return result1
	}
	result1, updateResult = reUpdateMaxValue4AdjacentTail(arr, result1, updateResult)
	if len(sortKeyList) == len(result1.tailList)+len(result1.adjacentTailList) {
		result1.maxValue = updateResult.maxValueTail
		result1.maxList = updateResult.maxListTail
		return result1
	}
	result1 = reUpdateHeadList(arr, result1, updateResult)
	if DEBUG {
		if !slices.Contains(result2.tailList, result1.maxList[len(result1.maxList)-1]) &&
			len(result2.sortedKeyList) >= 3 {
			log.Println(result1.maxList)
			log.Println(result2.maxList)
			log.Fatal("Error the tailf of maxList does not conntain tailList when joining 2 result")
		}
		if !checkUnique(result1.maxList) {
			log.Println(result1.maxList)
			log.Println(result2.maxList)
			log.Fatal("Error Unique in reUpdateHeadList")
		}
		if !checkInvalid(result1.maxList) {
			log.Println(result1.maxList)
			log.Println(result2.maxList)
			log.Fatal("Error Invalid in reUpdateHeadList")
		}
	}
	return result1
}

func reUpdateMaxValueForTailOfResult1(arr []int32, result1 Result) (Result, UpdateResult) {
	tailList := result1.tailList
	updateResult := UpdateResult{}
	for _, key := range tailList {
		result1 = calcKeyFromNonAdjacentList(arr, result1, key)
		if result1.mapMaxValue[key] > updateResult.maxValueTail {
			updateResult.maxValueTail = result1.mapMaxValue[key]
			updateResult.maxListTail = result1.mapNextListMax[key]
			updateResult.maxTailKey = key
		}
	}
	return result1, updateResult
}

func reUpdateMaxValue4AdjacentTail(arr []int32, result1 Result,
	updateResult UpdateResult) (Result, UpdateResult) {
	adjacentTailList := result1.adjacentTailList
	for _, key := range adjacentTailList {
		result1 = calcKeyFromNonAdjacentList(arr, result1, key)
		if result1.mapMaxValue[key] > updateResult.maxValueTail {
			updateResult.maxValueTail = result1.mapMaxValue[key]
			updateResult.maxListTail = result1.mapNextListMax[key]
			updateResult.maxTailKey = key
		}
	}
	return result1, updateResult
}

func reUpdateHeadList(arr []int32, result1 Result, updateResult UpdateResult) Result {
	tailKeyList := result1.tailList
	adjacentTailList := result1.adjacentTailList
	maxTailKey := updateResult.maxTailKey
	currentMaxList := result1.maxList
	var oldMaxTailKey int
	var oldMaxTailKeyIndex int

	if slices.Contains(tailKeyList, maxTailKey) {
		for oldMaxTailKeyIndex = len(currentMaxList) - 1; oldMaxTailKeyIndex >= 0; oldMaxTailKeyIndex-- {
			oldMaxTailKey = currentMaxList[oldMaxTailKeyIndex]
			if oldMaxTailKey == maxTailKey {
				break
			}
			if !slices.Contains(tailKeyList, oldMaxTailKey) {
				oldMaxTailKeyIndex++
				oldMaxTailKey = currentMaxList[oldMaxTailKeyIndex]
				break
			}
		}
	}
	if slices.Contains(adjacentTailList, maxTailKey) {
		for oldMaxTailKeyIndex = len(currentMaxList) - 1; oldMaxTailKeyIndex >= 0; oldMaxTailKeyIndex-- {
			oldMaxTailKey = currentMaxList[oldMaxTailKeyIndex]
			if oldMaxTailKey == maxTailKey {
				break
			}
			if !slices.Contains(tailKeyList, oldMaxTailKey) && !slices.Contains(adjacentTailList, oldMaxTailKey) {
				oldMaxTailKeyIndex++
				oldMaxTailKey = currentMaxList[oldMaxTailKeyIndex]
				break
			}
		}
	}

	if oldMaxTailKey == maxTailKey {
		return reCalculateMaxList(arr, result1, oldMaxTailKeyIndex-1)
	}

	oldMaxTailKeyIndex--
	if oldMaxTailKeyIndex < 0 { //no parent of oldMaxTailKey
		maxKey := maxTailKey
		parentMaxTailKey := findParentElement(result1.sortedKeyList, maxTailKey)
		if parentMaxTailKey > 0 {
			result1.mapMaxValue[parentMaxTailKey] = arr[parentMaxTailKey] + result1.mapMaxValue[maxTailKey]
			result1.mapNextListMax[parentMaxTailKey] = append([]int{parentMaxTailKey}, result1.mapNextListMax[maxTailKey]...)
			maxKey = parentMaxTailKey
		}
		result1.maxValue = result1.mapMaxValue[maxKey]
		result1.maxList = result1.mapNextListMax[maxKey]
		return result1
	}
	prevOldMaxTailKey := currentMaxList[oldMaxTailKeyIndex]

	if prevOldMaxTailKey+1 < maxTailKey {
		result1.mapMaxValue[prevOldMaxTailKey] = arr[prevOldMaxTailKey] + result1.mapMaxValue[maxTailKey]
		result1.mapNextListMax[prevOldMaxTailKey] = append([]int{prevOldMaxTailKey}, result1.mapNextListMax[maxTailKey]...)
		return reCalculateMaxList(arr, result1, oldMaxTailKeyIndex-1)
	}
	maxKey := maxTailKey
	result1 = calcKeyFromNonAdjacentList(arr, result1, prevOldMaxTailKey)
	if result1.mapMaxValue[prevOldMaxTailKey] > result1.mapMaxValue[maxKey] {
		maxKey = prevOldMaxTailKey
	}
	adjacentNewMaxTailKey := findParentElement(result1.sortedKeyList, maxTailKey)
	if adjacentNewMaxTailKey > 0 {
		result1.mapMaxValue[adjacentNewMaxTailKey] = arr[adjacentNewMaxTailKey] + result1.mapMaxValue[maxTailKey]
		result1.mapNextListMax[adjacentNewMaxTailKey] = append([]int{adjacentNewMaxTailKey}, result1.mapNextListMax[maxTailKey]...)
		if result1.mapMaxValue[adjacentNewMaxTailKey] > result1.mapMaxValue[maxKey] {
			maxKey = adjacentNewMaxTailKey
		}
	}
	for {
		oldMaxTailKeyIndex--
		if oldMaxTailKeyIndex < 0 {
			result1.maxValue = result1.mapMaxValue[maxKey]
			result1.maxList = result1.mapNextListMax[maxKey]
			return result1
		}
		selectedParent := currentMaxList[oldMaxTailKeyIndex]
		result1 = calcKeyFromNonAdjacentList(arr, result1, selectedParent)
		if result1.mapMaxValue[selectedParent] > result1.mapMaxValue[maxKey] {
			//keep ((0...1)th parent, selectedParent, prevOldMaxTailKey) > ((0...1)th parent, nearest parentKey, maxTailKey)
			// OR ((0...1)th parent, nearest parentKey, maxTailKey)
			maxKey = selectedParent
			break
		}
		newSelectedParent := findParentElement(result1.sortedKeyList, maxKey)
		if newSelectedParent > 0 {
			result1.mapMaxValue[newSelectedParent] = arr[newSelectedParent] + result1.mapMaxValue[maxKey]
			result1.mapNextListMax[newSelectedParent] = append([]int{newSelectedParent}, result1.mapNextListMax[maxKey]...)
			maxKey = newSelectedParent
		}
	}
	if oldMaxTailKeyIndex-1 >= 0 {
		result1 = reCalculateMaxList(arr, result1, oldMaxTailKeyIndex-1)
		maxKey = currentMaxList[0]
	}
	result1.maxValue = result1.mapMaxValue[maxKey]
	result1.maxList = result1.mapNextListMax[maxKey]
	return result1
}

func reCalculateMaxList(arr []int32, result1 Result, startIndex int) Result {
	currentMaxList := result1.maxList
	for j := startIndex; j >= 0; j-- {
		currentKey := currentMaxList[j]
		afterKey := currentMaxList[j+1]
		result1.mapMaxValue[currentKey] = arr[currentKey] + result1.mapMaxValue[afterKey]
		result1.mapNextListMax[currentKey] = append([]int{currentKey}, result1.mapNextListMax[afterKey]...)
	}
	result1.maxValue = result1.mapMaxValue[currentMaxList[0]]
	result1.maxList = result1.mapNextListMax[currentMaxList[0]]
	return result1
}

func findIndexElement(sortedKeyList []int, searchElm int) int {
	for i := len(sortedKeyList) - 1; i >= 0; i-- {
		if sortedKeyList[i] == searchElm {
			return i
		}
	}
	return -1
}

func findParentElement(sortedKeyList []int, searchElm int) int {
	searchIndex := findIndexElement(sortedKeyList, searchElm)
	for i := searchIndex - 1; i >= 0; i-- {
		candidateKey := sortedKeyList[i]
		if candidateKey+2 <= searchElm {
			return candidateKey
		}
	}
	return -1
}

func updateOtherHead(arr []int32, result Result, otherHeadNeedUpdate int, startIndex int) Result {
	if len(result.mapNonAdjacentList[otherHeadNeedUpdate]) == 0 {
		return result
	}
	if slices.Contains(result.tailList, otherHeadNeedUpdate+2) ||
		slices.Contains(result.adjacentTailList, otherHeadNeedUpdate+2) {
		result = calcKeyFromNonAdjacentList(arr, result, otherHeadNeedUpdate)
		return result
	}
	maxList := result.maxList
	nonAdjKeyList := []int{}
	excludeNonAdjCalc := []int{}
	for _, key := range result.mapNonAdjacentList[otherHeadNeedUpdate] {
		if len(excludeNonAdjCalc) == 0 || excludeNonAdjCalc[0] > key+2 {
			excludeNonAdjCalc = []int{key + 2}
		}
		if len(excludeNonAdjCalc) > 0 && key >= excludeNonAdjCalc[0] {
			break
		}
		if slices.Contains(maxList, key) {
			continue
		}
		nonAdjKeyList = append(nonAdjKeyList, key)
	}
	var leadMaxPathKey int
	i := startIndex
	for {
		if i == len(maxList) {
			break
		}
		leadMaxPathKey = maxList[i]
		if otherHeadNeedUpdate+2 <= leadMaxPathKey {
			break
		}
		i++
	}
	if otherHeadNeedUpdate+2 <= leadMaxPathKey && len(nonAdjKeyList) == 0 {
		result.mapMaxValue[otherHeadNeedUpdate] = arr[otherHeadNeedUpdate] + result.mapMaxValue[leadMaxPathKey]
		result.mapNextListMax[otherHeadNeedUpdate] = append([]int{otherHeadNeedUpdate}, result.mapNextListMax[leadMaxPathKey]...)
		return result
	}
	if i == len(maxList) || slices.Contains(result.tailList, leadMaxPathKey) { //noway to lead max tail
		//otherHeadNeedUpdate contains (adjanceList, tail list)
		result = calcKeyFromNonAdjacentList(arr, result, otherHeadNeedUpdate)
		return result
	}
	nonAdjKeyListNeedUpdate := []int{}
	for _, key := range nonAdjKeyList {
		if slices.Contains(result.tailList, key) ||
			slices.Contains(result.adjacentTailList, key) {
			continue
		}
		nonAdjKeyListNeedUpdate = append(nonAdjKeyListNeedUpdate, key)
	}
	for j := len(nonAdjKeyListNeedUpdate) - 1; j >= 0; j-- {
		result = updateOtherHead(arr, result, nonAdjKeyListNeedUpdate[j], i)
		if DEBUG {
			if !checkUnique(result.mapNextListMax[nonAdjKeyListNeedUpdate[j]]) {
				log.Println(result.mapNextListMax[nonAdjKeyListNeedUpdate[j]])
				log.Fatal("Error Unique in updateOtherHead")
			}
			if !checkInvalid(result.mapNextListMax[nonAdjKeyListNeedUpdate[j]]) {
				log.Println(result.mapNextListMax[nonAdjKeyListNeedUpdate[j]])
				log.Fatal("Error Invalid in updateOtherHead")
			}
		}
	}
	result = calcKeyFromNonAdjacentList(arr, result, otherHeadNeedUpdate)
	return result
}

func checkUnique(keyList []int) bool {
	i := 0
	j := 1
	for {
		if i == len(keyList) || j == len(keyList) {
			return true
		}
		if keyList[i] == keyList[j] {
			return false
		}
		i++
		j++
	}
}

func checkInvalid(keyList []int) bool {
	i := 0
	j := 1
	for {
		if i == len(keyList) || j == len(keyList) {
			return true
		}
		if keyList[i]+1 == keyList[j] {
			return false
		}
		i++
		j++
	}
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
		mapNextListMax:     make(map[int][]int),
		mapMaxValue:        make(map[int]int32),
		mapNonAdjacentList: mapNonAdjacentList,
		sortedKeyList:      sortedKeyList,
	}
	adjacentTailList := []int{}
	for i := len(sortedKeyList) - 1; i >= 0; i-- {
		key := sortedKeyList[i]
		if len(mapNonAdjacentList[key]) == 0 {
			tailList = append(tailList, key)
		}
		if len(mapNonAdjacentList[key]) == 1 || len(mapNonAdjacentList[key]) == len(tailList) {
			adjacentTailList = append(adjacentTailList, key)
		}
		calcResult = calcKeyFromNonAdjacentList(arr, calcResult, key)
		value, tempMaxList := calcResult.mapMaxValue[key], calcResult.mapNextListMax[key]

		if value > maxValue {
			maxValue = value
			maxList = tempMaxList
			calcResult.maxList = maxList
		}
	}
	calcResult.maxValue = maxValue
	calcResult.maxList = maxList
	calcResult.tailList = tailList
	calcResult.adjacentTailList = adjacentTailList
	return calcResult
}

func calcKeyFromNonAdjacentList(arr []int32, calcResult Result, key int) Result {
	var maxValue int32
	var maxKey int
	nonAdjacentList := calcResult.mapNonAdjacentList[key]
	if len(nonAdjacentList) == 0 {
		calcResult.mapNextListMax[key] = []int{key}
		calcResult.mapMaxValue[key] = arr[key]
		return calcResult
	}
	upperExcludeCalc := []int{}
	for _, key2 := range nonAdjacentList {
		if len(upperExcludeCalc) > 0 && key2 >= upperExcludeCalc[0] {
			break
		}
		if arr[key]+calcResult.mapMaxValue[key2] > maxValue {
			maxValue = arr[key] + calcResult.mapMaxValue[key2]
			maxKey = key2
		}
		if len(upperExcludeCalc) == 0 || upperExcludeCalc[0] > key2+2 {
			upperExcludeCalc = []int{key2 + 2}
		}
	}
	calcResult.mapNextListMax[key] = append([]int{key}, calcResult.mapNextListMax[maxKey]...)
	calcResult.mapMaxValue[key] = maxValue
	return calcResult
}

func main() {
	file, err := os.Open("input07.txt")
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
