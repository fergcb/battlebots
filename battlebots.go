package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const BOUTSPERMATCH = 5
const ROUNDSPERBOUT = 1000
const MACFILENAMESIZE = 100
const MAXWEAPONS = 100
const DISPLAYBOUTS = false
const WIDTH = 10
const HEIGHT = 10

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func max(a int, b int) int {
	if a < b {
		return b
	}

	return a
}

type Vector struct {
	x int
	y int
}

var directions = map[string]Vector{
	"N":  {0, -1},
	"NE": {1, -1},
	"E":  {1, 0},
	"SE": {1, 1},
	"S":  {0, 1},
	"SW": {-1, 1},
	"W":  {-1, 0},
	"NW": {-1, -1},
}

func (v1 Vector) equals(v2 Vector) bool {
	return v1.x == v2.x && v1.y == v2.y
}

func (v1 Vector) dist(v2 Vector) int {
	return max(abs(v1.x-v2.x), abs(v1.y-v2.y))
}

func (v1 Vector) isAdjacentTo(v2 Vector) bool {
	return v1.dist(v2) < 2
}

func (v1 Vector) add(v2 Vector) Vector {
	return Vector{v1.x + v2.x, v1.y + v2.y}
}

func (v Vector) inBounds() bool {
	return v.x > 0 && v.y > 0 && v.x < WIDTH && v.y < HEIGHT
}

type Bot struct {
	name   string
	cmd    string
	pos    Vector
	hp     int
	action string
}

func newBot(name string, cmd string) *Bot {
	return &Bot{name, cmd, Vector{}, 10, "NO"}
}

type Landmine struct {
	pos  Vector
	dead bool
}
type Projectile struct {
	pos  Vector
	dir  string
	dead bool
}

type Weapons struct {
	bullets   []*Projectile
	missiles  []*Projectile
	landmines []*Landmine
}

func newWeapons() *Weapons {
	return &Weapons{
		make([]*Projectile, 0),
		make([]*Projectile, 0),
		make([]*Landmine, 0),
	}
}

func main() {
	bots := []*Bot{
		newBot("huggy", "python ./bots/huggy.py"),
		newBot("nop", "java bots.Nop"),
	}

	numBots := len(bots)

	botDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	boutsWon := make([]int, numBots)
	matchesWon := make([]int, numBots)

	for i := 0; i < numBots; i++ {
		boutsWon[i] = 0
		matchesWon[i] = 0
	}

	for b1 := 0; b1 < numBots-1; b1++ {
		for b2 := b1 + 1; b2 < numBots; b2++ {
			bot1, bot2 := bots[b1], bots[b2]
			bot1Bouts, bot2Bouts := 0, 0

			fmt.Printf("%s vs %s\n", bots[b1].name, bots[b2].name)

			for bout := 0; bout < BOUTSPERMATCH; bout++ {
				fmt.Printf("%d ", bout)

				bot1.pos, bot1.hp = Vector{1, 1}, 10
				bot2.pos, bot2.hp = Vector{WIDTH - 2, HEIGHT - 2}, 10

				weapons := newWeapons()

				paralyzedRoundsRemaining := 0

				for round := 0; round < ROUNDSPERBOUT; round++ {
					bot1Arena := drawArena(bot1, bot2, weapons)
					bot2Arena := drawArena(bot2, bot1, weapons)

					bot1.action = runBot(bot1, bot1Arena, botDir)
					bot2.action = runBot(bot2, bot2Arena, botDir)

					if DISPLAYBOUTS {
						fmt.Printf("Round: %d\n", round)
						fmt.Println(bot1Arena)
					}

					if paralyzedRoundsRemaining == 0 {
						moveBots(bot1, bot2)
						checkLandmines(bot1, bot2, weapons)
					} else {
						paralyzedRoundsRemaining -= 1
						if DISPLAYBOUTS {
							fmt.Println("The bots are paralyzed.")
						}
					}

					if bot1.action == "P" {
						paralyzedRoundsRemaining = 2
						bot1.hp -= 1
					} else if bot2.action == "P" {
						paralyzedRoundsRemaining = 2
						bot2.hp -= 1
					}

					deployWeapons(bot1, bot2, weapons)
					deployWeapons(bot2, bot1, weapons)

					moveBullets(bot1, bot2, weapons)
					moveMissiles(bot1, bot2, weapons)

					if bot1.hp < 1 || bot2.hp < 1 {
						break
					}
				}

				if bot1.hp < bot2.hp {
					bot2Bouts += 1
					boutsWon[b2] += 1
				} else if bot2.hp < bot1.hp {
					bot1Bouts += 1
					boutsWon[b1] += 1
				}
			}

			if bot1Bouts < bot2Bouts {
				matchesWon[b2] += 1
			} else if bot2Bouts < bot1Bouts {
				matchesWon[b1] += 1
			}

			fmt.Println()
		}
	}

	fmt.Println("\nResults:")
	fmt.Println("Bot\tMatches\tBouts")
	for i := 0; i < numBots; i++ {
		fmt.Printf("%s\t%d\t%d\n", bots[i].name, matchesWon[i], boutsWon[i])
	}
}

