package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const NUMBOTS = 2
const BOUTSPERMATCH = 5
const ROUNDSPERBOUT = 1000
const MACFILENAMESIZE = 100
const MAXWEAPONS = 100
const DISPLAYBOUTS = false
const WIDTH = 10
const HEIGHT = 10

var directions = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

type Bot struct {
	name   string
	cmd    string
	x      int
	y      int
	hp     int
	action string
}

type Landmine struct {
	x    int
	y    int
	dead bool
}

type Projectile struct {
	x    int
	y    int
	dir  int
	dead bool
}

type Weapons struct {
	bullets   []*Projectile
	missiles  []*Projectile
	landmines []*Landmine
}

func main() {
	bots := [NUMBOTS]*Bot{
		newBot("huggy", "python ./bots/huggy.py"),
		newBot("nop", "java bots.Nop"),
	}

	botDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	boutsWon := make([]int, NUMBOTS)
	matchesWon := make([]int, NUMBOTS)

	for i := 0; i < NUMBOTS; i++ {
		boutsWon[i] = 0
		matchesWon[i] = 0
	}

	for b1 := 0; b1 < NUMBOTS-1; b1++ {
		for b2 := b1 + 1; b2 < NUMBOTS; b2++ {
			bot1, bot2 := bots[b1], bots[b2]
			bot1Bouts, bot2Bouts := 0, 0
			fmt.Printf("%s vs %s\n", bots[b1].name, bots[b2].name)
			for bout := 0; bout < BOUTSPERMATCH; bout++ {
				fmt.Printf("%d ", bout)
				bot1.x, bot1.y, bot1.hp = 1, 1, 10
				bot2.x, bot2.y, bot2.hp = WIDTH-2, HEIGHT-2, 10
				weapons := initWeapons()
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
						// Move bot 1
						bot1Moved := false
						dx, dy := 0, 0
						dx = getXMove(bot1.action)
						dy = getYMove(bot1.action)
						if newPosInBounds(bot1.x, bot1.y, dx, dy) {
							if !(bot1.x+dx == bot2.x) || !(bot1.y+dy == bot2.y) {
								bot1Moved = true
								bot1.x += dx
								bot1.y += dy
							}
						}

						// Move bot 2
						dx = getXMove(bot2.action)
						dy = getYMove(bot2.action)
						if newPosInBounds(bot2.x, bot2.y, dx, dy) {
							if !(bot2.x+dx == bot1.x) || !(bot2.y+dy == bot1.y) {
								bot2.x += dx
								bot2.y += dy
							}
						}

						// If bot2 was in the way the first time, let bot 1 try again
						if !bot1Moved {
							dx = getXMove(bot1.action)
							dy = getYMove(bot1.action)
							if newPosInBounds(bot1.x, bot1.y, dx, dy) {
								if !(bot1.x+dx == bot2.x) || !(bot1.y+dy == bot2.y) {
									bot1.x += dx
									bot1.y += dy
								}
							}
						}

						for _, l := range weapons.landmines {
							if directHit(bot1, l.x, l.y) {
								bot1.hp -= 2
								if inShrapnelRange(bot2, l.x, l.y) {
									bot2.hp -= 1
								}
								l.dead = true
							} else if directHit(bot2, l.x, l.y) {
								bot2.hp -= 2
								if inShrapnelRange(bot1, l.x, l.y) {
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

					for _, b := range weapons.bullets {
						dx := getXMove(directions[b.dir])
						dy := getYMove(directions[b.dir])

						for moves := 0; moves < 3; moves++ {
							if newPosInBounds(b.x, b.y, dx, dy) {
								b.x += dx
								b.y += dy

								if directHit(bot1, b.x, b.y) {
									bot1.hp -= 1
									b.dead = true
									break
								}

								if directHit(bot2, b.x, b.y) {
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

					for _, m := range weapons.missiles {
						dx := getXMove(directions[m.dir])
						dy := getYMove(directions[m.dir])

						for moves := 0; moves < 2; moves++ {
							if newPosInBounds(m.x, m.y, dx, dy) {
								m.x += dx
								m.y += dy

								if directHit(bot1, m.x, m.y) {
									bot1.hp -= 3
									if inShrapnelRange(bot2, m.x, m.y) {
										bot2.hp -= 1
									}
									m.dead = true
									break
								}

								if directHit(bot2, m.x, m.y) {
									bot2.hp -= 3
									if inShrapnelRange(bot1, m.x, m.y) {
										bot1.hp -= 1
									}
									m.dead = true
									break
								}
							} else {
								if inShrapnelRange(bot1, m.x, m.y) {
									bot1.hp -= 1
								}
								if inShrapnelRange(bot2, m.x, m.y) {
									bot2.hp -= 1
								}
								m.dead = true
							}
						}
					}

					n = 0
					for _, m := range weapons.missiles {
						if !m.dead {
							weapons.missiles[n] = m
							n += 1
						}
					}
					weapons.missiles = weapons.missiles[:n]

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
	for i := 0; i < NUMBOTS; i++ {
		fmt.Printf("%s\t%d\t%d\n", bots[i].name, matchesWon[i], boutsWon[i])
	}
}

func newBot(name string, cmd string) *Bot {
	return &Bot{name, cmd, 0, 0, 10, "NO"}
}

func initWeapons() *Weapons {
	return &Weapons{
		make([]*Projectile, 0),
		make([]*Projectile, 0),
		make([]*Landmine, 0),
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
		grid[b.y][b.x] = 'B'
		stats.WriteString(fmt.Sprintf("B x=%d y=%d dir=%s\n", b.x, b.y, directions[b.dir]))
	}

	for i := range weapons.missiles {
		m := weapons.missiles[i]
		grid[m.y][m.x] = 'M'
		stats.WriteString(fmt.Sprintf("M x=%d y=%d dir=%s\n", m.x, m.y, directions[m.dir]))
	}

	grid[bot.y][bot.x] = 'Y'
	grid[enemy.y][enemy.x] = 'X'
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

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func getXMove(dir string) int {
	if dir == "NE" {
		return 1
	} else if dir == "E" {
		return 1
	} else if dir == "SE" {
		return 1
	} else if dir == "SW" {
		return -1
	} else if dir == "W" {
		return -1
	} else if dir == "NW" {
		return -1
	}
	return 0
}

func getYMove(dir string) int {
	if dir == "S" {
		return 1
	} else if dir == "SE" {
		return 1
	} else if dir == "SW" {
		return 1
	} else if dir == "N" {
		return -1
	} else if dir == "NE" {
		return -1
	} else if dir == "NW" {
		return -1
	}
	return 0
}

func newPosInBounds(ox int, oy int, dx int, dy int) bool {
	return ox+dx >= 0 && ox+dx < WIDTH && oy+dy >= 0 && oy+dy < HEIGHT
}

func directHit(b *Bot, x int, y int) bool {
	return b.x == x && b.y == y
}

func landmineCollision(l1 *Landmine, l2 *Landmine) bool {
	return l1.x == l2.x && l1.y == l2.y
}

func inShrapnelRange(b *Bot, x int, y int) bool {
	return abs(b.x-x) < 2 && abs(b.y-y) < 2
}

func directionToInt(dir string) int {
	for i := range directions {
		if dir == directions[i] {
			return i
		}
	}

	return -1
}

func newProjectile(x int, y int, dir string) *Projectile {
	return &Projectile{x, y, directionToInt(dir), false}
}

func deployWeapons(bot *Bot, enemy *Bot, weapons *Weapons) {
	fields := strings.Fields(bot.action)
	weapon := fields[0]
	if weapon == "B" {
		bullet := newProjectile(bot.x, bot.y, fields[1])
		if bullet.dir != -1 {
			weapons.bullets = append(weapons.bullets, bullet)
		}
	} else if weapon == "M" {
		missile := newProjectile(bot.x, bot.y, fields[1])
		if missile.dir != 1 {
			weapons.missiles = append(weapons.missiles, missile)
		}
	} else if weapon == "L" {
		dir := fields[1]
		if newPosInBounds(bot.x, bot.y, getXMove(dir), getYMove(dir)) {
			var landmine *Landmine
			landmine.x = bot.x + getXMove(dir)
			landmine.y = bot.y + getYMove(dir)
			var collision bool
			for i := range weapons.landmines {
				l := weapons.landmines[i]
				if landmineCollision(landmine, l) {
					if inShrapnelRange(bot, l.x, l.y) {
						bot.hp -= 1
					}
					if inShrapnelRange(enemy, l.x, l.y) {
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

func clearArena(arena [][]rune) {
	for r := range arena {
		for c := range arena[r] {
			arena[r][c] = '.'
		}
	}
}
