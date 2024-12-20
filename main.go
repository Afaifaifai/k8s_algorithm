package main

import (
	ABC "k8s_algorithm/artificial_bee_colony"
	HGA "k8s_algorithm/hybrid_genetic_algorithm"
	GA "k8s_algorithm/new_genetic_algorithm"
	PSO "k8s_algorithm/particle_swarm_optimization"

	// NSGA "k8s_algorithm/nondominated_sorting_genetic_algorithm"
	"fmt"

	"k8s_algorithm/tools"
)

func main() {
	// tools.Transform_1dim()

	var items_weights [][tools.DIMENSION]float64 = tools.Read_data(tools.ITEMS_WEIGHT_FILE)
	var previous_state_of_knapsack [][tools.DIMENSION]float64 = tools.Read_data(tools.PREVIOUS_WEIGHT_FILE)
	var capacity_of_knapsack [][tools.DIMENSION]float64 = tools.Read_data(tools.WEIGHT_LIMIT_FILE)

	// GA
	fmt.Println("\nGA is running ...")
	ga_best_fitness, ga_best_fitness_solution, _ := GA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	fmt.Printf("GA Best fitness %f\n", ga_best_fitness)
	fmt.Printf("GA Best solution %v\n", ga_best_fitness_solution)

	// PSO
	fmt.Println("\nPSO is running ...")
	pso_best_fitness, pso_best_fitness_solution, _ := PSO.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	fmt.Printf("PSO Best fitness %f\n", pso_best_fitness)
	fmt.Printf("PSO Best solution %v\n", pso_best_fitness_solution)

	// HGA
	fmt.Println("\nHGA is running ...")
	hga_best_fitness, hga_best_fitness_solution, _ := HGA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	fmt.Printf("HGA Best fitness %f\n", hga_best_fitness)
	fmt.Printf("HGA Best solution %v\n", hga_best_fitness_solution)

	// ABC
	fmt.Println("\nABC is running ...")
	abc_best_fitness, abc_best_fitness_solution, _ := ABC.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	fmt.Printf("ABC Best fitness %f\n", abc_best_fitness)
	fmt.Printf("ABC Best solution %v\n", abc_best_fitness_solution)

	// NSGA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	tools.Write_data([][]int{ga_best_fitness_solution, pso_best_fitness_solution, hga_best_fitness_solution, abc_best_fitness_solution}, []string{"GA", "PSO", "HGA", "ABC"})
}

//
