import numpy as np


class Data_factory:
    def __init__(self, data_num: int, data_percent: list, func, txt_file_name):
        # 原始數據
        a = func
        self.original_data = a.answer()[0] # Max weights
        self.data_num = data_num
        self.data_percent = data_percent
        self.dim = len(self.original_data)
        self.txt_file_name = txt_file_name

    def get_data(self):
        data = []
        for dim in range(self.dim): # self.dim = len(weights) = 3
            # 計算原始數據的總和
            total_sum = int(sum(self.original_data[dim])) * self.data_percent[dim]
            # 方法 2: 隨機分配
            random_data = np.random.rand(self.data_num)
            # print(random_data, end="\t")
            random_data /= random_data.sum()  # 正規化為 1
            # print(random_data, end="\t")
            random_data *= total_sum  # 調整數據使總和等於原始總和
            # print(random_data)
            data.append(random_data)

        data = np.array(data)
        data = data.T # 每個row 

        # 使用 open() 寫入 txt 檔案，去除括號
        with open(self.txt_file_name, "w") as file:
            for i, value in enumerate(data):
                line = ' '.join(map(str, value))  # 將數據轉換為以空格分隔的字串
                if i < len(data) - 1:
                    file.write(f"{line}\n")  # 除最後一行外，其他行換行
                else:
                    file.write(f"{line}")    # 最後一行無換行
