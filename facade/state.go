package facade

import (
	"sort"
)

// ================================================================
// 状态
// ================================================================
type FSMState interface {
	OnEnter(obj GameObject, dt int64)
	OnUpdate(obj GameObject, dt int64) int64
	OnExit(obj GameObject, dt int64)
	GetTransitions() []FSMTransition
	AddTransition(FSMTransition)
}

type BaseFSMState struct {
	trans []FSMTransition
}

func NewFSMState() FSMState {
	return &BaseFSMState{
		trans: make([]FSMTransition, 0),
	}
}

func (s *BaseFSMState) OnEnter(obj GameObject, dt int64)        {}
func (s *BaseFSMState) OnUpdate(obj GameObject, dt int64) int64 { return 0 }
func (s *BaseFSMState) OnExit(obj GameObject, dt int64)         {}

func (s *BaseFSMState) GetTransitions() []FSMTransition {
	return s.trans
}

func (s *BaseFSMState) AddTransition(tran FSMTransition) {
	s.trans = append(s.trans, tran)
}

// ================================================================
// 状态转换
// ================================================================
type FSMTransition interface {
	IsValid(obj GameObject) bool
	NextState() FSMState
	OnTransition(obj GameObject)
}

type BaseFSMTransition struct {
}

func NewFSMTransition() FSMTransition {
	return &BaseFSMTransition{}
}

func (t *BaseFSMTransition) IsValid(obj GameObject) bool { return false }
func (t *BaseFSMTransition) NextState() FSMState         { return nil }
func (t *BaseFSMTransition) OnTransition(obj GameObject) {}

// ================================================================
//
// ================================================================
type FSM interface {
	Update(dt int64)
}

type BaseFSM struct {
	currState FSMState
	obj       GameObject
}

func NewFSM(obj GameObject, initState FSMState) FSM {
	return &BaseFSM{
		currState: initState,
		obj:       obj,
	}
}

func (f *BaseFSM) Update(dt int64) {
	t := dt
	for t > 0 {
		nextState, ok := f.switchStateValid()
		if ok {
			f.currState.OnExit(f.obj, t)
			f.currState = nextState
			f.currState.OnEnter(f.obj, t)
		}
		t = f.currState.OnUpdate(f.obj, t)
	}
}

func (f *BaseFSM) switchStateValid() (FSMState, bool) {
	for _, tran := range f.currState.GetTransitions() {
		if tran.IsValid(f.obj) {
			tran.OnTransition(f.obj)
			return tran.NextState(), true
		}
	}
	return nil, false
}

// ================================================================
//  fsm set/get 接口
// ================================================================

/*****************************************
* 状态修改器
*****************************************/
type StateModifier interface {
	OnModify(StateMgr)
	OnRemove(StateMgr)
	Update(dt int64)
	GetTime() int64
}

type BaseStateModifier struct {
	time int64
}

func (m *BaseStateModifier) OnModify(mgr StateMgr) {

}
func (m *BaseStateModifier) OnRemove(mgr StateMgr) {}
func (m *BaseStateModifier) Update(dt int64)       { m.time -= dt }
func (m *BaseStateModifier) GetTime() int64        { return m.time }

/* 移动 */
type MoveStateModifier struct {
	StateModifier
}

func NewMoveStateModifier(time int64) StateModifier {
	return &MoveStateModifier{
		StateModifier: &BaseStateModifier{
			time: time,
		},
	}
}

func (m *MoveStateModifier) OnModify(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.moveState++
	// fmt.Println("增加移动状态：", base.moveState)
}

func (m *MoveStateModifier) OnRemove(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.moveState--
	// fmt.Println("减少移动状态：", base.moveState)
}

/* 施法 */
type SpellStateModifier struct {
	StateModifier
}

func NewSpellStateModifier(time int64) StateModifier {
	return &SpellStateModifier{
		StateModifier: &BaseStateModifier{
			time: time,
		},
	}
}

func (m *SpellStateModifier) OnModify(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.spellState++
}

func (m *SpellStateModifier) OnRemove(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.spellState--
}

/* 攻击 */
type AttackeStateModifier struct {
	StateModifier
}

func NewAttackeStateModifier(time int64) StateModifier {
	return &AttackeStateModifier{
		StateModifier: &BaseStateModifier{
			time: time,
		},
	}
}

