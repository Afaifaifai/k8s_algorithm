package main

const DIMENSION int = 3             // the dimantion of parameters
const WORKER_NODES_QUANTITY int = 4 // the number of worker nodes

const ITEMS_WEIGHT_FILE string = "data/values.txt"        // each item weight (each pod cost)
const PREVIOUS_WEIGHT_FILE string = "data/old_values.txt" // the weight in the previous state of knapsack (previous load of the nodes)
const WEIGHT_LIMIT_FILE string = "data/knapsack.txt"      // the weight limit of knapsack (the capacity of the nodes)
const DATA_SPLIT string = " "                             // the split character that seperate the data in a line

// Genetic algorithm parameters
const KNAPSACK_QUANTITY int = WORKER_NODES_QUANTITY // the quantity of worker nodes is mapped to the quantity of knapsacks
const SOLUTION_SIZE int = 100                       // the number of solutions
const ITEM_QUANTITY int = 1000
const POPULATION_SIZE = ITEM_QUANTITY // the number of solutions
const GA_ITERATIONS int = 100
const MUTATION_RATE float64 = 0.1
const CROSSOVER_RATE float64 = 0.9
const SELECTION_RATE float64 = 0.4


