package gamelogic

import (
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Bullet struct {
	X, Y   float64
	VX, VY float64
	Active bool
}

type Zombie struct {
	X, Y   float64
	Health int
	Active bool
	Big    bool
}

type Game struct {
	PlayerImg    *ebiten.Image
	BulletImg    *ebiten.Image
	ZombieImg    *ebiten.Image
	BigZombieImg *ebiten.Image

	X, Y       float64
	DirX, DirY float64

	PlayerHealth int

	Bullets []*Bullet
	Zombies []*Zombie

	ShootTimer         float64
	ZombieTimer        float64
	SpawnTimeLeft      float64
	TimeUntilBigZombie float64
}

func NewGame(playerImg, bulletImg, zombieImg, bigZombieImg *ebiten.Image) *Game {
	if playerImg == nil {
		panic("playerImg is nil")
	}
	if bulletImg == nil {
		panic("bulletImg is nil")
	}
	if zombieImg == nil {
		panic("zombieImg is nil")
	}
	if bigZombieImg == nil {
		panic("bigZombieImg is nil")
	}

	return &Game{
		PlayerImg:          playerImg,
		BulletImg:          bulletImg,
		ZombieImg:          zombieImg,
		BigZombieImg:       bigZombieImg,
		X:                  (640 - 32) / 2,
		Y:                  (480 - 32) / 2,
		DirX:               1,
		DirY:               0,
		PlayerHealth:       5,
		SpawnTimeLeft:      99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999,
		TimeUntilBigZombie: 30.0,
	}
}

func (g *Game) Update() error {
	const playerSpeed = 2.5

	// Player movement.
	dx, dy := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		dy++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		dy--
	}

	if dx != 0 || dy != 0 {
		length := math.Hypot(dx, dy)
		dx /= length
		dy /= length

		g.X += dx * playerSpeed
		g.Y += dy * playerSpeed

		g.DirX = dx
		g.DirY = dy
	}

	// Keep player on screen.
	if g.X < 0 {
		g.X = 0
	}
	if g.X > 640-32 {
		g.X = 640 - 32
	}
	if g.Y < 0 {
		g.Y = 0
	}
	if g.Y > 480-32 {
		g.Y = 480 - 32
	}

	// Update bullets.
	for _, b := range g.Bullets {
		if !b.Active {
			continue
		}

		b.X += b.VX
		b.Y += b.VY

		if b.X < -50 || b.X > 690 || b.Y < -50 || b.Y > 530 {
			b.Active = false
		}
	}

	// Shoot every second.
	g.ShootTimer += 2.5 / 60.0
	if g.ShootTimer >= 1.0 {
		g.ShootTimer = 0
		g.spawnBullet()
	}

	// Spawn regular zombies for the first two minutes.
	if g.SpawnTimeLeft > 0 {
		g.SpawnTimeLeft -= 1.0 / 60.0
		g.ZombieTimer += 1.25 / 60.0

		if g.ZombieTimer >= 2.0 {
			g.ZombieTimer = 0
			g.spawnZombie()
		}
	}

	// Spawn a big zombie every 30 seconds.
	g.TimeUntilBigZombie -= 1.0 / 60.0
	if g.TimeUntilBigZombie <= 0 {
		g.spawnBigZombie()
		g.TimeUntilBigZombie = 30.0
	}

	// Zombies move toward the player.
	const zombieSpeed = 0.8

	for _, z := range g.Zombies {
		if !z.Active {
			continue
		}

		dx := (g.X + 16) - (z.X + 16)
		dy := (g.Y + 16) - (z.Y + 16)
		dist := math.Hypot(dx, dy)

		if dist > 0 {
			dx /= dist
			dy /= dist

			z.X += dx * zombieSpeed
			z.Y += dy * zombieSpeed
		}
	}

	g.checkCollisions()

	return nil
}

func (g *Game) spawnBullet() {
	const bulletSpeed = 10.0
	const scale = 2.0

	playerCenterX := g.X + 16
	playerCenterY := g.Y + 16

	bulletW := float64(g.BulletImg.Bounds().Dx()) * scale
	bulletH := float64(g.BulletImg.Bounds().Dy()) * scale

	startX := playerCenterX - bulletW/2
	startY := playerCenterY - bulletH/2

	for _, b := range g.Bullets {
		if !b.Active {
			b.X = startX
			b.Y = startY
			b.VX = g.DirX * bulletSpeed
			b.VY = g.DirY * bulletSpeed
			b.Active = true
			return
		}
	}

	g.Bullets = append(g.Bullets, &Bullet{
		X:      startX,
		Y:      startY,
		VX:     g.DirX * bulletSpeed,
		VY:     g.DirY * bulletSpeed,
		Active: true,
	})
}

