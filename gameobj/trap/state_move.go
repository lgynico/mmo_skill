package trap

import (
	"log"

	"github.com/lgynico/mmo_skill/facade"
)

var MoveState facade.FSMState = nil
var MoveToTriggerTransition facade.FSMTransition = nil

func init() {
	MoveState = &moveState{facade.NewFSMState()}

	MoveToTriggerTransition = &moveToTriggerTransition{facade.NewFSMTransition()}
	MoveState.AddTransition(MoveToTriggerTransition)

}

type moveState struct {
	facade.FSMState
}

func (s *moveState) OnEnter(obj facade.GameObject, dt int64) {
	log.Printf("[%.2f] 陷阱移动: %s\n", float64(obj.GetContext().GetTimeMillis())/1000, obj.(*Trap).GetConfig().Name)
	s.OnUpdate(obj, dt)
}

func (s *moveState) OnUpdate(obj facade.GameObject, dt int64) int64 {
	trap := obj.(*Trap)
	// 技能位移
	// common.UpdateSkillMove(trap, dt)
	trap.UpdateSkillMove(trap, dt)

	return 0
}

// =====================================================
// -> 生效
// =====================================================
type moveToTriggerTransition struct {
	facade.FSMTransition
}

func (t *moveToTriggerTransition) IsValid(obj facade.GameObject) bool {
	trap := obj.(*Trap)
	return !trap.IsSkillMove()
}

func (t *moveToTriggerTransition) NextState() facade.FSMState { return TriggerState }
