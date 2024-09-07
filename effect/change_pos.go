package effect

import (
	"log"

	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
)

// 改变位置
type ChangePos struct {
}

func NewChangePos() facade.SkillEffect {
	return &ChangePos{}
}

func (e *ChangePos) OnEffect(skill facade.Skill, conf *conf.SkillEffectRow, target facade.BattleUnit) {
	log.Printf("改变位置：方向=%d, 速度=%d, 距离=%d, 是否闪现=%d\n", conf.Param1, conf.Param2, conf.Param3, conf.Param4)
	forward := conf.Param1 == 0
	speed := conf.Param2
	dist := conf.Param3
	flash := conf.Param4 == 1

	// target := skill.GetUnit()
	heading := target.GetHeading().Copy()
	if !forward {
		heading = heading.Inverse()
	}

	if flash {
		heading.Mul(float64(dist))
		target.GetPosition().Add(heading)
	} else {
		target.SetHeading(heading)
		target.SetSkillMove(true)
		target.SetSkillMoveSpeed(float64(speed))
		// target.SetSkillMoveDist(dist)
		// target.SetSkillMovedDist(0)
		stateModifier := facade.NewMoveStateModifier(int64(1000.0 * dist / speed))
		target.AddStateModifier(stateModifier)
	}
}
