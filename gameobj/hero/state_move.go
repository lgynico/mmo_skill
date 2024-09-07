package hero

import (
	"log"

	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/proto"
)

var MoveState facade.FSMState = nil
var MoveToSpellTransition facade.FSMTransition = nil
var MoveToDeadTransition facade.FSMTransition = nil

func init() {
	MoveState = &moveState{facade.NewFSMState()}

	MoveToSpellTransition = &moveToSpellTransition{facade.NewFSMTransition()}
	MoveState.AddTransition(MoveToSpellTransition)

	MoveToDeadTransition = &moveToDeadTransition{facade.NewFSMTransition()}
	MoveState.AddTransition(MoveToDeadTransition)
}

type moveState struct {
	facade.FSMState
}

func (s *moveState) OnEnter(obj facade.GameObject, dt int64) {
	id := obj.GetId()
	log.Printf("[%.2f] 单位 %d 进入 move 状态：pos = %v, heading = %v\n", float64(obj.GetContext().GetTimeMillis())/1000, id, obj.GetPosition(), obj.GetHeading())
}

func (s *moveState) OnUpdate(obj facade.GameObject, dt int64) int64 {
	flag := false
	hero := obj.(*Hero)
	target := hero.GetTarget()
	ctx := hero.GetContext()
	if target == nil || target.IsDead() {
		target = common.FindNearestEnemy(hero, ctx.GetUnits())
		if target == nil {
			return 0
		}

		hero.SetTarget(target)
		flag = true
	}

	if !hero.CanMove() {
		return 0
	}

	if flag {
		action := &proto.BattleAction{
			Time:      ctx.GetTimeMillis(),
			Unit:      hero.GetId(),
			Key:       int32(consts.ACTION_MOVE),
			UnitPos:   facade.ConvertVector2d(hero.GetPosition()),
			TargetPos: facade.ConvertVector2d(target.GetPosition()),
		}
		ctx.RecordAction(action)
	}

	AttackDistance := hero.GetAttrs().Get(consts.AttackDistance)
	adi := float64(AttackDistance) / 100

	moveSpeed := hero.GetAttrs().Get(consts.MoveSpeed)
	mov := float64(moveSpeed) / 100

	pos, remainTime := common.MoveNoCross(hero.GetPosition(), target.GetPosition(), mov, dt, adi)
	hero.SetPosition(pos)

	// facade.Logger.D(fmt.Sprintf("单位 %d 移动到 %v", hero.GetId(), pos))

	return remainTime
}

// =====================================================
// -> 放技能
// =====================================================
type moveToSpellTransition struct {
	facade.FSMTransition
}

func (t *moveToSpellTransition) IsValid(obj facade.GameObject) bool {
	hero := obj.(*Hero)

	skillMgr := hero.GetSkillMgr()
	if !skillMgr.CastOver() {
		return true
	}

	return skillMgr.TryCastSkill()
}

func (t *moveToSpellTransition) NextState() facade.FSMState { return SpellState }

// =====================================================
// -> 死亡
// =====================================================

type moveToDeadTransition struct {
	facade.FSMTransition
}

func (t *moveToDeadTransition) IsValid(obj facade.GameObject) bool {
	hero := obj.(*Hero)
	return hero.IsDead()
}

func (t *moveToDeadTransition) NextState() facade.FSMState { return DeadState }
