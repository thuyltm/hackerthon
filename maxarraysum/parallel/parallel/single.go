package main

import (
	"log"
	"slices"
)

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
	oldSortedKeyList1 := result1.sortedKeyList
	result1 = join2Result(result1, result2)
	lastElmMaxList1 := result1.maxList[len(result1.maxList)-1]
	firstElmMaxList2 := result2.maxList[0]
	maxValueList := []int32{}
	maxListList := [][]int{}
	var maxValue3, maxValue4, maxValue5, maxValue6, maxValue7 int32
	var maxList3, maxList4, maxList5, maxList6, maxList7 []int
	if lastElmMaxList1+2 <= firstElmMaxList2 {
		maxValue := result1.maxValue + result2.maxValue
		maxList := make([]int, len(result1.maxList))
		copy(maxList, result1.maxList)
		tailMaxList := make([]int, len(result2.maxList))
		copy(tailMaxList, result2.maxList)
		maxList = append(maxList, tailMaxList...)
		maxValueList = append(maxValueList, maxValue)
		maxListList = append(maxListList, maxList)
	} else if lastElmMaxList1 == firstElmMaxList2 {
		tailValue2 := result2.maxValue - arr[firstElmMaxList2]
		maxValue := result1.maxValue + tailValue2
		maxList := make([]int, len(result1.maxList))
		copy(maxList, result1.maxList)
		tailMaxList := make([]int, len(result2.maxList)-1)
		copy(tailMaxList, result2.maxList[1:])
		maxList = append(maxList, tailMaxList...)
		maxValueList = append(maxValueList, maxValue)
		maxListList = append(maxListList, maxList)
	} else {
		maxValue3, maxList3 = reCalcForward(arr, result1, result2)
		maxValue4, maxList4 = reCalcBackward(arr, result1, result2)
		maxValue5, maxList5 = reCalcFowardWith2ndChild(arr, result1, result2)
		maxValue6, maxList6 = reCalcBackwardWith2ndChild(arr, result1, result2)
		maxValueList = append(maxValueList, []int32{maxValue3, maxValue4, maxValue5, maxValue6}...)
		maxListList = append(maxListList, [][]int{maxList3, maxList4, maxList5, maxList6}...)
	}
	anotherValue, anotherPath := choosePossible2ndMaxResult1(arr, result1, oldSortedKeyList1)
	if anotherValue > 0 {
		result3 := Result{
			sortedKeyList: result1.sortedKeyList,
			maxValue:      anotherValue,
			maxList:       anotherPath,
		}
		maxValue7, maxList7 = reCalcForward(arr, result3, result2)
		maxValueList = append(maxValueList, maxValue7)
		maxListList = append(maxListList, maxList7)
	}
	result1.maxValue, result1.maxList = getLargestValue(maxValueList, maxListList)
	return result1
}

