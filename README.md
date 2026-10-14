# 2dphyslib

A 2D rigid-body physics engine written in Go, with an Ebiten-based renderer for real-time simulation.

## What Is This?

`2dphyslib` is a from-scratch 2D physics engine that simulates rigid-body dynamics — gravity, collisions, rotation, friction, and restitution — using constraint-based impulse resolution. It ships with an example app (`main.go`) that lets you drop shapes into a walled arena and watch them interact in real time.

## Features

- **Rigid body simulation** — linear and angular velocity/acceleration integration per tick
- **Three shape types** — circles, axis-aligned/rotated rectangles, and arbitrary convex polygons
- **Collision detection**
  - Circle–circle
  - Circle–polygon (including rectangles)
  - Polygon–polygon via the Separating Axis Theorem (SAT)
- **Constraint-based collision resolution** — penetration constraints solved with a Gauss-Seidel iterative solver (5 iterations per frame)
- **Joint constraints** — anchor two bodies together at a point
- **Friction and restitution** — per-body coefficients influence post-collision impulses
- **Static bodies** — zero-mass bodies act as immovable walls/floors
- **Debug overlay** — press `D` to highlight collision contact points

## Controls (example app)

| Key | Action |
|-----|--------|
| `R` | Spawn a rectangle at the center |
| `C` | Spawn a circle at the center |
| `P` | Spawn a convex polygon (triangle) at the center |
| `D` | Toggle debug / collision visualization |
| `Esc` | Quit |

## Getting Started

```bash
git clone https://github.com/the-kada/2dphys
cd 2dphys
go run .
```

Requires Go 1.24+ and the [Ebiten](https://ebitengine.org/) game library (pulled automatically via `go mod`).

## Architecture

```
2dphys/
├── physics/
│   ├── body.go          # Rigid body: position, velocity, mass, inertia, impulse application
│   ├── shape.go         # Shape types: Circle, Rectangle, ConvexPolygon + moment-of-inertia formulas
│   ├── collision.go     # Narrow-phase collision detection, ContactInfo generation
│   ├── constraints.go   # JointConstraint & PenetrationConstraint; Gauss-Seidel solver
│   ├── forces.go        # Force accumulators (gravity applied in WorldPhysics.UpdatePhysics)
│   ├── worldphysics.go  # Simulation loop: broad/narrow phase, constraint solve, integration
│   ├── vec2.go          # 2D vector math
│   ├── matmn.go         # M×N matrix (used for Jacobian arithmetic)
│   └── vecn.go          # N-dimensional vector (used for lambda/impulse solving)
├── ebitendrawer/
│   └── draw.go          # Ebiten helpers: draw circles, polygons, contact highlights
└── main.go              # Example game loop wiring physics + renderer together
```

### Simulation loop (per frame)

1. Detect all pairwise collisions → build `ContactInfo` list
2. Create a `PenetrationConstraint` for each contact
3. Apply gravity force to every dynamic body
4. Integrate acceleration → update velocities
5. Pre-solve constraints (warm-start with cached lambda)
6. Solve constraints iteratively (5 passes of Gauss-Seidel)
7. Integrate velocities → update positions

### Using the physics package directly

```go
world := &physics.WorldPhysics{}

floor := physics.NewBody(
    physics.Vec2{400, 550},  // position
    physics.Vec2{},          // velocity
    physics.Vec2{},          // acceleration
    physics.Vec2{},          // initial forces
    0,                       // mass=0 → static
    0, 0, 0, 0,
    2, 1.3,                  // friction, restitution
    physics.NewRectangle(800, 40),
)
world.AddBody(floor)

ball := physics.NewBody(
    physics.Vec2{400, 100},
    physics.Vec2{}, physics.Vec2{}, physics.Vec2{},
    1.0,            // mass
    0, 0, 0, 0,
    2, 1.3,
    &physics.Circle{Radius: 30},
)
world.AddBody(ball)

// Each game tick:
world.UpdatePhysics(ebiten.ActualTPS())
```
