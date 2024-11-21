package genetic_algorithm

import (
	"fmt"
	s "k8s_algorithm/tools"
	"math"
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

type Data struct {
	Items_weights              [][s.DIMENSION]float64
	Previous_state_of_knapsack [][s.DIMENSION]float64
	Limit_of_knapsack          [][s.DIMENSION]float64
}

func Run(items_weights [][s.DIMENSION]float64,
	previous_state_of_knapsack [][s.DIMENSION]float64,
	limit_of_knapsack [][s.DIMENSION]float64,
) (float64, [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, []float64) {
	data := Data{
		Items_weights:              items_weights,
		Previous_state_of_knapsack: previous_state_of_knapsack,
		Limit_of_knapsack:          limit_of_knapsack,
	}
	matrix_iw := Transform_matrix(data.Items_weights)

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
		var elite_parents [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = Selection(solutions, fitness_of_solutions)
		for i := 0; i < s.ELITE_QUANTITY; i++ {
			for j := i + 1; j < s.ELITE_QUANTITY; j++ {
				var child1, child2 *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = Crossover(&elite_parents[i], &elite_parents[j])
				Mutate(child1)
				Mutate(child2)

				var child1_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(child1, matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
				if Is_under_limit(child1_weights, data.Limit_of_knapsack) {
					solutions = append(solutions, *child1)
					fitness_of_solutions = append(fitness_of_solutions, data.Calculate_fitness(child1_weights))
				}

				var child2_weights [][s.KNAPSACK_QUANTITY]float64 = data.Calculate_knapsack_weights(child2, matrix_iw) // DIMENSION * KNAPSACK_QUANTITY
				if Is_under_limit(child2_weights, data.Limit_of_knapsack) {
					solutions = append(solutions, *child2)
					fitness_of_solutions = append(fitness_of_solutions, data.Calculate_fitness(child2_weights))
				}
			}
		}

		var fitness_sort_in_idx []int = Argsort(fitness_of_solutions)
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
		// cmd := exec.Command("cmd", "/c", "cls")
		// cmd.Stdout = os.Stdout
		// cmd.Run()
		fmt.Printf("Best Fitness: %.2f, Iteration: %d\n", best_fitness, iter)
	}
	return best_fitness, best_fitness_solution, best_fitness_in_iterations
}

func Selection(solutions [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, fitness_of_solutions []float64) [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 {
	var fitness_sort_in_idx []int = Argsort(fitness_of_solutions)
	var elite_parents [][s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64 = make([][4][1000]float64, s.ELITE_QUANTITY)
	for i := 0; i < s.ELITE_QUANTITY; i++ {
		elite_parents[i] = solutions[fitness_sort_in_idx[i]]
	}
	return elite_parents
}

func Crossover(
	parent1 *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64,
	parent2 *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64,
) (*[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64, *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64) {

	var crossover_point = int(distuv.Uniform{Min: 0, Max: float64(s.ITEM_QUANTITY)}.Rand())
	var child1 [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64
	var child2 [s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64
	for i := 0; i < s.KNAPSACK_QUANTITY; i++ {
		for j := 0; j < s.ITEM_QUANTITY; j++ {
			if j < crossover_point {
				child1[i][j] = parent1[i][j]
				child2[i][j] = parent2[i][j]
			} else {
				child1[i][j] = parent2[i][j]
				child2[i][j] = parent1[i][j]
			}
		}
	}
	return &child1, &child2

}

func Mutate(child *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64) {
	for item_idx := 0; item_idx < s.ITEM_QUANTITY; item_idx++ {
		if rand.Float64() < s.MUTATION_RATE {
			var random_knapsack_idx int = int(distuv.Uniform{Min: 0, Max: float64(s.KNAPSACK_QUANTITY)}.Rand())
			for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
				child[ks_idx][item_idx] = 0
			}
			child[random_knapsack_idx][item_idx] = 1
		}
	}
}

func (data *Data) Calculate_knapsack_weights(
	solution *[s.KNAPSACK_QUANTITY][s.ITEM_QUANTITY]float64,
	matrix_iw *mat.Dense,
) [][s.KNAPSACK_QUANTITY]float64 {

	var knapsack_fitness = make([][s.KNAPSACK_QUANTITY]float64, s.DIMENSION) // each knapsack's weight in n dimension
	for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
		chromosome := mat.NewDense(1, s.ITEM_QUANTITY, solution[ks_idx][:]) // to slice then to matrix
		knapsack_previous_weights := mat.NewDense(1, s.DIMENSION, data.Previous_state_of_knapsack[ks_idx][:])
		knapsack_capacity := mat.NewDense(1, s.DIMENSION, nil)

		// matrix multiplication : (1*num_items) x (num_items*dimensions) --> (1*dimensions) --> displacement of weights of n dim in this knapsack
		knapsack_capacity.Mul(chromosome, matrix_iw)
		knapsack_capacity.Add(knapsack_capacity, knapsack_previous_weights)
		for dim_idx := 0; dim_idx < s.DIMENSION; dim_idx++ {
			knapsack_fitness[dim_idx][ks_idx] = knapsack_capacity.At(0, dim_idx)
		}
	}
	return knapsack_fitness

}

func (data *Data) Calculate_fitness(knapsack_weights [][s.KNAPSACK_QUANTITY]float64) float64 { // DIMENSION * KNAPSACK_QUANTITY
	var fitness float64 = 0
	for i := 0; i < len(knapsack_weights); i++ {
		for j := 0; j < len(knapsack_weights[0]); j++ {
			knapsack_weights[i][j] /= data.Limit_of_knapsack[j][i] //////////////////
		}
	}
	for _, weights_percestage := range knapsack_weights {
		// fmt.Println(weights_percestage)
		fitness += stat.StdDev(weights_percestage[:], nil)
	}
	return fitness
}

func Transform_matrix(ary [][s.DIMENSION]float64) *mat.Dense {
	row := len(ary)
	col := len(ary[0])
	var data []float64
	for _, i := range ary {
		for _, j := range i {
			data = append(data, j)
		}
	}
	// transpose := mat.NewDense(col, row, nil)
	// transpose.Copy(mat.Transpose{Matrix: })
	return mat.NewDense(row, col, data)
}

func Argsort(ary []float64) []int {
	var results []int = make([]int, len(ary))
	for i := 0; i < len(ary); i++ {
		results[i] = i
	}
	sort.Slice(results, func(i, j int) bool {
		return ary[results[i]] < ary[results[j]] // increasing order
	})
	return results
}

func Is_under_limit(knapsack_weights [][s.KNAPSACK_QUANTITY]float64, limit_of_knapsack [][s.DIMENSION]float64) bool {
	for i := 0; i < len(knapsack_weights); i++ {
		for j := 0; j < s.KNAPSACK_QUANTITY; j++ {
			if knapsack_weights[i][j] > limit_of_knapsack[j][i] {
				return false
			}
		}
	}
	return true
}