func choosePossible2ndMaxResult1(arr []int32, result1 Result, oldSortedKeyList1 []int) (int32, []int) {
	maxList1 := result1.maxList
	if len(oldSortedKeyList1) == len(maxList1) {
		return -1, []int{}
	}
	if len(oldSortedKeyList1) == 2 && len(maxList1) == 1 {
		for i := 0; i < len(oldSortedKeyList1); i++ {
			if oldSortedKeyList1[i] != maxList1[0] {
				return arr[oldSortedKeyList1[i]], []int{oldSortedKeyList1[i]}
			}
		}
	}
	var maxValue int32
	var maxList []int
	maxList1KeyList := []int{maxList1[len(maxList1)-1]}
	if len(maxList1) >= 2 {
		maxList1KeyList = append([]int{maxList1[len(maxList1)-2]}, maxList1KeyList...)
	}
	selectedKeyList := []int{}
	for i := len(oldSortedKeyList1) - 3; i <= len(oldSortedKeyList1)-1; i++ {
		if !slices.Contains(maxList1KeyList, oldSortedKeyList1[i]) {
			selectedKeyList = append(selectedKeyList, oldSortedKeyList1[i])
		}
	}
	if len(maxList1) == 1 && len(selectedKeyList) == 2 && selectedKeyList[0]+1 < selectedKeyList[1] {
		return arr[selectedKeyList[0]] + arr[selectedKeyList[1]], append([]int{selectedKeyList[0]}, selectedKeyList[1])
	}
	if len(maxList1) >= 2 && len(selectedKeyList) >= 1 {
		selectedTail := selectedKeyList[0]
		if len(selectedKeyList) > 1 {
			selectedTail = selectedKeyList[1]
		}
		if (maxList1KeyList[1]+1 == selectedTail || selectedTail+1 == maxList1KeyList[1]) &&
			maxList1KeyList[0]+1 < selectedTail {
			maxValue = result1.maxValue - arr[maxList1KeyList[1]] + arr[selectedTail]
			maxList = make([]int, len(maxList1)-1)
			copy(maxList, maxList1[:len(maxList1)-1])
			maxList = append(maxList, selectedTail)
			return maxValue, maxList
		}
		if (maxList1KeyList[0]+1 == selectedKeyList[0] || selectedKeyList[0]+1 == maxList1KeyList[0]) &&
			selectedKeyList[0]+1 < maxList1KeyList[1] {
			result5 := Result{
				sortedKeyList: []int{selectedKeyList[0], maxList1KeyList[1]},
				maxValue:      arr[selectedKeyList[0]] + arr[maxList1KeyList[1]],
				maxList:       []int{selectedKeyList[0], maxList1KeyList[1]},
			}
			result3 := Result{
				sortedKeyList: oldSortedKeyList1,
				maxValue:      result1.maxValue,
				maxList:       result1.maxList,
			}
			return reCalcBackward(arr, result3, result5)
		}
		var result4 Result
		if len(selectedKeyList) == 1 {
			result4 = Result{
				sortedKeyList: selectedKeyList,
				maxValue:      arr[selectedKeyList[0]],
				maxList:       selectedKeyList,
			}
		} else if len(selectedKeyList) == 2 && selectedKeyList[0]+1 < selectedKeyList[1] {
			result4 = Result{
				sortedKeyList: selectedKeyList,
				maxValue:      arr[selectedKeyList[0]] + arr[selectedKeyList[1]],
				maxList:       selectedKeyList,
			}
		} else {
			return -1, []int{}
		}
		result3 := Result{
			sortedKeyList: oldSortedKeyList1,
			maxValue:      result1.maxValue,
			maxList:       result1.maxList,
		}
		return reCalcBackward(arr, result3, result4)
	}
	return -1, []int{}
}

func getLargestValue(maxValueList []int32, maxListList [][]int) (int32, []int) {
	var maxValue int32
	var maxIndex int
	var i int
	for i = 0; i < len(maxValueList); i++ {
		if maxValueList[i] > maxValue {
			maxValue = maxValueList[i]
			maxIndex = i
		}
	}
	return maxValueList[maxIndex], maxListList[maxIndex]
}

func reCalcFowardWith2ndChild(arr []int32, result1 Result, result2 Result) (int32, []int) {
	maxList1 := result1.maxList
	maxList2 := result2.maxList
	lastElmMaxList1 := maxList1[len(maxList1)-1]
	newMaxList := make([]int, len(maxList1)-1)
	copy(newMaxList, maxList1[:len(maxList1)-1])
	tailNewMaxList := make([]int, len(maxList2))
	copy(tailNewMaxList, maxList2)
	newMaxList = append(newMaxList, tailNewMaxList...)
	newMaxValue := result1.maxValue - arr[lastElmMaxList1] + result2.maxValue
	return newMaxValue, newMaxList
}

