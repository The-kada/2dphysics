package gameerrors

import "errors"

var (
	GameExitedError = errors.New("Game exited") // probably an antipattern ro handle exits as errors, TODO: see a more ebiten standard way of handling errors
)
