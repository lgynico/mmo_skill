package skill

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
)

var SkillCreator = &skillCreator{}

type skillCreator struct{}

func NewSkillCreator() facade.SkillCreator {
	return &skillCreator{}
}

func (sc *skillCreator) CreateSkill(cfg *conf.SkillConfig) facade.Skill {
	return NewSkill(cfg)
}