func reCalcBackwardWith2ndChild(arr []int32, result1 Result, result2 Result) (int32, []int) {
	maxList1 := result1.maxList
	maxList2 := result2.maxList
	firstElmMaxList2 := maxList2[0]
	newMaxList := make([]int, len(maxList1))
	copy(newMaxList, maxList1)
	tailNewMaxList := make([]int, len(maxList2)-1)
	copy(tailNewMaxList, maxList2[1:])
	newMaxList = append(maxList1, tailNewMaxList...)
	newMaxValue := result1.maxValue + result2.maxValue - arr[firstElmMaxList2]
	return newMaxValue, newMaxList
}

func join2Result(result1 Result, result2 Result) Result {
	sortedKeyList1 := result1.sortedKeyList
	sortedKeyList2 := result2.sortedKeyList
	if sortedKeyList1[len(sortedKeyList1)-1] != sortedKeyList2[len(sortedKeyList2)-1] {
		startIndexJoinResult2 := 0
		if sortedKeyList1[len(sortedKeyList1)-1] == sortedKeyList2[0] {
			startIndexJoinResult2 = 1
		}
		result1.sortedKeyList = append(result1.sortedKeyList, sortedKeyList2[startIndexJoinResult2:]...)
	}
	return result1
}

func debug(tempResult Result, result1 Result, result2 Result, arr []int32, start int, end int, middle int) {
	if tempResult.maxValue != result1.maxValue {
		log.Println("==============debug===================")
		log.Println(result1.maxList, "value=", result1.maxValue)
		log.Println(tempResult.maxList, "value=", tempResult.maxValue)
		log.Println(result2.maxList, "value=", result2.maxValue)
		log.Println(result2.sortedKeyList)
		for i := 0; i < len(result1.maxList); i++ {
			if result1.maxList[i] != tempResult.maxList[i] {
				log.Printf("Diff result1 %d, temp %d", result1.maxList[i], tempResult.maxList[i])
				log.Println(result1.maxList[i], arr[result1.maxList[i]])
				log.Println(tempResult.maxList[i], arr[tempResult.maxList[i]])
				if i-1 >= 0 {
					log.Printf("Parent %d", result1.maxList[i-1])
					log.Println(result1.mapNonAdjacentList[result1.maxList[i-1]])
				}
				break
			}
		}
		log.Println("==============end debug===================")
		log.Fatal("Error")
	}
}

func reCalcBackward(arr []int32, result1 Result, result2 Result) (int32, []int) {
	sortedKeyList := result1.sortedKeyList
	maxValue := result2.maxValue
	maxList := make([]int, len(result2.maxList))
	copy(maxList, result2.maxList)
	maxKey := result2.maxList[0]
	startSearchIndex := len(sortedKeyList) - len(result2.sortedKeyList)
	var parentMaxTailKey, parentMaxTailKeyIndex int
	var otherParentMaxTailKey int
	parentMaxTailKey, parentMaxTailKeyIndex = findNearestParent(sortedKeyList, maxKey, startSearchIndex)
	if parentMaxTailKey < 0 {
		return maxValue, maxList
	}
	if parentMaxTailKeyIndex == 0 {
		maxValue = maxValue + arr[parentMaxTailKey]
		maxList = append([]int{parentMaxTailKey}, maxList...)
		return maxValue, maxList
	}
	startSearchIndex = parentMaxTailKeyIndex
	otherParentMaxTailKey = sortedKeyList[parentMaxTailKeyIndex-1]
	parentKeyList := []int{parentMaxTailKey}
	if otherParentMaxTailKey+1 == parentMaxTailKey {
		parentKeyList = append(parentKeyList, otherParentMaxTailKey)
	}
	var newMaxValue int32
	var newMaxList []int
	for _, key2 := range parentKeyList {
		var currentValue int32
		var currentList []int
		if slices.Contains(result1.maxList, key2) {
			headMaxValue, headMaxList := extractHeadValue(arr, result1, key2)
			currentValue = maxValue + headMaxValue
			maxList2 := make([]int, len(result2.maxList))
			copy(maxList2, result2.maxList)
			currentList = append(headMaxList, maxList2...)
		} else {
			maxList2 := make([]int, len(result2.maxList))
			copy(maxList2, result2.maxList)
			result3 := Result{
				sortedKeyList: append([]int{key2}, result2.sortedKeyList...),
				maxValue:      maxValue + arr[key2],
				maxList:       append([]int{key2}, maxList2...),
			}
			currentValue, currentList = reCalcBackward(arr, result1, result3)
		}

		if currentValue > newMaxValue {
			newMaxValue = currentValue
			newMaxList = currentList
		}
	}
	return newMaxValue, newMaxList
}

