package facade

import (
	"log"

	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/geom"
	"github.com/lgynico/mmo_skill/proto"
)

// ===============================================================================
// 技能处理器
// =========================“======================================================
type SkillHandler interface {
	EventHandler
	GetSkill() Skill

	IsPassive() bool
	IsActive() bool
	IsToggle() bool

	// 技能通用处理器
	OnSkillInit()  // 初始化
	OnSkillStart() // 开始施法
	OnSkillCast()  // 施法
	OnSkillEnd()   // 施法结束

	// 持续引导技能处理器
	OnChannelStart()         // 开始施法
	OnChannelInterval(int64) // 施法间隔
	OnChannelEnd()           // 施法结束

	// 开关类技能
	OnToggleOn()  // 技能开
	OnToggleOff() // 技能关

	OnHandleEffect(effCfg *conf.SkillEffectRowEx)
}

type BaseSkillHandler struct {
	skill Skill
}

func NewSkillHandler(skill Skill) SkillHandler {
	return &BaseSkillHandler{
		skill: skill,
	}
}

// ========================= 事件方法 ================================
func (h *BaseSkillHandler) OnEnterScene(event *GameEvent) bool    { return false }
func (h *BaseSkillHandler) OnLeaveScene(event *GameEvent) bool    { return false }
func (h *BaseSkillHandler) OnTakeDamage(event *GameEvent) bool    { return false }
func (h *BaseSkillHandler) OnGiveDamage(event *GameEvent) bool    { return false }
func (h *BaseSkillHandler) OnTakeHealing(event *GameEvent) bool   { return false }
func (h *BaseSkillHandler) OnGiveHealing(event *GameEvent) bool   { return false }
func (h *BaseSkillHandler) OnDead(event *GameEvent) bool          { return false }
func (h *BaseSkillHandler) OnReborn(event *GameEvent) bool        { return false }
func (h *BaseSkillHandler) OnAddBuff(event *GameEvent) bool       { return false }
func (h *BaseSkillHandler) OnImmediately(event *GameEvent) bool   { return false }
func (h *BaseSkillHandler) OnKillTarget(event *GameEvent) bool    { return false }
func (h *BaseSkillHandler) OnOverHeal(event *GameEvent) bool      { return false }
func (h *BaseSkillHandler) OnGiveHealFront(event *GameEvent) bool { return false }

// ========================= 其它方法 ================================
func (h *BaseSkillHandler) GetSkill() Skill {
	return h.skill
}

func (sh *BaseSkillHandler) IsPassive() bool {
	skill := sh.GetSkill()
	skillType := skill.GetConfig().Type
	return (skillType & SKILL_TYPE_PASSIVE) == SKILL_TYPE_PASSIVE
}

func (sh *BaseSkillHandler) IsActive() bool {
	skillType := sh.skill.GetConfig().Type
	return (skillType & SKILL_TYPE_ACTIVE) == SKILL_TYPE_ACTIVE
}

func (sh *BaseSkillHandler) IsToggle() bool {
	skillType := sh.skill.GetConfig().Type
	return (skillType & SKILL_TYPE_TOGGLE) == SKILL_TYPE_TOGGLE
}

// ========================= 触发方法 ================================
func (sh *BaseSkillHandler) OnSkillInit()  {}
func (sh *BaseSkillHandler) OnSkillStart() {}

func (sh *BaseSkillHandler) OnSkillCast() {
	// 角色回能
	skill := sh.GetSkill()
	if skill.GetConfig().IsAtk {
		EnergyRecover(skill.GetUnit(), skill, consts.AttackEnergy)
	}
}

func (sh *BaseSkillHandler) OnSkillEnd() {}

func (sh *BaseSkillHandler) OnChannelStart()            {}
func (sh *BaseSkillHandler) OnChannelInterval(dt int64) {}
func (sh *BaseSkillHandler) OnChannelEnd()              {}

func (sh *BaseSkillHandler) OnToggleOn()  {}
func (sh *BaseSkillHandler) OnToggleOff() {}

func (sh *BaseSkillHandler) OnHandleEffect(effCfg *conf.SkillEffectRowEx) {}

