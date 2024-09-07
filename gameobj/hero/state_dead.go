package hero

import "github.com/lgynico/mmo_skill/facade"

var DeadState facade.FSMState = nil
var DeadToMoveTransition facade.FSMTransition

func init() {
	DeadState = &deadState{facade.NewFSMState()}

	DeadToMoveTransition = &deadToMoveTransition{facade.NewFSMTransition()}
	DeadState.AddTransition(DeadToMoveTransition)
}

type deadState struct {
	facade.FSMState
}

func (s *deadState) OnEnter(obj facade.GameObject, dt int64) {
	// fmt.Println("hero is dead:", obj.GetId())
}

// =====================================================
// -> 移动
// =====================================================
type deadToMoveTransition struct {
	facade.FSMTransition
}
