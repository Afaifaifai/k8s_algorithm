package new_genetic_algorithm

import (
	"fmt"
	s "k8s_algorithm/tools"
	"math"
	"math/rand"
	"sort"
)

var Items_weights [][s.DIMENSION]float64
var Previous_state_of_knapsack [][s.DIMENSION]float64
var Limit_of_knapsack [][s.DIMENSION]float64

func Run(
	items_weights [][s.DIMENSION]float64,
	previous_state_of_knapsack [][s.DIMENSION]float64,
	limit_of_knapsack [][s.DIMENSION]float64,
) (float64, []int, []float64) {

	Items_weights = items_weights
	Previous_state_of_knapsack = previous_state_of_knapsack
	Limit_of_knapsack = limit_of_knapsack

	var best_fitness float64 = math.MaxFloat64
	var best_fitness_solution []int
	var best_fitness_in_iterations []float64

	var chromosomes []Chromosome = initialize_chromosomes()
	chromosomes = argsort(chromosomes) // sort chromosomes by fitness

	for iter := 0; iter < s.GA_ITERATIONS; iter++ {
		var elite_parents []Chromosome = selection(chromosomes) // the chrmosomes is already argsorterd
		for i := 0; i < s.ELITE_QUANTITY; i++ {
			for j := i + 1; j < s.ELITE_QUANTITY; j++ {
				var child1_genes, child2_genes []int = crossover(elite_parents[i], elite_parents[j])
				mutate(child1_genes)
				mutate(child2_genes)
				chromosomes = append(chromosomes, born(child1_genes), born(child2_genes))
			}
		}

		chromosomes = replacement(chromosomes)

		if chromosomes[0].fitness < best_fitness {
			best_fitness = chromosomes[0].fitness
			best_fitness_solution = chromosomes[0].genes
		}
		best_fitness_in_iterations = append(best_fitness_in_iterations, chromosomes[0].fitness)
		if s.PRINT_PERMIT {
			fmt.Printf("Best Fitness: %.6f, Iteration: %d\n", best_fitness, iter)
		}
	}

	return best_fitness, best_fitness_solution, best_fitness_in_iterations
}

func selection(chromosomes []Chromosome) []Chromosome {
	var elite_parents []Chromosome = make([]Chromosome, s.ELITE_QUANTITY)
	for i := 0; i < s.ELITE_QUANTITY; i++ {
		elite_parents[i] = chromosomes[i]
	}
	return elite_parents
}

func crossover(parent1 Chromosome, parent2 Chromosome) ([]int, []int) {
	var crossover_point int = rand.Intn(s.ITEM_QUANTITY)
	var child1_genes []int = make([]int, s.ITEM_QUANTITY)
	var child2_genes []int = make([]int, s.ITEM_QUANTITY)
	for i := 0; i < s.ITEM_QUANTITY; i++ {
		if i < crossover_point {
			child1_genes[i] = parent1.genes[i]
			child2_genes[i] = parent2.genes[i]
		} else {
			child1_genes[i] = parent2.genes[i]
			child2_genes[i] = parent1.genes[i]
		}
	}
	return child1_genes, child2_genes
}

func mutate(genes []int) {
	for i := 0; i < s.ITEM_QUANTITY; i++ {
		if rand.Float64() < s.MUTATION_RATE {
			var original_knapsack int = genes[i]
			for genes[i] == original_knapsack {
				genes[i] = rand.Intn(s.KNAPSACK_QUANTITY)
			}
		}
	}
}

func born(genes []int) Chromosome {
	var knapsack_dim_weights *[s.KNAPSACK_QUANTITY][s.DIMENSION]float64 = calculate_weights(genes)
	if under_limit(knapsack_dim_weights) {
		return Chromosome{
			genes:   genes,
			fitness: calculate_fitness(knapsack_dim_weights),
		}
	} else {
		return Chromosome{
			genes:   genes,
			fitness: math.MaxFloat64,
		}
	}
}

func replacement(chromosomes []Chromosome) []Chromosome {
	chromosomes = argsort(chromosomes)
	var new_generation []Chromosome = make([]Chromosome, s.POPULATION_SIZE)
	for i := 0; i < s.POPULATION_SIZE; i++ {
		new_generation[i] = chromosomes[i]
	}
	return new_generation
}

func argsort(chromosomes []Chromosome) []Chromosome {
	var argsort_chromosomes []Chromosome = make([]Chromosome, len(chromosomes))
	var index []int = make([]int, len(chromosomes))
	for i := 0; i < len(chromosomes); i++ {
		index[i] = i
	}

	sort.Slice(index, func(i, j int) bool {
		return chromosomes[index[i]].fitness < chromosomes[index[j]].fitness
	})

	for i, idx := range index {
		argsort_chromosomes[i] = chromosomes[idx]
	}

	return argsort_chromosomes
}
