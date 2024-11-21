package particle_swarm_optimization

import (
	"fmt"
	s "k8s_algorithm/tools"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/stat"
)

type Particle struct {
	Velocity       []float64 // length : ITEM_QUANTITY
	Position       []int     // length : ITEM_QUANTITY
	pBest_position []int     // Persional best postion
	fitness        float64   // The sum of the standard deviation of n dimensions
	pBest_fitness  float64   // Persional best fitness
}

var Items_weights [][s.DIMENSION]float64
var Previous_state_of_knapsack [][s.DIMENSION]float64
var Limit_of_knapsack [][s.DIMENSION]float64

func Run(
	items_weights [][s.DIMENSION]float64,
	previous_state_of_knapsack [][s.DIMENSION]float64,
	limit_of_knapsack [][s.DIMENSION]float64,
) (float64, []int, []float64) {

	/*
		Step 1 Update velocity : v = v + c1 * rand(0,1) * (pBest - p) + c2 * rand(0,1) * (gBest - p)
		Step 2 Update position : p = p + v
		Step 3 Update fitness : f = sum(std(items_weights * (previous_state_of_knapsack + p)))
	*/

	Items_weights = items_weights
	Previous_state_of_knapsack = previous_state_of_knapsack
	Limit_of_knapsack = limit_of_knapsack

	var particles []Particle
	var gBest_position []int                  // Global best position
	var gBest_fitness float64                 // Global best fitness
	var gBest_fitness_in_iterations []float64 // Best fitness in all iterations
	particles, gBest_fitness, gBest_position = init_particles()

	for iter := 0; iter < s.PSO_ITERATIONS; iter++ {
		for i := 0; i < s.PARTICLE_QUANTITY; i++ {
			var particle *Particle = &particles[i]
			particle.update_velocity(gBest_position)
			particle.update_position()
			particle.update_fitness_and_pBest()
			if particle.pBest_fitness < gBest_fitness && particle.pBest_fitness > 0 {
				gBest_fitness = particle.pBest_fitness
				copy(gBest_position, particle.pBest_position)
			}
			gBest_fitness_in_iterations = append(gBest_fitness_in_iterations, gBest_fitness)
			// fmt.Printf("%.2f ", particle.fitness)
		}
		fmt.Println("\nIteration", iter, "Best fitness", gBest_fitness)
	}
	return gBest_fitness, gBest_position, gBest_fitness_in_iterations
}

func (particle *Particle) update_fitness_and_pBest() {
	var knapsack_weights [s.DIMENSION][s.KNAPSACK_QUANTITY]float64

	for item_idx, ks_idx := range particle.Position { // Current weight of n dimensions
		for dim := 0; dim < s.DIMENSION; dim++ {
			knapsack_weights[dim][ks_idx] += Items_weights[item_idx][dim]
		}
	}

	// var over_weight bool = false
	var penalty float64 = 0.0
	for dim := 0; dim < s.DIMENSION; dim++ {
		for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ { // Plus previous weight of n dimensions
			knapsack_weights[dim][ks_idx] += Previous_state_of_knapsack[ks_idx][dim]
			// if !is_under_limit(&knapsack_weights) {
			// 	// knapsack_weights[dim][ks_idx] = 1.5 * Limit_of_knapsack[ks_idx][dim]
			// 	over_weight = true
			// }
			knapsack_weights[dim][ks_idx] /= Limit_of_knapsack[ks_idx][dim] // Calculate the percentage of each knapsack in n dimensions
			penalty += knapsack_weights[dim][ks_idx] - 1                    // (Weight / Weight_Limit) - 1 : if Weight > Weight_Limit, then (Weight / Weight_Limit) > 1
		}
	}

	var fitness float64 = 0.0
	for dim := 0; dim < s.DIMENSION; dim++ {
		fitness += stat.StdDev(knapsack_weights[dim][:], nil) // fitness is the sum of the standard deviation of n dimensions
	}

	if is_under_limit(&knapsack_weights) {
		particle.fitness = fitness // Punishment strategy is to make fitness larger
	} else {
		particle.fitness = fitness + penalty*s.LAMBDA
	}

	// Update personal best position and fitness
	if particle.fitness < particle.pBest_fitness {
		particle.pBest_fitness = particle.fitness
		copy(particle.pBest_position, particle.Position)
	}
}

