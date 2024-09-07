package projectile

import (
	"log"

	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/effect"
	"github.com/lgynico/mmo_skill/facade"
)

var FlyState facade.FSMState = nil

var FlyToEndTransition facade.FSMTransition = nil

func init() {
	FlyState = &flyState{facade.NewFSMState()}

	FlyToEndTransition = &flyToEndTransition{facade.NewFSMTransition()}
	FlyState.AddTransition(FlyToEndTransition)
}

type flyState struct {
	facade.FSMState
}

// func (s *flyState) OnUpdate(obj facade.GameObject, dt int64) int64 {
// 	proj := obj.(*Projectile)
// 	proj.stateCd -= int(dt)
// 	speed := proj.config.Speed
// 	pos := common.Move(proj.GetPosition(), proj.GetHeading(), float64(speed), dt)
// 	proj.SetPosition(pos)

// 	// 碰撞检测
// 	// bCreateTrap := (proj.config.ProcessType == common.FLYER_PROCESS_TYPE_CREATE_TRAP)
// 	bCollision := false
// 	shape := proj.GetBoundingShape(true)
// 	for _, unit := range proj.GetContext().GetUnits() {
// 		if unit.IsDead() || !facade.IsTarget(proj.GetUnit(), unit, proj.config.TargetType) {
// 			continue
// 		}
// 		if !geom.TwoShapeCollision(shape, unit.GetBoundingShape(true)) {
// 			continue
// 		}

// 		bCollision = true
// 		log.Printf("技能(%d) 命中单位(%d) \n", proj.config.Id, unit.GetId())
// 		config := proj.config
// 		for _, eff := range config.Effects {
// 			skillEffect, ok := effect.NewSkillEffect(eff.Type)
// 			if !ok {
// 				// logger.I("unknown effect type: ", eff.Type)
// 				continue
// 			}

// 			skillEffectDecorator := facade.NewSkillEffectDecorator(skillEffect, eff.Deco.Type, eff.Deco.Params)
// 			skillEffectDecorator.OnEffect(proj.GetSkill(), eff.SkillEffectRow)
// 		}
// 		// if bCreateTrap {
// 		// 	common.HandleTrap(proj.GetUnit(), config.TrapRow, config.TrapPos, 0)
// 		// 	break
// 		// } else {
// 		// 	s.addBuff(proj.GetUnit(), unit, config.Buffs)
// 		// }
// 	}

// 	ctx := proj.GetContext()
// 	if bCollision {
// 		ctx.RemoveObj(proj)
// 	} else {
// 		proj.stateCd -= facade.DeltaTime
// 	}
// 	return 0
// }

func (s *flyState) OnUpdate(obj facade.GameObject, dt int64) int64 {
	proj := obj.(*Projectile)
	proj.stateCd -= int(dt)
	speed := proj.config.Speed
	target := proj.GetObjTarget()
	pos := common.Move(proj.GetPosition(), target.GetPosition(), float64(speed), dt)
	proj.SetPosition(pos)

	ctx := proj.GetContext()
	unit := target.(facade.BattleUnit)
	if unit == nil || unit.IsDead() {
		ctx.RemoveObj(proj)
		return 0
	}

	if proj.GetBoundingShape(true).IsPointInside(unit.GetPosition()) {
		log.Printf("[%.2f] 技能 %d(%s) 命中单位 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, proj.config.Id, proj.config.Name, unit.GetId())

		config := proj.config
		for _, eff := range config.Effects {
			skillEffect, ok := effect.NewSkillEffect(eff.Type)
			if !ok {
				// logger.I("unknown effect type: ", eff.Type)
				continue
			}

			skillEffectDecorator := facade.NewSkillEffectDecorator(skillEffect, eff.Deco.Type, eff.Deco.Params)
			skillEffectDecorator.OnEffect(proj.GetSkill(), eff.SkillEffectRow, unit)
		}
		ctx.RemoveObj(proj)

	} else {
		proj.stateCd -= facade.DeltaTime

	}

	return 0
}

// func (s *flyState) addBuff(target facade.BattleUnit, buffs []int) {
// 	buffMgr := target.GetBuffMgr()
// 	for _, buffId := range buffs {
// 		buff := buff.NewBuff(buffId, src, 0) // todo
// 		buffMgr.AddBuff(buff)
// 	}
// }

// =====================================================
// -> 结束
// =====================================================
type flyToEndTransition struct {
	facade.FSMTransition
}

func (t *flyToEndTransition) IsValid(obj facade.GameObject) bool {
	proj := obj.(*Projectile)
	return proj.stateCd <= 0
}

func (t *flyToEndTransition) NextState() facade.FSMState {
	return EndState
}
