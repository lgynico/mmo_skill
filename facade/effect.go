package facade

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/utils"
)

/******************************************************
* 技能效果
******************************************************/
type SkillEffect interface {
	OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit)
}

/******************************************************
* 技能装饰器
******************************************************/
type SkillEffectDecorator interface {
	SkillEffect
}

// Base
type BaseSkillEffectDecorator struct {
	SkillEffect
}

func NewBaseSkillEffectDecorator(effect SkillEffect) SkillEffectDecorator {
	return &BaseSkillEffectDecorator{
		SkillEffect: effect,
	}
}

func (d *BaseSkillEffectDecorator) OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit) {
	d.SkillEffect.OnEffect(skill, conf, target)
}

// 重复执行N次
type RepeatSED struct {
	cnt int
	SkillEffect
}

func NewRepeatSkillEffectDecorator(effect SkillEffect, count int) SkillEffectDecorator {
	return &RepeatSED{
		SkillEffect: effect,
		cnt:         count,
	}
}

func (d *RepeatSED) OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit) {
	for i := 0; i < d.cnt; i++ {
		d.SkillEffect.OnEffect(skill, conf, target)
	}
}

// 条件为真一直执行
type WhileSED struct {
	cond func(skill Skill) bool
	SkillEffect
}

func NewWhileSkillEffectDecorator(effect SkillEffect, condFunc func(skill Skill) bool) SkillEffectDecorator {
	return &WhileSED{
		SkillEffect: effect,
		cond:        condFunc,
	}
}

func (d *WhileSED) OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit) {
	for d.cond(skill) {
		d.SkillEffect.OnEffect(skill, conf, target)
	}
}

// 延迟执行
type WaitSED struct {
	delay int
	SkillEffect
}

func NewWaitSkillEffectDecorator(effect SkillEffect, delay int) SkillEffectDecorator {
	return &WaitSED{
		SkillEffect: effect,
		delay:       delay,
	}
}

func (d *WaitSED) OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit) {
	delay := d.delay
	task := NewDelayEffectTask(d.SkillEffect, skill, conf, target, int64(delay))
	skill.GetUnit().GetContext().AddTask(task)
	// go func() {
	// 	// time.Sleep(time.Duration(d.delay) * time.Millisecond)
	// 	t := d.delay * 1000
	// 	for {
	// 		if t > 0 {
	// 			t -= DeltaTime
	// 			continue
	// 		}
	// 		task := NewEffTask(d.SkillEffect, skill, conf, target)
	// 		skill.GetUnit().GetContext().AddTask(task)
	// 	}
	// }()
}

// 前提
type IfSED struct {
	cond func(skill Skill) bool
	SkillEffect
}

func NewIfSkillEffectDecorator(effect SkillEffect, condFunc func(skill Skill) bool) SkillEffectDecorator {
	return &IfSED{
		SkillEffect: effect,
		cond:        condFunc,
	}
}

func (d *IfSED) OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit) {
	if d.cond(skill) {
		d.SkillEffect.OnEffect(skill, conf, target)
	}
}

// 概率
type ProbSED struct {
	prob float32
	SkillEffect
}

func NewProbSkillEffectDecorator(effect SkillEffect, prob float32) SkillEffectDecorator {
	return &ProbSED{
		SkillEffect: effect,
		prob:        prob,
	}
}

func (d *ProbSED) OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit) {
	i := utils.RandIntByCrypto(1, 10000)
	prob := float32(i) / 10000
	if prob <= d.prob {
		d.SkillEffect.OnEffect(skill, conf, target)
	}
}

// 重复并等待
type RepeatWithDelaySED struct {
	count int
	delay int
	SkillEffect
}

func NewRepeatWithDelaySkillEffectDecorator(effect SkillEffect, count int, delay int) SkillEffectDecorator {
	return &RepeatWithDelaySED{
		SkillEffect: effect,
		count:       count,
		delay:       delay,
	}
}

func (d *RepeatWithDelaySED) OnEffect(skill Skill, conf *conf.SkillEffectRow, target BattleUnit) {
	for i := 0; i < d.count; i++ {
		delay := d.delay * 1000 * (i + 1)
		task := NewDelayEffectTask(d.SkillEffect, skill, conf, target, int64(delay))
		skill.GetUnit().GetContext().AddTask(task)
	}
	// go func() {
	// 	t := d.delay * 1000
	// 	ctx := skill.GetUnit().GetContext()
	// 	for i := 0; i < d.count; {
	// 		if t > 0 {
	// 			t -= DeltaTime
	// 			continue
	// 		}

	// 		task := NewEffTask(d.SkillEffect, skill, conf, target)
	// 		ctx.AddTask(task)
	// 		t = d.delay * 1000
	// 		i++
	// 	}
	// }()
}

func NewSkillEffectDecorator(effect SkillEffect, typ int, params []interface{}) SkillEffectDecorator {
	switch typ {
	case conf.EFF_DECO_REPEAT: // repeat count
		cnt := params[0].(int)
		return NewRepeatSkillEffectDecorator(effect, cnt)
	case conf.EFF_DECO_WHILE: // while conditon
		cond := params[0].(*conf.SkillCondition)
		return NewWhileSkillEffectDecorator(effect, newCondFunc(cond))
	case conf.EFF_DECO_WAIT: // wait delay
		delay := params[0].(int)
		return NewWaitSkillEffectDecorator(effect, delay)
	case conf.EFF_DECO_IF: // if condition
		cond := params[0].(*conf.SkillCondition)
		return NewIfSkillEffectDecorator(effect, newCondFunc(cond))
	case conf.EFF_DECO_PROB: // prob probability
		prob := params[0].(float32)
		return NewProbSkillEffectDecorator(effect, prob)
	case conf.EFF_DECO_REPEAT_WITH_DELAY: // repeat count withdelay delay
		cnt := params[0].(int)
		delay := params[1].(int)
		return NewRepeatWithDelaySkillEffectDecorator(effect, cnt, delay)
	default:
		return NewBaseSkillEffectDecorator(effect)
	}
}

func newCondFunc(cond *conf.SkillCondition) func(Skill) bool {
	c := cond
	return func(s Skill) bool {
		return testCondition(s, c)
	}
}

func testCondition(s Skill, condition *conf.SkillCondition) bool {
	var (
		unit  = s.GetUnit() // TODO oppoTarget
		value int
	)

	baseValue := float64(condition.Value)
	switch condition.AttrName {
	case "moving":
		value = int(unit.GetMoveVector().Length())
	case "hp":
		value = unit.GetHp()
		if condition.IsPerc {
			baseValue = float64(unit.GetAttrs().Get(consts.MaxHp)*condition.Value) / 100.0
		}
	default:
		key := getAttrKey(condition.AttrName)
		if key == -1 {
			return false
		}
		value = unit.GetAttrs().Get(key)
		if condition.IsPerc {
			baseValue = float64(value*condition.Value) / 100.0
		}
	}
	return compare(value, int(baseValue), condition.Comparator)
}

func compare(value, baseValue int, comparator string) bool {
	switch comparator {
	case "=":
		return value == baseValue
	case ">":
		return value > baseValue
	case "<":
		return value < baseValue
	case ">=":
		return value >= baseValue
	case "<=":
		return value <= baseValue
	case "!=":
		return value != baseValue
	}
	return false
}

func getAttrKey(attrName string) int {
	switch attrName {
	case "maxhp":
		return consts.MaxHp
	case "atk":
		return consts.Attack
	case "mov":
		return consts.MoveSpeed
	case "asp":
		return consts.AttackSpeed
	case "adi":
		return consts.AttackDistance
	}
	return -1
}
