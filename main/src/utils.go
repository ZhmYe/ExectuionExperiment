package src

import (
	"math/rand"
	"time"
)

func getDegree(DAG [][]int, index int) int {
	// i->j DAG[i][j] = 1
	degree := 0
	for i := 0; i < len(DAG); i++ {
		if DAG[i][index] == 1 && i != index {
			degree += 1
		}
	}
	return degree
}
func TopologicalOrder(DAG [][]int) []int {
	sortResult := make([]int, 0)
	for k := 0; k < len(DAG); k++ {
		for i := 0; i < len(DAG); i++ {
			if getDegree(DAG, i) == 0 {
				// 取出度数为0的点加入到结果中，并将它连的所有点取消连接
				sortResult = append(sortResult, i)
				for j := 0; j < len(DAG); j++ {
					DAG[i][j] = 0
				}
				break
			}
		}
	}
	return sortResult
}

func generateRandomAddress() string {
	n := 16
	var letters = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)

}

func generateRandomTxhash() string {
	n := 16
	var letters = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)

}
func checkInPath(n int, path []int) bool {
	for _, p := range path {
		if p == n {
			return true
		}
	}
	return false
}
func findCycle(graph Graph, target int, index int, path []int, result *[][]int) {
	//fmt.Print(index)
	//fmt.Print(" ")
	if graph[index][target] == 1 {
		tmp := sort(append(path, index))
		exist := false
		for _, r := range *result {
			if len(r) == len(tmp) {
				same := true
				for k, _ := range r {
					if r[k] != tmp[k] {
						same = false
						break
					}
				}
				exist = same
			}
			if exist {
				break
			}
		}
		if !exist {
			*result = append(*result, tmp)
		}
		//fmt.Println(append(path, index))
	} else {
		for i, _ := range graph[index] {
			if graph[index][i] == 1 && !checkInPath(i, path) {
				findCycle(graph, target, i, append(path, index), result)
			}
		}
	}
}
func findCycles(graph Graph) [][]int {
	results := make([][]int, 0)
	for i, _ := range graph {
		findCycle(graph, i, i, *new([]int), &results)
	}
	return results
}
func sort(a []int) []int {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
	return a
}
func getMaxFromCounter(m map[int]int) int {
	maxCount := 0
	maxid := -1
	for txid, count := range m {
		if count > maxCount {
			maxCount = count
			maxid = txid
		}
	}
	return maxid
}
func checkStillCycle(m map[int]bool) bool {
	for _, flag := range m {
		if !flag {
			return false
		}
	}
	return true
}
