package main

import (
	"flag"
	"fmt"
	"log"
	"squava2/mover"
)

func main() {
	computerFirstPtr := flag.Bool("C", false, "Computer takes first move (default false)")
	flag.Parse()

	markers := []rune{'O', '_', 'X'}

	board := make([][]rune, 5)
	for i := 0; i < 5; i++ {
		board[i] = []rune{'_', '_', '_', '_', '_'}
	}

	mvr := mover.NewFromFile(flag.Arg(0))
	mvr.NextPlayer(-1)
	if *computerFirstPtr {
		mvr.NextPlayer(1)
	}

	for {

		player, n, m, counter, useIt := mvr.Next()
		if !useIt || counter > 24 {
			break
		}

		fmt.Printf("%c move %d,%d\n", markers[player+1], n, m)
		board[n][m] = markers[player+1]

		for i := 0; i < 5; i++ {
			for _, marker := range board[i] {
				fmt.Printf("%c ", marker)
			}
			fmt.Println()
		}
		fmt.Println()
		_, err := fmt.Scanf("\n")
		if err != nil {
			log.Print(err)
		}
	}
}
