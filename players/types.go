package players

// Player interface describes something that has an internal representation of
// a squava game and can choose a move based on that internal representation.
// The "board" type isn't specified externally, but that means that each implementation
// of Player has to have its own way to find a "winner" (loser), and to make a human
// readable representation of its internal board state.
// MakeMove has player type so that a driver program can set a board to some desired
// config before letting the Player choose a move.
type Player interface {
	Name() string
	MakeMove(int, int, int)           // x,y coords, type of player (MAXIMIZER, MINIMIZER)
	ChooseMove() (int, int, int, int) // x,y coords of move, value, leaf node count
	FindWinner() int
	String() string // human readable formatted board
	// Options(...string) // name=value pairs particular to an implementation
}

// Manifest constants to improve understanding
const (
	MAXIMIZER = 1
	MINIMIZER = -1
	UNSET     = 0
)
