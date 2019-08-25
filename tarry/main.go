package main

import "fmt"

//Message ...
type Message struct {
	Sender string
}

//Token ...
type Token struct {
	Sender string
}

//Neighbour ...
type Neighbour struct {
	ID     string
	Weight uint32
	From   chan Token // Canal que envia
	To     chan Token // Canal que recebe
}

func redirect(in chan Token, neigh Neighbour) {
	token := <-neigh.From
	in <- token
}

func tarry(id string, msg Message, neighs ...Neighbour) {
	in := make(chan Token, 10)
	for _, neigh := range neighs {
		go redirect(in, neigh)
	}
	if msg.Sender == "init" {
		fmt.Printf("A <-- B indica que A é pai de B\n\n")
		fmt.Printf("O iniciador é: %s\n", id)
		token := Token{id}
		neighs[0].To <- token
		for index, neigh := range neighs {
			if index != 0 {
				token = <-in
				fmt.Printf("%s -> %s\n", neigh.ID, id)
				token.Sender = id
				neigh.To <- token
			}
		}
		token = <-in
		fmt.Printf("%s -> %s\n", token.Sender, id)
	} else {
		token := <-in
		fmt.Printf("%s -> %s\n", token.Sender, id)
		fmt.Printf("%s <-- %s\n", token.Sender, id)
		parentIndex := -1
		for index, neigh := range neighs {
			if neigh.ID != token.Sender {
				token.Sender = id
				neigh.To <- token
				token := <-in
				fmt.Printf("%s -> %s\n", token.Sender, id)
			} else {
				parentIndex = index
			}
		}
		token.Sender = id
		neighs[parentIndex].To <- token
	}
}

func main() {
	pR := make(chan Token, 1)
	tS := make(chan Token, 1)
	sT := make(chan Token, 1)
	rQ := make(chan Token, 1)
	rP := make(chan Token, 1)
	qR := make(chan Token, 1)
	pQ := make(chan Token, 1)
	qP := make(chan Token, 1)
	rT := make(chan Token, 1)
	tR := make(chan Token, 1)
	rS := make(chan Token, 1)
	sR := make(chan Token, 1)

	go tarry("T", Message{}, Neighbour{"R", 4, rT, tR}, Neighbour{"S", 1, sT, tS})
	go tarry("S", Message{}, Neighbour{"R", 1, rS, sR}, Neighbour{"T", 1, tS, sT})
	go tarry("R", Message{}, Neighbour{"Q", 1, qR, rQ}, Neighbour{"P", 3, pR, rP}, Neighbour{"S", 1, sR, rS}, Neighbour{"T", 4, tR, rT})
	go tarry("Q", Message{}, Neighbour{"R", 1, rQ, qR}, Neighbour{"P", 1, pQ, qP})
	tarry("P", Message{"init"}, Neighbour{"Q", 1, qP, pQ}, Neighbour{"R", 3, rP, pR})
}
