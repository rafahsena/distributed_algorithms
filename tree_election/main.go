package main

import "fmt"

type Message struct {
	id string
}

type Token struct {
	Sender string
	value  int32
}

type Children struct {
	ID   string
	From chan Token // Canal que envia
	To   chan Token // Canal que recebe
}

func redirect(in chan Token, child Children) {
	token := <-child.From
	in <- token
}

func node(id string, msg Message, value int32, children ...Children) {

	in := make(chan Token, 10)
	for _, child := range children {
		go redirect(in, child)
	}
	if msg.id == "init" {
		fmt.Printf("O código considera um valor como chave para escolher novo iniciador\n\n")
		fmt.Printf("O iniciador é: %s\n", id)
		token := Token{id, -1}
		newInitiator := token
		for _, child := range children {
			fmt.Printf("%s -> %s\n", child.ID, id)
			child.To <- token
			tmp := <-in
			if tmp.value > newInitiator.value {
				newInitiator = tmp
			}
		}
		fmt.Printf("O novo iniciador é: %s\n", newInitiator.Sender)
	} else {
		parent := <-in
		token := Token{id, value}
		greatest := token
		parentIndex := -1
		for index, child := range children {
			if child.ID != parent.Sender {
				child.To <- token
				tmp := <-in
				if tmp.value > greatest.value {
					greatest = tmp
				}
				fmt.Printf("%s -> %s\n", token.Sender, id)
			} else {
				parentIndex = index
			}
		}
		children[parentIndex].To <- greatest
	}
}

func main() {
	pR := make(chan Token, 1)
	rP := make(chan Token, 1)
	pQ := make(chan Token, 1)
	qP := make(chan Token, 1)
	qT := make(chan Token, 1)
	tQ := make(chan Token, 1)
	sP := make(chan Token, 1)
	pS := make(chan Token, 1)

	go node("T", Message{}, 40, Children{"Q", qT, tQ})
	go node("S", Message{}, 6, Children{"P", pS, sP})
	go node("R", Message{}, 10, Children{"P", pR, rP})
	go node("Q", Message{}, 70, Children{"T", tQ, qT}, Children{"P", pQ, qP})
	node("P", Message{"init"}, 9, Children{"Q", qP, pQ}, Children{"R", rP, pR}, Children{"S", sP, pS})
}