func findNearestParent(sortedKeyList []int, searchKey int, startIndex int) (int, int) {
	if startIndex == -1 {
		startIndex = 0
	}
	for i := startIndex; i >= 0; i-- {
		key := sortedKeyList[i]
		if key+2 <= searchKey {
			return key, i
		}
	}
	return -1, -1
}

func findNearestChild(sortedKeyList []int, searchKey int, startIndex int) (int, int) {
	if startIndex == -1 {
		startIndex = 0
	}
	for i := startIndex; i < len(sortedKeyList); i++ {
		key := sortedKeyList[i]
		if searchKey+2 <= key {
			return key, i
		}
	}
	return -1, -1
}

func reCalcForward(arr []int32, result1 Result, result2 Result) (int32, []int) {
	sortedKeyList := result1.sortedKeyList
	maxList1 := result1.maxList
	lastKeyMaxList1 := maxList1[len(maxList1)-1]
	startSearchIndex := len(sortedKeyList) - len(result2.sortedKeyList)
	childKey, childIndex := findNearestChild(sortedKeyList, lastKeyMaxList1, startSearchIndex)
	var maxChildValue int32
	var maxChildList []int
	var maxValue int32
	var maxList []int
	if childKey > 0 && childIndex+1 < len(sortedKeyList) {
		otherChildKey := sortedKeyList[childIndex+1]
		startSearchIndex = childIndex
		traverseChildList := []int{childKey}
		if childKey+1 == otherChildKey {
			traverseChildList = append(traverseChildList, otherChildKey)
			startSearchIndex = childIndex + 1
		}
		for _, currentKey := range traverseChildList {
			if slices.Contains(result2.maxList, currentKey) {
				maxChildValue, maxChildList = extractTailValue(arr, result2, currentKey)
			} else {
				maxChildValue, maxChildList = findChildrenMaxList(arr, result1, result2, currentKey, startSearchIndex)
			}
			currentValue := result1.maxValue + maxChildValue
			if currentValue > maxValue {
				maxValue = currentValue
				maxList = make([]int, len(result1.maxList))
				copy(maxList, result1.maxList)
				maxList = append(maxList, maxChildList...)
			}
		}
		return maxValue, maxList
	} else if childKey > 0 {
		maxList := make([]int, len(result1.maxList))
		copy(maxList, result1.maxList)
		maxList = append(maxList, childKey)
		return result1.maxValue + arr[childKey], maxList
	}
	maxList = make([]int, len(result1.maxList))
	copy(maxList, result1.maxList)
	return result1.maxValue, maxList
}

func extractTailValue(arr []int32, result Result, extractedKey int) (int32, []int) {
	maxList := result.maxList
	maxValue := result.maxValue
	var excludeKeyList []int
	var i int
	for i = 0; i < len(maxList); i++ {
		currentKey := maxList[i]
		if currentKey == extractedKey {
			break
		}
		excludeKeyList = append(excludeKeyList, currentKey)
	}
	var currentKey int
	for j := 0; j < len(excludeKeyList); j++ {
		currentKey = excludeKeyList[j]
		maxValue = maxValue - arr[currentKey]
	}
	newMaxList := make([]int, len(maxList)-i)
	copy(newMaxList, maxList[i:])
	return maxValue, newMaxList
}

