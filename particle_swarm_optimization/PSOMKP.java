import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Random;

// 定義物品類
class Item {
    int id;
    int weight;
    int value;

    public Item(int id, int weight, int value) {
        this.id = id;
        this.weight = weight;
        this.value = value;
    }
}

// 定義背包類
class Knapsack {
    int id;
    int capacity;

    public Knapsack(int id, int capacity) {
        this.id = id;
        this.capacity = capacity;
    }
}

// 定義粒子類
class Particle {
    int[] position; // 分配向量，長度為物品數
    double[] velocity; // 速度向量
    double fitness;
    int[] pBestPosition;
    double pBestFitness;

    public Particle(int numItems) {
        position = new int[numItems];
        velocity = new double[numItems];
        pBestPosition = new int[numItems];
        pBestFitness = Double.NEGATIVE_INFINITY;
    }

    // 深拷貝
    public Particle clone() {
        Particle clone = new Particle(position.length);
        clone.position = Arrays.copyOf(this.position, this.position.length);
        clone.velocity = Arrays.copyOf(this.velocity, this.velocity.length);
        clone.fitness = this.fitness;
        clone.pBestPosition = Arrays.copyOf(this.pBestPosition, this.pBestPosition.length);
        clone.pBestFitness = this.pBestFitness;
        return clone;
    }
}

// 定義 PSO 演算法類
public class PSOMKP {
    // PSO 參數
    int swarmSize = 30;
    int maxIter = 1000;
    double w = 0.729; // 慣性權重
    double c1 = 1.49445; // 認知係數
    double c2 = 1.49445; // 社會係數
    double lambda = 1000; // 懲罰係數

    int numItems;
    int numKnapsacks;
    List<Item> items;
    List<Knapsack> knapsacks;
    List<Particle> swarm;
    Particle gBest;

    Random rand;

    public PSOMKP(List<Item> items, List<Knapsack> knapsacks) {
        this.items = items;
        this.knapsacks = knapsacks;
        this.numItems = items.size();
        this.numKnapsacks = knapsacks.size();
        this.swarm = new ArrayList<>();
        this.gBest = new Particle(numItems).clone();
        this.rand = new Random();
    }

    // 初始化粒子群
    public void initializeSwarm() {
        for (int i = 0; i < swarmSize; i++) {
            Particle p = new Particle(numItems);
            for (int j = 0; j < numItems; j++) {
                p.position[j] = rand.nextInt(numKnapsacks + 1); // 0 ~ numKnapsacks
                p.velocity[j] = rand.nextDouble() * 2 - 1; // -1 ~ 1
            }
            p.fitness = evaluateFitness(p.position);
            p.pBestPosition = Arrays.copyOf(p.position, numItems);
            p.pBestFitness = p.fitness;

            // 更新全局最佳
            if (p.fitness > gBest.pBestFitness) {
                gBest = p.clone();
            }

            swarm.add(p);
        }
    }

    // 計算適應度
    public double evaluateFitness(int[] position) {
        int totalValue = 0;
        int[] knapsackWeights = new int[numKnapsacks];
        int[] itemAssignment = new int[numItems]; // 用於檢查物品是否被多個背包選擇

        for (int i = 0; i < numItems; i++) {
            int assignedKnapsack = position[i];
            if (assignedKnapsack > 0) { // 如果物品被分配到某個背包
                totalValue += items.get(i).value;
                knapsackWeights[assignedKnapsack - 1] += items.get(i).weight;
                itemAssignment[i]++;
            }
        }

        // 計算違反約束的懲罰
        double penalty = 0.0;

        // 背包容量違反
        for (int j = 0; j < numKnapsacks; j++) {
            if (knapsackWeights[j] > knapsacks.get(j).capacity) {
                penalty += (knapsackWeights[j] - knapsacks.get(j).capacity);
            }
        }

        // 物品被多個背包選擇
        for (int i = 0; i < numItems; i++) {
            if (itemAssignment[i] > 1) {
                penalty += (itemAssignment[i] - 1);
            }
        }

        return totalValue - lambda * penalty;
    }

    // 更新速度和位置
    public void updateParticles() {
        for (Particle p : swarm) {
            for (int j = 0; j < numItems; j++) {
                // 更新速度
                double r1 = rand.nextDouble();
                double r2 = rand.nextDouble();
                p.velocity[j] = w * p.velocity[j]
                        + c1 * r1 * (p.pBestPosition[j] - p.position[j])
                        + c2 * r2 * (gBest.pBestPosition[j] - p.position[j]);

                // 使用 Sigmoid 函數將速度轉換為概率
                double sigmoid = 1.0 / (1.0 + Math.exp(-p.velocity[j]));

                // 更新位置
                if (rand.nextDouble() < sigmoid) {
                    // 隨機選擇一個背包或不放入
                    p.position[j] = rand.nextInt(numKnapsacks + 1);
                }
                // 否則，保持原來的分配
            }

            // 計算新的適應度
            double newFitness = evaluateFitness(p.position);

            // 更新個人最佳
            if (newFitness > p.pBestFitness) {
                p.pBestFitness = newFitness;
                p.pBestPosition = Arrays.copyOf(p.position, numItems);
            }

            // 更新全局最佳
            if (newFitness > gBest.pBestFitness) {
                gBest = p.clone();
            }
        }
    }

    // 執行 PSO
    public void runPSO() {
        initializeSwarm();

        for (int iter = 0; iter < maxIter; iter++) {
            updateParticles();

            // 可選：打印當前最佳適應度
            if (iter % 100 == 0 || iter == maxIter - 1) {
                System.out.println("Iteration " + iter + ": Best Fitness = " + gBest.pBestFitness);
            }

            // 可選：提前停止條件，例如連續若干迭代沒有改善
            // 這裡簡化，不實現提前停止
        }

        // 打印最終結果
        System.out.println("\n最終最佳適應度: " + gBest.pBestFitness);
        System.out.println("最佳分配方案:");
        for (int i = 0; i < numItems; i++) {
            System.out.println("物品 " + items.get(i).id + " -> 背包 " + gBest.pBestPosition[i]);
        }
    }

    // 主函數
    public static void main(String[] args) {
        // 定義物品（id, weight, value）
        List<Item> items = new ArrayList<>();
        items.add(new Item(1, 10, 60));
        items.add(new Item(2, 20, 100));
        items.add(new Item(3, 30, 120));
        items.add(new Item(4, 25, 90));
        items.add(new Item(5, 15, 70));

        // 定義背包（id, capacity）
        List<Knapsack> knapsacks = new ArrayList<>();
        knapsacks.add(new Knapsack(1, 50));
        knapsacks.add(new Knapsack(2, 60));

        // 創建 PSO 實例並運行
        PSOMKP pso = new PSOMKP(items, knapsacks);
        pso.runPSO();
    }
}
