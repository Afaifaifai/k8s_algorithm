package hybrid_genetic_algorithm

import (
	"fmt"
	ga "k8s_algorithm/genetic_algorithm"
	s "k8s_algorithm/tools"
	"math"
)

var Items_weights [][s.DIMENSION]float64
var Previous_state_of_knapsack [][s.DIMENSION]float64
var Limit_of_knapsack [][s.DIMENSION]float64

func Run(items_weights [][s.DIMENSION]float64,
	previous_state_of_knapsack [][s.DIMENSION]float64,
	limit_of_knapsack [][s.DIMENSION]float64,
) (float64, [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, []float64) {

	data := ga.Data{
		Items_weights:              items_weights,
		Previous_state_of_knapsack: previous_state_of_knapsack,
		Limit_of_knapsack:          limit_of_knapsack,
	}

	Items_weights = items_weights
	Previous_state_of_knapsack = previous_state_of_knapsack
	Limit_of_knapsack = limit_of_knapsack

	matrix_iw := ga.Transform_matrix(data.Items_weights)

	var best_fitness float64 = math.MaxFloat64
	var best_fitness_solution [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64
	var best_fitness_in_iterations []float64

	var solutions [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = data.Generate_solutions()
	for iter := 0; iter < s.GA_ITERATIONS; iter++ {
		// The sum of the weights of the items in each knapsack is used to measure the value of a solution
		// fmt.Println(len(solutions))
		var fitness_of_solutions []float64 = make([]float64, s.SOLUTION_SIZE)
		for i := 0; i < s.SOLUTION_SIZE; i++ {
			var knapsack_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(&solutions[i], matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
			fitness_of_solutions[i] = data.Calculate_fitness(knapsack_weights)
		}

		// Select  -->  Crossover  -->  Mutate
		var elite_parents [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = ga.Selection(solutions, fitness_of_solutions)
		for i := 0; i < s.ELITE_QUANTITY; i++ {
			for j := i + 1; j < s.ELITE_QUANTITY; j++ {
				var child1, child2 *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = ga.Crossover(&elite_parents[i], &elite_parents[j])
				ga.Mutate(child1)
				ga.Mutate(child2)

				var child1_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(child1, matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
				if ga.Is_under_limit(child1_weights, data.Limit_of_knapsack) {
					solutions = append(solutions, *child1)
					fitness_of_solutions = append(fitness_of_solutions, data.Calculate_fitness(child1_weights))
				}

				var child2_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(child2, matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
				if ga.Is_under_limit(child2_weights, data.Limit_of_knapsack) {
					solutions = append(solutions, *child2)
					fitness_of_solutions = append(fitness_of_solutions, data.Calculate_fitness(child2_weights))
				}
			}
		}

		var fitness_sort_in_idx []int = ga.Argsort(fitness_of_solutions)
		for idx := 0; idx < s.SOLUTION_SIZE; idx++ {
			if idx%10 == 0 {
				fitness := Simulated_Annealing(&solutions[fitness_sort_in_idx[idx]])
				fitness_of_solutions[fitness_sort_in_idx[idx]] = fitness
			}
		}

		fitness_sort_in_idx = ga.Argsort(fitness_of_solutions)
		var new_gen_solutions [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = make([][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, s.SOLUTION_SIZE)
		for idx := 0; idx < s.SOLUTION_SIZE; idx++ {
			new_gen_solutions[idx] = solutions[fitness_sort_in_idx[idx]]
		}

		if fitness_of_solutions[fitness_sort_in_idx[0]] < best_fitness {
			best_fitness = fitness_of_solutions[fitness_sort_in_idx[0]]
			best_fitness_solution = solutions[fitness_sort_in_idx[0]]
		}
		best_fitness_in_iterations = append(best_fitness_in_iterations, best_fitness)
		solutions = new_gen_solutions

		fmt.Printf("Best Fitness: %.2f, Iteration: %d\n", best_fitness, iter)
	}
	return best_fitness, best_fitness_solution, best_fitness_in_iterations
}
