package arrays

func SumArray(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAll(arrays ...[]int) []int {
	lengthOfNumbers := len(arrays)
	sums := make([]int, lengthOfNumbers)

	for i, numbers := range arrays {
		sums[i] = SumArray(numbers)
	}

	return sums
}