// ===================================================================================
// 技能
// ===================================================================================
type Skill interface {
	EventLinstener
	SkillHandler
	// SkillImpl
	Cast() bool
	Update(int64)

	GetLv() int
	GetConfig() *conf.SkillConfig
	GetPhase() int
	SetPhase(int)
	// GetTargetPos() *geom.Vector2d
	Cooldown(dt int)
	SetCd(time int)
	GetCd() int

	SetValue(int)
	GetValue() int

	// SetImpl(SkillImpl)
	// GetImpl() SkillImpl
	SetHandler(SkillHandler)
	GetUnit() BattleUnit
	SetUnit(BattleUnit)

	// 目标
	SetObjTarget(GameObject)
	GetObjTarget() GameObject

	SetPosTarget(*geom.Vector2d)
	GetPosTarget() *geom.Vector2d
}

type BaseSkill struct {
	SkillHandler
	unit   BattleUnit
	id     int
	lv     int
	config *conf.SkillConfig

	phase     int
	phaseTime int
	cd        int
	val       int // 用来传递伤害等数值

	singTime  int
	spellTime int
	shakeTime int

	objTarget GameObject
	posTarget *geom.Vector2d
}

func NewSkill(cfg *conf.SkillConfig) Skill {
	skill := &BaseSkill{
		id:     cfg.SkillEntryRow.Id,
		config: cfg,
		phase:  SKILL_PHASE_READY,
		cd:     0,
	}
	// skill.SetImpl(skill)
	skill.SetHandler(NewSkillHandler(skill))
	return skill
}

func (s *BaseSkill) Cast() bool {
	if s.cd > 0 || s.GetPhase() != SKILL_PHASE_READY {
		return false
	}

	if s.GetConfig().Unique && s.GetUnit().GetEp() < consts.UNIQUE_SKILL_ENERGY {
		return false
	}

	dist := s.GetConfig().Range
	if s.GetConfig().IsAtk {
		dist = s.GetUnit().GetAttrs().Get(consts.AttackDistance)
	}

	target := FindTarget(s.GetUnit(), s.GetConfig().TargetType, dist, s.GetUnit().GetTarget())
	if target == nil {
		return false
	}
	s.SetObjTarget(target)
	s.SetPosTarget(target.GetPosition())

	// log.Printf("释放技能：%s(%d)\n", s.config.Name, s.id)

	// s.phase = SKILL_PHASE_READY
	// s.phaseTime = 0
	if s.GetConfig().Unique {
		ep := s.GetUnit().GetEp() - consts.UNIQUE_SKILL_ENERGY
		s.GetUnit().SetEp(ep)
	}

	s.singTime = CalSkillSingTime(s)
	s.spellTime = CalSkillSpellTime(s)
	s.shakeTime = CalSkillShakeTime(s)

	return true
}

func (s *BaseSkill) Update(dt int64) {
	// 选择目标
	if s.phase == SKILL_PHASE_READY {
		// s.targetPos = s.selectPos()
		s.phase = SKILL_PHASE_SING
		s.phaseTime = 0
		s.OnSkillStart()
	}
	// 前摇
	if s.phase == SKILL_PHASE_SING {
		if dt = s.handlePhaseSing(dt); dt < 0 {
			return
		}
		s.phase = SKILL_PHASE_SPELL
		s.phaseTime = 0
		// CastPoint
		s.OnSkillCast()

		// if s.config.GuideType > 0 {
		// 	s.phase = SKILL_PHASE_GUIDE
		// 	s.OnChannelStart()
		// } else {
		// 	s.phase = SKILL_PHASE_SHAKE
		// }
	}
	// 施法
	if s.phase == SKILL_PHASE_SPELL {
		if dt = s.handlePhaseSpell(dt); dt < 0 {
			return
		}
		s.phase = SKILL_PHASE_SHAKE
		s.phaseTime = 0
	}
	// 引导
	// if s.phase == SKILL_PHASE_GUIDE {
	// 	if dt = s.handlePhaseGuide(dt); dt < 0 {
	// 		return
	// 	}
	// 	s.OnChannelEnd()
	// 	s.phase = SKILL_PHASE_SHAKE
	// 	s.phaseTime = 0
	// }
	// 后摇
	if s.phase == SKILL_PHASE_SHAKE {
		if dt = s.handlePhaseShake(dt); dt < 0 {
			return
		}
		s.skillEnd()
	}
}

