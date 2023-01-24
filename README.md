# Tic Tac Toe!
commandline demo app built using [Bubbletea](https://github.com/charmbracelet/bubbletea) for playing tic tac toe in the command line!

Currently, I am working on making it playable with a friend, over a simple TCP connection.

## TODO:
- [x] initialize connection in main before starting the game
- [x] randomly decide between X and O
- [ ] handle the TCP moveMsg in Update
- [x] block playing when it's not your turn
- [x] add sending of the message upon playing a move

### Advanced features
- [ ] enable standalone mode
- [ ] add the other players cursor
- [ ] send cursor moves to the other player
- [ ] receive cursor moves and move the cursor

### Refactoring
- [ ] make View use StringBuilder
- [x] extract "x" and "o" as constants
- [ ] decide which functions should be model's methods and which should be standalone.
