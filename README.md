# ☄️ ASTRO

> A love letter to the 1979 arcade classic — rebuilt from scratch with modern web tech.

![Game Status](https://img.shields.io/badge/status-in%20development-yellow)
![License](https://img.shields.io/badge/license-MIT-blue)

---

## What Is This?

A 2D physics-based space shooter inspired by Atari's *Asteroids*. You pilot a lone ship through an endless field of drifting rocks, blasting them into smaller and smaller fragments until the screen is clear — then doing it all over again, faster, meaner, with less mercy.

No health bars. No checkpoints. Pure arcade survival.

---

## Gameplay

- **Rotate** your ship and fire thrust to drift through space
- **Shoot** asteroids — large ones split into medium, medium into small
- **Survive** each wave to face a harder one
- **Don't stop moving** — momentum is your friend and your enemy

### Controls

| Action         | Key              |
|----------------|------------------|
| Rotate Left    | `←` / `A`        |
| Rotate Right   | `→` / `D`        |
| Thrust         | `↑` / `W`        |
| Fire           | `Space`          |
| Pause          | `Esc` / `P`      |

---

## Features

### Core
- Fixed-timestep game loop with delta time — smooth at any frame rate
- Newtonian-style inertia and drag — the ship drifts, not teleports
- Screen wrapping — exit left, enter right. Classic.
- Circle and polygon collision detection with broadphase optimisation

### Asteroids
- Three sizes: Large → Medium → Small, each with randomised polygon shapes
- Fragments inherit parent velocity with randomised spread on split
- Wave difficulty scales with each cleared screen

### Player
- Respawn with invincibility window after death
- Extra lives awarded at score milestones
- Classic 3-life start

### Enemies
- UFO saucer spawns periodically — Large (inaccurate fire) and Small (aims at you)

### Feel
- Particle effects on all destruction events
- Screen shake on impact
- Full sound design: thrust, fire, explosions, heartbeat pulse that tracks remaining asteroids
- Persistent high score

---

## Milestones

Development is tracked as GitHub Issues across 8 milestones:

| # | Milestone | Status |
|---|-----------|--------|
| 1 | Core Engine Foundation | 🔵 Planned |
| 2 | Player Ship | 🔵 Planned |
| 3 | Weapons & Projectiles | 🔵 Planned |
| 4 | Asteroids | 🔵 Planned |
| 5 | Scoring & Progression | 🔵 Planned |
| 6 | Game Feel & Polish | 🔵 Planned |
| 7 | Menus & Game States | 🔵 Planned |
| 8 | Stretch Goals | 🔵 Planned |

See [Issues](../../issues) for the full breakdown.

---

## Architecture

```
asteroids/
├── main.go               # Entry point, Ebitengine bootstrap
├── go.mod
├── go.sum
├── game/
│   ├── game.go           # Core Game struct — implements ebiten.Game
│   ├── state.go          # Game state machine (title, playing, paused, etc.)
│   └── wave.go           # Wave progression logic
├── entity/
│   ├── ship.go
│   ├── asteroid.go
│   ├── bullet.go
│   ├── ufo.go
│   └── particle.go
├── physics/
│   ├── vec2.go           # 2D vector math
│   ├── body.go           # Rigid body state & integration
│   └── collision.go      # Detection + event dispatch
├── input/
│   └── input.go          # Keyboard polling via ebiten
├── audio/
│   └── audio.go          # Sound effects & heartbeat
├── ui/
│   ├── hud.go
│   ├── title.go
│   └── gameover.go
└── assets/
    ├── sounds/
    └── fonts/
```

---

## Getting Started

```bash
# Clone the repo
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
cd YOUR_REPO

# Run directly
go run .

# Or build a binary
go build -o asteroids .
./asteroids
```

**Requirements:** Go 1.21+ — Ebitengine's dependencies are handled by `go mod`.

> On Linux you may need a few system libs for audio/graphics: `libgl1-mesa-dev`, `libasound2-dev`, `xorg-dev`. See the [Ebitengine install guide](https://ebitengine.org/en/documents/install.html) for your OS.

---

## Tech Stack

- **Language:** Go
- **Engine:** [Ebitengine](https://ebitengine.org) — a simple 2D game engine for Go
- **Audio:** Ebitengine's built-in audio (via `github.com/hajimehoshi/ebiten/v2/audio`)
- **Storage:** OS-native file system for high score persistence

No JavaScript. No browser. Just Go, a game loop, and some vectors.

---

## Stretch Goals

Once the core is solid, the backlog includes:

- 📱 Mobile / touch controls
- 👾 Two-player local co-op
- 🏆 Top-10 leaderboard with initials entry
- 🎨 Alternate visual themes (neon, pixel art, wireframe)

---

## Contributing

This project is being built issue-by-issue. If you want to contribute:

1. Pick an open issue
2. Branch off `main` — name it `issue/N-short-description`
3. Open a PR referencing the issue when done

---

## License

MIT — do whatever you want with it.

---

*"The rocks don't care how good you are. They just keep coming."*
