package skill

import (
	"log"

	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/effect"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/skill/impl"
)

type Skill struct {
	facade.Skill
}

func NewSkill(cfg *conf.SkillConfig) facade.Skill {
	skill := &Skill{
		Skill: facade.NewSkill(cfg),
	}
	// skill.SetImpl(skill)
	skill.setSkillEventHandler()

	return skill
}

func (s *Skill) setSkillEventHandler() {
	handler := NewSkillHandler(s)
	switch s.GetConfig().Id {
	case facade.SKILL_BURNING_RAGE:
		handler = impl.NewBurningRage(handler)
	case facade.SKILL_OVER_HEAL:
		handler = impl.NewOverHealing(handler)
	case facade.SKILL_HEAL_FRONT:
		handler = impl.NewHealBefore(handler)
	case facade.SKILL_NO_DIE:
		handler = impl.NewUnDead(handler)
	default:
	}
	s.SetHandler(handler)
}

type SkillHandler struct {
	facade.SkillHandler
}

func NewSkillHandler(s facade.Skill) facade.SkillHandler {
	return &SkillHandler{
		SkillHandler: facade.NewSkillHandler(s),
	}
}

func (s *SkillHandler) OnSkillCast() {
	s.SkillHandler.OnSkillCast()
	cfg := s.GetSkill().GetConfig()
	for _, eff := range cfg.Effects {
		s.doHandleEffect(eff)
	}
}

func (s *SkillHandler) OnHandleEffect(effCfg *conf.SkillEffectRowEx) {
	s.doHandleEffect(effCfg)
}

func (s *SkillHandler) doHandleEffect(effCfg *conf.SkillEffectRowEx) {
	skillEffect, ok := effect.NewSkillEffect(effCfg.Type)
	if !ok {
		log.Println("unknown effect type: ", effCfg.Type)
		return
	}

	skillEffectDecorator := facade.NewSkillEffectDecorator(skillEffect, effCfg.Deco.Type, effCfg.Deco.Params)
	target := s.GetSkill().GetObjTarget().(facade.BattleUnit)
	skillEffectDecorator.OnEffect(s.GetSkill(), effCfg.SkillEffectRow, target)
}

// func (s *SkillHandler) OnImmediately(event *facade.GameEvent) {
// 	cfg := s.GetSkill().GetConfig()
// 	for _, eff := range cfg.Effects {
// 		skillEffect, ok := effect.NewSkillEffect(eff.Type)
// 		if !ok {
// 			logger.I("unknown effect type: ", eff.Type)
// 			continue
// 		}

// 		skillEffectDecorator := facade.NewSkillEffectDecorator(skillEffect, eff.Deco.Type, eff.Deco.Params)

// 		target := s.GetSkill().GetObjTarget().(facade.BattleUnit)
// 		skillEffectDecorator.OnEffect(s.GetSkill(), eff.SkillEffectRow, target)
// 	}
// }
