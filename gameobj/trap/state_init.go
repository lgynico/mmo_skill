package trap

import "github.com/lgynico/mmo_skill/facade"

var InitState facade.FSMState = nil

// var InitToDetectTransition facade.FSMTransition = nil
var InitToTriggerTransition facade.FSMTransition = nil
var InitToMoveTransition facade.FSMTransition = nil

func init() {
	InitState = &initState{facade.NewFSMState()}

	// InitToDetectTransition = &initToDetectTransition{}
	// InitState.AddTransition(InitToDetectTransition)

	InitToTriggerTransition = &initToTriggerTransition{facade.NewFSMTransition()}
	InitState.AddTransition(InitToTriggerTransition)

	InitToMoveTransition = &initToMoveTransition{facade.NewFSMTransition()}
	InitState.AddTransition(InitToMoveTransition)
}

type initState struct {
	facade.FSMState
}

func (s *initState) OnExit(obj facade.GameObject, dt int64) {
	ctx := obj.GetContext()
	trap := obj.(*Trap)
	trap.createTime = int(ctx.GetTimeMillis())
	config := trap.GetConfig()
	// if config.Condition == common.TRAP_CONDITION_DETECT {
	// 	trap.triggerTime = -1
	// } else {
	trap.triggerTime = int(ctx.GetTimeMillis()) + config.Delay
	// }
}

// =====================================================
// -> 警戒
// =====================================================
// type initToDetectTransition struct {
// 	facade.FSMTransition
// }

// func (t *initToDetectTransition) IsValid(obj facade.GameObject) bool {
// 	trap := obj.(*Trap)
// 	return trap.GetConfig().Condition == 1
// }

// func (t *initToDetectTransition) NextState() facade.FSMState {
// 	return DetectState
// }

// =====================================================
// -> 触发
// =====================================================
type initToTriggerTransition struct {
	facade.FSMTransition
}

func (t *initToTriggerTransition) IsValid(obj facade.GameObject) bool {
	trap := obj.(*Trap)
	return !trap.IsSkillMove()
}

func (t *initToTriggerTransition) NextState() facade.FSMState {
	return TriggerState
}

// =====================================================
// -> 移动
// =====================================================
type initToMoveTransition struct {
	facade.FSMTransition
}

func (t *initToMoveTransition) IsValid(obj facade.GameObject) bool {
	trap := obj.(*Trap)
	return trap.IsSkillMove()
}

func (t *initToMoveTransition) NextState() facade.FSMState {
	return MoveState
}
