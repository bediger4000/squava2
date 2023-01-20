package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"squava2/players"
)

const (
	MAXIMIZER = 1
	MINIMIZER = -1
)

func main() {

	maxDepthPtr := flag.Int("d", 10, "maximum lookahead depth (alpha/beta)")
	deterministic := flag.Bool("D", false, "Play deterministically")
	firstType := flag.String("1", "A", "first player type, A: alphabeta, G: A/B+avoid bad positions, M: MCTS")
	secondType := flag.String("2", "M", "second player type, A: alphabeta, G: A/B+avoid bad positions, M: MCTS")
	nonInteractive := flag.Int("n", 1, "play <number> games non-interactively")
	i1 := flag.Int("i1", 500000, "MCTS iterations, player 1")
	i2 := flag.Int("i2", 500000, "MCTS iterations, player 2")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	if *nonInteractive > 1 {
		nonInteractiveGames(*nonInteractive, *firstType, *secondType, *maxDepthPtr)
		return
	}

	var winner int

	moveCounter := 0

	first, second := createPlayers(*firstType,
		*secondType, *maxDepthPtr, *deterministic)

	if *firstType == "M" {
		first.(*players.MCTS).SetIterations(*i1)
	}

	if *secondType == "M" {
		// second.(*players.MCTS).SetUCTK(*u2)
		second.(*players.MCTS).SetIterations(*i2)
	}

	gameStart := time.Now()
	for moveCounter < 25 {

		before := time.Now()
		i, j, value, leafCount := first.ChooseMove()
		et := time.Since(before)
		second.MakeMove(i, j, MINIMIZER)

		moveCounter++
		fmt.Printf("X (%s) <%d,%d> (%d) [%d] %v\n", first.Name(), i, j, value, leafCount, et)

		winner = first.FindWinner() // main() thinks first is maximizer
		if winner != 0 || moveCounter >= 25 {
			break
		}

		before = time.Now()
		i, j, value, leafCount = second.ChooseMove()
		et = time.Since(before)
		first.MakeMove(i, j, MINIMIZER)

		moveCounter++
		fmt.Printf("O (%s) <%d,%d> (%d) [%d] %v\n", second.Name(), i, j, value, leafCount, et)

		fmt.Printf("%s\n", first)

		winner1 := first.FindWinner()
		winner2 := -second.FindWinner() // main thinks second is minimizer
		if winner1 != winner2 {
			fmt.Printf("Winner disagreement. First %d, second %d\n", winner1, winner2)
		}
		if winner2 != 0 {
			winner = winner2
			break
		}

	}
	gameET := time.Since(gameStart)

	switch winner {
	case 1:
		fmt.Printf("player 1 X (%s) wins, %v\n", first.Name(), gameET)
	case -1:
		fmt.Printf("player 2 O (%s) wins, %v\n", second.Name(), gameET)
	default:
		fmt.Printf("Cat wins\n")
	}

	fmt.Printf("%s\n", first)

}

func nonInteractiveGames(gameCount int, firstType, secondType string, maxDepth int) {

	for i := 0; i < gameCount; i++ {

		moveCounter := 0

		first, second := createPlayers(firstType, secondType, maxDepth, false)

		fmt.Printf("%d %s %s %d ", i, first.Name(), second.Name(), maxDepth)

		var moves [25][2]int
		var values [25][2]int
		var winner int

		for moveCounter < 25 {

			i, j, value, _ := first.ChooseMove()
			moves[moveCounter][0], moves[moveCounter][1] = i, j
			values[moveCounter][0] = value
			second.MakeMove(i, j, MINIMIZER)
			moveCounter++
			winner = first.FindWinner()
			if winner != 0 || moveCounter >= 25 {
				break
			}

			i, j, value, _ = second.ChooseMove()
			moves[moveCounter][0], moves[moveCounter][1] = i, j
			values[moveCounter][1] = value
			first.MakeMove(i, j, MINIMIZER)
			moveCounter++
			winner = -second.FindWinner() // main thinks second is minimizer
			if winner != 0 {
				break
			}
		}

		fmt.Printf("%d %d", moveCounter, winner)

		for i := 0; i < moveCounter; i++ {
			marker := [2]string{"", ""}
			for j := 0; j < 2; j++ {
				if values[i][j] > 9000 {
					marker[j] = "+"
				}
				if values[i][j] < -9000 {
					marker[j] = "-"
				}
			}
			fmt.Printf(" %d%s,%d%s", moves[i][0], marker[0], moves[i][1], marker[1])
		}

		fmt.Printf("\n")
	}
}

func createPlayers(firstType, secondType string, maxDepth int, deterministic bool) (players.Player, players.Player) {

	firstType = strings.ToUpper(firstType)
	secondType = strings.ToUpper(secondType)

	return createPlayer(firstType, maxDepth, deterministic, 500000),
		createPlayer(secondType, maxDepth, deterministic, 500000)
}

func createPlayer(typ string, maxDepth int, deterministic bool, iterations int) players.Player {

	typ = strings.ToUpper(typ)

	switch typ {
	case "A":
		return players.NewAlphaBeta(deterministic, maxDepth)
	case "G":
		ab := players.NewAlphaBeta(deterministic, maxDepth)
		ab.SetAvoid()
		return ab
	case "M":
		mcts := players.NewMCTS(iterations)
		return mcts
	case "U":
		mcts := players.NewMCTS(iterations)
		mcts.SetUCB1()
		return mcts
	}

	return nil
}
