package group_process

import (
	"context"
	"fmt"
	"math"
	"sync"
)

/*
*
创建多协程任务
*/
type groupProcess struct {
	processNum int // 协程数量
	handleNum  int // 处理数量
	tasks      []int

	cancel context.CancelFunc
	ctx    context.Context

	runCompletedChan chan int   //完成通道
	runStartChan     chan []int //执行通道

	handleFunc HandleFunc
	endFunc    EndFunc
}

// HandleFunc 处理任务
type HandleFunc func(w *sync.WaitGroup, index int)
type EndFunc func()

// New
// processNum 协程数量
// handleNum 处理数量
// handleFunc 处理方法
// endFunc 处理结束
func New(processNum int, handleNum int, handleFunc HandleFunc, endFunc EndFunc) *groupProcess {
	process := &groupProcess{
		processNum: processNum,
		handleNum:  handleNum,
		tasks:      make([]int, 0),
		handleFunc: handleFunc,
		endFunc:    endFunc,
	}
	return process
}

// Start 开始
func (g *groupProcess) Start() {

	g.ctx, g.cancel = context.WithCancel(context.Background())

	//生成任务数量
	for i := 0; i < g.handleNum; i++ {
		g.tasks = append(g.tasks, i)
	}
	g.runCompletedChan = make(chan int, 0)
	g.runStartChan = make(chan []int, 0)
	go g.handle()
	var lists = g.arrayChunk(g.tasks, g.processNum)
	fmt.Println(lists)
	w := &sync.WaitGroup{}
	for i := 0; i < len(lists); i++ {
		w.Add(1)
		var items = make([]int, 0)
		for _, index := range lists[i] {
			items = append(items, index)
		}
		g.runStartChan <- items
		select {
		case <-g.ctx.Done():
			close(g.runStartChan)
			goto BREAK
		case <-g.runCompletedChan:
			w.Done()
		}

	}

BREAK:
	g.endFunc()
	w.Wait()

}

// handle( 处理任务
func (g *groupProcess) handle() {
	for {
		select {
		case <-g.ctx.Done():
			fmt.Println("---------")
			close(g.runCompletedChan)
			return
		case list := <-g.runStartChan:
			w := &sync.WaitGroup{}
			w.Add(len(list))
			for i := 0; i < len(list); i++ {
				go g.handleFunc(w, list[i])
			}
			w.Wait()
			//todo 处理完成
			g.runCompletedChan <- 1
		}
	}
}

// Close 关闭所有协程
func (g *groupProcess) Close() {
	g.cancel()
}
func (g *groupProcess) arrayChunk(arr []int, size int) [][]int {
	if size < 1 {
		panic("[ArrayChunk]`size cannot be less than 1")
	}
	length := len(arr)
	if length == 0 {
		return nil
	}
	chunkNum := int(math.Ceil(float64(length) / float64(size)))
	var res [][]int
	var item []int
	var start int
	for i, end := 0, 0; chunkNum > 0; chunkNum-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}

		item = nil
		start = i * size
		for ; start < end; start++ {
			item = append(item, arr[start])
		}
		if item != nil {
			res = append(res, item)
		}

		i++
	}
	return res
}