func initArena() [][]rune {
	arena := make([][]rune, HEIGHT)
	for r := range arena {
		arena[r] = make([]rune, WIDTH)
	}
	clearArena(arena)
	return arena
}

func drawArena(bot *Bot, enemy *Bot, weapons *Weapons) string {
	var arena strings.Builder
	var stats strings.Builder

	grid := initArena()

	stats.WriteString(fmt.Sprintf("Y hp=%d\n", bot.hp))
	stats.WriteString(fmt.Sprintf("X hp=%d\n", enemy.hp))

	clearArena(grid)
	for i := range weapons.bullets {
		b := weapons.bullets[i]
		grid[b.pos.y][b.pos.x] = 'B'
		stats.WriteString(fmt.Sprintf("B x=%d y=%d dir=%s\n", b.pos.x, b.pos.y, b.dir))
	}

	for i := range weapons.missiles {
		m := weapons.missiles[i]
		grid[m.pos.y][m.pos.x] = 'M'
		stats.WriteString(fmt.Sprintf("M x=%d y=%d dir=%s\n", m.pos.x, m.pos.y, m.dir))
	}

	grid[bot.pos.y][bot.pos.x] = 'Y'
	grid[enemy.pos.y][enemy.pos.x] = 'X'
	// fmt.Println(bot.x, bot.y)
	// fmt.Println(enemy.x, enemy.y)

	for r := range grid {
		for c := range grid[r] {
			arena.WriteRune(grid[r][c])
		}
		arena.WriteRune('\n')
	}

	return arena.String() + stats.String()
}

func runBot(bot *Bot, arena string, dir string) string {
	fields := strings.Fields(bot.cmd)
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Dir = dir
	cmd.Args = append(cmd.Args, arena)
	stdout, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(stdout))
}

func moveBots(bot1 *Bot, bot2 *Bot) {
	// Move bot 1
	bot1Moved := false
	v := directions[bot1.action]
	newPos := bot1.pos.add(v)
	if newPos.inBounds() && !newPos.equals(bot2.pos) {
		bot1Moved = true
		bot1.pos = newPos
	}

	// Move bot 2
	v = directions[bot2.action]
	newPos = bot2.pos.add(v)
	if newPos.inBounds() && !newPos.equals(bot1.pos) {
		bot2.pos = newPos
	}

	// If bot2 was in the way the first time, let bot 1 try again
	if !bot1Moved {
		v := directions[bot1.action]
		newPos := bot1.pos.add(v)
		if newPos.inBounds() && !newPos.equals(bot2.pos) {
			bot1.pos = newPos
		}
	}
}

