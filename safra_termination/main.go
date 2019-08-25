// Implementação do algoritmo de travessia de Tarry
package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Message struct {
	Sender string
	dist   uint32
}

type Token struct {
	counter int32
	Sender  string
}

type Neighbour struct {
	ID     string
	Weight uint32
	From   chan interface{} // Canal que envia
	To     chan interface{} // Canal que recebe
}

func redirect(inChandy chan Message, inSafra chan Token, neigh Neighbour) {
	for value := range neigh.From {
		switch msg := value.(type) {
		case Token:
			inSafra <- msg
		case Message:
			inChandy <- msg
		default:
			fmt.Println("Deu ruim")
		}
	}
}

func safra(id string, msg Message, nmap map[string]Neighbour, wg *sync.WaitGroup, cont *int, inSafra chan Token, neighs []Neighbour) {

	if msg.Sender == "init" {
		tk := Token{0, id}
		wg.Wait()
		for _, neigh := range neighs {
			// Entrega o msg para os vizinhos
			neigh.To <- tk
		}
	}
	for token := range inSafra {
		tk := token
		println("------------------")
		println(tk.counter)
		println("------------------")
		if msg.Sender == "init" {
			if tk.counter == 0 {
				break
			} else {
				tk.counter = 0
			}
		}
		tk.counter += int32(*cont)
		wg.Wait()
		for _, neigh := range neighs {
			// Entrega o msg para o vizinho se ele não for o pai
			if token.Sender != neigh.ID {
				tk.Sender = id
				neigh.To <- tk
			}
		}
		nmap[token.Sender].To <- tk
	}
	fmt.Println("Cabou")

}

func process(id string, msg Message, neighs ...Neighbour) {
	var pai Neighbour
	var dist uint32 = math.MaxUint32
	cont := 0
	var wg sync.WaitGroup

	// Redirecionando todos os canais de entrada para um único canal "in_chandy" de entrada
	inChandy := make(chan Message)
	inSafra := make(chan Token)
	nmap := make(map[string]Neighbour)

	for _, neigh := range neighs {
		nmap[neigh.ID] = neigh
		go redirect(inChandy, inSafra, neigh)
	}

	// Algoritmo de Safra
	go safra(id, msg, nmap, &wg, &cont, inSafra, neighs)

	// Chandy-Misra algorithm
	if msg.Sender == "init" {
		// Processo iniciador
		wg.Add(1)
		dist = 0
		message := Message{id, dist}
		for _, neigh := range neighs {
			neigh.To <- message // Envia o msg do iniciador para o vizinho
			cont++
		}
		wg.Done()
		for range inChandy {
			cont--
		}
	} else {
		for message := range inChandy {
			wg.Add(1)
			cont--
			// Processo não iniciador
			fmt.Printf("De %s para %s\n", message.Sender, id)
			neigh := nmap[message.Sender]
			if message.dist+neigh.Weight < dist {
				dist = message.dist + neigh.Weight
				pai = neigh
				fmt.Printf("* %s é pai de %s, e a dist até o iniciador é: %d\n", pai.ID, id, dist)

				for _, neigh := range neighs {
					// Entrega o msg para o vizinho se ele não for o pai
					if pai.ID != neigh.ID {
						message.Sender = id
						message.dist = dist
						neigh.To <- message
						cont++
					}
				}
			}
			fmt.Printf("%s: dist -> %d, pai: %s\n", id, dist, pai.ID)
			wg.Done()
		}
	}

}

func main() {

	pR := make(chan interface{})
	tS := make(chan interface{})
	sT := make(chan interface{})
	rQ := make(chan interface{})
	rP := make(chan interface{})
	qR := make(chan interface{})
	pQ := make(chan interface{})
	qP := make(chan interface{})
	rT := make(chan interface{})
	tR := make(chan interface{})
	rS := make(chan interface{})
	sR := make(chan interface{})

	go process("T", Message{}, Neighbour{"R", 4, rT, tR}, Neighbour{"S", 1, sT, tS})
	go process("S", Message{}, Neighbour{"R", 1, rS, sR}, Neighbour{"T", 1, tS, sT})
	go process("R", Message{}, Neighbour{"Q", 1, qR, rQ}, Neighbour{"P", 3, pR, rP}, Neighbour{"S", 1, sR, rS}, Neighbour{"T", 4, tR, rT})
	go process("Q", Message{}, Neighbour{"R", 1, rQ, qR}, Neighbour{"P", 1, pQ, qP})
	process("P", Message{"init", 0}, Neighbour{"Q", 1, qP, pQ}, Neighbour{"R", 3, rP, pR})
	time.Sleep(2 * time.Second)

}
