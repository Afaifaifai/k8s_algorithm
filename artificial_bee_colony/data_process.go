package artificial_bee_colony

import (
	s "k8s_algorithm/tools"
	"math/rand"

	"gonum.org/v1/gonum/stat"
)

func initialize_food_sources() []FoodSource {
	var food_sources []FoodSource
	for i := 0; i < s.COLONY_SIZE; i++ {
		population, fitness := generate_uniform_valid_population_fitness()
		food_sources = append(food_sources, FoodSource{
			population: population,
			fitness:    fitness,
			trials:     0,
		})
	}
	return food_sources
}

func generate_uniform_valid_population_fitness() ([]int, float64) {
	var population []int = make([]int, s.ITEM_QUANTITY)
	var knapsack_dim_weights [][]float64
	for {
		for i := 0; i < s.POPULATION_SIZE; i++ {
			population[i] = rand.Intn(s.KNAPSACK_QUANTITY)
		}
		knapsack_dim_weights = calculate_weights(population)
		if under_limit(knapsack_dim_weights) {
			break
		}
	}

	return population, calculate_fitness(knapsack_dim_weights)
}

func calculate_weights(genes []int) [][]float64 {
	var knapsack_dim_weights [][]float64 = make([][]float64, s.KNAPSACK_QUANTITY)
	for i := 0; i < s.KNAPSACK_QUANTITY; i++ {
		knapsack_dim_weights[i] = make([]float64, s.DIMENSION)
	}

	for item_idx, ks_idx := range genes {
		for dim_idx := 0; dim_idx < s.DIMENSION; dim_idx++ {
			knapsack_dim_weights[ks_idx][dim_idx] += Items_weights[item_idx][dim_idx]
		}
	}

	for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
		for dim_idx := 0; dim_idx < s.DIMENSION; dim_idx++ {
			knapsack_dim_weights[ks_idx][dim_idx] += Previous_state_of_knapsack[ks_idx][dim_idx]
		}
	}

	return knapsack_dim_weights
}

func under_limit(knapsack_dim_weights [][]float64) bool {
	for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
		for dim_idx := 0; dim_idx < s.DIMENSION; dim_idx++ {
			if knapsack_dim_weights[ks_idx][dim_idx] > Limit_of_knapsack[ks_idx][dim_idx] {
				return false
			}
		}
	}

	return true
}

func calculate_fitness(knapsack_dim_weights [][]float64) float64 {
	var dim_knapsack_weights [][]float64 = make([][]float64, s.DIMENSION)
	for i := 0; i < s.DIMENSION; i++ {
		dim_knapsack_weights[i] = make([]float64, s.KNAPSACK_QUANTITY)
	}

	var fitness, penalty float64 = 0.0, 0.0
	for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
		for dim_idx := 0; dim_idx < s.DIMENSION; dim_idx++ {
			dim_knapsack_weights[dim_idx][ks_idx] = knapsack_dim_weights[ks_idx][dim_idx] / Limit_of_knapsack[ks_idx][dim_idx]
			if dim_knapsack_weights[dim_idx][ks_idx] > 1.0 {
				penalty += dim_knapsack_weights[dim_idx][ks_idx] - 1.0
			}
		}
	}

	for _, weights_percestage := range dim_knapsack_weights {
		// fmt.Println(weights_percestage)
		fitness += stat.StdDev(weights_percestage[:], nil)
	}

	return fitness + penalty*s.ABC_LAMBDA
}
