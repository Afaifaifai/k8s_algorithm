package tools

const DIMENSION int = 3             // the dimantion of parameters
const WORKER_NODES_QUANTITY int = 4 // the number of worker nodes

const ITEMS_WEIGHT_FILE string = "data/values.txt"        // each item weight (each pod cost)
const PREVIOUS_WEIGHT_FILE string = "data/old_values.txt" // the weight in the previous state of knapsack (previous load of the nodes)
const WEIGHT_LIMIT_FILE string = "data/knapsack.txt"      // the weight limit of knapsack (the capacity of the nodes)
const DATA_SPLIT string = " "                             // the split character that seperate the data in a line

const KNAPSACK_QUANTITY int = WORKER_NODES_QUANTITY // the quantity of worker nodes is mapped to the quantity of knapsacks
const SOLUTION_SIZE int = 100                       // the number of solutions
const ITEM_QUANTITY int = 1000                      // the number of items

// Genetic algorithm parameters
const POPULATION_SIZE = ITEM_QUANTITY // each item pick a  knapsack
const GA_ITERATIONS int = 100
const MUTATION_RATE float64 = 0.1
const CROSSOVER_RATE float64 = 0.9
const SELECTION_RATE float64 = 0.4
const ELITE_QUANTITY = int(float64(SOLUTION_SIZE) * SELECTION_RATE)

// PSO parameters
const PSO_ITERATIONS int = 100000
const PARTICLE_QUANTITY int = 100
const LAMBDA float64 = 5

const W = 0.729
const C1 = 1.49445
const C2 = 1.49445

const EPSILON float64 = 1 / float64(ITEM_QUANTITY) * 5 // basic exploration rate

// Hybrid genetic algorithm parameters
// GA
const HGA_ITERATIONS int = 100
const HGA_MUTATION_RATE float64 = 0.1
const HGA_CROSSOVER_RATE float64 = 0.9
const HGA_SELECTION_RATE float64 = 0.4
const HGA_ELITE_QUANTITY = int(float64(SOLUTION_SIZE) * SELECTION_RATE)

// SA
const SA_PERCENTAGE float64 = 0.15
const SA_SOLUTION_QUANTITY int = int(float64(SOLUTION_SIZE) * SA_PERCENTAGE)
const SA_ITERATIONS int = 100
const INITIAL_TEMPERATURE float64 = 1000
const COOLING_RATE float64 = 0.95
const MINIMUM_TEMPERATURE float64 = 1e-3
