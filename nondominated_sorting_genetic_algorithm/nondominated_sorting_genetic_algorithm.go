package nondominated_sorting_genetic_algorithm

import (
	"fmt"
	s "k8s_algorithm/tools"
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
) {

	Items_weights = items_weights
	Previous_state_of_knapsack = previous_state_of_knapsack
	Limit_of_knapsack = limit_of_knapsack

	// var best_fitness float64 = math.MaxFloat64
	// var best_fitness_solution []int
	// var best_fitness_in_iterations []float64

	var chromosomes []Chromosome = initialize_chromosomes()
	// chromosomes = argsort(chromosomes) // sort chromosomes by fitness
	chromosomes = nondominated_sorting(chromosomes, s.ELITE_QUANTITY)

	for iter := 0; iter < s.GA_ITERATIONS; iter++ {
		fmt.Println(iter)
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

		// if chromosomes[0].fitness < best_fitness {
		// 	best_fitness = chromosomes[0].fitness
		// 	best_fitness_solution = chromosomes[0].genes
		// }
		// best_fitness_in_iterations = append(best_fitness_in_iterations, chromosomes[0].fitness)
		// fmt.Printf("Best Fitness: %.2f, Iteration: %d\n", best_fitness, iter)
	}

	// return best_fitness, best_fitness_solution, best_fitness_in_iterations
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
			fitness: *calculate_fitness(knapsack_dim_weights),
		}
	} else {
		return Chromosome{
			genes:   genes,
			fitness: s.NSGA_MAX_FITNESSES,
		}
	}
}

func replacement(chromosomes []Chromosome) []Chromosome {
	chromosomes = nondominated_sorting(chromosomes, s.DIMENSION)
	var new_generation []Chromosome = make([]Chromosome, s.POPULATION_SIZE)
	for i := 0; i < s.POPULATION_SIZE; i++ {
		new_generation[i] = chromosomes[i]
	}
	return new_generation
}

// func argsort(chromosomes []Chromosome) []Chromosome {
// 	var argsort_chromosomes []Chromosome = make([]Chromosome, len(chromosomes))
// 	var index []int = make([]int, len(chromosomes))
// 	for i := 0; i < len(chromosomes); i++ {
// 		index[i] = i
// 	}

// 	sort.Slice(index, func(i, j int) bool {
// 		return chromosomes[index[i]].fitness < chromosomes[index[j]].fitness
// 	})

// 	for i, idx := range index {
// 		argsort_chromosomes[i] = chromosomes[idx]
// 	}

// 	return argsort_chromosomes
// }

func nondominated_sorting(chromosomes []Chromosome, selecting_quantity int) []Chromosome {
	var dominated_count []int = make([]int, len(chromosomes))  // The number of times the i-th population is dominated
	var dominate_idx [][]int = make([][]int, len(chromosomes)) // Indices of population's index dominated by the i-th population
	var front [][]int
	front = append(front, []int{})
	for i := 0; i < len(chromosomes); i++ {
		// Calculate the i-th population dominated relation
		// 1. The i-th population is dominated by the j-th population 	--> count the total times being dominated
		// 2. The i-th population dominates the j-th population 		--> record the j-th population index
		for j := 0; j < len(chromosomes); j++ {
			population1, population2 := chromosomes[i].fitness, chromosomes[j].fitness
			if population1[0] <= population2[0] && population1[1] <= population2[1] && population1[2] <= population2[2] {
				if population1[0] < population2[0] || population1[1] < population2[1] || population1[2] < population2[2] {
					dominate_idx[i] = append(dominate_idx[i], j)
				}
			} else if population2[0] <= population1[0] && population2[1] <= population1[1] && population2[2] <= population1[2] {
				if population2[0] < population1[0] || population2[1] < population1[1] || population2[2] < population1[2] {
					dominated_count[i] += 1
				}
			}
		}
		if dominated_count[i] == 0 {
			front[0] = append(front[0], i)
		}
	}

	i := 0
	for len(front[i]) != 0 {
		var next_front []int
		for _, population_idx := range front[i] {
			for _, dominated_idx := range dominate_idx[population_idx] {
				dominated_count[dominated_idx] -= 1
				if dominated_count[dominated_idx] == 0 {
					next_front = append(next_front, dominated_idx)
				}
			}
		}
		i += 1
		front = append(front, next_front)
	}
	front = front[:len(front)-1]

	for i = 0; i < len(front); i++ {
		if selecting_quantity < len(front[i]) {
			front[i] = crowding_distance_sorting(front[i], chromosomes)
		}
		selecting_quantity -= len(front[i])
	}

	var results []Chromosome
	for _, f := range front {
		for _, population_idx := range f {
			results = append(results, chromosomes[population_idx])
		}
	}

	return results
}

type Data struct {
	idx                   int
	population            [s.DIMENSION]float64
	crowding_distance     [s.DIMENSION]float64
	crowding_distance_sum float64
}

func crowding_distance_sorting(population_idxes []int, chromosomes []Chromosome) []int {
	var data []Data
	for _, population_idx := range population_idxes {
		data = append(data, Data{
			idx:                   population_idx,
			population:            chromosomes[population_idx].fitness,
			crowding_distance:     [s.DIMENSION]float64{0, 0, 0},
			crowding_distance_sum: 0.0,
		})
	}

	for dim := 0; dim < s.DIMENSION; dim++ {
		sort.Slice(data, func(i, j int) bool {
			return data[i].population[dim] > data[j].population[dim] // decreasing order
		})

		data[0].crowding_distance[dim] = s.NSGA_MAX_DISTANCE
		data[len(data)-1].crowding_distance[dim] = s.NSGA_MAX_DISTANCE

		var f_min, f_max float64 = data[0].population[dim], data[len(data)-1].population[dim]
		var f_max_min float64 = f_max - f_min

		for i := 1; i < len(data)-1; i++ {
			data[i].crowding_distance[dim] = (data[i+1].population[dim] - data[i-1].population[dim]) / f_max_min
		}

		for i := 0; i < len(data); i++ {
			data[i].crowding_distance_sum += data[i].crowding_distance[dim]
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].crowding_distance_sum > data[j].crowding_distance_sum // decreasing order
	})
	var results []int
	for _, d := range data {
		results = append(results, d.idx)
	}

	return results
}
