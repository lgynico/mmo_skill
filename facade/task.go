package facade

import "github.com/lgynico/mmo_skill/conf"

type Task interface {
	Exec(ctx Battle) bool
}

/* 效果 */
type EffectTask struct {
	eff    SkillEffect
	effCfg *conf.SkillEffectRow
	skill  Skill
	target BattleUnit
}

func NewEffTask(eff SkillEffect, skill Skill, effCfg *conf.SkillEffectRow, target BattleUnit) Task {
	return &EffectTask{
		eff:    eff,
		skill:  skill,
		effCfg: effCfg,
		target: target,
	}
}

func (t *EffectTask) Exec(ctx Battle) bool {
	t.eff.OnEffect(t.skill, t.effCfg, t.target)
	return true
}

/* 延迟效果事件 */
type DelayEffectTask struct {
	EffectTask
	delay int64
}

func NewDelayEffectTask(eff SkillEffect, skill Skill, effCfg *conf.SkillEffectRow, target BattleUnit, delay int64) Task {
	task := NewEffTask(eff, skill, effCfg, target)
	effTask := task.(*EffectTask)
	return &DelayEffectTask{
		EffectTask: *effTask,
		delay:      delay,
	}
}

func (t *DelayEffectTask) Exec(ctx Battle) bool {
	t.delay -= ctx.GetDeltaTime()
	if t.delay > 0 {
		return false
	}

	t.eff.OnEffect(t.skill, t.effCfg, t.target)
	return true
}
