Given an array of integers, find the subset of non-adjacent elements with the maximum sum. Calculate the sum of that subset. It is possible that the maximum sum is , the case when all elements are negative.

Example

arr=[-2, 1, 3, -4, 5]

The following subsets with more than 1 element exist. These exclude the empty subset and single element subsets which are also valid.

Subset      Sum
[-2, 3, 5]   6
[-2, 3]      1
[-2, -4]    -6
[-2, 5]      3
[1, -4]     -3
[1, 5]       6
[3, 5]       8
The maximum subset sum is 8. Note that any individual element is a subset as well.

arr=[-2, -3, -1]
In this case, it is best to choose no element: return 0.

input00.txt=>151598486
input06.txt=>7412694
input07.txt=>114846322
input08.txt=>27902749
input09.txt=>175041179