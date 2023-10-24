package src

import (
	"math/rand"
	"time"
)

func getDegree(DAG [][]int, index int) int {
	// i->j DAG[i][j] = 1
	degree := 0
	abort := 0
	for i := 0; i < len(DAG); i++ {
		if DAG[i][index] == 1 && i != index {
			degree += 1
		}
		if DAG[index][i] == -1 {
			abort += 1
		}
	}
	if abort == len(DAG) {
		return -1
	}
	return degree
}
func TopologicalOrder(DAG [][]int) []int {
	degrees := make([]int, len(DAG))
	for i, _ := range degrees {
		degrees[i] = getDegree(DAG, i)
	}
	sortResult := make([]int, 0)
	visited := make(map[int]bool, 0)
	for k := 0; k < len(DAG); k++ {
		for i := 0; i < len(DAG); i++ {
			_, flag := visited[i]
			if flag {
				continue
			}
			if degrees[i] == 0 {
				// 取出度数为0的点加入到结果中，并将它连的所有点取消连接
				sortResult = append(sortResult, i)
				visited[i] = true
				for j := 0; j < len(DAG); j++ {
					if DAG[i][j] == 1 {
						degrees[j] -= 1
					}
				}
				break
			}
		}
	}
	//fmt.Println(len(sortResult))
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

func computeCascade(acgs []ACG) (map[string]map[string][]Unit, map[string]int) {
	// nextReadNumberInAddress 统计每笔交易在每个地址后面直接相连的读操作个数
	nextReadNumberInAddress := make(map[string]map[string][]Unit, 0) // map[tx_hash] -> map[address] nextReadSet
	for i, acg := range acgs {
		// 遍历acg中的每个address得到对应的stateset
		// 每个stateset的writeset的最后一个unit的tx_hash -> map[address] -> 下一个stateset中readset的长度
		// 下一个stateset通过继续向后面的acg寻找key得到，如果不存在key则向后寻找
		// 后续需考虑判断tx_hash是否出现在address后
		for address, stateset := range acg {
			writeSet := stateset.WriteSet
			tmpDistance := 0
			if len(writeSet) == 0 {
				continue
			}
			lastElement := writeSet[len(writeSet)-1] // lastElement.tx_hash
			// 可能是在块内并发已经被abort的，要取到没有被abort的最后一个
			tmpFlag := true
			for {
				if lastElement.tx.abort {
					tmpDistance += 1
					if len(writeSet)-1-tmpDistance < 0 {
						tmpFlag = false
						break
					}
					lastElement = writeSet[len(writeSet)-1-tmpDistance]
				} else {
					break
				}
			}
			// 说明该地址下所有的写集都已经被abort，无需进行后续的讨论
			if !tmpFlag {
				continue
			}
			_, inMap := nextReadNumberInAddress[lastElement.tx.txHash]
			if !inMap {
				nextReadNumberInAddress[lastElement.tx.txHash] = make(map[string][]Unit, 0)
			}
			flag := false
			for j := i + 1; j < len(acgs); j++ {
				nextStateSet, exist := acgs[j][address]
				// 如果下一个hashtable里包含了address，那么就记录其读集长度然后结束
				if exist {
					nextReadNumberInAddress[lastElement.tx.txHash][address] = nextStateSet.ReadSet
					flag = true
					break
				}
			}
			// 如果没有后续hashtable
			if !flag {
				nextReadNumberInAddress[lastElement.tx.txHash][address] = make([]Unit, 0)
			}
		}
	}
	//record := nextReadNumberInAddress
	// end for nextReadNumberInAddress
	cascade := make(map[string]int, 0)
	// 计算instance每个address上第一个读集的级联度，所有交易在nextReadNumberInAddress中的所有address相加
	inRecord := make(map[string]bool) // 判断是否是每个address第一个读集
	for _, hashtable := range acgs {
		for address, stateset := range hashtable {
			_, haveRecord := inRecord[address]
			// 是每个address的第一个读集
			if !haveRecord {
				inRecord[address] = true // 标记
				// 如果是第一个读集 才需要更新cascade变量
				cascade[address] = getReadSetNumber(stateset.ReadSet, nextReadNumberInAddress)
			}
		}
	}
	return nextReadNumberInAddress, cascade
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
