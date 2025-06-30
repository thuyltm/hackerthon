package main

import (
	"log"
	"slices"
)

func maxSubsetSum(arr []int32) int32 {
	result := divideAndConquer(arr, 0, len(arr)-1)
	log.Println(result.maxValue)
	//printResult(result)
	//result2 := maxSubsetSumWithoutChannel(arr, 0, len(arr)-1)
	//log.Println(result2.maxValue)
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
}
func printValue(arr []int32, index []int) []int32 {
	var result []int32
	for i := 0; i < len(index); i++ {
		result = append(result, arr[index[i]])
	}
	return result
}

func divideAndConquer(arr []int32, start int, end int) Result {
	if end-start <= 100 {
		return maxSubsetSumWithoutChannel(arr, start, end)
	}
	middle := start + (end-start)/2
	result1 := divideAndConquer(arr, start, middle)
	result2 := divideAndConquer(arr, middle, end)
	result1 = combine(arr, result1, result2)
	result1.start = start
	result1.end = end
	return result1
}

func combine(arr []int32, result1 Result, result2 Result) Result {
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
	result1 = reUpdateMaxValue4All(result1, arr)
	//log.Println("==============Start Special Case result2.maxValue > result1.maxCase ===================")
	if result2.maxValue > result1.maxValue {
		result1.maxValue = result2.maxValue
		result1.maxList = result2.maxList
	}
	//log.Println("==============End Special Corner Case===================")
	//log.Println(result1)
	if result1.sortedKeyList[len(result1.sortedKeyList)-1] == result2.sortedKeyList[0] {
		result1.sortedKeyList = append(result1.sortedKeyList, result2.sortedKeyList[1:]...)
	} else {
		result1.sortedKeyList = append(result1.sortedKeyList, result2.sortedKeyList...)
	}
	result1.tailList = result2.tailList
	result1.adjacentTailList = result2.adjacentTailList
	result1.end = result2.end
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
		log.Println(result1.mapNonAdjacentList[4324])
		log.Println(result1.mapMaxValue[4327])
		log.Println(result1.mapMaxValue[4326])
		log.Println(result1.adjacentTailList)
		log.Println(result1.tailList)
		//printResult(result1)
		//printResult(result2)
		log.Println("==============end debug===================")
		log.Fatal("Error")
	}
}

