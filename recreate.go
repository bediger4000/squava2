package main

import (
	"fmt"
	"log"
	"os"
	"squava2/mover"
)

func main() {

	markers := []rune{'X', '_', 'O'}

	board := make([][]rune, 5)
	for i := 0; i < 5; i++ {
		board[i] = []rune{'_', '_', '_', '_', '_'}
	}

	mvr := mover.NewFromFile(os.Args[1])
	mvr.NextPlayer(1)

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
