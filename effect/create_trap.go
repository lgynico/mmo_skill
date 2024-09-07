package effect

import (
	"log"
	"math"

	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/geom"
)

// 创建陷阱
type CreateTrap struct{}

func NewCreateTrap() facade.SkillEffect {
	return &CreateTrap{}
}

func (e *CreateTrap) OnEffect(skill facade.Skill, config *conf.SkillEffectRow, target facade.BattleUnit) {
	if trapRow, ok := conf.ConfMgr.SkillConfAdapter.TrapRows[config.Param1]; ok {
		e.createTrap(skill, trapRow, []int{config.Param2, config.Param3, config.Param4, config.Param5})
		log.Printf("技能效果-创建陷阱: %s(%d)\n", trapRow.Name, config.Param1)
		return
	}

	log.Printf("config not exists: name=%s, row=%d\n", conf.CONF_SKILL_EFFECT, config.Param1)
}

func (e *CreateTrap) createTrap(skill facade.Skill, config *conf.SkillTrapRowEx, params []int) {
	unit := skill.GetUnit()
	ctx := unit.GetContext()
	trap := ctx.CreateTrap(config, unit)
	trap.SetValue(skill.GetValue())
	trap.SetSkill(skill)

	// posChooser := GetPosChooser(dir)
	// posChooser.Choose(ctx, t, unit, posParams, unit.GetTarget())
	dir := params[0]
	dist := params[1]
	genType := params[2]
	genParam := params[3]
	heading, pos := e.chooseTrapPos(trap, dir, dist)
	trap.SetHeading(heading)
	if genType == 0 {
		trap.SetPosition(pos)
	} else {
		trap.SetPosition(unit.GetPosition().Copy())
		trap.SetSkillMove(true)
		// trap.SetSkillMovedDist(0)
		// dist := int(pos.DistTo(trap.GetPosition()))
		// trap.SetSkillMoveDist(dist)
		speed := 0
		if genType == 1 {
			speed = genParam
		} else if genType == 2 {
			speed = dist / genParam * 100
		}
		trap.SetSkillMoveSpeed(float64(speed))
	}

	// 延迟时间为0则立刻执行
	if !trap.IsSkillMove() && config.Delay == 0 {
		trap.Update(facade.DeltaTime)
		if config.Times > 1 {
			ctx.AddObj(trap)
		}
	} else {
		ctx.AddObj(trap)
	}

	s := trap.GetBoundingShape(true)
	log.Printf("陷阱 %s 形状：%v, 单位 %d 位置：%v\n", trap.GetConfig().Name, s, unit.GetId(), unit.GetPosition())
}

// func ensureTrapPos(trap facade.Trap, unit facade.BattleUnit, dir int, dist int) {
// 	var heading, pos *geom.Vector2d
// 	if dir == 99 {
// 		pos = trap.GetSkill().GetPosTarget()
// 		heading = pos.SubN(unit.GetPosition())
// 		heading.Normalize()
// 	} else {
// 		heading, pos = chooseTrapPos(unit, dir, dist)
// 	}

// 	trap.SetHeading(heading)
// 	trap.SetPosition(pos)

// 	angle := heading.Angle()
// 	mat := geom.MatrixRotate(angle)
// 	mat = mat.Mul(geom.MatrixTranslateV(pos))
// 	shape := trap.GetBoundingShape(true)
// 	shape.Transform(mat)
// }

func (e *CreateTrap) chooseTrapPos(trap facade.Trap, dir int, dist int) (*geom.Vector2d, *geom.Vector2d) {
	skill := trap.GetSkill()
	unit := trap.GetUnit()

	if dir == 99 {
		pos := skill.GetPosTarget()
		heading := pos.SubN(unit.GetPosition())
		heading.Normalize()
		return pos, heading
	}

	heading := unit.GetHeading().Copy()
	pos := unit.GetPosition().Copy()

	radian := float64(0)
	if dir >= 0 && dir <= 12 {
		if dir == 12 {
			dir = 0
		}
		radian = math.Pi / 6 * float64(dir)
	} else if dir <= 16 {
		radian = math.Pi / 4 * float64(dir-1)
	}
	heading.Rotate(radian)

	if dist > 0 {
		v := heading.NormalizeN()
		v.Mul(float64(dist))
		pos.Translate(v.GetX(), v.GetY())
	}

	return heading, pos
}
