# Algorithmic Players for the Game of Squava

I wrote some algorithmic players for the game of [Squava]().
I've written two variations on Alpha-beta minimaxing,
and two variations on Monte Carlo Tree Search.

This is my second attempt, my [first attempt](https://github.com/bediger4000/squava) is full of
cruft and mistakes.
I did not implement Monte Carlo Tree Search correctly, for example.

## Rules of the game

Squava is a tic-tac-toe variant. Moves are made like tic-tac-toe, except on a
5x5 grid of cells.
Players alternate marking cells, conventionally with `X` or `O`.
Four cells of the same mark in a row (vertical, horizontal or diagonal)
wins for the player with that mark. Three cells in a row loses. That is, a
player can win outright, or lose outright.

The rules have an ambiguity, in that it isn't clear what to do if a single marker
fills in a row of 3, say, and a diagonal of 4. Does that player win or lose?

I chose "win", mainly because it's computationally easier to check for 4-in-a-row
as a win separately from 3-in-a-row as a loss. After all, every 4-in-a-row has
3-in-a-row inside it.

Neither player can win until the 7th move (4 for starting player, 3 for the other).
The starting player can win on odd-numbered moves by winning with 4-in-a-row.
The starting player can lose on even-numbered moves by losing with 3-in-a-row.

Similarly, the second player wins on even-numbered moves by getting 4-in-a-row,
or loses on odd-numbered moves with 3-in-a-row.

A game has a maximum of 25 moves.
I believe that the starting player always loses 25-move games, but I can't prove it.

I've not see it written down anywhere, but "squava" is probably "square yavalath",
[Yavalath](http://cambolbro.com/games/yavalath/)
being the inspiration  for squava.

### Cat Games

Can "cat" get the game, as in ordinary tic-tac-toe?
I don't know.
I've got lots of 25-move example games, all of them lost by the first player.

## Playing the game

I wrote an interactive player,
and a program that matches two algorithmic players against each other.

### Interactive game

```
$ go build sqv.go
$ ./sqv -t M -C
X (MCTS/Plain) <3,2> (9255) [500000] 1.940846297s
   0 1 2 3 4
0  _ _ _ _ _ 
1  _ _ _ _ _ 
2  _ _ _ _ _ 
3  _ _ X _ _ 
4  _ _ _ _ _ 

Your move: 1 4
   0 1 2 3 4
0  _ _ _ _ _ 
1  _ _ _ _ O 
2  _ _ _ _ _ 
3  _ _ X _ _ 
4  _ _ _ _ _ 
```

In the above game fragment, a freshly-compiled `sqv` program ran the plain
Monte Carlo Tree Search player to find the first move, 3,2, the `X`.

The human chose the next move, 1,4, signified by an 'O'

### Inter-algorithm games

```
$ go build playoff.go
$ ./playoff -1 U -2 M -n 2
0    MCTS/UCB1    MCTS/Plain   20   -1   40.23  2,2 3,1 3,2 0,2 1,1 0,1 4,3 0,4 0,3 3,4 4,4 3,3 1,0 2,1+ 4,1 4,2+ 2,4 1,2+ 4,0 2,3+ 
1    MCTS/UCB1    MCTS/Plain   9    1    28.08  1,1 2,2 4,1 0,4 1,4 4,3 3+,1 0,1 2+,1 
```

Above, we have 2 games, first move by MCTS with UCB1 move selection
versus plain old MCTS.
The first game took 20 moves, the plain MCTS 2nd player won in 40.23 seconds
The second game took 9 moves, the first player won in 28.08 seconds.
The series of moves are listed after the per-game statistics.

You can re-use the series of moves in two ways:

1. The `recreate` program accepts either a file name with the string of
moves as contents, or the string of moves on the command line:
   * `./recreate '1,1 2,2 4,1 0,4 1,4 4,3 3+,1 0,1 2+,1'`
2. The  `sqv` program can accept a partial game on the command line,
setting up the board for the algorithmic player:
   * `./sqv -t U -p '1,1 2,2 4,1 0,4 1,4'`

You can investigate which move the algorithmic players make in a given
situation with the `-p 'x,y x,y...'` partial game.

## Software Engineering

#### `Player` interface

The `Player` interface describes code that has an internal representation of
a squava game and can choose a move based on that internal representation.
The "board" isn't specified externally, each algorithm can have an internal
representation of the game board customized for itself.
That means that each implementation
of `Player` has to have its own way to find a "winner" (or loser),
and to make a human readable representation of its internal board state.
Method `MakeMove` has player type so that driver programs
(like `sqv.go`)
can set a board to some desired
configuration before letting the algorithm choose a move.

```go
type Player interface {
    Name() string
    MakeMove(int, int, int)           // x,y coords, type of player (MAXIMIZER, MINIMIZER
    ChooseMove() (int, int, int, int) // x,y coords of move, value, leaf node count
    FindWinner() int
    String() string // human readable formatted board
```

Code for all algorithmic players lives in the same package.
I had them in separate packages in my first attempt at algorithmic players.
It seemed like that arrangement required lots of redundant code,
so this time around I put them all in the same package.
This, too, seems like a mistake, as I had to rename some of the variables
holding things like "which sets of slots make a 4-in-a-row".
I don't know what the solution for this is.

The Alpha-beta minimaxing algorithms have an internal representation
of the board that looks like this:

```go
type board [5][5]int
```

The Monte Carlo Tree Search variants have an internal board representation
that looks like this:

```go
var board  [25]int
```

Go the programming language really doesn't have 2-D arrays,
so I worried about the efficiency of `[5][5]int` in the Alpha-beta
algorithms.

There's a mismatch between what a human (me) wants to give the interactive
program (a row and a column number) versus what the Monte Carlo Tree Search
programs use (`[25]int` array.
It was also harder to specify the array indexes that make a 4-in-a-row win,
or a 3-in-a-row loss.
Again, I don't  know which alternative really is better.

## Other Investigations

Peiyan Yang has put together a [program](https://github.com/iForgot321/Squava)
that has solved the game for the first player.
