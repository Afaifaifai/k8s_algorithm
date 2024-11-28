package main

import (
	// GA "k8s_algorithm/new_genetic_algorithm"
	// HGA "k8s_algorithm/hybrid_genetic_algorithm"
	// NSGA "k8s_algorithm/nondominated_sorting_genetic_algorithm"
	"fmt"
	ABC "k8s_algorithm/artificial_bee_colony"
	"k8s_algorithm/tools"
	// "k8s_algorithm/tools"
)

func main() {
	var items_weights [][tools.DIMENSION]float64 = tools.Read_data(tools.ITEMS_WEIGHT_FILE)
	var previous_state_of_knapsack [][tools.DIMENSION]float64 = tools.Read_data(tools.PREVIOUS_WEIGHT_FILE)
	var capacity_of_knapsack [][tools.DIMENSION]float64 = tools.Read_data(tools.WEIGHT_LIMIT_FILE)

	// best_fitness, _, _ := GA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	// fmt.Println(best_fitness)

	// best_fitness_in_iterations, _, _ := PSO.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	// fmt.Println(best_fitness_in_iterations)

	// _, _, best_fitness_in_iterations := HGA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	// fmt.Println(best_fitness_in_iterations)

	// NSGA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)

	_, _, best_fitness_in_iterations := ABC.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	fmt.Println(best_fitness_in_iterations)
}

//
