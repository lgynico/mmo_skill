package facade

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/geom"
)

const DEFAULT_AOI_RANGE = 10

// 地图物体基础接口
type GameObject interface {
	Transform
	// EventSystem
	ContextHolder

	GetId() int32
	SetId(int32)
	GetRoleId() int64
	GetCfgId() int32

	GetObjType() consts.GameObjectType
	GetTeamId() int32
	SetTeamId(teamId int32)

	GetEventSystem() *EventSystem
	SetEventSystem(*EventSystem)

	Update(dt int64)
	SetImpl(GameObject)

	GetFSM() FSM
	SetFSM(FSM)
}

type BaseGameObject struct {
	Transform
	ContextHolder
	impl         GameObject
	eventSystem  *EventSystem
	id           int32                 // 唯一id
	cfgId        int32                 // 配置id
	roleId       int64                 // 属于的玩家id
	teamId       int32                 // 阵营id
	objectType   consts.GameObjectType // 类型
	stateMachine FSM                   // 状态机
}

func NewGameObject(id int32, roleId int64, objType consts.GameObjectType, cfgId int) GameObject {
	return &BaseGameObject{
		Transform:     NewTransform(),
		ContextHolder: NewContextHolder(),
		id:            id,
		roleId:        roleId,
		objectType:    objType,
		cfgId:         int32(cfgId),
	}
}

func (b *BaseGameObject) GetId() int32 {
	return b.id
}

func (b *BaseGameObject) SetId(id int32) {
	b.id = id
}

func (b *BaseGameObject) GetRoleId() int64 {
	return b.roleId
}

func (b *BaseGameObject) GetCfgId() int32 {
	return b.cfgId
}

func (b *BaseGameObject) GetTeamId() int32 {
	return b.teamId
}

func (m *BaseGameObject) SetTeamId(teamId int32) {
	m.teamId = teamId
}

func (b *BaseGameObject) GetObjType() consts.GameObjectType {
	return b.objectType
}

func (b *BaseGameObject) OnEnterAOI(obj GameObject) {

}

func (b *BaseGameObject) OnLeaveAOI(obj GameObject) {

}

func (b *BaseGameObject) OnSync() {

}

func (m *BaseGameObject) SetEventSystem(es *EventSystem) {
	m.eventSystem = es
}

func (m *BaseGameObject) GetEventSystem() *EventSystem {
	return m.eventSystem
}

func (m *BaseGameObject) SetImpl(obj GameObject) {
	m.impl = obj
}

func (m *BaseGameObject) Update(dt int64) {
	if fsm := m.GetFSM(); fsm != nil {
		fsm.Update(dt)
	}
}

func (m *BaseGameObject) GetFSM() FSM {
	return m.stateMachine
}

func (m *BaseGameObject) SetFSM(fsm FSM) {
	m.stateMachine = fsm
}

// func (m *BaseGameObject) GetInitPos() *geom.Vector2d {
// 	return m.initPos
// }

// func (m *BaseGameObject) SetInitPos(initPos *geom.Vector2d) {
// 	m.initPos = initPos
// }

/*************************************
* 战斗接口
*************************************/
type BattleUnit interface {
	GameObject
	BattleUnitProps
	StateMgr
	SkillMover

	GetBuffMgr() BuffMgr
	SetBuffMgr(BuffMgr)

	GetSkillMgr() SkillMgr
	SetSkillMgr(SkillMgr)

	GetHaloMgr() HaloMgr
	SetHaloMgr(HaloMgr)

	GetAttrs() Attrs
	SetAttrs(Attrs)

	GetMoveVector() *geom.Vector2d
	SetMoveVector(*geom.Vector2d)

	GetMaxHp() int
	GetLvl() int
}

type BaseBattleUnit struct {
	GameObject
	BattleUnitProps
	StateMgr
	SkillMover
	lvl       int
	attrs     Attrs
	skillMgr  SkillMgr
	buffMgr   BuffMgr
	haloMgr   HaloMgr
	movVector *geom.Vector2d
}

func NewBattleUnit(id int32, roleId int64, objType consts.GameObjectType, cfgId int, lvl int) BattleUnit {
	b := &BaseBattleUnit{
		GameObject:      NewGameObject(id, roleId, objType, cfgId),
		BattleUnitProps: NewBattleUnitProps(),
		StateMgr:        NewStateMgr(),
		movVector:       geom.NewVector2d(0, 0),
		SkillMover:      NewSkillMover(),
		lvl:             lvl,
	}

	b.SetBuffMgr(NewBuffMgr(b))
	b.SetEventSystem(NewEventSystem())
	b.SetHaloMgr(NewHaloMgr())
	return b
}

func (b *BaseBattleUnit) GetAttrs() Attrs {
	return b.attrs
}

func (b *BaseBattleUnit) SetAttrs(attrs Attrs) {
	b.attrs = attrs
}

func (b *BaseBattleUnit) GetSkillMgr() SkillMgr {
	return b.skillMgr
}

func (b *BaseBattleUnit) SetSkillMgr(skillMgr SkillMgr) {
	b.skillMgr = skillMgr
	// skillMgr.AddPassiveSkillsListener()
}

func (b *BaseBattleUnit) GetBuffMgr() BuffMgr {
	return b.buffMgr
}

func (b *BaseBattleUnit) SetBuffMgr(buffMgr BuffMgr) {
	b.buffMgr = buffMgr
}

func (b *BaseBattleUnit) GetHaloMgr() HaloMgr    { return b.haloMgr }
func (b *BaseBattleUnit) SetHaloMgr(mgr HaloMgr) { b.haloMgr = mgr }

