package projectile

import "github.com/lgynico/mmo_skill/facade"

var EndState facade.FSMState = nil

func init() {
	EndState = &endState{facade.NewFSMState()}
}

type endState struct {
	facade.FSMState
}