func (s *BaseSkill) handlePhaseSing(dt int64) int64 {
	s.phaseTime += int(dt)
	return int64(s.phaseTime - s.singTime)
}

func (s *BaseSkill) handlePhaseSpell(dt int64) int64 {
	s.phaseTime += int(dt)
	return int64(s.phaseTime - s.spellTime)
}

// func (s *BaseSkill) handlePhaseGuide(dt int64) int64 {
// 	s.OnChannelInterval(dt)
// 	s.phaseTime += int(dt)
// 	return int64(s.phaseTime - s.singTime)
// }

func (s *BaseSkill) handlePhaseShake(dt int64) int64 {
	s.phaseTime += int(dt)
	return int64(s.phaseTime - s.shakeTime)
}

func (s *BaseSkill) skillEnd() {
	s.phase = SKILL_PHASE_END
	s.SetCd(CalSkillCd(s))
	s.OnSkillEnd()
}

// func (s *BaseSkill) selectPos() *geom.Vector2d {
// 	if s.targetPos != nil {
// 		return s.targetPos
// 	}
// 	unit := s.GetUnit()
// 	if target := unit.GetTarget(); target != nil {
// 		return target.GetPosition()
// 	}

//		// TODO 选择朝向前的一个位置
//		return unit.GetPosition()
//	}
func (s *BaseSkill) GetLv() int {
	return s.lv
}

func (s *BaseSkill) GetConfig() *conf.SkillConfig {
	return s.config
}

func (s *BaseSkill) GetPhase() int {
	return s.phase
}

func (s *BaseSkill) SetPhase(phase int) {
	s.phase = phase
}

func (s *BaseSkill) Cooldown(dt int) {
	s.cd -= dt
	if s.cd < 0 {
		s.cd = 0
	}

	if s.cd == 0 && s.GetPhase() == SKILL_PHASE_END {
		s.SetPhase(SKILL_PHASE_READY)
	}
}

func (s *BaseSkill) SetCd(time int) {
	s.cd = time
}

func (s *BaseSkill) GetCd() int {
	return s.cd
}

func (s *BaseSkill) SetValue(value int) {
	s.val = value
}

func (s *BaseSkill) GetValue() int {
	return s.val
}

func (s *BaseSkill) GetPasEffect(event int) []*conf.SkillEffectRowEx {
	effects := make([]*conf.SkillEffectRowEx, 0)
	for _, pasEffect := range s.GetConfig().PasEffects {
		if pasEffect.Event == event {
			effects = append(effects, pasEffect)
		}
	}

	return effects
}

func (s *BaseSkill) OnEventTrigger(eventType EventType, event *GameEvent) {
	effects := s.GetPasEffect(int(eventType))
	for _, effect := range effects {
		if !s.handleEvent(eventType, event) {
			s.OnHandleEffect(effect)
		}
	}
}

func (s *BaseSkill) handleEvent(evtType EventType, event *GameEvent) bool {
	isHandled := false
	switch evtType {
	case EVENT_IMMEDIATELY:
		isHandled = s.OnImmediately(event)
	case EVENT_ENTER_SCENE:
		isHandled = s.OnEnterScene(event)
	case EVENT_LEAVE_SCENE:
		isHandled = s.OnLeaveScene(event)
	case EVENT_GIVE_DAMAGE:
		isHandled = s.OnGiveDamage(event)
	case EVENT_TAKE_DAMAGE:
		isHandled = s.OnTakeDamage(event)
	case EVENT_GIVE_HEAL:
		isHandled = s.OnGiveHealing(event)
	case EVENT_TAKE_HEAL:
		isHandled = s.OnTakeHealing(event)
	case EVENT_DEAD:
		isHandled = s.OnDead(event)
	case EVENT_REBORN:
		isHandled = s.OnReborn(event)
	case EVENT_ADD_BUFF:
		isHandled = s.OnAddBuff(event)
	case EVENT_KILL_TARGET:
		isHandled = s.OnKillTarget(event)
	case EVENT_OVER_HEAL:
		isHandled = s.OnOverHeal(event)
	case EVENT_GIVE_HEAL_FRONT:
		isHandled = s.OnGiveHealFront(event)
	}

	return isHandled
}

