package artificial_bee_colony

import (
	"fmt"
	s "k8s_algorithm/tools"
	"math"
	"math/rand"
)

type FoodSource struct {
	population []int
	fitness    float64
	trials     int
}

var Items_weights [][s.DIMENSION]float64
var Previous_state_of_knapsack [][s.DIMENSION]float64
var Limit_of_knapsack [][s.DIMENSION]float64

var best_fitness float64 = math.MaxFloat64
var best_fitness_solution []int

func Run(
	items_weights [][s.DIMENSION]float64,
	previous_state_of_knapsack [][s.DIMENSION]float64,
	limit_of_knapsack [][s.DIMENSION]float64,
) (float64, []int, []float64) {

	Items_weights = items_weights
	Previous_state_of_knapsack = previous_state_of_knapsack
	Limit_of_knapsack = limit_of_knapsack

	best_fitness = math.MaxFloat64
	best_fitness_solution = nil
	var best_fitness_in_iterations []float64

	var food_sources []FoodSource = initialize_food_sources()

	for iter := 0; iter < s.ABC_MAX_ITERATIONS; iter++ {
		employed_bees(food_sources)
		on_looker_bees(food_sources)
		scout_bees(food_sources)

		best_fitness_in_iterations = append(best_fitness_in_iterations, best_fitness)
		fmt.Printf("Best Fitness: %.10f, Iteration: %d\n", best_fitness, iter)
	}

	return best_fitness, best_fitness_solution, best_fitness_in_iterations
}

func employed_bees(food_sources []FoodSource) {
	for i := 0; i < s.EMPLOYED_BEE_QUANTITY; i++ {
		var neighbor_food_source = generate_neighbor_food_source(food_sources[i])
		if neighbor_food_source.fitness < food_sources[i].fitness { // neighbor is better
			food_sources[i] = neighbor_food_source
		} else {
			food_sources[i].trials++
		}

		if food_sources[i].fitness < best_fitness {
			best_fitness = food_sources[i].fitness
			best_fitness_solution = food_sources[i].population
		}
	}
}

func on_looker_bees(food_sources []FoodSource) {
	var fitness_reciprocal [s.COLONY_SIZE]float64
	var sum_fitness_reciprocal float64 = 0.0
	for i, food_sorce := range food_sources {
		fitness_reciprocal[i] += 1 / food_sorce.fitness
		sum_fitness_reciprocal += fitness_reciprocal[i]
	}

	var probabilities []float64
	for _, fr := range fitness_reciprocal {
		probabilities = append(probabilities, fr/sum_fitness_reciprocal)
	}

	for i := 0; i < s.ONLOOKER_BEE_QUANTITY; i++ {
		var selected_food_source_idx = roulette_wheel_selection(probabilities)
		var neighbor_food_source = generate_neighbor_food_source(food_sources[selected_food_source_idx])
		if neighbor_food_source.fitness < food_sources[selected_food_source_idx].fitness { // neighbor is better
			food_sources[selected_food_source_idx] = neighbor_food_source
		} else {
			food_sources[selected_food_source_idx].trials++
		}

		if food_sources[i].fitness < best_fitness {
			best_fitness = food_sources[i].fitness
			best_fitness_solution = food_sources[i].population
		}
	}
}

func scout_bees(food_sources []FoodSource) {
	for i := 0; i < len(food_sources); i++ {
		if food_sources[i].trials >= s.SCOUT_LIMIT {
			population, fitness := generate_uniform_valid_population_fitness()
			food_sources[i] = FoodSource{population, fitness, 0}
		}

		if food_sources[i].fitness < best_fitness {
			best_fitness = food_sources[i].fitness
			best_fitness_solution = food_sources[i].population
		}
	}
}

func generate_neighbor_food_source(food_source FoodSource) FoodSource {
	var population []int = make([]int, s.ITEM_QUANTITY)
	copy(population, food_source.population)
	var rand_idx = rand.Intn(s.ITEM_QUANTITY)
	var original_ks_idx = population[rand_idx]
	for population[rand_idx] == original_ks_idx {
		population[rand_idx] = rand.Intn(s.KNAPSACK_QUANTITY)
	}

	var knapsack_dim_weights *[s.KNAPSACK_QUANTITY][s.DIMENSION]float64 = calculate_weights(population)
	return FoodSource{population, calculate_fitness(knapsack_dim_weights), 0}
}

func roulette_wheel_selection(probabilities []float64) int {
	var rand_num = rand.Float64()
	var cumulative_probability float64 = 0.0
	for i := 0; i < len(probabilities); i++ {
		cumulative_probability += probabilities[i]
		if cumulative_probability > rand_num {
			return i
		}
	}
	return len(probabilities) - 1
}
