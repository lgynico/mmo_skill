package projectile

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/geom"
)

type Projectile struct {
	facade.Projectile
	config  *conf.SkillProjectileRowEx
	stateCd int
}

func NewProjectile(id int32, roleId int64, skill facade.Skill, config *conf.SkillProjectileRowEx) facade.Projectile {
	p := &Projectile{
		Projectile: facade.NewProjectile(id, roleId),
		stateCd:    config.Time,
		config:     config,
	}
	p.SetSkill(skill)

	unit := skill.GetUnit()
	p.SetUnit(unit)
	p.SetFSM(facade.NewFSM(p, FlyState))
	p.SetBoundingShape(geom.NewShapeInt(geom.ShapeType(config.Shape[0]), config.Shape[1:]...))

	// heading := targetPos.SubN(unit.GetPosition())
	// heading.Normalize()
	heading := unit.GetHeading().Copy()
	heading.Normalize()
	p.SetHeading(heading)

	// radius := unit.GetBoundingShape(false).Radius() // 人物是一个点了
	radius := float64(unit.GetAttrs().Get(consts.AttackDistance)) / 100
	offset := heading.MulN(radius)
	pos := unit.GetPosition().TranslateN(offset.GetX(), offset.GetY())
	p.SetPosition(pos)

	p.SetContext(unit.GetContext())

	p.SetImpl(p)
	return p
}

func (p *Projectile) GetConfigId() int32 {
	return int32(p.config.Id)
}
