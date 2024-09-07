package hero

import (
	"log"

	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/utils"
)

var SpellState facade.FSMState = nil
var SpellToMoveTransition facade.FSMTransition = nil
var SpellToDeadTransition facade.FSMTransition = nil

func init() {
	SpellState = &spellState{facade.NewFSMState()}

	SpellToMoveTransition = &spellToMoveTransition{facade.NewFSMTransition()}
	SpellState.AddTransition(SpellToMoveTransition)

	SpellToDeadTransition = &spellToDeadTransition{facade.NewFSMTransition()}
	SpellState.AddTransition(SpellToDeadTransition)
}

type spellState struct {
	facade.FSMState
}

func (s *spellState) OnEnter(obj facade.GameObject, dt int64) {
	log.Printf("[%.2f] 单位 %d 进入 spell 状态\n", float64(obj.GetContext().GetTimeMillis())/1000, obj.GetId())

	// 不响应移动
	// TODO 如果是普攻？
	hero := obj.(*Hero)
	hero.stateModifier = facade.NewMoveStateModifier(-1)
	hero.AddStateModifier(hero.stateModifier)

	// if hero.spellRequest != nil {
	// 	req := hero.spellRequest
	// 	targetObj, _ := hero.GetContext().GetObj(req.TargetId)
	// 	hero.GetSkillMgr().CastSkill(int(req.SkillId), req.Hold, geom.NewVector2d(req.CastPos.X, req.CastPos.Z), targetObj)
	// 	hero.spellRequest = nil
	// }

}

func (s *spellState) OnUpdate(obj facade.GameObject, dt int64) int64 {
	hero := obj.(*Hero)
	skillMgr := hero.GetSkillMgr()

	if skillMgr.CastOver() && !skillMgr.TryCastSkill() {
		return 0
	}

	skillMgr.Update(dt)
	// 技能位移
	// common.UpdateSkillMove(hero, dt)
	return 0
}

func (s *spellState) OnExit(obj facade.GameObject, dt int64) {
	hero := obj.(*Hero)
	if hero.stateModifier != nil {
		hero.RemoveStateModifier(hero.stateModifier)
		hero.stateModifier = nil
	}
}

// =====================================================
// -> 移动
// =====================================================
type spellToMoveTransition struct {
	facade.FSMTransition
}

func (t *spellToMoveTransition) IsValid(obj facade.GameObject) bool {
	hero := obj.(*Hero)
	if hero.IsDead() {
		return false
	}

	if !hero.GetSkillMgr().CastOver() {
		return false
	}

	if hero.GetTarget() == nil || hero.GetTarget().IsDead() {
		return true
	}

	adi := float64(hero.GetAttrs().Get(consts.AttackDistance)) / 100
	distSq := hero.GetPosition().DistSqTo(hero.GetTarget().GetPosition())
	// return distSq > adi*adi
	return utils.IsFloat64Gt(distSq, adi*adi)
}

func (t *spellToMoveTransition) NextState() facade.FSMState {
	return MoveState
}

// =====================================================
// -> 死亡
// =====================================================

type spellToDeadTransition struct {
	facade.FSMTransition
}

func (t *spellToDeadTransition) IsValid(obj facade.GameObject) bool {
	hero := obj.(*Hero)
	return hero.IsDead()
}

func (t *spellToDeadTransition) NextState() facade.FSMState {
	return DeadState
}
