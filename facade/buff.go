package facade

import (
	"log"

	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/proto"
)

type BuffHandler interface {
	OnHandleBuff(Buff, BattleUnit)
	OnRemoveBuff(Buff, BattleUnit)
	OnCalValue(Buff, BattleUnit) int
}

type BaseBuffHandler struct {
}

func NewBuffHandler() BuffHandler {
	return &BaseBuffHandler{}
}

func (b *BaseBuffHandler) OnHandleBuff(buff Buff, unit BattleUnit) {}
func (b *BaseBuffHandler) OnRemoveBuff(buff Buff, unit BattleUnit) {
	for _, subBuff := range buff.GetSubBuffs() {
		subBuff.OnRemoveBuff(subBuff, unit)
	}
}

func (b *BaseBuffHandler) OnCalValue(buff Buff, unit BattleUnit) int {
	return buff.GetValue()
}

// ==================================================================
// ==================================================================

type Buff interface {
	BuffHandler
	HandleBuff(unit BattleUnit, dt int64, isSub bool)

	ResetEffTimes()
	IncrLv(maxLv int)
	IsRemove() bool
	CalValue(BattleUnit) int

	SetStartTime(time int)
	GetStartTime() int

	GetId() int
	GetSrc() BattleUnit

	GetValue() int
	SetValue(int)
	GetConfig() *conf.SkillBuffRow
	GetLv() int
	SetBuffHandler(BuffHandler)

	GetSkill() Skill

	SetSubBuffs([]Buff)
	GetSubBuffs() []Buff
}

type BaseBuff struct {
	BuffHandler
	id     int // 配置
	config *conf.SkillBuffRow
	// srcUnit BattleUnit // 来源
	// uid       int              // 唯一id
	skill             Skill
	startTime         int  // 开始时间
	lastTime          int  // 上次生效时间
	effTimes          int  // 生效次数
	lv                int  // 层次
	remove            bool // 是否移除
	perpetualEffected bool // 永久buff生效
	val               int  // 数值

	subBuffs []Buff // 子buff
}

func NewBuff(id int, skill Skill, val int) Buff {
	config, ok := conf.ConfMgr.SkillBuff.GetRowByInt(id)
	if !ok {
		return nil
	}

	buff := &BaseBuff{
		id:          id,
		config:      config,
		skill:       skill,
		startTime:   0,
		lastTime:    0,
		effTimes:    0,
		lv:          1,
		val:         val,
		BuffHandler: NewBuffHandler(),
	}

	return buff
}

func (b *BaseBuff) HandleBuff(unit BattleUnit, dt int64, isSub bool) {
	if unit.IsDead() {
		return
	}
	// 永久buff
	if b.isPerpetualBuff() {
		if !b.perpetualEffected {
			b.OnHandleBuff(b, unit)
			b.perpetualEffected = true
		}
		return
	}

	ctx := unit.GetContext()
	if b.effTimes < b.config.EffTimes {
		if b.lastTime+b.config.EffInterval <= int(ctx.GetTimeMillis()) {
			log.Printf("[%.2f] 单位 %d 处理buff (%s), 次数: %d\n", float64(ctx.GetTimeMillis())/1000, unit.GetId(), b.config.Name, b.effTimes+1)
			b.OnHandleBuff(b, unit)
			b.effTimes++
			b.lastTime = int(ctx.GetTimeMillis())
		}
	}

	if b.config.Duration >= 0 {
		if !isSub && b.startTime+b.config.Duration <= int(ctx.GetTimeMillis()) {
			log.Printf("[%.2f] 单位 %d 移除buff (%s)\n", float64(ctx.GetTimeMillis())/1000, unit.GetId(), b.config.Name)
			b.OnRemoveBuff(b, unit)
			b.remove = true
		}
	}
}

// func (b *BaseBuff) HandleRemove(unit BattleUnit, ctx Scene) {
// 	b.OnRemoveBuff(b, unit, ctx)
// }

func (b *BaseBuff) GetConfig() *conf.SkillBuffRow {
	return b.config
}

func (b *BaseBuff) isPerpetualBuff() bool {
	return b.config.Duration == -1 && b.config.EffTimes == -1 && b.config.EffInterval == -1
}

func (b *BaseBuff) SetBuffHandler(handler BuffHandler) {
	b.BuffHandler = handler
}

func (b *BaseBuff) SetStartTime(time int) {
	b.startTime = time
}

func (b *BaseBuff) GetStartTime() int {
	return b.startTime
}

func (b *BaseBuff) IncrLv(maxLv int) {
	if b.lv < maxLv {
		b.lv++
		// b.perpetualEffected = false
	}
}

func (b *BaseBuff) ResetEffTimes() {
	b.effTimes = 0
	b.lastTime = 0
}

func (b *BaseBuff) GetId() int {
	return b.id
}

func (b *BaseBuff) GetSrc() BattleUnit {
	return b.skill.GetUnit()
}

func (b *BaseBuff) IsRemove() bool {
	return b.remove
}

func (b *BaseBuff) CalValue(unit BattleUnit) int {
	return b.BuffHandler.OnCalValue(b, unit)
}

// 不能实现
// func (b *BaseBuff) OnCalValue(buff Buff, unit BattleUnit) int {
// 	return buff.GetValue()
// }

// func (b *BaseBuff) OnHandleBuff(unit BattleUnit, ctx Scene) {}

// func (b *BaseBuff) OnRemoveBuff(unit BattleUnit, ctx Scene) {}

func (b *BaseBuff) SetValue(val int) {
	b.val = val
}

