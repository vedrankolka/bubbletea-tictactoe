# Tic Tac Toe!
Command line demo app built using [Bubbletea](https://github.com/charmbracelet/bubbletea) for playing tic tac toe in the command line!

Currently, it is playable with a friend, over a simple TCP connection.

![](https://github.com/vedrankolka/bubbletea-tictactoe/blob/develop/gifs/demo-wait.gif)

![](https://github.com/vedrankolka/bubbletea-tictactoe/blob/develop/gifs/demo-dail.gif)

## Getting started
TODO: write this part.

## TODO:
- [x] initialize connection in main before starting the game
- [x] randomly decide between X and O
- [x] handle the TCP moveMsg in Update
- [x] block playing when it's not your turn
- [x] add sending of the message upon playing a move
- [x] add bottom panel to View printing the player and playerTurn
- [x] detect and report a tie
- [ ] add getting started

### Advanced features
- [ ] enable standalone mode
- [ ] add the other players cursor
- [ ] send cursor moves to the other player
- [ ] receive cursor moves and move the cursor
- [ ] have real random X, O assignment*
- [ ] detect a tie even sooner (if it not possible for either side to win)

*Each player chooses a nonce, sends the hash to the other, once the player receives the hash, sends its nonce, receives the opponent's nonce, it checks if the hash is from the received nonce, (nonce1 + nonce2) % 2 == 0 -> X else O

### Refactoring
- [ ] make View use StringBuilder
- [x] extract "x" and "o" as constants
- [ ] decide which functions should be model's methods and which should be standalone
- [ ] make each model's method change only one model property (meaning extract switchPlayer as a method and don't make handlePlayerEnter do it)
- [ ] add more log statements
- [ ] don't send moves when no action was taken
- [x] capitalize x and o to X and O
