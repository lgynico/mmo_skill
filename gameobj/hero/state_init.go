package hero

import "github.com/lgynico/mmo_skill/facade"

var InitState facade.FSMState = nil
var InitToMoveTransition facade.FSMTransition = nil

func init() {
	InitState = &initState{facade.NewFSMState()}

	InitToMoveTransition = &initToMoveTransition{facade.NewFSMTransition()}
	InitState.AddTransition(InitToMoveTransition)
}

type initState struct {
	facade.FSMState
}

// =====================================================
// -> 移动
// =====================================================
type initToMoveTransition struct {
	facade.FSMTransition
}

func (t *initToMoveTransition) IsValid(obj facade.GameObject) bool { return true }
func (t *initToMoveTransition) NextState() facade.FSMState         { return MoveState }