func checkLandmines(bot1 *Bot, bot2 *Bot, weapons *Weapons) {
	for _, l := range weapons.landmines {
		if bot1.pos.equals(l.pos) {
			bot1.hp -= 2
			if bot2.pos.isAdjacentTo(l.pos) {
				bot2.hp -= 1
			}
			l.dead = true
		} else if bot2.pos.equals(l.pos) {
			bot2.hp -= 2
			if bot1.pos.isAdjacentTo(l.pos) {
				bot1.hp -= 1
			}
			l.dead = true
		}
	}

	n := 0
	for _, l := range weapons.landmines {
		if !l.dead {
			weapons.landmines[n] = l
			n += 1
		}
	}
	weapons.landmines = weapons.landmines[:n]
}

func deployWeapons(bot *Bot, enemy *Bot, weapons *Weapons) {
	fields := strings.Fields(bot.action)

	if len(fields) != 2 {
		return
	}

	if _, ok := directions[fields[1]]; !ok {
		return
	}

	weapon := fields[0]
	dir := fields[1]

	if weapon == "B" {
		bullet := &Projectile{bot.pos, dir, false}
		weapons.bullets = append(weapons.bullets, bullet)
	} else if weapon == "M" {
		missile := &Projectile{bot.pos, dir, false}
		weapons.missiles = append(weapons.missiles, missile)
	} else if weapon == "L" {
		targetPos := bot.pos.add(directions[dir])
		if targetPos.inBounds() {
			var landmine *Landmine

			landmine.pos = targetPos

			var collision bool
			for i := range weapons.landmines {
				l := weapons.landmines[i]
				if landmine.pos.equals(l.pos) {
					if bot.pos.isAdjacentTo(l.pos) {
						bot.hp -= 1
					}
					if enemy.pos.isAdjacentTo(l.pos) {
						enemy.hp -= 1
					}

					collision = true
					weapons.landmines[i].dead = true
					break
				}
			}

			if !collision {
				weapons.landmines = append(weapons.landmines, landmine)
			}
		}
	}
}

func moveBullets(bot1 *Bot, bot2 *Bot, weapons *Weapons) {
	for _, b := range weapons.bullets {
		v := directions[b.dir]
		for moves := 0; moves < 3; moves++ {
			newPos := b.pos.add(v)
			if newPos.inBounds() {
				b.pos = newPos

				if bot1.pos.equals(b.pos) {
					bot1.hp -= 1
					b.dead = true
					break
				}

				if bot2.pos.equals(b.pos) {
					bot2.hp -= 1
					b.dead = true
					break
				}
			} else {
				b.dead = true
			}
		}
	}

	n := 0
	for _, b := range weapons.bullets {
		if !b.dead {
			weapons.bullets[n] = b
			n += 1
		}
	}
	weapons.bullets = weapons.bullets[:n]
}

func moveMissiles(bot1 *Bot, bot2 *Bot, weapons *Weapons) {
	for _, m := range weapons.missiles {
		v := directions[m.dir]

		for moves := 0; moves < 2; moves++ {
			newPos := m.pos.add(v)
			if newPos.inBounds() {
				m.pos = newPos

				if bot1.pos.equals(m.pos) {
					bot1.hp -= 3
					if bot2.pos.isAdjacentTo(m.pos) {
						bot2.hp -= 1
					}
					m.dead = true
					break
				}

				if bot2.pos.equals(m.pos) {
					bot2.hp -= 3
					if bot1.pos.isAdjacentTo(m.pos) {
						bot1.hp -= 1
					}
					m.dead = true
					break
				}
			} else {
				if bot1.pos.isAdjacentTo(m.pos) {
					bot1.hp -= 1
				}
				if bot2.pos.isAdjacentTo(m.pos) {
					bot2.hp -= 1
				}
				m.dead = true
			}
		}
	}

	n := 0
	for _, m := range weapons.missiles {
		if !m.dead {
			weapons.missiles[n] = m
			n += 1
		}
	}
	weapons.missiles = weapons.missiles[:n]

}

func clearArena(arena [][]rune) {
	for r := range arena {
		for c := range arena[r] {
			arena[r][c] = '.'
		}
	}
}
