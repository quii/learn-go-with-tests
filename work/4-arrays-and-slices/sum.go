package arrays

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAll(numbersToSum ...[]int) []int {
	// lenOfNumbers := len(numbersToSum)
	// sums := make([]int, lenOfNumbers)

	var sums []int

	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
		// sums[index] = Sum(numbers)
	}

	return sums
}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
			continue
		}
		sums = append(sums, Sum(numbers[1:]))
	}

	return sums
}
