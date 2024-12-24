package tools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Read_data(filename string) [][]float64 {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Cannot locate the file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var results [][]float64

	for scanner.Scan() {
		line := scanner.Text()

		values := strings.Split(line, " ")
		var nums_dim []float64
		for _, value := range values {
			value = strings.TrimSuffix(value, "m")

			var num float64
			num, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err != nil {
				fmt.Println("Invalid format of float64:", err)
				continue
			}
			nums_dim = append(nums_dim, num)
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

func Write_data(record [][]int, order []string) {
	// 1. 將每個切片轉換為字串，並在每行前加上行號
	mkdir(RESULT_DIR)
	var lines []string
	for i, row := range record {
		// 將每個 row 轉為字串，例如 [1, 2, 3] => "1, 2, 3"
		line := strings.Trim(strings.Replace(fmt.Sprint(row), " ", ", ", -1), "[]")
		lines = append(lines, fmt.Sprintf("%s : %s", order[i], line))
	}

	// 合併所有行為單一字串，用換行符號分隔
	content := strings.Join(lines, "\n")

	// 2. 取得當前時間戳
	timestamp := time.Now().Format("20060102_150405") // 格式：YYYYMMDD_HHMMSS

	// 3. 使用時間戳作為檔案名稱
	fileName := fmt.Sprintf("%s%s", RESULT_DIR, timestamp)

	// 4. 將內容寫入檔案
	err := os.WriteFile(fileName, []byte(content), 0644) // 0644 表示可讀寫權限
	if err != nil {
		fmt.Println("無法寫入檔案:", err)
		return
	}

	fmt.Printf("Writing successfully : %s\n", fileName)
}

func Transform_1dim() {
	// 原始檔案名稱
	var input_file []string = []string{"data/values.txt", "data/old_values.txt", "data/knapsack.txt"}
	for _, file_name := range input_file {
		var output_file string = strings.Replace(file_name, ".txt", "_1dim.txt", -1)
		// 讀取原始檔案
		file, err := os.Open(file_name)
		if err != nil {
			fmt.Println("無法開啟檔案:", err)
			return
		}
		defer file.Close()

		// 開啟輸出檔案
		outFile, err := os.Create(output_file)
		if err != nil {
			fmt.Println("無法建立檔案:", err)
			return
		}
		defer outFile.Close()

		// 逐行處理檔案
		scanner := bufio.NewScanner(file)
		writer := bufio.NewWriter(outFile)
		for scanner.Scan() {
			// 取得每一行並分割為欄位
			line := scanner.Text()
			fields := strings.Fields(line) // 將行按空格分割

			// 如果至少有一個欄位，取第一個欄位
			if len(fields) > 0 {
				_, err := writer.WriteString(fields[0] + "\n") // 寫入第一個欄位
				if err != nil {
					fmt.Println("寫入檔案時發生錯誤:", err)
					return
				}
			}
		}

		// 確保所有內容寫入檔案
		writer.Flush()

		// 確認是否成功處理
		if err := scanner.Err(); err != nil {
			fmt.Println("讀取檔案時發生錯誤:", err)
		} else {
			fmt.Printf("處理完成，輸出檔案為: %s\n", output_file)
		}
	}

}

func mkdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 若目錄不存在，創建目錄
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("無法創建目錄: %v\n", err)
			return err
		}
		fmt.Println("目錄已成功創建:", dir)
	} else {
		fmt.Println("目錄已存在:", dir)
	}
	return nil
}
