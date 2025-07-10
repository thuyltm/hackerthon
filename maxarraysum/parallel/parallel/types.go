package main

type Result struct {
	mapNextListMax     map[int][]int
	mapMaxValue        map[int]int32
	maxList            []int
	maxValue           int32
	mapNonAdjacentList map[int][]int
	sortedKeyList      []int
	start              int
	end                int
}

type Input struct {
	start int
	end   int
	batch int
}

type SortedInput struct {
	result Result
	batch  int
}

type CombineSortedInput struct {
	result []Result
	batch  int
}

type UpdateResult struct {
	maxValueTail int32
	maxListTail  []int
	maxTailKey   int
}
