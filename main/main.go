package main

import (
	"main/src"
	"math/rand"
	"time"
)

func test() {
	//src.TestSmallbank(true)
	//src.TestGenerateTransaction(true)
	//src.TestGetACG(true)
	//src.TestBuildConflictGraph(true)
	//src.TestTarjan(true)
	src.TestFindCycles(true)
	//src.TestNezha(true)
	//src.TestFabric(true)
}
func main() {
	rand.Seed(time.Now().UnixNano())
	test()
}