func (b *BaseBuff) GetValue() int {
	return b.val
}

func (b *BaseBuff) GetLv() int {
	return b.lv
}

func (b *BaseBuff) GetSkill() Skill {
	return b.skill
}

func (b *BaseBuff) SetSubBuffs(buffs []Buff) {
	b.subBuffs = buffs
}

func (b *BaseBuff) GetSubBuffs() []Buff {
	return b.subBuffs
}

// ====================================================================
// ====================================================================

type BuffMgr interface {
	UpdateBuffs(unit BattleUnit, dt int64)
	GetBuffs() []Buff
	AddBuff(Buff)
	RemoveBuff(id int, src BattleUnit)
	ForceRemoveBuff(Buff)
	Clear()
	GetUnit() BattleUnit
}

type BaseBuffMgr struct {
	BuffMgr
	buffs []Buff
	unit  BattleUnit
}

func NewBuffMgr(unit BattleUnit) BuffMgr {
	return &BaseBuffMgr{
		buffs: make([]Buff, 0),
		unit:  unit,
	}
}

func (m *BaseBuffMgr) AddBuff(buff Buff) {
	if m.immuneBuff(buff) { // 免疫
		return
	}

	ctx := m.GetUnit().GetContext()
	buff.SetStartTime(int(ctx.GetTimeMillis()))
	config := buff.GetConfig()
	if config.Type == int(consts.BUFF_DAMAGE) ||
		config.Type == int(consts.BUFF_HEALING) ||
		config.Type == int(consts.BUFF_ENERGY) {
		m.buffs = append(m.buffs, buff)
		log.Printf("[%.2f] 单位 %d 增加buff %s 来源Id %d\n", float64(ctx.GetTimeMillis())/1000, m.GetUnit().GetId(), buff.GetConfig().Name, buff.GetSrc().GetId())
		SendEvent(EVENT_ADD_BUFF, m.GetUnit(), buff)
		//ctx.PushBuff(m.GetUnit().GetId(), int32(buff.GetId()), int32(buff.GetValue()))

		return
	}

	// TODO buff如何叠加
	for _, b := range m.buffs {
		if b.GetId() == buff.GetId() && b.GetSrc().GetId() == buff.GetSrc().GetId() {
			maxLayer := buff.GetConfig().Layer
			if maxLayer > 1 {
				buff.OnRemoveBuff(buff, m.unit)
				b.IncrLv(maxLayer)
				buff.OnHandleBuff(buff, m.unit)
			}

			buff.SetStartTime(int(ctx.GetTimeMillis()))
			log.Printf("[%.2f] 单位 %d 刷新buff %s\n", float64(ctx.GetTimeMillis())/1000, m.GetUnit().GetId(), buff.GetConfig().Name)
			return
		}
	}
	// 驱散buff
	// for k := range buff.GetConfig().DispelBuffSet {
	// 	m.RemoveBuff(k, buff.GetSrc())
	// }
	m.buffs = append(m.buffs, buff)
	log.Printf("[%.2f] 单位 %d 增加buff %s\n", float64(ctx.GetTimeMillis())/1000, m.GetUnit().GetId(), buff.GetConfig().Name)
	SendEvent(EVENT_ADD_BUFF, m.GetUnit(), ctx, buff)
	//ctx.PushBuff(m.GetUnit().GetId(), int32(buff.GetId()), int32(buff.GetValue()))

	action := &proto.BattleAction{
		Time:  ctx.GetTimeMillis(),
		Unit:  m.GetUnit().GetId(),
		Key:   int32(consts.ACTION_ADD_BUFF),
		Value: int64(buff.GetId()),
	}
	ctx.RecordAction(action)

}

func (m *BaseBuffMgr) immuneBuff(buff Buff) bool {
	for _, b := range m.buffs {
		if len(buff.GetConfig().Tag) > 0 && buff.GetConfig().Tag == b.GetConfig().ImmuneTag {
			return true
		}
	}
	return false
}

func (m *BaseBuffMgr) RemoveBuff(id int, src BattleUnit) {
	for i, buff := range m.buffs {
		if buff.GetId() == id && buff.GetSkill().GetUnit().GetId() == src.GetId() {
			m.buffs = append(m.buffs[:i], m.buffs[i+1:]...)
			buff.OnRemoveBuff(buff, src)
			return
		}
	}
}

func (m *BaseBuffMgr) ForceRemoveBuff(buff Buff) {
	m.buffs = BuffFilter(m.buffs, func(b Buff) bool {
		return b != buff
	})
}

func (m *BaseBuffMgr) Clear() {
	m.buffs = make([]Buff, 0)
}

func (m *BaseBuffMgr) GetBuffs() []Buff {
	return m.buffs
}

func (m *BaseBuffMgr) UpdateBuffs(unit BattleUnit, dt int64) {
	m.buffs = BuffFilter(m.buffs, func(b Buff) bool {
		return !b.IsRemove()
	})

	// TODO 隐藏一bug：size 大于 len(m.buffs)
	// TODO 可能导致某些buff不被执行
	// size := len(m.buffs)
	// for i := 0; i < size; i++ {
	for i := 0; i < len(m.buffs); i++ {
		buff := m.buffs[i]
		buff.HandleBuff(unit, dt, false)

		if buff.GetSubBuffs() != nil {
			for _, subBuff := range buff.GetSubBuffs() {
				subBuff.HandleBuff(unit, dt, true)
			}
		}
	}

}

func (m *BaseBuffMgr) GetUnit() BattleUnit {
	return m.unit
}
