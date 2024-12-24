package nondominated_sorting_genetic_algorithm

import (
	s "k8s_algorithm/tools"
	"math/rand"

	"gonum.org/v1/gonum/stat"
)

type Chromosome struct {
	genes   []int
	fitness []float64
}

func initialize_chromosomes() []Chromosome {
	var chromosomes []Chromosome
	for i := 0; i < s.CHROMOSOME_QUANTITY; i++ {
		genes, fitness := generate_uniform_genes_fitness()
		chromosomes = append(chromosomes, Chromosome{
			genes:   genes,
			fitness: fitness,
		})
	}
	return chromosomes
}

func generate_uniform_genes_fitness() ([]int, []float64) {
	var genes []int = make([]int, s.POPULATION_SIZE)
	var knapsack_dim_weights [][]float64
	for {
		for i := 0; i < s.POPULATION_SIZE; i++ {
			genes[i] = rand.Intn(s.KNAPSACK_QUANTITY)
		}
		knapsack_dim_weights = calculate_weights(genes)
		if under_limit(knapsack_dim_weights) {
			break
		}
	}

	return genes, calculate_fitness(knapsack_dim_weights)
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

func calculate_fitness(knapsack_dim_weights [][]float64) []float64 {
	var dim_knapsack_weights [][]float64 = make([][]float64, s.DIMENSION)
	for i := 0; i < s.DIMENSION; i++ {
		dim_knapsack_weights[i] = make([]float64, s.KNAPSACK_QUANTITY)
	}

	var fitness []float64 = make([]float64, s.DIMENSION)
	for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
		for dim_idx := 0; dim_idx < s.DIMENSION; dim_idx++ {
			dim_knapsack_weights[dim_idx][ks_idx] = knapsack_dim_weights[ks_idx][dim_idx] / Limit_of_knapsack[ks_idx][dim_idx]
		}
	}

	for idx, weights_percestage := range dim_knapsack_weights {
		// fmt.Println(weights_percestage)
		fitness[idx] = stat.StdDev(weights_percestage[:], nil)
	}

	return fitness
}
