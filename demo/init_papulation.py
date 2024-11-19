import numpy as np


class Init_population:
    def __init__(self, dim, item_num, values: np.array, max_weight: np.array, knapsack_num: int, old_value: np.array,
                 particle_num: int):
        self.values = values  # p的transpose矩陣
        self.max_weight = max_weight
        self.knapsack_num = knapsack_num
        self.old_value = old_value
        self.item_num = item_num
        self.particle_num = particle_num
        self.dim = dim

    def init_population(self, particle_num: int) -> list:
        if int(particle_num) == 0:
            return []
        return [np.random.randint(self.knapsack_num, size=self.item_num) for _ in range(particle_num)]
        # 每個物品在哪個背包 假設5個物品 3個背包 -->  [背包1 背包0 背包0 背包1 背包2]
        ''' [
                array([1, 2, 0, 1, 0]),  # 第 1 個粒子的分配方案
                array([2, 0, 1, 1, 2])   # 第 2 個粒子的分配方案
            ] '''

    def init_population_dim(self, population: list, particle_num: int) -> list:
        # population 就是物品在哪個背包  [背包1 背包0 背包3 背包3 背包2]
        # particle_num 幾組解
        if int(particle_num) == 0:
            return []
        sol = [np.zeros(shape=[self.knapsack_num, self.item_num]) for _ in range(particle_num)]

        ''' [
            array([[0., 0., 0., 0., 0.],  # 第 1 個粒子的零矩陣
                   [0., 0., 0., 0., 0.],
                   [0., 0., 0., 0., 0.]]),
            array([[0., 0., 0., 0., 0.],  # 第 2 個粒子的零矩陣
                   [0., 0., 0., 0., 0.],
                   [0., 0., 0., 0., 0.]])
            ] '''

        for chromosome in range(len(population)):
            item = 0
            for knapsack in population[chromosome]:
                sol[chromosome][knapsack, item] = 1 # 第item個物品 有沒有在 第knapsack個背包  有就是1
                item += 1
        return sol


    ''''''''''''''''''
    def fitness_value(self, solution: list, knapsack: int, dim: int):
        total_value = np.sum(np.array(self.values[dim]) * np.array(solution)) + self.old_value[dim][knapsack]
        return total_value
    




    

    def fix_population(self, population: list) -> list:
        ''' 這時候的population : 變成01 每個array的 row代表有無物品 column代表第幾個knapsack '''
        ''' 反正一橫行就是一個背包 '''
        ''' [
            array([[0., 0., 0., 0., 0.],   # 第 1 個背包
                   [0., 0., 0., 0., 0.],   # 第 2 個背包
                   [0., 0., 0., 0., 0.]]), # 第 3 個背包
            array([[0., 0., 0., 0., 0.],   # 第 1 個背包
                   [0., 0., 0., 0., 0.],   # 第 2 個背包
                   [0., 0., 0., 0., 0.]])  # 第 3 個背包
            ] '''
        while True:
            is_remove = False
            index = 0
            while True:
                if index >= len(population):
                    break
                for dim in range(self.dim):
                    for each_knapsack in range(self.knapsack_num):
                        if self.fitness_value(population[index][each_knapsack], each_knapsack, dim) > \
                                self.max_weight[dim][each_knapsack]:
                            del population[index]
                            is_remove = True
                            break
                        else:
                            is_remove = False
                    if is_remove:
                        break

                if not is_remove:
                    index += 1

            new_sol = self.init_population_dim(self.init_population(self.particle_num - len(population)),
                                               self.particle_num - len(population))

            population = population + new_sol

            if len(new_sol) == 0:
                break

        return population
