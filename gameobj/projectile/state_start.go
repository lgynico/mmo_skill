package projectile

import "github.com/lgynico/mmo_skill/facade"

var StartState facade.FSMState = nil

var StartToFlyTransition facade.FSMTransition = nil

func init() {
	StartState = &startState{facade.NewFSMState()}

	StartToFlyTransition = &startToFlyTransition{facade.NewFSMTransition()}
	StartState.AddTransition(StartToFlyTransition)
}

type startState struct {
	facade.FSMState
}

// =====================================================
// -> 飞行
// =====================================================
type startToFlyTransition struct {
	facade.FSMTransition
}

func (t *startToFlyTransition) IsValid(obj facade.GameObject) bool {
	proj := obj.(facade.Projectile)
	return proj.GetObjTarget() != nil || proj.GetPosTarget() != nil
}

func (t *startToFlyTransition) NextState() facade.FSMState {
	return FlyState
}
