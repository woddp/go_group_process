# go_group_process
go Concurrent processing of multiple processes

```
	var itme = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "10", "10", "12310"}
  //processNum 协程数量
  processNum:=2;
  // handleNum 处理数量
  handleNum:=len(itme)
	process1 := group_process.New(processNum,handleNum, func(w *sync.WaitGroup, index int) {
		defer w.Done()

		time.Sleep(time.Second * 1)
		fmt.Println(itme[index])

	}, func() {
		//end todo
		fmt.Println("end")
	})

	go process1.Start()

	time.Sleep(time.Second*3)
	process1.Close()
```
