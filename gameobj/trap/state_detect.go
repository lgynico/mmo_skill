package trap

/*
import (
	"fmt"
	"scenesrv/service/scene"
	"scenesrv/service/scene/common"
	"scenesrv/service/scene/geom"
)

var DetectState battle.FSMState = nil
var DetectToTriggerTransition battle.FSMTransition = nil
var DetectToEndTransition battle.FSMTransition = nil

func init() {
	DetectState = &detectState{battle.NewFSMState()}

	DetectToTriggerTransition = &detectToTriggerTransition{battle.NewFSMTransition()}
	DetectState.AddTransition(DetectToTriggerTransition)

	DetectToEndTransition = &detectToEndTransition{battle.NewFSMTransition()}
	DetectState.AddTransition(DetectToEndTransition)
}

type detectState struct {
	battle.FSMState
}

// =====================================================
// -> 触发
// =====================================================
type detectToTriggerTransition struct {
	battle.FSMTransition
}

func (t *detectToTriggerTransition) IsValid(obj battle.GameObject) bool {
	trap := obj.(*Trap)
	ctx := trap.GetContext()
	for _, obj := range ctx.GetObjs() {
		if unit, ok := obj.(battle.BattleUnit); ok {
			if unit.IsDead() || !common.IsTarget(trap.GetUnit(), unit, trap.GetConfig().TargetType) {
				continue
			}
			if geom.TwoShapeCollision(trap.GetDetectCircle(), unit.GetBoundingShape(true)) {
				log.Printf("有人撞到了陷阱。。。。。。\n")
				return true
			}
		}
	}
	return false
}

func (t *detectToTriggerTransition) NextState() battle.FSMState {
	return TriggerState
}

func (t *detectToTriggerTransition) OnTransition(obj battle.GameObject) {
	trap := obj.(*Trap)
	ctx := trap.GetContext()

	trap.SetDetect(false)
	config := trap.GetConfig()
	trap.triggerTime = int(ctx.GetTimeMillis()) + config.Delay

	// battle.PushTrapTrigger(ctx.GetMServer(), trap.CurrAreaId, trap, trap.triggerTime)
}

// =====================================================
// -> 结束
// =====================================================
type detectToEndTransition struct {
	battle.FSMTransition
}

func (t *detectToEndTransition) IsValid(obj battle.GameObject) bool {
	trap := obj.(*Trap)
	ctx := trap.GetContext()
	// 时间结束了都没有人进入范围
	return (trap.createTime + trap.GetConfig().Duration) <= int(ctx.GetTimeMillis())
}

func (t *detectToEndTransition) NextState() battle.FSMState {
	return EndState
}

*/