func (s *BaseSkill) SetHandler(handler SkillHandler) {
	s.SkillHandler = handler
}

func (s *BaseSkill) GetUnit() BattleUnit {
	return s.unit
}

func (s *BaseSkill) SetUnit(unit BattleUnit) {
	s.unit = unit
}

func (s *BaseSkill) SetObjTarget(obj GameObject) {
	s.objTarget = obj
}

func (s *BaseSkill) GetObjTarget() GameObject {
	return s.objTarget
}

func (s *BaseSkill) SetPosTarget(pos *geom.Vector2d) {
	s.posTarget = pos
}

func (s *BaseSkill) GetPosTarget() *geom.Vector2d {
	return s.posTarget
}

// ======================================================================================
// 技能创建器
// ======================================================================================

type SkillCreator interface {
	CreateSkill(cfg *conf.SkillConfig) Skill
}

// ======================================================================================
// 技能管理器
// ======================================================================================
type SkillMgr interface {
	CastSkill(skillId int, targetPos *geom.Vector2d, targetObj GameObject) bool
	TryCastSkill() bool
	Update(dt int64)

	Countdown(dt int64)
	Reset()

	GetAttackSkill() Skill
	GetActSkills() []Skill
	GetCastingSkill() Skill
	SetCastingSkill(Skill)
	// GetLastSkill() (skid int, t int64)

	// AddSkill(cfg *conf.SkillConfig)
	AddSkills(cfgList []*conf.SkillConfig)
	RemoveSkill(int)

	GetUnit() BattleUnit
	SetUnit(BattleUnit)

	CastOver() bool
}

type BaseSkillMgr struct {
	SkillCreator
	attackSkill Skill // 普攻技能
	// mapId2ActSkill    map[int]Skill // 主动技能
	// mapId2PasSkill    map[int]Skill // 被动技能
	actSkills         []Skill // 主动技能
	pasSkills         []Skill // 被动技能
	castingSkill      Skill   // 正在释放的技能
	skillToCast       Skill   // 下一个要释放的技能
	lastSkill         int     // 处理连招用
	lastSkillCastTime int64   // 处理连招用
	unit              BattleUnit
}

func NewSkillMgr(unit BattleUnit, skillCreator SkillCreator, cfg *conf.SkillConfigs) SkillMgr {
	mgr := &BaseSkillMgr{
		attackSkill: nil,
		// mapId2ActSkill: make(map[int]Skill),
		// mapId2PasSkill: make(map[int]Skill),
		actSkills:    make([]Skill, 0),
		pasSkills:    make([]Skill, 0),
		castingSkill: nil,
		unit:         unit,
	}

	if skillCreator == nil {
		skillCreator = mgr
	}
	mgr.SkillCreator = skillCreator

	// mgr.atkSkill = skillCreator.CreateSkill(atkSkill)
	skillCfg := make([]*conf.SkillConfig, 0)
	skillCfg = append(skillCfg, cfg.Attack)
	for _, cfg := range cfg.Normals {
		// if cfg.Id != 132021 && cfg.Id != 132022 {
		// 	continue
		// }

		skillCfg = append(skillCfg, cfg)
	}
	//skillCfg = append(skillCfg, cfg.Normals...)
	mgr.AddSkills(skillCfg)
	return mgr
}