func (m *AttackeStateModifier) OnModify(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.attackState++
}

func (m *AttackeStateModifier) OnRemove(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.attackState--
}

/* 魅惑 */
type CharmStateModifier struct {
	StateModifier
}

func NewCharmStateModifier(time int64) StateModifier {
	return &CharmStateModifier{
		StateModifier: &BaseStateModifier{
			time: time,
		},
	}
}

func (m *CharmStateModifier) OnModify(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.charmStage++
}

func (m *CharmStateModifier) OnRemove(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.charmStage--
}

/* 释放大招 */
type UniqueSpellStateModifier struct {
	StateModifier
}

func NewUniqueSpellStateModifier(time int64) StateModifier {
	return &UniqueSpellStateModifier{
		StateModifier: &BaseStateModifier{
			time: time,
		},
	}
}

func (m *UniqueSpellStateModifier) OnModify(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.banUniqueState++
}

func (m *UniqueSpellStateModifier) OnRemove(mgr StateMgr) {
	base := mgr.(*BaseStateMgr)
	base.banUniqueState--
}

/*****************************************
* 状态管理器
*****************************************/
type StateMgr interface {
	CanMove() bool
	CanSpell() bool
	CanAttack() bool
	CanSelect() bool
	IsCharmState() bool
	CanSpellUnique() bool
	AddStateModifier(StateModifier)
	RemoveStateModifier(StateModifier)
	Update(dt int64)
}

type BaseStateMgr struct {
	modifiers        []StateModifier
	persistModifiers []StateModifier
	moveState        int
	spellState       int
	attackState      int
	selectState      int
	charmStage       int
	banUniqueState   int
}

func NewStateMgr() StateMgr {
	return &BaseStateMgr{
		modifiers:        make([]StateModifier, 0),
		persistModifiers: make([]StateModifier, 0),
		moveState:        0,
		spellState:       0,
		attackState:      0,
		selectState:      0,
	}
}

func (m *BaseStateMgr) CanMove() bool {
	// fmt.Println("判断是否可以移动：", m.moveState)
	return m.moveState == 0
}

func (m *BaseStateMgr) CanSpell() bool {
	return m.spellState == 0
}

func (m *BaseStateMgr) CanAttack() bool {
	return m.attackState == 0
}

func (m *BaseStateMgr) CanSelect() bool {
	return m.selectState == 0
}

func (m *BaseStateMgr) IsCharmState() bool {
	return m.charmStage != 0
}

// 允许释放大招
func (m *BaseStateMgr) CanSpellUnique() bool {
	return m.banUniqueState == 0
}

func (m *BaseStateMgr) AddStateModifier(modifier StateModifier) {
	modifier.OnModify(m)
	if modifier.GetTime() == -1 {
		m.persistModifiers = append(m.persistModifiers, modifier)
	} else {
		m.modifiers = append(m.modifiers, modifier)
		m.sortModifiers()
	}
}

func (m *BaseStateMgr) sortModifiers() {
	sort.Slice(m.modifiers, func(i, j int) bool {
		m1 := m.modifiers[i]
		m2 := m.modifiers[j]
		return m1.GetTime() > m2.GetTime() // 从大到小排列
	})
}

func (m *BaseStateMgr) RemoveStateModifier(modifier StateModifier) {
	if modifier.GetTime() == -1 {
		for i := len(m.persistModifiers) - 1; i >= 0; i-- {
			mod := m.persistModifiers[i]
			if mod == modifier {
				mod.OnRemove(m)
				m.persistModifiers = append(m.persistModifiers[:i], m.persistModifiers[i+1:]...)
				return
			}
		}

	}

	for i := len(m.modifiers) - 1; i >= 0; i-- {
		mod := m.modifiers[i]
		if mod == modifier {
			mod.OnRemove(m)
			m.modifiers = append(m.modifiers[:i], m.modifiers[i+1:]...)
			return
		}
	}
}

func (m *BaseStateMgr) Update(dt int64) {
	for i := len(m.modifiers) - 1; i >= 0; i-- {
		mod := m.modifiers[i]
		mod.OnRemove(m)
		m.modifiers = append(m.modifiers[:i], m.modifiers[i+1:]...)
	}
}
