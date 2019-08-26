package main

import (
	"fmt"
	"sync"
)

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
	for {
		token := <-child.From
		in <- token
	}
}

func node(id string, msg Message, value int32, wg *sync.WaitGroup, children ...Children) {

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
			child.To <- token
			tmp := <-in
			if tmp.value > newInitiator.value {
				newInitiator = tmp
			}
		}
		fmt.Printf("O novo iniciador é: %s\n", newInitiator.Sender)
		// Avisa aos filhos quem é o novo iniciador
		for _, child := range children {
			child.To <- newInitiator
		}
		wg.Wait()
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
			} else {
				parentIndex = index
			}
		}
		//Envia o maior entre ele e seus filhos para o pai
		children[parentIndex].To <- greatest
		//Recebe o novo iniciador do pai
		greatest = <-in
		fmt.Printf("%s foi avisado que %s é o novo iniciador\n", id, greatest.Sender)
		//Avisa aos filhos quem é o novo iniciador
		for _, child := range children {
			if child.ID != parent.Sender {
				child.To <- greatest
			}
		}
		wg.Done()
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

	var wg sync.WaitGroup
	wg.Add(4)
	go node("T", Message{}, 40, &wg, Children{"Q", qT, tQ})
	go node("S", Message{}, 6, &wg, Children{"P", pS, sP})
	go node("R", Message{}, 10, &wg, Children{"P", pR, rP})
	go node("Q", Message{}, 70, &wg, Children{"T", tQ, qT}, Children{"P", pQ, qP})
	node("P", Message{"init"}, 9, &wg, Children{"Q", qP, pQ}, Children{"R", rP, pR}, Children{"S", sP, pS})
}
