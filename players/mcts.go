package players

import (
	"fmt"
	"math/rand"
	"strings"
)

type gameState struct {
	player int
	board  [25]int
}

type Node struct {
	move       int
	player     int
	parent     *Node
	childNodes []*Node
	wins       float64
	visits     float64
	score      float64
	// score should be 0 for a losing move,
	// 1 for a winning move
	untriedMoves []int
	winner       int
}

type MCTS struct {
	board      [25]int
	iterations int
}

func NewMCTS(iterations int) *MCTS {
	return &MCTS{iterations: iterations}
}

func (p *MCTS) Name() string {
	return "MCTS/Plain"
}

func (p *MCTS) SetIterations(iterations int) {
	p.iterations = iterations
}

func (p *MCTS) MakeMove(x, y int, player int) {
	p.board[5*x+y] = player
}

// ChooseMove should choose computer's next move and
// return x,y coords of move and its score.
func (p *MCTS) ChooseMove() (xcoord int, ycoord int, value int, leafcount int) {

	var best int
	var score float64

	best, score, leafcount = bestMove(p.board, p.iterations)

	p.board[best] = MAXIMIZER

	// Since this implementations's "board" is a plain array, a move has to
	// translate to <x,y> coords
	xcoord = best / 5
	ycoord = best % 5

	value = int(score * 1000.)

	return
}

func bestMove(board [25]int, iterations int) (move int, score float64, leafCount int) {

	root := &Node{
		player: MINIMIZER, // opponent made the last move
	}
	root.untriedMoves = make([]int, 0, 25)
	for i := range board {
		if board[i] == UNSET {
			root.untriedMoves = append(root.untriedMoves, i)
		}
	}

	state := &gameState{}

	for iters := 0; iters < iterations; iters++ {

		// reset state
		for j := 0; j < 25; j++ {
			state.board[j] = board[j]
		}
		state.player = MINIMIZER

		node := root

		// Selection
		for len(node.untriedMoves) == 0 && len(node.childNodes) > 0 {
			node = node.selectBestChild()
			state.makeMove(node.move)
		}

		// node points to a Node struct that has no child nodes
		// OR
		// node points to a struct Node that has untried moves.
		//
		// state should represent the board resulting from following
		// the "best child" nodes.

		var win bool

		// Expansion will pick an untried move on the struct Node
		// pointed to by Node, if it has untried moves. If node points to a
		// struct Node reached by following "best child" nodes,
		// node may not have untried moves.
		if len(node.untriedMoves) > 0 {
			mv := node.untriedMoves[rand.Intn(len(node.untriedMoves))]

			state.makeMove(mv)

			node = node.AddChild(mv, state) // AddChild take mv out of untriedMoves slice
			node.winner = findWinner(&(state.board))
			if node.winner == MAXIMIZER {
				node.score = 1.0
				win = true
			}
			// node represents mv, the previously untried move
		}

		// Simulation
		if node.winner == UNSET {
			moves := state.remainingMoves()

			for len(moves) > 0 {
				// Whoever can make a winning move for them should make it
				m := chooseAWinner(&(state.board), moves, 0-state.player)
				if m < 0 {
					// no winning move for current player
					// Whoever can avoid a losing move for them should avoid it
					acceptableMoves := removeLosses(&(state.board), moves, 0-state.player)
					m = acceptableMoves[rand.Intn(len(acceptableMoves))]
				}

				state.makeMove(m)
				winner := findWinner(&(state.board))
				if winner != UNSET {
					if winner == MAXIMIZER {
						win = true
					}
					break
				}
				cutElement(&moves, m)
			}
		}

		leafCount++

		winIncr := 0.0
		if win {
			winIncr = 1.0
		}

		for node != nil {
			node.visits += 1.0
			node.wins += winIncr
			node.score = node.wins / node.visits
			node = node.parent
		}
	}

	fmt.Printf("after iterations root node %.0f/%.0f\n", root.wins, root.visits)

	fmt.Println("Child nodes:")
	for _, c := range root.childNodes {
		fmt.Printf("\tmove %d, player %d, %.0f/%.0f/%.3f\n", c.move, c.player, c.wins, c.visits, c.score)
	}

	moveNode := root.selectBestChild()
	fmt.Printf("\nbest move node move %d, player %d, %.0f/%.0f/%.3f\n", moveNode.move, moveNode.player, moveNode.wins, moveNode.visits, moveNode.score)
	move = moveNode.move
	score = moveNode.score

	return
}

