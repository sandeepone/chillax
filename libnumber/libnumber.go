package libnumber

import (
    "sort"
)

func LargestInt(numbers []int) int {
    sort.Ints(numbers)
    return numbers[len(numbers) - 1]
}

func FirstGapIntSlice(numbers []int) int {
    sort.Sort(sort.Reverse(sort.IntSlice(numbers)))

    for index, value := range numbers {
        if index + 1 < len(numbers) {
            nextValue := numbers[index + 1]

            if (value - 1) != nextValue {
                return (value - 1)
            }
        }
    }
    return -1
}