func (m *BaseSkillMgr) CastSkill(skillId int, targetPos *geom.Vector2d, targetObj GameObject) bool {
	var skillToBeCast Skill
	for _, skill := range m.actSkills {
		skillToBeCast = skill
	}

	if skillToBeCast == nil {
		return false
	}

	// log.Printf("准备释放技能：%s(%d)\n", skillToBeCast.GetConfig().Name, skillToBeCast.GetConfig().Id)
	if skillToBeCast.GetCd() > 0 {
		// log.Printf("技能 %s 正在cd中：%dms\n", skillToBeCast.GetConfig().Name, skillToBeCast.GetCd())
		return false
	}

	if m.castingSkill != nil {
		// log.Printf("已有正在释放的技能：%s(%d), 阶段%d\n", m.castingSkill.GetConfig().Name, m.castingSkill.GetConfig().Id, m.castingSkill.GetPhase())
		// if m.castingSkill.GetPhase() <= SKILL_PHASE_SING {
		// 	return false
		// }
		return m.castingSkill.GetPhase() > SKILL_PHASE_SING

		// 处理连招
		// if skillToBeCast.GetConfig().PreSkill == m.castingSkill.GetConfig().Id {
		// 	skillToBeCast.SetPosTarget(targetPos)
		// 	skillToBeCast.SetObjTarget(targetObj)
		// 	m.skillToCast = skillToBeCast
		// }
		// return true
	}

	// if skillToBeCast = m.ensurePreSkill(skillToBeCast); skillToBeCast == nil {
	// 	return false
	// }

	skillToBeCast.SetPosTarget(targetPos)
	skillToBeCast.SetObjTarget(targetObj)
	skillToBeCast.Cast()
	m.lastSkill = skillToBeCast.GetConfig().Id
	m.lastSkillCastTime = m.GetUnit().GetContext().GetTimeMillis()
	m.SetCastingSkill(skillToBeCast)

	return true
}

func (m *BaseSkillMgr) TryCastSkill() bool {
	if m.GetUnit().CanSpell() {
		for _, skill := range m.GetActSkills() {
			if !m.GetUnit().CanSpellUnique() && skill.GetConfig().Unique {
				continue
			}

			if m.doTryCastSkill(skill) {
				return true
			}
		}
	}

	return m.GetUnit().CanAttack() && m.doTryCastSkill(m.GetAttackSkill())
}

func (m *BaseSkillMgr) doTryCastSkill(skill Skill) bool {
	if skill.Cast() {
		m.SetCastingSkill(skill)

		unit := m.GetUnit()
		ctx := unit.GetContext()
		action := &proto.BattleAction{
			Time:      ctx.GetTimeMillis(),
			Unit:      unit.GetId(),
			Key:       int32(consts.ACTION_SPELL),
			Value:     int64(skill.GetConfig().Id),
			Target:    skill.GetObjTarget().GetId(),
			UnitPos:   ConvertVector2d(unit.GetPosition()),
			TargetPos: ConvertVector2d(skill.GetPosTarget()),
		}
		ctx.RecordAction(action)

		log.Printf("[%.2f] 单位 %d 释放技能 %d(%s), 目标: %d %v\n", float64(ctx.GetTimeMillis())/1000, unit.GetId(), skill.GetConfig().Id, skill.GetConfig().Name,
			skill.GetObjTarget().GetId(), skill.GetPosTarget())
		return true
	}
	return false
}

// func (m *BaseSkillMgr) ensurePreSkill(curSkill Skill) Skill {
// 	preSkill := curSkill.GetConfig().PreSkill
// 	if preSkill == 0 {
// 		return curSkill
// 	}

// 	ctx := m.GetUnit().GetContext()
// 	if preSkill == m.lastSkill && curSkill.GetConfig().ComboTime >= int(ctx.GetCurTime()-m.lastSkillCastTime) {
// 		return curSkill
// 	}

// 	lastRow := conf.SkillConfAdapter.EntryRows[preSkill]
// 	var sid int
// 	for lastRow != nil {
// 		sid = lastRow.Id
// 		lastRow = conf.SkillConfAdapter.EntryRows[lastRow.PreSkill]
// 	}

// 	if skill, ok := m.mapId2ActSkill[sid]; ok {
// 		return skill
// 	}

// 	return nil
// }

