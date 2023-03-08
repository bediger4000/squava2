package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	MAXIMIZER = 1
	MINIMIZER = -1
	UNSET     = 0
)

func main() {
	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))

	marks := [2]int{-1, 1}

	for i := 0; true; i++ {
		var board [25]int
		var winner, move int
		bits := rand.Uint64()
		for move = 0; move < 25; move++ {
			x := ((bits & (1 << move)) >> move)
			mark := marks[x]
			board[move] = mark
			winner = findWinner(&board)
			if winner != UNSET {
				break
			}
		}

		if move == 25 {
			who := [3]string{"O", "cat", "X"}[board[winner]+1]
			fmt.Printf("%s won:\n%s\n\n", who, boardString(board))
		}

	}
}

// findWinner will return MAXIMIZER or MINIMIZER if somebody won,
// UNSET if nobody wins based on argument board.
// Pointer to [25]int to avoid creating copies of array that
// don't get used.
func findWinner(board *[25]int) int {
	for _, i := range importantCells {
		if (*board)[i] != UNSET {
			for _, quad := range mctsWinningQuads[i] {
				sum := (*board)[quad[0]] + (*board)[quad[1]] + (*board)[quad[2]] + (*board)[quad[3]]
				switch sum {
				case 4:
					return MAXIMIZER
				case -4:
					return MINIMIZER
				}
			}
		}
	}
	for _, i := range importantCells {
		if (*board)[i] != UNSET {
			for _, triplet := range mctsLosingTriplets[i] {
				sum := (*board)[triplet[0]] + (*board)[triplet[1]] + (*board)[triplet[2]]
				switch sum {
				case 3:
					return MINIMIZER
				case -3:
					return MAXIMIZER
				}
			}
		}
	}
	return UNSET
}

// boardString exists as a separate function so that if
// printf-style debugging is necessary, gameState.board
// can also get printed.
func boardString(board [25]int) string {
	buf := &strings.Builder{}
	buf.WriteString("   0 1 2 3 4\n")
	for i := 0; i < 25; i++ {
		if (i % 5) == 0 {
			fmt.Fprintf(buf, "%c  ", rune(i/5)+'0')
		}
		fmt.Fprintf(buf, "%c ", "O_X"[board[i]+1])
		if (i % 5) == 4 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

var importantCells = [9]int{2, 7, 10, 11, 12, 13, 14, 17, 22}

// 25 rows only to make looping easier. The filled-in
// rows are the only quads you actually have to check
// to find out if there's a win
var mctsWinningQuads = [25][][]int{
	{}, {},
	{{0, 1, 2, 3}, {1, 2, 3, 4}, {2, 7, 12, 17}},
	{}, {}, {}, {},
	{{5, 6, 7, 8}, {6, 7, 8, 9}, {7, 12, 17, 22}},
	{}, {},
	{{0, 5, 10, 15}, {5, 10, 15, 20}},
	{{1, 6, 11, 16}, {6, 11, 16, 21}, {3, 7, 11, 15}, {5, 11, 17, 23}},
	{{10, 11, 12, 13}, {11, 12, 13, 14}, {4, 8, 12, 16}, {8, 12, 16, 20}, {0, 6, 12, 18}, {6, 12, 18, 24}},
	{{3, 8, 13, 18}, {8, 13, 18, 23}, {1, 7, 13, 19}, {9, 13, 17, 21}},
	{{4, 9, 14, 19}, {9, 14, 19, 24}},
	{}, {},
	{{15, 16, 17, 18}, {16, 17, 18, 19}},
	{}, {}, {}, {},
	{{20, 21, 22, 23}, {21, 22, 23, 24}},
	{}, {},
}

// 25 rows only to make looping easier. The filled-in
// rows are the only triplets you actually have to check
// to find out if there's a loss.
var mctsLosingTriplets = [][][]int{
	{}, {},
	{{0, 1, 2}, {1, 2, 3}, {2, 3, 4}, {2, 7, 12}, {2, 6, 10}, {14, 8, 2}},
	{}, {}, {}, {},
	{{5, 6, 7}, {6, 7, 8}, {7, 8, 9}, {2, 7, 12}, {7, 12, 17}, {3, 7, 11}, {7, 11, 15}, {1, 7, 13}, {7, 13, 19}},
	{}, {},
	{{10, 11, 12}, {0, 5, 10}, {5, 10, 15}, {10, 15, 20}, {2, 6, 10}, {10, 16, 22}},
	{{10, 11, 12}, {11, 12, 13}, {1, 6, 11}, {6, 11, 16}, {11, 16, 21}, {3, 7, 11}, {7, 11, 15}, {5, 11, 17}, {11, 17, 23}},
	{{10, 11, 12}, {11, 12, 13}, {12, 13, 14}, {2, 7, 12}, {7, 12, 17}, {12, 17, 22}, {0, 6, 12}, {6, 12, 18}, {12, 18, 24}, {4, 8, 12}, {8, 12, 16}, {12, 16, 20}},
	{{11, 12, 13}, {12, 13, 14}, {3, 8, 13}, {8, 13, 18}, {13, 18, 23}, {1, 7, 13}, {7, 13, 19}, {21, 17, 13}, {17, 13, 9}},
	{{12, 13, 14}, {4, 9, 14}, {9, 14, 19}, {14, 19, 24}, {22, 18, 14}, {14, 8, 2}},
	{}, {},
	{{15, 16, 17}, {16, 17, 18}, {17, 18, 19}, {7, 12, 17}, {12, 17, 22}, {5, 11, 17}, {11, 17, 23}, {21, 17, 13}, {17, 13, 9}},
	{}, {}, {}, {},
	{{20, 21, 22}, {21, 22, 23}, {22, 23, 24}, {12, 17, 22}, {10, 16, 22}, {22, 18, 14}},
	{}, {},
}
