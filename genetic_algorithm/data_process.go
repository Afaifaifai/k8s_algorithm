package genetic_algorithm

import (
	s "k8s_algorithm/settings"

	"gonum.org/v1/gonum/stat/distuv"
)

func (data *Data) generate_uniform_population() ([]int, [][s.DIMENSION]float64) {
	uniform := distuv.Uniform{ // [ Min, Max )
		Min: 0,
		Max: float64(s.KNAPSACK_QUANTITY),
	}

	var population []int = make([]int, s.POPULATION_SIZE) // items pick knapsack
	var knapsack_weights [][s.DIMENSION]float64 = make([][s.DIMENSION]float64, s.KNAPSACK_QUANTITY)
	for i := 0; i < s.ITEM_QUANTITY; i++ {
		knapsack_idx := int(uniform.Rand())                              // ITEM[i] pick the knapsack
		add(&(knapsack_weights[knapsack_idx]), &(data.Items_weights[i])) // 3 dim means 3 parameters
		population[i] = knapsack_idx
	}
	return population, knapsack_weights
}

func (data *Data) generate_solutions() [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 {
	var solutions [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = make([][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, s.SOLUTION_SIZE)
	for i := 0; i < s.SOLUTION_SIZE; i++ {
		var population []int                        // items pick knapsack
		var knapsack_weights [][s.DIMENSION]float64 // KNAPSACK_QUANTITY * DIMENSION

		for {
			population, knapsack_weights = data.generate_uniform_population()
			if under_capacity(knapsack_weights, data.Limit_of_knapsack) {
				break
			}
		}

		/*
			popilation = [ks2, ks0, ks0, ks1, ks0, ...]  --> each item pick a knapsack
		*/

		var chromosome_formats [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64
		for item_idx, knapsack_idx := range population {
			chromosome_formats[knapsack_idx][item_idx] = 1
		}
		/*
			chromosome = [1, 0, 0, 1, 0, ...]  --> ks0  // each knapsack pick items
						 [0, 1, 1, 0, 0, ...]  --> ks1
						 ...
		*/
		solutions[i] = chromosome_formats
	}
	return solutions
}

func add(ary1 *[s.DIMENSION]float64, ary2 *[s.DIMENSION]float64) {
	length := len(ary1)
	for i := 0; i < length; i++ {
		(*ary1)[i] += (*ary2)[i]
	}
}

func under_capacity(knapsack_weights [][s.DIMENSION]float64, capacity_of_knapsack [][s.DIMENSION]float64) bool {
	for i := 0; i < s.KNAPSACK_QUANTITY; i++ {
		for j := 0; j < s.DIMENSION; j++ {
			if knapsack_weights[i][j] > capacity_of_knapsack[i][j] {
				return false
			}
		}
	}
	return true
}
