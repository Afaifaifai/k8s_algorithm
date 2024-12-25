package main

import (
	ABC "k8s_algorithm/artificial_bee_colony"
	HGA "k8s_algorithm/hybrid_genetic_algorithm"
	GA "k8s_algorithm/new_genetic_algorithm"
	PSO "k8s_algorithm/particle_swarm_optimization"
	"time"

	// NSGA "k8s_algorithm/nondominated_sorting_genetic_algorithm"

	// tools "k8s_algorithm/tools"

	"flag"
	"fmt"
	"k8s_algorithm/tools"
)

func main() {
	// tools.Transform_1dim()
	items_weight_file := flag.String("items", tools.ITEMS_WEIGHT_FILE, "file path for item_weight")
	previous_weight_file := flag.String("previous", tools.PREVIOUS_WEIGHT_FILE, "file path for previous weights of knapsacks")
	weight_limit_file := flag.String("limit", tools.WEIGHT_LIMIT_FILE, "file path for weight limit of knapsacks")
	flag.Parse()
	// fmt.Println(*items_weight_file, *previous_weight_file, *weight_limit_file)
	if *items_weight_file == "" || *previous_weight_file == "" || *weight_limit_file == "" {
		fmt.Println("Error: Missing required file paths.")
		flag.Usage() // 顯示使用方法
		return
	}

	var items_weights [][]float64 = tools.Read_data(*items_weight_file)
	var previous_state_of_knapsack [][]float64 = tools.Read_data(*previous_weight_file)
	var capacity_of_knapsack [][]float64 = tools.Read_data(*weight_limit_file)

	var dim, item_quantity, knapsack_quantity int = len(items_weights[0]), len(items_weights), len(capacity_of_knapsack)
	tools.Setup(dim, item_quantity, knapsack_quantity)
	// fmt.Println(tools.DIMENSION, tools.ITEM_QUANTITY, knapsack_quantity)

	// GA
	fmt.Println("\nGA is running ...")
	start := time.Now()
	ga_best_fitness, ga_best_fitness_solution, _ := GA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	elapsed := time.Since(start)
	fmt.Printf("GA Best fitness %f\n", ga_best_fitness)
	fmt.Printf("GA Best solution %v\n", ga_best_fitness_solution)
	fmt.Printf("GA execution time : %.6f ms\n\n", elapsed.Seconds()*1000)

	// PSO
	fmt.Println("\nPSO is running ...")
	start = time.Now()
	pso_best_fitness, pso_best_fitness_solution, _ := PSO.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	elapsed = time.Since(start)
	fmt.Printf("PSO Best fitness %f\n", pso_best_fitness)
	fmt.Printf("PSO Best solution %v\n", pso_best_fitness_solution)
	fmt.Printf("PSO execution time : %.6f ms\n\n", elapsed.Seconds()*1000)

	// HGA
	fmt.Println("\nHGA is running ...")
	start = time.Now()
	hga_best_fitness, hga_best_fitness_solution, _ := HGA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	elapsed = time.Since(start)
	fmt.Printf("HGA Best fitness %f\n", hga_best_fitness)
	fmt.Printf("HGA Best solution %v\n", hga_best_fitness_solution)
	fmt.Printf("HGA execution time : %.6f ms\n\n", elapsed.Seconds()*1000)

	// ABC
	fmt.Println("\nABC is running ...")
	start = time.Now()
	abc_best_fitness, abc_best_fitness_solution, _ := ABC.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	elapsed = time.Since(start)
	fmt.Printf("ABC execution time : %.6f ms\n\n", elapsed.Seconds()*1000)
	fmt.Printf("ABC Best fitness %f\n", abc_best_fitness)
	fmt.Printf("ABC Best solution %v\n", abc_best_fitness_solution)

	// start = time.Now()
	// NSGA.Run(items_weights, previous_state_of_knapsack, capacity_of_knapsack)
	// elapsed = time.Since(start)
	// fmt.Printf("NSGA execution time : %.6f ms\n\n", elapsed.Seconds()*1000)

	tools.Write_data([][]int{ga_best_fitness_solution, pso_best_fitness_solution, hga_best_fitness_solution, abc_best_fitness_solution}, []string{"GA", "PSO", "HGA", "ABC"})
}