func (b *BaseBattleUnit) GetMoveVector() *geom.Vector2d {
	return b.movVector
}

func (b *BaseBattleUnit) SetMoveVector(vec *geom.Vector2d) {
	b.movVector = vec
}

func (b *BaseBattleUnit) TakeHealing(val int) {
	b.BattleUnitProps.TakeHealing(val)

	maxHp := b.GetAttrs().Get(consts.MaxHp)
	if b.GetHp() > maxHp {
		b.SetHp(maxHp)
	}
}

func (b *BaseBattleUnit) Update(dt int64) {
	b.GameObject.Update(dt)
	if skillMgr := b.GetSkillMgr(); skillMgr != nil {
		// skillMgr.Update(dt) // Update in
		skillMgr.Countdown(dt)
	}
	b.UpdateSkillMove(b, dt)
	b.StateMgr.Update(dt)
	b.Countdown(dt) // 回蓝限制
}

func (t *BaseBattleUnit) GetMaxHp() int {
	return t.attrs.Get(consts.MaxHp)
}

func (bu *BaseBattleUnit) GetLvl() int {
	return bu.lvl
}

func (m *BaseBattleUnit) MoveTo(pos *geom.Vector2d) {
	mov := m.attrs.Get(consts.MoveSpeed)
	msp := float64(mov) / 100

	heading := pos.SubN(m.GetPosition())
	heading.Normalize()

	p := heading.MulN(msp * float64(DeltaTime) / 1000)
	m.GetPosition().Add(p)

	nowHeading := pos.SubN(m.GetPosition())
	nowHeading.Normalize()

	// 位置修正
	if int(nowHeading.LengthSq()) != int(heading.LengthSq()) {
		m.SetPosition(pos)
	}
}

/*************************************************
* 陷阱
**************************************************/
type Trap interface {
	GameObject
	SkillMover
	SkillRelative
	GetConfig() *conf.SkillTrapRowEx
	// GetDetectCircle() *geom.Circle
	GetCreateTime() int
	GetTriggerTime() int
	// IsDetect() bool
	// SetDetect(bool)
	GetValue() int
	SetValue(int)
}

type BaseTrap struct {
	GameObject
	SkillMover
	SkillRelative
	// config
	config      *conf.SkillTrapRowEx
	createTime  int // 创建时间
	triggerTime int // 触发时间
	val         int // 用来传递伤害等数值
	// detectCircle *geom.Circle // 检测圆
	// bDetect      bool         // 是否检测
	// collisionTarget battle.BattleUnit // 碰撞到的单位
}

func NewTrap(id int32, roleId int64, config *conf.SkillTrapRowEx) Trap {
	trap := &BaseTrap{
		GameObject:    NewGameObject(id, roleId, consts.GAMEOBJECT_TRAP, config.Id),
		SkillMover:    NewSkillMover(),
		SkillRelative: NewSkillRelative(),
		config:        config,
		// bDetect:    config.Condition == 1,
	}
	trap.SetBoundingShape(geom.NewShape(geom.ShapeType(config.Shape[0]), config.Shape[1:]...))
	return trap
}

func (t *BaseTrap) SetCreateTime(time int) {
	t.createTime = time
}

func (t *BaseTrap) GetConfig() *conf.SkillTrapRowEx {
	return t.config
}

// func (t *BaseTrap) GetDetectCircle() *geom.Circle {
// 	if t.detectCircle == nil {
// 		t.detectCircle = geom.NewCircle(float64(t.config.DetectRaduis))
// 	}
// 	mat := geom.MatrixTranslateV(t.GetPosition())
// 	t.detectCircle.Transform(mat)
// 	return t.detectCircle
// }

func (t *BaseTrap) GetCreateTime() int {
	return t.createTime
}

func (t *BaseTrap) GetTriggerTime() int {
	return t.triggerTime
}

// func (t *BaseTrap) IsDetect() bool {
// 	return t.bDetect
// }

// func (t *BaseTrap) SetDetect(b bool) {
// 	t.bDetect = b
// }

func (t *BaseTrap) SetValue(value int) {
	t.val = value
}

func (t *BaseTrap) GetValue() int {
	return t.val
}

// func (t *BaseTrap) GetContext() Scene {
// 	return t.unit.GetContext()
// }

/*************************************************
* 抛射物
**************************************************/
type Projectile interface {
	GameObject
	SkillRelative
	GetObjTarget() GameObject
	SetObjTarget(GameObject)
	GetPosTarget() *geom.Vector2d
	SetPosTarget(*geom.Vector2d)
}

type BaseProjectile struct {
	GameObject
	SkillRelative
	posTarget *geom.Vector2d
	objTarget GameObject
}

func NewProjectile(id int32, roleId int64) Projectile {
	return &BaseProjectile{
		GameObject:    NewGameObject(id, roleId, consts.GAMEOBJECT_PROJECTILE, int(id)),
		SkillRelative: NewSkillRelative(),
	}
}

func (p *BaseProjectile) GetObjTarget() GameObject {
	return p.objTarget
}

func (p *BaseProjectile) SetObjTarget(obj GameObject) {
	p.objTarget = obj
}

func (p *BaseProjectile) GetPosTarget() *geom.Vector2d {
	return p.posTarget
}

func (p *BaseProjectile) SetPosTarget(pos *geom.Vector2d) {
	p.posTarget = pos
}

// func (p *BaseProjectile) GetContext() Scene {
// 	return p.GetUnit().GetContext()
// }