func updateWhenMerge2Result(result1 Result, result2 Result) Result {
	tailList1 := result1.tailList
	sortedKeyList2 := result2.sortedKeyList
	maxList2 := result2.maxList
	loopSortedKeyList2 := []int{maxList2[0]}
	if len(maxList2) >= 2 {
		comparableKey2 := maxList2[1]
		comparableValue2 := result2.mapMaxValue[comparableKey2]
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
			if result2.mapMaxValue[key] >= comparableValue2 {
				loopSortedKeyList2 = append([]int{maxList2[0]}, key)
				comparableValue2 = result2.mapMaxValue[key]
			}
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
	for _, key := range result1.tailList {
		for _, key2 := range loopSortedKeyList2 {
			if key+1 < key2 && !slices.Contains(tailList1, key2) {
				result1.mapNonAdjacentList[key] = append(result1.mapNonAdjacentList[key], key2)
			}
		}
	}
	for _, key := range result1.adjacentTailList {
		for _, key2 := range loopSortedKeyList2 {
			if key+1 < key2 && !slices.Contains(tailList1, key2) {
				result1.mapNonAdjacentList[key] = append(result1.mapNonAdjacentList[key], key2)
			}
		}
	}
	return result1
}

func reUpdateMaxValue4All(result1 Result, arr []int32) Result {
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
	return reUpdateHeadList(arr, result1, updateResult)
}

func reUpdateMaxValueForTailOfResult1(arr []int32, result1 Result) (Result, UpdateResult) {
	tailList := result1.tailList
	updateResult := UpdateResult{}
	for _, key := range tailList {
		var maxValueKey int32
		nonAdjList := result1.mapNonAdjacentList[key]
		if len(nonAdjList) == 0 {
			maxValueKey = result1.mapMaxValue[key]
		} else {
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
		}
		if maxValueKey > updateResult.maxValueTail {
			updateResult.maxValueTail = maxValueKey
			updateResult.maxListTail = result1.mapNextListMax[key]
			updateResult.maxTailKey = key
		}
	}
	return result1, updateResult
}

func reUpdateMaxValue4AdjacentTail(arr []int32, result1 Result,
	updateResult UpdateResult) (Result, UpdateResult) {
	adjacentTailList := result1.adjacentTailList
	maxTailKey := updateResult.maxTailKey
	tailList := result1.tailList
	for _, key := range adjacentTailList {
		nonAdjacentKeyList := result1.mapNonAdjacentList[key]
		//AdjacentTail contains more key from result2
		isContailAllTailKeyList := true
		for _, nonAdjacentKey := range nonAdjacentKeyList {
			if !slices.Contains(tailList, nonAdjacentKey) {
				isContailAllTailKeyList = false
			}
		}
		if isContailAllTailKeyList &&
			(slices.Contains(nonAdjacentKeyList, maxTailKey) || key+1 < maxTailKey) {
			result1.mapMaxValue[key] = arr[key] + result1.mapMaxValue[maxTailKey]
			result1.mapNextListMax[key] = append([]int{key}, result1.mapNextListMax[maxTailKey]...)
			if !slices.Contains(nonAdjacentKeyList, maxTailKey) {
				result1.mapNonAdjacentList[key] = append([]int{maxTailKey}, result1.mapNonAdjacentList[key]...)
			}
			if result1.mapMaxValue[key] > updateResult.maxValueTail {
				updateResult.maxValueTail = result1.mapMaxValue[key]
				updateResult.maxListTail = result1.mapNextListMax[key]
				updateResult.maxTailKey = key
			}
			continue
		}
		var newMaxValue int32
		var newMaxKey int
		for _, nonAdjacentKey := range nonAdjacentKeyList {
			currentValue := arr[key] + result1.mapMaxValue[nonAdjacentKey]
			if currentValue > newMaxValue {
				newMaxValue = currentValue
				newMaxKey = nonAdjacentKey
			}
		}
		result1.mapMaxValue[key] = newMaxValue
		result1.mapNextListMax[key] = append([]int{key}, result1.mapNextListMax[newMaxKey]...)
		if result1.mapMaxValue[key] > updateResult.maxValueTail {
			updateResult.maxValueTail = result1.mapMaxValue[key]
			updateResult.maxListTail = result1.mapNextListMax[key]
			updateResult.maxTailKey = key
		}
	}
	return result1, updateResult
}

func reUpdateHeadList(arr []int32, result1 Result, updateResult UpdateResult) Result {
	sortedKeyList := result1.sortedKeyList
	tailKeyList := result1.tailList
	adjacentTailList := result1.adjacentTailList
	maxTailKey := updateResult.maxTailKey
	var maxList []int
	var maxValue int32
	for i := len(sortedKeyList) - len(tailKeyList) - len(adjacentTailList) - 1; i >= 0; i-- {
		key := sortedKeyList[i]
		var maxValueKey int32
		var maxKey int
		reUpdateValueNonAdjKey := []int{}
		nonAdjList := result1.mapNonAdjacentList[key]

		for _, nonAdjKey := range nonAdjList {
			if !slices.Contains(tailKeyList, nonAdjKey) && !slices.Contains(adjacentTailList, nonAdjKey) {
				reUpdateValueNonAdjKey = append(reUpdateValueNonAdjKey, nonAdjKey)
			}
			if slices.Contains(adjacentTailList, nonAdjKey) && nonAdjKey == maxTailKey {
				reUpdateValueNonAdjKey = append(reUpdateValueNonAdjKey, nonAdjKey)
				break
			} else if slices.Contains(adjacentTailList, nonAdjKey) {
				reUpdateValueNonAdjKey = append(reUpdateValueNonAdjKey, nonAdjKey)
			}
			if slices.Contains(tailKeyList, nonAdjKey) && len(reUpdateValueNonAdjKey) > 0 {
				lastElm := reUpdateValueNonAdjKey[len(reUpdateValueNonAdjKey)-1]
				if lastElm+1 == nonAdjKey {
					reUpdateValueNonAdjKey = append(reUpdateValueNonAdjKey, nonAdjKey)
				}
				break
			}
		}
		for _, nonAdjKey := range reUpdateValueNonAdjKey {
			newMaxValue := arr[key] + result1.mapMaxValue[nonAdjKey]
			if newMaxValue > maxValueKey {
				maxValueKey = newMaxValue
				maxKey = nonAdjKey
			}
		}
		result1.mapMaxValue[key] = maxValueKey
		result1.mapNextListMax[key] = append([]int{key}, result1.mapNextListMax[maxKey]...)
		/* if key == 4324 && result1.mapMaxValue[4326] == 177708 {
			log.Println(nonAdjList)
			log.Println(reUpdateValueNonAdjKey)
			log.Println(result1.tailList)
			log.Println(result1.adjacentTailList)
			log.Println(maxTailKey)
			log.Fatal("Error")
		} */
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
