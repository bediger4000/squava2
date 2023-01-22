package main

import (
	"fmt"
	"log"
	"os"
	"squava2/mover"
)

func main() {

	markers := []rune{'X', 'O'}

	board := make([][]rune, 5)
	for i := 0; i < 5; i++ {
		board[i] = []rune{'_', '_', '_', '_', '_'}
	}

	mvr := mover.NewFromFile(os.Args[1])

	for moveCounter := 0; moveCounter < 25; moveCounter++ {

		n, m, useIt := mvr.Next()
		if !useIt {
			continue
		}

		fmt.Printf("%c move %d,%d\n", markers[moveCounter%2], n, m)
		board[n][m] = markers[moveCounter%2]

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
