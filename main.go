package main

import (
	"fmt"
	GA "k8s_algorithm/genetic_algorithm"

	// PSO "k8s_algorithm/particle_swarm_optimization"
	// HGA "k8s_algorithm/hybrid_genetic_algorithm"
	"k8s_algorithm/tools"
)

func main() {
	var items_weights [][tools.DIMENSION]float64 = tools.Read_data(tools.ITEMS_WEIGHT_FILE)
	var previous_state_of_knapsack [][tools.DIMENSION]float64 = tools.Read_data(tools.PREVIOUS_WEIGHT_FILE)
	var capacity_of_knapsack [][tools.DIMENSION]float64 = tools.Read_data(tools.WEIGHT_LIMIT_FILE)

	_, _, best_fitness_in_iterations := GA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	fmt.Println(best_fitness_in_iterations)

	// best_fitness_in_iterations, _, _ := PSO.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	// fmt.Println(best_fitness_in_iterations)

	// _, _, best_fitness_in_iterations := HGA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	// fmt.Println(best_fitness_in_iterations)

}

//
