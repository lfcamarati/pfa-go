package main

import "time"

func worker(workerId int, msg chan int) {
	for res := range msg {
		println("worker", workerId, "recebeu", res)
		time.Sleep(time.Second)
	}
}

func main() {
	channel := make(chan int)

	for i := 0; i < 5; i++ {
		go worker(i, channel)
	}

	for i := 1; i <= 10; i++ {
		channel <- i
	}
}
