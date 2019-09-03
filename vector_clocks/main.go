package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Body      string
	Timestamp [3]int
	mux       sync.Mutex
}

func event(pid int, counter [3]int) {
	counter[pid-1]++
	fmt.Printf("Event in process pid=%v. Contador = %v\n", pid, counter)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func calcVectorTimestamp(recvTimestamp [3]int, counter [3]int, pid int) [3]int {
	var newTimestampVector [3]int
	for i := 0; i < len(newTimestampVector); i++ {
		if i != pid-1 {
			newTimestampVector[i] = max(recvTimestamp[i], counter[i])
		} else {
			newTimestampVector[i] = counter[i] + 1
		}
	}
	return newTimestampVector
}

func sendMessage(ch chan Message, pid int, counter [3]int) [3]int {
	counter[pid-1]++
	ch <- Message{Body: "Test msg!!!", Timestamp: counter}
	fmt.Printf("Mensagem enviada por pid=%v. Contador = %v\n", pid, counter)
	return counter

}

func receiveMessage(ch chan Message, pid int, counter [3]int) [3]int {
	message := <-ch
	message.mux.Lock()
	defer message.mux.Unlock()
	counter = calcVectorTimestamp(message.Timestamp, counter, pid)
	fmt.Printf("Mensagem recebida no pid=%v. Contador = %v\n", pid, counter)
	return counter
}

func processOne(ch12, ch21 chan Message) {
	pid := 1
	counter := [3]int{0, 0, 0}
	counter = sendMessage(ch12, pid, counter)
	counter = receiveMessage(ch21, pid, counter)

}

func processTwo(ch12, ch21, ch23, ch32 chan Message) {
	pid := 2
	counter := [3]int{0, 0, 0}
	counter = receiveMessage(ch12, pid, counter)
	counter = sendMessage(ch21, pid, counter)
	counter = sendMessage(ch23, pid, counter)
	counter = receiveMessage(ch32, pid, counter)

}

func processThree(ch23, ch32 chan Message) {
	pid := 3
	counter := [3]int{0, 0, 0}
	counter = receiveMessage(ch23, pid, counter)
	counter = sendMessage(ch32, pid, counter)

}

func main() {
	oneTwo := make(chan Message, 100)
	twoOne := make(chan Message, 100)
	twoThree := make(chan Message, 100)
	threeTwo := make(chan Message, 100)

	go processOne(oneTwo, twoOne)
	go processTwo(oneTwo, twoOne, twoThree, threeTwo)
	go processThree(twoThree, threeTwo)

	time.Sleep(5 * time.Second)
}
