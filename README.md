# BattleBots
**A robot coding tournament.** 

## The Game
The aim of the game is to write a bot that will destroy all the other bots in the tournament.

Bots will be pitted against one another in a 1v1 round-robin tournament. Each match will consist of 5 bouts, with the winning bot being that with the highest HP at the end of the bout. The bot with the most bout wins takes the match, and the bot with the most match wins at the end of the rounament takes the title of champion.

A bout is turn based. At the top of the each turn, the current state of the battlefield will be rendered as an ASCII string, which will be passed to each bot as a command-line argument. Each bot must respond by printing their desired action to `stdout`, to be read by the control program, which will attempt to execute these actions one after the other. Once the bots have moved, weapons have been deployed, and HP updated, the control program will proceed to the next turn. There is a limit of 1000 rounds, but the bout will end early if a bot reaches 0 HP. If both bots have equal HP by the 1000th round, the bout is a draw, and no points are scored.

### The Arena
At the start of each round, each bot will recieve an input string formatted as follows:

```
X.....LLL.
..........
..........
..........
M.........
..........
..........
..........
..........
...B.....Y
Y hp=10
X hp=7
B x=3 y=9 dir=W
B x=3 y=9 dir=S
M x=0 y=4 dir=S
L x=6 y=0
L x=7 y=0
L x=8 y=0
```

The first ten lines (ending with a newline character, `\n`) are a graphical representation of the arena the bots are fighting in. The meanings of the symbols are as follows:

* `.` = an empty square
* `Y` = you - the bot recieving this input
* `X` = the opponent bot
* `L` = a landmine
* `B` = a bullet
* `M` = a missile

If several entities are overlapping, only one symbol will be shown, buit the positions of _all_ projectiles, including those not visible on the grid, are included in the stats readout below the grid. Projectiles and mines will never overlop with bots, as they would have already been destroyed, and the bot damaged.

The stats of each entity in the arena are listed below the grid. Each entity is represented by a line of space-separated fields. The first field is the name of the entity (`Y` = you, `X` = the enemy, `L`/`B`/`M` = a weapon), and each following field is a key/value pair denoting an attribute of that entity. For the bots, this is just their HP. For weapons, attributes include their coordinates and direction of travel.

### Actions

Once your bot has parsed the input string and decided its next move, it should print a string to `stdout` indicating its desired action. The actions available are as follows:

#### Movement
* `N` - Move one square north.
* `NE` - Move one square north east.
* `E` - Move one square east.
* `SE` - Move one square south east.
* `S` - Move one square south.
* `SW` - Move one square south west.
* `W` - Move one square west.
* `NW` - Move one square north west.

#### Weapons
* `L <dir>` - Deploy a landmine in the given direction. The `L` should be followed by a single space, then a movement direction as above. The landmine will be placed in the same square that the bot would have arrived in by moving in the given direction. **A landmine deals 2 points of damage to a bot that stands on it, and 1 point to a bot standing in an adjacent square.** A landmine will also be triggered if a second landmine is deployed in the same square. The landmine will remain stationary until it is triggered, and it is destroyed.
* `B <dir>` - Fire a bullet travelling in the given direction. **The bullet will move 3 squares in the given direction each round, dealing 1 point of damage to any bot it hits.** The bullet is destroyed when it hits a bot or the edge of the arena.
* `M <dir>` - Fire a missile travelling in the given direction. **The missile will move 2 squares in the given direction each round, dealing 2 points of damage to any bot it hits, and 1 point of splash damage to any bots in adjacent squares.** The missile is destroyed when it hits a bot or the edge of the arena.
* `P` - Deploy an electro-magnetic pulse (EMP). **The EMP disables the movement circuits of each bot for 2 rounds, meaning the bots can't move during that time. However, both bots can still deploy weapons while they're paralyzed.** _Deploying an EMP will cost 1 HP for the bot that uses it._

#### Do Nothing
* `NO` - A bot that outputs `NO` will take no action. In fact, any output that does not satisfy the requirements of the actions listed above will result in no action being carried out by the bot that turn.

## Competing

To enter the tournament, all you need to do is write your bot, and create a pull request to add the bot to the `bots` folder in this repository and the list of bots in `battlebots.go` (line 50). If your bot has a build step, you'll also need to add the appropriate steps to `build.sh`.

Each bot should consist of a single source file with a unique name. See `huggy.py` and `Nop.java` for examples. The bot can be written in any language, as long as it can be executed from an Ubuntu command-line environment with a single command, and can take command-line arguments. It should be noted that extra configuration is required to support new languages. The list of supported languages is provided below. If you wish to write a bot in another language, get in touch and support for that language can be arranged.

If your bot requires a build step, you should add the commands required to build the source files to `build.sh`. The compilation process may involve the creation of additional files, but only the compiled program should remain, in the 'bots' folder, once `build.sh` is run. The Java bot "Nop" is provided as an example of a compiled battle bot.

The final step in registering a bot (before opening a PR), is to add a new Bot object to the `bots` array in `battlebots.go` (line 50). You should provide a unique name, and a single command that can be used to execute your bot. The commands will be executed in the project root (i.e. the same directory as `battlebots.go`), and therefore any paths should be relative to this.

### Currently Supported Languages
* Go (1.18.5)
* JavaScript (Node.js, version 18.10.x)
* Python (3.10)
* Java (18, Azul Zulu OpenJDK)

### Notes
* Only one submission is allowed per competitor.
* You may update your bot at any time, by submitting additional pull requests.
* Your bot may not use any additional external dependencies. You must use only the core libraries and built-ins of your language of choice.
* Your bot may not read, write or execute other files, except its own source code.
* The example bots, huggy and Nop, have been included by way of example. They will remain in the `bots` folder, but will not be compete in the tournament once at least two player-made bots have been submitted.

## A Note on Changes
Until a number of entries have been submitted, its hard to verify the control program for balance, fairness and fun. There may be minor changes to the rules of the game, surrounding the weapons and actions available to the bots, the size and contents of the arena, and the contents of the input string. The general format of the input string will not change. Keep an eye on the repo for such changes, but they will be announced where possible and owners of broken bots will be contacted to request fixes.

## The Results
The tournament will be run every time a new bot is submitted. A Github Action is used to install the necessary language interpreters/compilers, the `build.sh` script is executed, and the control program starts simulating the tournament.

In a future update, the results table will automatically be available for viewing online.