// Implementação do algoritmo de travessia de Tarry
package main

import (
	"fmt"
	"math"
	"sync"
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

func safra(id string, wg2 *sync.WaitGroup, msg Message, nmap map[string]Neighbour, wg *sync.WaitGroup, cont *int32, inSafra chan Token, neighs []Neighbour) {
	var parent Neighbour
	contador := 1
	if msg.Sender == "init" {
		for {
			fmt.Printf("    $ Safra está executando no ciclo: %d\n", contador)
			fmt.Printf("    $ Contador do iniciador: %d\n", (*cont))
			token := Token{(*cont), id}
			wg.Wait()
			neighs[0].To <- token
			for index, neigh := range neighs {
				if index != 0 {
					token = <-inSafra
					token.Sender = id
					neigh.To <- token
				}
			}
			token = <-inSafra
			fmt.Printf("    $ Contador do token após a travessia: %d\n", token.counter)
			if token.counter == 0 {
				fmt.Println("\n    $ Safra detectou o término")
				wg2.Done()
				break
			}
			contador++
		}
	} else {
		for token := range inSafra {
			token.counter += (*cont)
			fmt.Printf("    $ %s -> %s, soma = %d, counter = %d, ciclo: %d\n", token.Sender, id, token.counter, (*cont), contador)
			for _, neigh := range neighs {
				if parent.ID == "" {
					parent = nmap[token.Sender]
				}
				if neigh.ID != parent.ID {
					token.Sender = id
					wg.Wait()
					neigh.To <- token
					token = <-inSafra
				}
				token.Sender = id
				wg.Wait()
				parent.To <- token
			}
		}
	}

}

func process(id string, wg2 *sync.WaitGroup, msg Message, neighs ...Neighbour) {
	var pai Neighbour
	var dist uint32 = math.MaxUint32
	cont := int32(0)
	var wg sync.WaitGroup

	// Redirecionando todos os canais de entrada para um único canal "in_chandy" de entrada
	inChandy := make(chan Message, 100)
	inSafra := make(chan Token, 100)
	nmap := make(map[string]Neighbour)

	for _, neigh := range neighs {
		nmap[neigh.ID] = neigh
		go redirect(inChandy, inSafra, neigh)
	}

	// Algoritmo de Safra

	// Chandy-Misra algorithm
	if msg.Sender == "init" {
		// Processo iniciador
		dist = 0
		fmt.Printf("* O iniciador envia para os filhos uma mensagem com distância igual a zero\n")
		message := Message{id, dist}
		for _, neigh := range neighs {
			neigh.To <- message // Envia o msg do iniciador para o vizinho
			cont++
		}
		go safra(id, wg2, msg, nmap, &wg, &cont, inSafra, neighs)
		for range inChandy {
			cont--
		}
	} else {
		go safra(id, wg2, msg, nmap, &wg, &cont, inSafra, neighs)
		for message := range inChandy {
			wg.Add(1)
			cont--
			// Processo não iniciador
			neigh := nmap[message.Sender]
			if message.dist+neigh.Weight < dist {
				dist = message.dist + neigh.Weight
				pai = neigh
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
			fmt.Printf("* Nó %s tem distância até o iniciador igual a %d\n", id, dist)
			wg.Done()
		}
	}

}

func main() {

	pR := make(chan interface{}, 100)
	tS := make(chan interface{}, 100)
	sT := make(chan interface{}, 100)
	rQ := make(chan interface{}, 100)
	rP := make(chan interface{}, 100)
	qR := make(chan interface{}, 100)
	pQ := make(chan interface{}, 100)
	qP := make(chan interface{}, 100)
	rT := make(chan interface{}, 100)
	tR := make(chan interface{}, 100)
	rS := make(chan interface{}, 100)
	sR := make(chan interface{}, 100)

	var wg2 sync.WaitGroup

	wg2.Add(1)
	go process("T", &wg2, Message{}, Neighbour{"R", 4, rT, tR}, Neighbour{"S", 1, sT, tS})
	go process("S", &wg2, Message{}, Neighbour{"R", 1, rS, sR}, Neighbour{"T", 1, tS, sT})
	go process("R", &wg2, Message{}, Neighbour{"Q", 1, qR, rQ}, Neighbour{"P", 3, pR, rP}, Neighbour{"S", 1, sR, rS}, Neighbour{"T", 4, tR, rT})
	go process("Q", &wg2, Message{}, Neighbour{"R", 1, rQ, qR}, Neighbour{"P", 1, pQ, qP})
	go process("P", &wg2, Message{"init", 0}, Neighbour{"Q", 1, qP, pQ}, Neighbour{"R", 3, rP, pR})
	wg2.Wait()

}