// cutElement removes element from slice ary
// that has value v. Disorders ary.
func cutElement(ary *[]int, v int) {
	for i, m := range *ary {
		if m == v {
			(*ary)[i] = (*ary)[len(*ary)-1]
			*ary = (*ary)[:len((*ary))-1]
			break
		}
	}
}

// removeLosses takes out losing moves for players,
// returns slice of non-losing moves
func removeLosses(board *[25]int, moves []int, player int) []int {
	if len(moves) < 2 {
		return moves
	}
	acceptableMoves := make([]int, 0, len(moves))
	for _, m := range moves {
		(*board)[m] = player
		x := findWinner(board)
		if x == UNSET || x == player {
			acceptableMoves = append(acceptableMoves, m)
		}
		(*board)[m] = UNSET
	}
	if len(acceptableMoves) > 0 {
		return acceptableMoves
	}
	return moves
}

func (node *Node) AddChild(mv int, state *gameState) *Node {
	// fmt.Printf("node.AddChild(%d, %d)\n", mv, state.player)
	ch := &Node{
		move:         mv,
		parent:       node,
		player:       state.player,
		untriedMoves: state.remainingMoves(),
	}
	node.childNodes = append(node.childNodes, ch)
	// weed out mv as an untried move
	cutElement(&(node.untriedMoves), mv)

	/* fmt.Printf("Child nodes %d:\n", len(node.childNodes))
	for _, n := range node.childNodes {
		fmt.Printf("\tmove %d player %d, %.0f/%.0f/%.3f\n", n.move, n.player, n.wins, n.visits, n.score)
	}
	fmt.Printf("untried moves: %v\n", node.untriedMoves)
	*/

	return ch
}

func (node *Node) selectBestChild() *Node {
	best := node.childNodes[0]
	bestScore := node.childNodes[0].score

	for _, c := range node.childNodes {
		if c.score > bestScore {
			best = c
			bestScore = c.score
		}
	}

	return best
}

// remainingMoves returns an array of all moves left
// unmade on state.board
func (state *gameState) remainingMoves() []int {
	mvs := make([]int, 0, 25)
	j := 0
	for i := 0; i < 25; i++ {
		if state.board[i] == UNSET {
			mvs = append(mvs, i)
			j++
		}
	}
	return mvs
}

func (state *gameState) makeMove(mv int) {
	state.player = 0 - state.player
	state.board[mv] = state.player
}

func (p *MCTS) PrintBoard() {
	fmt.Printf("%s\n", p)
}

func (p *MCTS) SetScores(_ bool) {
}

// FindWinner will return MAXIMIZER or MINIMIZER if somebody won,
// UNSET if nobody wins based on current board.
func (p *MCTS) FindWinner() int {
	return findWinner(&(p.board))
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

// chooseAWinner picks a winning move for player, if there is one, from board.
// Returns -1 if there's no winning move
func chooseAWinner(board *[25]int, moves []int, player int) int {
	var winningMoves []int
	for _, mv := range moves {
		(*board)[mv] = player
		w := findWinner(board)
		(*board)[mv] = UNSET
		if w == player {
			winningMoves = append(winningMoves, mv)
		}
	}

	if len(winningMoves) > 0 {
		if len(winningMoves) == 1 {
			return winningMoves[0]
		}
		return winningMoves[rand.Intn(len(winningMoves))]
	}

	return -1
}

func (p *MCTS) String() string {
	return boardString(p.board)
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
