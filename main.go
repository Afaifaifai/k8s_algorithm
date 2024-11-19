package main

import (
	"fmt"
	"k8s_algorithm/genetic_algorithm"
	"k8s_algorithm/settings"
)

func main() {
	var items_weights [][settings.DIMENSION]float64 = read_data(settings.ITEMS_WEIGHT_FILE)
	var previous_state_of_knapsack [][settings.DIMENSION]float64 = read_data(settings.PREVIOUS_WEIGHT_FILE)
	var capacity_of_knapsack [][settings.DIMENSION]float64 = read_data(settings.WEIGHT_LIMIT_FILE)

	_, _, best_fitness_in_iterations := genetic_algorithm.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	fmt.Println(best_fitness_in_iterations)

}

//