func (m *BaseSkillMgr) Update(dt int64) {
	skill := m.GetCastingSkill()
	if skill == nil {
		if m.skillToCast == nil {
			return
		}

		skill = m.skillToCast
		m.skillToCast = nil

		skill.Cast()
		m.SetCastingSkill(skill)

		m.lastSkill = skill.GetConfig().Id
		m.lastSkillCastTime = m.GetUnit().GetContext().GetTimeMillis()
	}

	skill.Update(dt)
	if skill.GetPhase() == SKILL_PHASE_END {
		m.SetCastingSkill(nil)
	}
}

func (m *BaseSkillMgr) Reset() {
	m.SetCastingSkill(nil)
	for _, skill := range m.actSkills {
		skill.SetPhase(SKILL_PHASE_READY)
	}
	// for _, skill := range m.mapId2ActSkill {
	// 	skill.SetPhase(SKILL_PHASE_READY)
	// }
}

func (m *BaseSkillMgr) CreateSkill(cfg *conf.SkillConfig) Skill {
	return NewSkill(cfg)
}

func (m *BaseSkillMgr) GetCastingSkill() Skill {
	return m.castingSkill
}

func (m *BaseSkillMgr) SetCastingSkill(skill Skill) {
	m.castingSkill = skill
}

func (m *BaseSkillMgr) GetLastSkill() (skid int, t int64) {
	return m.lastSkill, m.lastSkillCastTime
}

func (m *BaseSkillMgr) Countdown(dt int64) {
	// for _, s := range m.mapId2ActSkill {
	for _, s := range m.actSkills {
		if s.GetPhase() != SKILL_PHASE_READY {
			s.Cooldown(int(dt))
		}
	}

	if m.attackSkill.GetPhase() != SKILL_PHASE_READY {
		m.attackSkill.Cooldown(int(dt))
	}
}

func (m *BaseSkillMgr) addSkill(cfg *conf.SkillConfig) {
	skill := m.doAddSkill(cfg)
	if skill == nil || !skill.IsPassive() {
		return
	}

	// if skill.GetConfig().FakeBuff != 0 {
	// 	unit.GetContext().PushBuff(unit.GetId(), int32(skill.GetConfig().FakeBuff), 0)
	// }
}

func (m *BaseSkillMgr) AddSkills(cfgList []*conf.SkillConfig) {
	for _, cfg := range cfgList {
		m.addSkill(cfg)
	}

	m.SetUnit(m.GetUnit())
	// for _, passSkill := range m.mapId2PasSkill {
	for _, passSkill := range m.pasSkills {
		m.addEventListener(passSkill)
	}
}

func (m *BaseSkillMgr) doAddSkill(cfg *conf.SkillConfig) Skill {
	skill := m.SkillCreator.CreateSkill(cfg)
	if skill == nil {
		return nil
	}

	// skill.SetUnit(m.GetUnit())
	if skill.GetConfig().IsAtk {
		m.attackSkill = skill
		return skill
	}

	id := cfg.SkillEntryRow.Id
	if skill.IsPassive() {
		// if _, ok := m.mapId2PasSkill[id]; !ok {
		// 	m.mapId2PasSkill[id] = skill
		// }
		exists := false
		for _, s := range m.pasSkills {
			if s.GetConfig().SkillEntryRow.Id == id {
				exists = true
				break
			}
		}
		if !exists {
			m.pasSkills = append(m.pasSkills, skill)
		}
	}
	if skill.IsActive() {
		skill.SetCd(cfg.EnterCoolDown)
		skill.SetPhase(SKILL_PHASE_END)
		// if _, ok := m.mapId2ActSkill[id]; !ok {
		// 	m.mapId2ActSkill[id] = skill
		// }
		exists := false
		for _, s := range m.actSkills {
			if s.GetConfig().SkillEntryRow.Id == id {
				exists = true
				break
			}
		}
		if !exists {
			m.actSkills = append(m.actSkills, skill)
		}
	}

	return skill
}