func (g *Game) spawnZombie() {
	const margin = 40

	var x, y float64

	switch rand.Intn(4) {
	case 0:
		x = float64(rand.Intn(640))
		y = -margin
	case 1:
		x = 640 + margin
		y = float64(rand.Intn(480))
	case 2:
		x = float64(rand.Intn(640))
		y = 480 + margin
	case 3:
		x = -margin
		y = float64(rand.Intn(480))
	}

	for _, z := range g.Zombies {
		if !z.Active {
			z.X = x
			z.Y = y
			z.Health = 3
			z.Active = true
			z.Big = false
			return
		}
	}

	g.Zombies = append(g.Zombies, &Zombie{
		X:      x,
		Y:      y,
		Health: 3,
		Active: true,
		Big:    false,
	})
}

func (g *Game) spawnBigZombie() {
	const margin = 40

	var x, y float64

	switch rand.Intn(4) {
	case 0:
		x = float64(rand.Intn(640))
		y = -margin
	case 1:
		x = 640 + margin
		y = float64(rand.Intn(480))
	case 2:
		x = float64(rand.Intn(640))
		y = 480 + margin
	case 3:
		x = -margin
		y = float64(rand.Intn(480))
	}

	for _, z := range g.Zombies {
		if !z.Active {
			z.X = x
			z.Y = y
			z.Health = 10
			z.Active = true
			z.Big = true
			return
		}
	}

	g.Zombies = append(g.Zombies, &Zombie{
		X:      x,
		Y:      y,
		Health: 10,
		Active: true,
		Big:    true,
	})
}

func (g *Game) checkCollisions() {
	const bulletScale = 2.0

	bulletW := float64(g.BulletImg.Bounds().Dx()) * bulletScale
	bulletH := float64(g.BulletImg.Bounds().Dy()) * bulletScale

	// Bullets versus zombies.
	for _, b := range g.Bullets {
		if !b.Active {
			continue
		}

		for _, z := range g.Zombies {
			if !z.Active {
				continue
			}

			zombieW, zombieH := g.zombieSize(z)

			if b.X < z.X+zombieW &&
				b.X+bulletW > z.X &&
				b.Y < z.Y+zombieH &&
				b.Y+bulletH > z.Y {

				b.Active = false
				z.Health--

				if z.Health <= 0 {
					z.Active = false
				}

				break
			}
		}
	}

	// Zombies versus player.
	const playerW, playerH = 32.0, 32.0

	for _, z := range g.Zombies {
		if !z.Active {
			continue
		}

		zombieW, zombieH := g.zombieSize(z)

		if g.X < z.X+zombieW &&
			g.X+playerW > z.X &&
			g.Y < z.Y+zombieH &&
			g.Y+playerH > z.Y {

			g.PlayerHealth--

			if g.PlayerHealth <= 0 {
				print("Game Over! You have been overwhelmed by the undead.\n")
				os.Exit(0)
			}

			// Knock the zombie back.
			dx := (z.X + zombieW/2) - (g.X + playerW/2)
			dy := (z.Y + zombieH/2) - (g.Y + playerH/2)
			dist := math.Hypot(dx, dy)

			if dist > 0 {
				z.X += (dx / dist) * 30
				z.Y += (dy / dist) * 30
			}
		}
	}
}

// zombieSize returns the actual dimensions of the appropriate sprite.
func (g *Game) zombieSize(z *Zombie) (float64, float64) {
	img := g.ZombieImg

	if z.Big {
		img = g.BigZombieImg
	}

	return float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 34,
		G: 139,
		B: 34,
		A: 255,
	})

	// Player.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.X, g.Y)
	screen.DrawImage(g.PlayerImg, op)

	// Bullets (2x).
	for _, b := range g.Bullets {
		if !b.Active {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(b.X, b.Y)
		screen.DrawImage(g.BulletImg, op)
	}

	// Regular and big zombies use separate sprites.
	for _, z := range g.Zombies {
		if !z.Active {
			continue
		}

		img := g.ZombieImg
		if z.Big {
			img = g.BigZombieImg
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(z.X, z.Y)
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}