func extractHeadValue(arr []int32, result Result, extractedKey int) (int32, []int) {
	maxList := result.maxList
	maxValue := result.maxValue
	var excludeKeyList []int
	var i int
	for i = len(maxList) - 1; i >= 0; i-- {
		currentKey := maxList[i]
		if currentKey == extractedKey {
			break
		}
		excludeKeyList = append(excludeKeyList, currentKey)
	}
	var currentKey int
	for j := 0; j < len(excludeKeyList); j++ {
		currentKey = excludeKeyList[j]
		maxValue = maxValue - arr[currentKey]
	}
	newMaxList := make([]int, i+1)
	copy(newMaxList, maxList[:i+1])
	return maxValue, newMaxList
}

func findChildrenMaxList(arr []int32, result1 Result, result2 Result, key int, startSearchIndex int) (int32, []int) {
	sortedKeyList := result1.sortedKeyList
	var childKey, childIndex int
	childKey, childIndex = findNearestChild(sortedKeyList, key, startSearchIndex)
	if childKey < 0 {
		return arr[key], []int{key}
	}
	if childIndex == len(sortedKeyList)-1 {
		return arr[key] + arr[childKey], append([]int{key}, childKey)
	}
	startSearchIndex = childIndex
	otherChildKey := sortedKeyList[childIndex+1]
	childKeyList := []int{childKey}
	if childKey+1 == otherChildKey && arr[childKey] < arr[otherChildKey] {
		childKeyList = append(childKeyList, otherChildKey)
		startSearchIndex = childIndex + 1
	}
	var maxChildValue int32
	var maxChildList []int
	var maxValue int32
	var maxList []int
	for _, currentKey := range childKeyList {
		if slices.Contains(result2.maxList, currentKey) {
			maxChildValue, maxChildList = extractTailValue(arr, result2, currentKey)
		} else {
			maxChildValue, maxChildList = findChildrenMaxList(arr, result1, result2, currentKey, startSearchIndex)
		}
		currentValue := arr[key] + maxChildValue
		if currentValue > maxValue {
			maxValue = currentValue
			maxList = append([]int{key}, maxChildList...)
		}
	}

	return maxValue, maxList
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
			if len(mapNonAdjacentList[int(left)]) == 1 &&
				mapNonAdjacentList[int(left)][0]+2 <= i {
				break
			}
			mapNonAdjacentList[int(left)] = append(mapNonAdjacentList[int(left)], int(i))
		}
		if len(mapNonAdjacentList[int(left)]) == 2 {
			break
		}
	}
	if len(mapNonAdjacentList[int(left)]) == 0 && arr[left] > 0 {
		mapNonAdjacentList[int(left)] = []int{}
	}
	return mapNonAdjacentList, sortedKeyList
}

func calcAndReturnListMaxSum(mapNonAdjacentList map[int][]int, arr []int32, sortedKeyList []int) Result {
	maxList := []int{}
	var maxValue int32
	calcResult := Result{
		mapNextListMax:     make(map[int][]int),
		mapMaxValue:        make(map[int]int32),
		mapNonAdjacentList: mapNonAdjacentList,
		sortedKeyList:      sortedKeyList,
	}
	for i := len(sortedKeyList) - 1; i >= 0; i-- {
		key := sortedKeyList[i]
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
	return calcResult
}

func calcKeyFromNonAdjacentList(arr []int32, calcResult Result, key int) Result {
	nonAdjacentList := calcResult.mapNonAdjacentList[key]
	if len(nonAdjacentList) == 0 {
		calcResult.mapNextListMax[key] = []int{key}
		calcResult.mapMaxValue[key] = arr[key]
		return calcResult
	}
	maxKey := nonAdjacentList[0]
	for _, key2 := range nonAdjacentList {
		if calcResult.mapMaxValue[key2] > calcResult.mapMaxValue[maxKey] {
			maxKey = key2
		}
	}
	calcResult.mapNextListMax[key] = append([]int{key}, calcResult.mapNextListMax[maxKey]...)
	calcResult.mapMaxValue[key] = arr[key] + calcResult.mapMaxValue[maxKey]
	return calcResult
}
