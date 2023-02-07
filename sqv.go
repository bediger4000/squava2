package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"squava2/mover"
	"squava2/players"
)

const (
	HUMAN    = -1
	COMPUTER = 1
)

func main() {

	computerFirstPtr := flag.Bool("C", false, "Computer takes first move (default false)")
	maxDepthPtr := flag.Int("d", 10, "maximum lookahead depth (alpha/beta)")
	typ := flag.String("t", "A", "player type, A: alphabeta, G: A/B+avoid bad positions, M: MCTS/Plain, U: MCTS/UCB1")
	i := flag.Int("i", 500000, "MCTS iterations")
	partialGame := flag.String("p", "", "partial game, filename or comma-sep move string")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	var winner int

	moveCounter := 0

	computerPlayer := createPlayer(*typ, *maxDepthPtr, *i)

	next := HUMAN
	if *computerFirstPtr {
		next = COMPUTER
	}

	// computerPlayer keeps track of the board internally,
	// but we'll keep track too, so the human can be informed
	// that an input move has already been taken.
	bd := new(Board)

	if *partialGame != "" {
		next = gameSoFar(next, *partialGame, bd, computerPlayer)
		playerPhrase := "human"
		if next == 1 {
			playerPhrase = "computer"
		}
		fmt.Printf("After reading in partial game, next player %s\n", playerPhrase)
		fmt.Printf("%s\n", computerPlayer)
	}

	for moveCounter < 25 {

		switch next {

		case HUMAN:
			l, m := bd.readMove()
			computerPlayer.MakeMove(l, m, HUMAN)
			next = COMPUTER

		case COMPUTER:
			before := time.Now()
			i, j, value, leafCount := computerPlayer.ChooseMove()
			et := time.Since(before)

			fmt.Printf("X (%s) <%d,%d> (%d) [%d] %v\n", computerPlayer.Name(), i, j, value, leafCount, et)

			bd.makeMove(i, j, COMPUTER)
			next = HUMAN
		}

		moveCounter++
		winner = computerPlayer.FindWinner()

		if winner != 0 || moveCounter >= 25 {
			break
		}

		fmt.Printf("%s\n", computerPlayer)
	}

	switch winner {
	case 1:
		fmt.Printf("player 1 X (%s) wins\n", computerPlayer.Name())
	case -1:
		fmt.Printf("player 2 O (human) wins\n")
	default:
		fmt.Printf("Cat wins\n")
	}

	fmt.Printf("%s\n", computerPlayer)
}

func createPlayer(typ string, maxDepth int, iterations int) players.Player {

	typ = strings.ToUpper(typ)

	switch typ {
	case "A":
		return players.NewAlphaBeta(false, maxDepth)
	case "G":
		ab := players.NewAlphaBeta(false, maxDepth)
		ab.SetAvoid()
		return ab
	case "M":
		mcts := players.NewMCTS(iterations)
		mcts.SetIterations(iterations)
		return mcts
	case "U":
		mcts := players.NewMCTS(iterations)
		mcts.SetIterations(iterations)
		mcts.SetUCB1()
		return mcts
	}

	return nil
}

type Board [5][5]int

func (bd *Board) makeMove(x, y, player int) {
	bd[x][y] = player
}

func (bd *Board) readMove() (x, y int) {
	readMove := false
	for !readMove {
		fmt.Printf("Your move: ")
		_, err := fmt.Scanf("%d %d\n", &x, &y)
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			fmt.Printf("Failed to read: %v\n", err)
			os.Exit(1)
		}
		switch {
		case x < 0 || x > 4 || y < 0 || y > 4:
			fmt.Printf("Choose two numbers between 0 and 4, try again\n")
		case bd[x][y] == 0:
			readMove = true
		case bd[x][y] != 0:
			fmt.Printf("Cell (%d, %d) already occupied, try again\n", x, y)
		}
	}
	bd.makeMove(x, y, HUMAN)
	return x, y
}

func gameSoFar(firstPlayer int, partial string, bd *Board, p players.Player) int {

	var moves *mover.Mvr

	// var partial could name a file, or be a string of x,y moves
	if _, err := os.Stat(partial); err == nil {
		moves = mover.NewFromFile(partial)
	} else {
		moves = mover.NewFromBuffer([]byte(partial))
	}

	moves.NextPlayer(firstPlayer)

	next := firstPlayer

	for {
		player, n, m, counter, useIt := moves.Next()
		if !useIt || counter > 24 {
			break
		}
		(*bd)[n][m] = player
		p.MakeMove(n, m, player)
		next = player
	}

	return 0 - next
}
