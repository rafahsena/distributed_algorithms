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
	black   bool
}

type Neighbour struct {
	ID     string
	Weight uint32
	From   chan Message // Canal que envia
	To     chan Message // Canal que recebe
	black  bool
}

func redirect(in chan Message, neigh Neighbour) {
	for {
		msg := <-neigh.From
		in <- msg
	}
}

func process(id string, msg Message, wg *sync.WaitGroup, neighs ...Neighbour) {
	wg.Add(1)
	var pai Neighbour
	var dist uint32 = math.MaxUint32

	// Redirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Message, 1)
	nmap := make(map[string]Neighbour)

	for _, neigh := range neighs {
		nmap[neigh.ID] = neigh
		go redirect(in, neigh)
	}

	if msg.Sender == "init" {
		// Processo iniciador
		dist = 0
		message := Message{id, dist}
		for _, neigh := range neighs {
			neigh.To <- message // Envia o msg do iniciador para o vizinho
		}
	} else {
		for message := range in {
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
					}
				}
			} else {
				//fmt.Printf("O caminho em %s já é o mais curto no momento\n", id)
			}
			fmt.Printf("%s: dist -> %d, pai: %s\n", id, dist, pai.ID)
		}
	}

	wg.Done()
}

func main() {

	pR := make(chan Message, 1)
	tS := make(chan Message, 1)
	sT := make(chan Message, 1)
	rQ := make(chan Message, 1)
	rP := make(chan Message, 1)
	qR := make(chan Message, 1)
	pQ := make(chan Message, 1)
	qP := make(chan Message, 1)
	rT := make(chan Message, 1)
	tR := make(chan Message, 1)
	rS := make(chan Message, 1)
	sR := make(chan Message, 1)

	var wg sync.WaitGroup

	go process("T", Message{}, &wg, Neighbour{"R", 4, rT, tR, false}, Neighbour{"S", 1, sT, tS, false})
	go process("S", Message{}, &wg, Neighbour{"R", 1, rS, sR, false}, Neighbour{"T", 1, tS, sT, false})
	go process("R", Message{}, &wg, Neighbour{"Q", 1, qR, rQ, false}, Neighbour{"P", 3, pR, rP, false}, Neighbour{"S", 1, sR, rS, false}, Neighbour{"T", 4, tR, rT, false})
	go process("Q", Message{}, &wg, Neighbour{"R", 1, rQ, qR, false}, Neighbour{"P", 1, pQ, qP, false})
	process("P", Message{"init", 0}, &wg, Neighbour{"Q", 1, qP, pQ, false}, Neighbour{"R", 3, rP, pR, false})
	time.Sleep(2 * time.Second)

}
