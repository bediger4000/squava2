package main

/*
 * Produce an Elo-style rating for the algorithmic players
 * by repeatedly playing one against another.
 */

import (
	"flag"
	"fmt"
	"math"
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

	gameCount := flag.Int("n", 1, "play <number> games non-interactively")
	aRating := flag.Float64("A", 1300., "Alpha-beta minimaxing player initial rating")
	aGames := flag.Float64("a", 14., "Alpha-beta minimaxing player effective games count")
	gRating := flag.Float64("G", 1300., "Better static valuation player (G) initial rating")
	gGames := flag.Float64("g", 14., "Better static valuation player (G) player effective games count")
	mRating := flag.Float64("M", 1300., "M player initial rating")
	mGames := flag.Float64("m", 14., "M player player effective games count")
	uRating := flag.Float64("U", 1300., "U player initial rating")
	uGames := flag.Float64("u", 14., "U player player effective games count")

	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	nonInteractiveGames(*gameCount, *aRating, *aGames, *gRating, *gGames, *mRating, *mGames, *uRating, *uGames)
}

type PlayerRating struct {
	name           string
	rating         float64
	effectiveGames float64
}

func nonInteractiveGames(gameCount int, aRating, aGames, gRating, gGames, mRating, mGames, uRating, uGames float64) {

	started := time.Now()

	var playerList [4]PlayerRating

	for i := 0; i < 4; i++ {
		playerList[i].rating = 1300.
		playerList[i].effectiveGames = 14.0
	}

	playerList[0].name = "A"
	playerList[0].rating = aRating
	playerList[0].effectiveGames = aGames

	playerList[1].name = "G"
	playerList[1].rating = gRating
	playerList[1].effectiveGames = gGames

	playerList[2].name = "M"
	playerList[2].rating = mRating
	playerList[2].effectiveGames = mGames

	playerList[3].name = "U"
	playerList[3].rating = uRating
	playerList[3].effectiveGames = uGames

	for i := 0; i < gameCount; i++ {

		firstChoice := rand.Intn(4)
		secondChoice := rand.Intn(4)
		for firstChoice == secondChoice {
			secondChoice = rand.Intn(4)
		}

		first, second := createPlayers(
			playerList[firstChoice].name,
			playerList[secondChoice].name,
			10,
			false,
		)

		moveCounter := 0

		var moves [25][2]int
		var values [25][2]int
		var winner int

		before := time.Now()

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
		elapsed := time.Since(before)

		// Either moveCounter == 25, or winner != 0, or both
		var firstScore, secondScore float64
		var winning string
		switch winner {
		case MAXIMIZER:
			winning = playerList[firstChoice].name
			firstScore = 1.0
		case MINIMIZER:
			winning = playerList[secondChoice].name
			secondScore = 1.0
		default:
			// cat can get a game, but maybe only if the players cooperate?
			winning = "cat"
			firstScore = 0.5
			secondScore = 0.5
		}

		previousFirstRating := playerList[firstChoice].rating
		playerList[firstChoice].effectiveGames++
		K := 800. / playerList[firstChoice].effectiveGames
		E := We(playerList[firstChoice].rating, playerList[secondChoice].rating)
		playerList[firstChoice].rating += K * (firstScore - E)

		previousSecondRating := playerList[secondChoice].rating
		playerList[secondChoice].effectiveGames++
		K = 800. / playerList[secondChoice].effectiveGames
		E = We(playerList[secondChoice].rating, previousFirstRating)
		playerList[secondChoice].rating += K * (secondScore - E)

		fmt.Printf("%d\t%.02f\t%s\t%s\t%s\t%.0f\t%.0f\t%.0f\t%.0f\t%.0f\t%.0f\n",
			i,
			elapsed.Seconds(),
			playerList[firstChoice].name,
			playerList[secondChoice].name,
			winning,
			previousFirstRating,
			playerList[firstChoice].rating,
			playerList[firstChoice].effectiveGames,
			previousSecondRating,
			playerList[secondChoice].rating,
			playerList[secondChoice].effectiveGames,
		)
	}
	for i := range playerList {
		fmt.Printf("# %s: %.0f, %.0f games\n",
			playerList[i].name,
			playerList[i].rating,
			playerList[i].effectiveGames,
		)
	}
	overallET := time.Since(started)
	fmt.Printf("# Overall elapsed time %.2f\n", overallET.Seconds())
}

func We(R, Ri float64) float64 {
	exponent := (Ri - R) / 400.
	return 1.0 / (1.0 + math.Pow(10., exponent))
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
