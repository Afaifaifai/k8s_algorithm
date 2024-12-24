package hybrid_genetic_algorithm

import (
	s "k8s_algorithm/tools"
	"math"
	"math/rand"
)

func Simulated_Annealing(solution []int, fitness float64) ([]int, float64) {
	var best_solution []int = make([]int, s.ITEM_QUANTITY)
	copy(best_solution, solution)
	var best_fitness float64 = fitness

	var temperature float64 = s.INITIAL_TEMPERATURE
	for i := 0; i < s.SA_ITERATIONS; i++ {
		// fmt.Println(i)
		var neighbor []int = make([]int, s.ITEM_QUANTITY)
		copy(neighbor, solution)
		var rand_item_idx int = rand.Intn(s.ITEM_QUANTITY)

		original_ks_idx := neighbor[rand_item_idx]
		for neighbor[rand_item_idx] == original_ks_idx {
			neighbor[rand_item_idx] = rand.Intn(s.KNAPSACK_QUANTITY)
		}

		var neighbor_weights [][]float64 = calculate_weights(neighbor)
		var neighbor_fitness float64 = calculate_fitness(neighbor_weights)
		var delta float64 = neighbor_fitness - fitness

		if delta < 0 { // neighbor fitness < solution fitness  --> neighbor is better
			solution = neighbor
			fitness = neighbor_fitness
		} else { // neighbor fitness >= solution fitness  --> solution is better
			var probability float64
			if temperature > 0 {
				probability = math.Exp(delta / temperature)
			} else {
				probability = 0
			}

			if rand.Float64() < probability {
				solution = neighbor
				fitness = neighbor_fitness
			}
		}

		temperature *= s.COOLING_RATE
		if temperature < s.MINIMUM_TEMPERATURE {
			break
		}

		if fitness < best_fitness {
			best_fitness = fitness
			copy(best_solution, solution)
		}
	}
	return best_solution, best_fitness
}

// func sa_fitness(solution *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64) float64 { // DIMENSION * KNAPSACK_QUANTITY
// 	var knapsack_weights [s.DIMENSION][s.KNAPSACK_QUANTITY]float64
// 	for ks_idx, knapsack := range *solution {
// 		for item_idx, item_exists := range knapsack {
// 			for dim := 0; dim < s.DIMENSION; dim++ {
// 				knapsack_weights[dim][ks_idx] += Items_weights[item_idx][dim] * item_exists // item_exists = 0 or 1
// 			}
// 		}
// 	}

// 	var fitness float64 = 0
// 	for dim := 0; dim < s.DIMENSION; dim++ {
// 		for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
// 			knapsack_weights[dim][ks_idx] = (knapsack_weights[dim][ks_idx] + Previous_state_of_knapsack[ks_idx][dim]) / Limit_of_knapsack[ks_idx][dim]
// 		}
// 	}

// 	for dim := 0; dim < s.DIMENSION; dim++ {
// 		for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
// 			if knapsack_weights[dim][ks_idx] > 1 {
// 				var max float64 = max_element(knapsack_weights[dim][:])
// 				knapsack_weights[dim][0] = max
// 				knapsack_weights[dim][1] = -max
// 				break
// 			}
// 		}
// 	}

// 	for _, weights_percestage := range knapsack_weights {
// 		fitness += stat.StdDev(weights_percestage[:], nil)
// 	}
// 	return fitness
// }

// func max_element(ary []float64) float64 {
// 	var max float64 = 0
// 	for _, element := range ary {
// 		if element > max {
// 			max = element
// 		}
// 	}
// 	return max
// }