func (particle *Particle) update_velocity(gBest_position []int) { // v
	// w : inertia weight
	// c1, c2 : cognitive and social parameters
	// r1, r2 : random numbers in [0, 1]
	var new_velocity []float64 = make([]float64, s.ITEM_QUANTITY)

	for i := 0; i < s.ITEM_QUANTITY; i++ {
		var w, c1, c2, r1, r2 float64 = s.W, s.C1, s.C2, rand.Float64(), rand.Float64()

		new_velocity[i] = w*particle.Velocity[i] +
			c1*r1*float64(particle.pBest_position[i]-particle.Position[i]) +
			c2*r2*float64(gBest_position[i]-particle.Position[i])
	}
	particle.Velocity = new_velocity
}

func (particle *Particle) update_position() {
	// x = sigmoid(v)
	var position []int = make([]int, s.ITEM_QUANTITY)
	for i := 0; i < s.ITEM_QUANTITY; i++ {
		var x float64 = particle.Velocity[i]
		var y float64 = inverse_function(x) // if x -> 0, then y -> 0, else y -> 1
		if rand.Float64() < y {             // y is the probability of staying in the same position
			position[i] = rand.Intn(s.KNAPSACK_QUANTITY)
		} else {
			position[i] = particle.Position[i]
		}
	}
	particle.Position = position
}

func init_particles() ([]Particle, float64, []int) { // return the number of solutions
	var particles []Particle
	var gBest_position []int = make([]int, s.ITEM_QUANTITY) // Global best position
	var gBest_fitness float64 = math.MaxFloat64
	for i := 0; i < s.PARTICLE_QUANTITY; i++ {
		var position []int = generate_valid_position()
		var velocity []float64 = make([]float64, s.ITEM_QUANTITY)
		for j := 0; j < s.ITEM_QUANTITY; j++ {
			velocity[j] = rand.Float64()*2 - 1 // (-1, 1)
		}

		// var knapsack_weights [s.DIMENSION][s.KNAPSACK_QUANTITY]float64
		// for _, ks_idx := range position { // Current weight of n dimensions
		// 	for dim := 0; dim < s.DIMENSION; dim++ {
		// 		knapsack_weights[dim][ks_idx] += Items_weights[ks_idx][dim] + Previous_state_of_knapsack[ks_idx][dim]
		// 		if knapsack_weights[dim][ks_idx] > Limit_of_knapsack[ks_idx][dim] {
		// 			knapsack_weights[dim][ks_idx] = 1.5 * Limit_of_knapsack[ks_idx][dim]
		// 		}
		// 	}
		// }
		// for is_under_limit(&position) {

		// }

		var particle Particle = Particle{
			Velocity:       velocity,
			Position:       position,
			pBest_position: position,
			fitness:        0.0,
			pBest_fitness:  math.MaxFloat64,
		}
		particle.update_fitness_and_pBest()
		if particle.pBest_fitness < gBest_fitness && particle.pBest_fitness > 0 {
			gBest_fitness = particle.pBest_fitness
			copy(gBest_position, particle.pBest_position)
		}

		particles = append(particles, particle)
	}
	return particles, gBest_fitness, gBest_position
}

// func sigmoid(x float64) float64 {
// 	return 1 / (1 + math.Exp(-x))
// }

func inverse_function(x float64) float64 { // if x -> 0, then y -> 0, else y -> 1
	return 1 / (1 + math.Abs(x))
}

func is_under_limit(knapsack_weights *[s.DIMENSION][s.KNAPSACK_QUANTITY]float64) bool {
	for dim := 0; dim < s.DIMENSION; dim++ {
		for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ {
			if knapsack_weights[dim][ks_idx] > Limit_of_knapsack[ks_idx][dim] {
				return false
			}
		}
	}
	return true
}

func generate_valid_position() []int {
	var position []int = make([]int, s.ITEM_QUANTITY)

Outer_loop:
	for true {
		var knapsack_weights [s.DIMENSION][s.KNAPSACK_QUANTITY]float64

		for i := 0; i < s.ITEM_QUANTITY; i++ {
			var ks_idx int = rand.Intn(s.KNAPSACK_QUANTITY)
			position[i] = ks_idx

			for dim := 0; dim < s.DIMENSION; dim++ {
				knapsack_weights[dim][ks_idx] += Items_weights[i][dim]
			}
		}

		for dim := 0; dim < s.DIMENSION; dim++ {
			for ks_idx := 0; ks_idx < s.KNAPSACK_QUANTITY; ks_idx++ { // Plus previous weight of n dimensions
				if (knapsack_weights[dim][ks_idx] + Previous_state_of_knapsack[ks_idx][dim]) > Limit_of_knapsack[ks_idx][dim] {
					continue Outer_loop
				}
			}
		}

		break
	}

	return position
}
