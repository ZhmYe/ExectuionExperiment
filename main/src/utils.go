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

func getReadSetNumber(readSet []Unit, record map[string]map[string][]Unit) int {
	repeatCheck := make(map[string]bool)
	total := 0
	if len(readSet) == 0 {
		return 0
	}
	for _, unit := range readSet {
		_, repeat := repeatCheck[unit.tx.txHash]
		// 如果当前交易有两笔读操作在同一个address或当前交易已经被abort,无需重复计算
		if repeat || unit.tx.abort {
			continue
		}
		total += 1
		repeatCheck[unit.tx.txHash] = true
		CascadeInAddress, haveCascade := record[unit.tx.txHash]
		if haveCascade {
			for _, eachReadSet := range CascadeInAddress {
				total += getReadSetNumber(eachReadSet, record)
			}
		}
	}
	return total
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
