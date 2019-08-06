// Implementação do algoritmo de travessia de Tarry
package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Token struct {
	Sender string
	dist   uint32
}

type Neighbour struct {
	ID     string
	Weight uint32
	From   chan Token // Canal que envia
	To     chan Token // Canal que recebe
}

func redirect(in chan Token, neigh Neighbour) {
	for {
		token := <-neigh.From
		in <- token
	}
}

func process(id string, token Token, wg *sync.WaitGroup, neighs ...Neighbour) {
	wg.Add(1)
	var pai Neighbour
	var dist uint32 = math.MaxUint32

	// Redirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Token, 1)
	nmap := make(map[string]Neighbour)

	for _, neigh := range neighs {
		nmap[neigh.ID] = neigh
		go redirect(in, neigh)
	}

	if token.Sender == "init" {
		// Processo iniciador
		dist = 0
		tk := Token{id, dist}
		for _, neigh := range neighs {
			neigh.To <- tk // Envia o token do iniciador para o vizinho
		}
	} else {
		for tk := range in {
			// Processo não iniciador
			fmt.Printf("De %s para %s\n", tk.Sender, id)
			neigh := nmap[tk.Sender]
			if tk.dist+neigh.Weight < dist {
				dist = tk.dist + neigh.Weight
				pai = neigh
				fmt.Printf("* %s é pai de %s, e a dist até o iniciador é: %d\n", pai.ID, id, dist)

				for _, neigh := range neighs {
					// Entrega o token para o vizinho se ele não for o pai
					if pai.ID != neigh.ID {
						tk.Sender = id
						tk.dist = dist
						neigh.To <- tk
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

	var wg sync.WaitGroup

	go process("T", Token{}, &wg, Neighbour{"R", 4, rT, tR}, Neighbour{"S", 1, sT, tS})
	go process("S", Token{}, &wg, Neighbour{"R", 1, rS, sR}, Neighbour{"T", 1, tS, sT})
	go process("R", Token{}, &wg, Neighbour{"Q", 1, qR, rQ}, Neighbour{"P", 3, pR, rP}, Neighbour{"S", 1, sR, rS}, Neighbour{"T", 4, tR, rT})
	go process("Q", Token{}, &wg, Neighbour{"R", 1, rQ, qR}, Neighbour{"P", 1, pQ, qP})
	process("P", Token{"init", 0}, &wg, Neighbour{"Q", 1, qP, pQ}, Neighbour{"R", 3, rP, pR})
	time.Sleep(2 * time.Second)

}