func (m *BaseSkillMgr) addEventListener(skill Skill) {
	unit := m.GetUnit()
	// for _, eventType := range skill.GetConfig().Events {
	// 	// if eventType == int(EVENT_IMMEDIATELY) {
	// 	// 	skill.OnImmediately(NewGameEvent(unit))
	// 	// 	continue
	// 	// }
	// 	eventName, ok := GetEventName(EventType(eventType))
	// 	if !ok {
	// 		continue
	// 	}

	// 	if eventSystem := unit.GetEventSystem(); eventSystem != nil {
	// 		eventSystem.AddListener(eventName, skill)
	// 	}
	// }

	for _, pasEffect := range skill.GetConfig().PasEffects {
		if eventSystem := unit.GetEventSystem(); eventSystem != nil {
			eventSystem.AddListener(EventType(pasEffect.Event), skill)
		}
	}
}

func (m *BaseSkillMgr) RemoveSkill(id int) {
	skill := m.doRemoveSkill(id)
	if skill == nil || !skill.IsPassive() {
		return
	}

	unit := m.GetUnit()
	eventSystem := unit.GetEventSystem()
	if eventSystem == nil {
		return
	}

	for _, eventType := range skill.GetConfig().Events {
		eventSystem.RemoveListener(EventType(eventType), skill)
	}
}

func (m *BaseSkillMgr) doRemoveSkill(id int) Skill {
	var skill Skill = nil
	// if _, ok := m.mapId2ActSkill[id]; ok {
	// 	skill = m.mapId2ActSkill[id]
	// 	delete(m.mapId2ActSkill, id)
	// }
	// if _, ok := m.mapId2PasSkill[id]; ok {
	// 	skill = m.mapId2PasSkill[id]
	// 	delete(m.mapId2PasSkill, id)
	// }
	for i, s := range m.actSkills {
		if s.GetConfig().SkillEntryRow.Id == id {
			m.actSkills = append(m.actSkills[:i], m.actSkills[i+1:]...)
			break
		}
	}
	for i, s := range m.pasSkills {
		if s.GetConfig().SkillEntryRow.Id == id {
			m.pasSkills = append(m.pasSkills[:i], m.pasSkills[i+1:]...)
			break
		}
	}
	return skill
}

func (m *BaseSkillMgr) GetUnit() BattleUnit {
	return m.unit
}

func (m *BaseSkillMgr) SetUnit(unit BattleUnit) {
	m.unit = unit

	m.attackSkill.SetUnit(unit)

	// for _, s := range m.mapId2ActSkill {
	for _, s := range m.actSkills {
		s.SetUnit(unit)
	}

	// for _, s := range m.mapId2PasSkill {
	for _, s := range m.pasSkills {
		s.SetUnit(unit)
		s.SetObjTarget(unit)
	}

}

func (m *BaseSkillMgr) CastOver() bool {
	if m.castingSkill != nil && m.castingSkill.GetPhase() != SKILL_PHASE_END {
		return false
	}

	if m.skillToCast != nil {
		return false
	}

	return true
}

func (m *BaseSkillMgr) GetAttackSkill() Skill {
	return m.attackSkill
}

func (m *BaseSkillMgr) GetActSkills() []Skill {
	// skills := make([]Skill, 0, len(m.mapId2ActSkill))
	// for _, s := range m.mapId2ActSkill {
	// 	skills = append(skills, s)
	// }

	return m.actSkills
}

// ======================================================================================
// 技能衍生物创建器
// ======================================================================================

type SkillDerivantCreator interface {
	CreateTrap(config *conf.SkillTrapRowEx, unit BattleUnit) Trap
	CreateProjectile(skill Skill, config *conf.SkillProjectileRowEx) Projectile
	CreateBuff(id int, src Skill, val int) Buff
	CreateHalo(id int, skill Skill) Halo
}

type BaseSkillDerivantCreator struct {
}

func NewSkillDerivantCreator() SkillDerivantCreator {
	return &BaseSkillDerivantCreator{}
}

func (c *BaseSkillDerivantCreator) CreateTrap(config *conf.SkillTrapRowEx, unit BattleUnit) Trap {
	return nil
}

func (c *BaseSkillDerivantCreator) CreateProjectile(skill Skill, config *conf.SkillProjectileRowEx) Projectile {
	return nil
}

func (c *BaseSkillDerivantCreator) CreateBuff(id int, src Skill, val int) Buff {
	return nil
}

func (c *BaseSkillDerivantCreator) CreateHalo(id int, skill Skill) Halo {
	return nil
}
