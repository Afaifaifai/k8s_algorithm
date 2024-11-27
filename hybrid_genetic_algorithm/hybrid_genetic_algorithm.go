package hybrid_genetic_algorithm

// import (
// 	"fmt"
// 	ga "k8s_algorithm/genetic_algorithm"
// 	s "k8s_algorithm/tools"
// 	"math"
// )

// var Items_weights [][s.DIMENSION]float64
// var Previous_state_of_knapsack [][s.DIMENSION]float64
// var Limit_of_knapsack [][s.DIMENSION]float64

// func Run(items_weights [][s.DIMENSION]float64,
// 	previous_state_of_knapsack [][s.DIMENSION]float64,
// 	limit_of_knapsack [][s.DIMENSION]float64,
// ) (float64, [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, []float64) {

// 	data := ga.Data{
// 		Items_weights:              items_weights,
// 		Previous_state_of_knapsack: previous_state_of_knapsack,
// 		Limit_of_knapsack:          limit_of_knapsack,
// 	}

// 	Items_weights = items_weights
// 	Previous_state_of_knapsack = previous_state_of_knapsack
// 	Limit_of_knapsack = limit_of_knapsack

// 	matrix_iw := ga.Transform_matrix(data.Items_weights)

// 	var best_fitness float64 = math.MaxFloat64
// 	var best_fitness_solution [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64
// 	var best_fitness_in_iterations []float64

// 	var solutions [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = data.Generate_solutions()
// 	for iter := 0; iter < s.GA_ITERATIONS; iter++ {
// 		// The sum of the weights of the items in each knapsack is used to measure the value of a solution
// 		// fmt.Println(len(solutions))
// 		var fitness_of_solutions []float64 = make([]float64, s.SOLUTION_SIZE)
// 		for i := 0; i < s.SOLUTION_SIZE; i++ {
// 			var knapsack_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(&solutions[i], matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
// 			fitness_of_solutions[i] = data.Calculate_fitness(knapsack_weights)
// 		}

// 		// Select  -->  Crossover  -->  Mutate
// 		var elite_parents [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = ga.Selection(solutions, fitness_of_solutions)
// 		for i := 0; i < s.ELITE_QUANTITY; i++ {
// 			for j := i + 1; j < s.ELITE_QUANTITY; j++ {
// 				var child1, child2 *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = ga.Crossover(&elite_parents[i], &elite_parents[j])
// 				ga.Mutate(child1)
// 				ga.Mutate(child2)

// 				var child1_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(child1, matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
// 				if ga.Is_under_limit(child1_weights, data.Limit_of_knapsack) {
// 					solutions = append(solutions, *child1)
// 					fitness_of_solutions = append(fitness_of_solutions, data.Calculate_fitness(child1_weights))
// 				}

// 				var child2_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(child2, matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
// 				if ga.Is_under_limit(child2_weights, data.Limit_of_knapsack) {
// 					solutions = append(solutions, *child2)
// 					fitness_of_solutions = append(fitness_of_solutions, data.Calculate_fitness(child2_weights))
// 				}
// 			}
// 		}

// 		var fitness_sort_in_idx []int = ga.Argsort(fitness_of_solutions)
// 		for idx := 0; idx < s.SOLUTION_SIZE; idx++ {
// 			if idx%10 == 0 {
// 				fitness := Simulated_Annealing(&solutions[fitness_sort_in_idx[idx]])
// 				fitness_of_solutions[fitness_sort_in_idx[idx]] = fitness
// 			}
// 		}

// 		fitness_sort_in_idx = ga.Argsort(fitness_of_solutions)
// 		var new_gen_solutions [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = make([][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, s.SOLUTION_SIZE)
// 		for idx := 0; idx < s.SOLUTION_SIZE; idx++ {
// 			new_gen_solutions[idx] = solutions[fitness_sort_in_idx[idx]]
// 		}

// 		if fitness_of_solutions[fitness_sort_in_idx[0]] < best_fitness {
// 			best_fitness = fitness_of_solutions[fitness_sort_in_idx[0]]
// 			best_fitness_solution = solutions[fitness_sort_in_idx[0]]
// 		}
// 		best_fitness_in_iterations = append(best_fitness_in_iterations, best_fitness)
// 		solutions = new_gen_solutions

// 		fmt.Printf("Best Fitness: %.2f, Iteration: %d\n", best_fitness, iter)
// 	}
// 	return best_fitness, best_fitness_solution, best_fitness_in_iterations
// }

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
	fmt.Println("Hybrid Genetic Algorithm")
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

		chromosomes = argsort(chromosomes)
		for i := 0; i < s.SA_SOLUTION_QUANTITY; i++ {
			genes, fitness := Simulated_Annealing(chromosomes[i].genes, chromosomes[i].fitness)
			chromosomes[i] = Chromosome{genes, fitness}
		}

		chromosomes = replacement(chromosomes)

		if chromosomes[0].fitness < best_fitness {
			best_fitness = chromosomes[0].fitness
			best_fitness_solution = chromosomes[0].genes
		}
		best_fitness_in_iterations = append(best_fitness_in_iterations, chromosomes[0].fitness)
		fmt.Printf("Best Fitness: %.2f, Iteration: %d\n", best_fitness, iter)
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
