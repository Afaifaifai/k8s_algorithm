package tools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Read_data(filename string) [][DIMENSION]float64 {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Cannot locate the file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var results [][DIMENSION]float64

	for scanner.Scan() {
		line := scanner.Text()

		values := strings.Split(line, " ")
		var nums_dim [DIMENSION]float64
		for dim, value := range values {
			// fmt.Println(value)
			var num float64
			num, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err != nil {
				fmt.Println("Invalid format of float64:", err)
				continue
			}
			nums_dim[dim] = num
		}
		results = append(results, nums_dim)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Read Error:", err)
	}
	// results := mat.NewDense(s.DIMENSION, len(nums)/s.DIMENSION, nil) // row : DIM 		col : len / DIM
	// results.Copy(&mat.Transpose{Matrix: mat.NewDense(len(nums)/s.DIMENSION, s.DIMENSION, nums)})
	// results := mat.NewDense(len(nums)/s.DIMENSION, s.DIMENSION, nums)
	/*
				[[dim1]
		result=	 [dim2]
			 	 [dim3]]
	*/
	return results
}
