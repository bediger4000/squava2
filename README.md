# The Game of Squava

## Rules

Squava is a tic-tac-toe variant. Moves are made like tic-tac-toe, except on a
5x5 grid of cells. Players alternate marking cells, conventionally with `X` or
`O`.  Four cells of the same mark in a row (verical, horizontal or diagonal)
wins for the player with that mark. Three cells in a row loses. That is, a
player can win outright, or lose.

You can [play the JavaScript version](https://bediger4000.github.io/squava.html) right now!

You can [play a Monte Carlo Tree Search version](https://bediger4000.github.io/squavam.html) right now!

The rules have an ambiguity, in that it isn't clear what to do if a single marker
fills in a row of 3, say, and a diagonal of 4. Does that player win or lose?

I chose "win", mainly because it's computationally easier to check for 4-in-a-row
as a win separately from 3-in-a-row as a loss. After all, every 4-in-a-row has
3-in-a-row inside it.

Neither player can win until the 7th move (4 for starting play, 3 for the other).
The starting player can win on odd-numbered moves by winning with 4-in-a-row.
The starting player can lose on even-numbered moves by losing with 3-in-a-row.

Similarly, the second player wins on even-numbered moves by getting 4-in-a-row,
or loses on odd-numbered moves with 3-in-a-row.

Peiyan Yang has put together a [program](https://github.com/iForgot321/Squava)
that has solved the game for the first player, if that player plays (2, 2).
