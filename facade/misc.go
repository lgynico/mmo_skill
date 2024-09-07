package facade

import (
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/geom"
	"github.com/lgynico/mmo_skill/utils"
)

// =========================================================================================================
// 属性
// =========================================================================================================
type Attrs interface {
	Get(key int) int
	GetBase(key int) int
	GetPerc(key int) int
	GetExtr(key int) int
	AddExtr(key int, val int)
	SubExtr(key int, val int)
	AddPerc(key int, val int)
	SubPerc(key int, val int)
	Reset()
	Range(func(key, val int))
}

type BaseAttrs struct {
	// base []int
	// perc []int
	// extr []int
	base  map[int]int
	perc  map[int]int
	extr  map[int]int
	final map[int]int
}

func NewAttr(attrs map[int]int) Attrs {
	// len := len(attrs)
	// a := &BaseAttrs{
	// 	base: make([]int, 0, len),
	// 	perc: make([]int, len),
	// 	extr: make([]int, len),
	// }
	// a.base = append(a.base, attrs...)

	attr := &BaseAttrs{
		base:  make(map[int]int),
		perc:  make(map[int]int),
		extr:  make(map[int]int),
		final: make(map[int]int),
	}

	for k, v := range attrs {
		attr.base[k] = v
		attr.updateFinal(k)
	}

	return attr
}

func (a *BaseAttrs) Reset() {
	len := len(a.base)
	for i := 0; i < len; i++ {
		a.perc[i] = 0
		a.extr[i] = 0
	}
}

func (a *BaseAttrs) Get(key int) int {
	return a.final[key]
}

func (a *BaseAttrs) GetBase(key int) int {
	return a.base[key]
}

func (a *BaseAttrs) GetPerc(key int) int {
	return a.perc[key]
}

func (a *BaseAttrs) GetExtr(key int) int {
	return a.extr[key]
}

func (a *BaseAttrs) AddExtr(key int, val int) {
	a.extr[key] += val
	a.updateFinal(key)
}

func (a *BaseAttrs) SubExtr(key int, val int) {
	a.extr[key] -= val
	a.updateFinal(key)
}

func (a *BaseAttrs) AddPerc(key int, val int) {
	a.perc[key] += val
	a.updateFinal(key)
}

func (a *BaseAttrs) SubPerc(key int, val int) {
	a.perc[key] -= val
	a.updateFinal(key)
}

func (a *BaseAttrs) updateFinal(key int) {
	base := float32(a.base[key])
	perc := float32(a.perc[key])
	extr := float32(a.extr[key])
	a.final[key] = int(base*(10000.0+perc)/10000.0 + extr)
}

func (a *BaseAttrs) Range(fn func(key, val int)) {
	for k, v := range a.final {
		fn(k, v)
	}
}

// =========================================================================================================
// BU属性
// =========================================================================================================
type BattleUnitProps interface {
	TakeDamage(val int) (hp, sp int)
	TakeHealing(val int)

	GetHp() int
	SetHp(int)
	IsDead() bool

	GetEp() int
	SetEp(int)
	AddEp(ep int, maxEp int, isFromBuff bool) int

	AddSp(buffId, value int)
	SubSp(value int) int
	GetSp() int
	RemoveSp(buffId int)

	SetTarget(BattleUnit)
	GetTarget() BattleUnit

	Countdown(dt int64)
}

type BaseBattleUnitProps struct {
	hp     int                         // 生命值
	ep     int                         // 能量值
	sp     []*utils.Pair[int32, int64] // 护盾值
	target BattleUnit

	epCounter int // 每秒ep计数
	epTimer   int // ep计时器
}

func NewBattleUnitProps() BattleUnitProps {
	return &BaseBattleUnitProps{
		sp: make([]*utils.Pair[int32, int64], 0),
	}
}

func (b *BaseBattleUnitProps) TakeDamage(val int) (hp, sp int) {
	remain := b.SubSp(val)
	sp = val - remain
	hp = remain

	if remain <= 0 {
		return
	}

	b.hp -= val
	if b.hp < 0 {
		b.hp = 0
	}

	return
}

func (b *BaseBattleUnitProps) TakeHealing(val int) {
	b.hp += val
}

func (b *BaseBattleUnitProps) GetHp() int                { return b.hp }
func (b *BaseBattleUnitProps) SetHp(hp int)              { b.hp = hp }
func (b *BaseBattleUnitProps) IsDead() bool              { return b.hp <= 0 }
func (b *BaseBattleUnitProps) GetEp() int                { return b.ep }
func (b *BaseBattleUnitProps) SetEp(ep int)              { b.ep = ep }
func (b *BaseBattleUnitProps) SetTarget(unit BattleUnit) { b.target = unit }
func (b *BaseBattleUnitProps) GetTarget() BattleUnit     { return b.target }

func (b *BaseBattleUnitProps) AddSp(buffId, value int) {
	for _, p := range b.sp {
		if p.Key() == int32(buffId) {
			p.SetValue(int64(value))
			return
		}
	}

	p := utils.NewPair(int32(buffId), int64(value))
	b.sp = append(b.sp, p)
}

func (b *BaseBattleUnitProps) SubSp(value int) int {
	remain := int64(value)
	for len(b.sp) > 0 && remain > 0 {
		p := b.sp[0]
		sheild := p.Value()
		if sheild >= remain {
			sheild -= remain
			remain = 0
		} else {
			remain -= sheild
			sheild = 0
		}

		if sheild <= 0 {
			b.sp = append(b.sp[:0], b.sp[1:]...)
		} else {
			p.SetValue(sheild)
		}
	}

	return int(remain)
}

func (b *BaseBattleUnitProps) GetSp() int {
	sp := 0
	for _, p := range b.sp {
		sp += int(p.Value())
	}
	return sp
}

func (b *BaseBattleUnitProps) RemoveSp(buffId int) {
	for i, p := range b.sp {
		if int(p.Key()) == buffId {
			b.sp = append(b.sp[:i], b.sp[i+1:]...)
			return
		}
	}
}

func (b *BaseBattleUnitProps) AddEp(ep int, maxEp int, isFromBuff bool) int {
	if isFromBuff {
		b.ep += ep

	} else {
		if b.epCounter >= consts.ENERGY_RECOVER_MAX_PER_SECOND {
			return 0
		}

		ep = utils.MaxInt(ep, consts.ENERGY_RECOVER_MAX_PER_SECOND-b.epCounter)
		b.ep += ep
		b.epCounter += ep
	}

	if b.ep > maxEp {
		b.ep = maxEp
	}

	return ep
}

func (b *BaseBattleUnitProps) Countdown(dt int64) {
	b.epTimer += int(dt)
	if b.epTimer >= 1000 {
		b.epTimer -= 1000
		b.epCounter = 0
	}
}

// =============================================================================
// 上下文管理器
// =============================================================================
type ContextHolder interface {
	GetContext() Battle
	SetContext(Battle)
}

type BaseContextHolder struct {
	ctx Battle
}

func NewContextHolder() ContextHolder {
	return &BaseContextHolder{}
}

func (ch *BaseContextHolder) GetContext() Battle {
	return ch.ctx
}

func (ch *BaseContextHolder) SetContext(ctx Battle) {
	ch.ctx = ctx
}

/***********************************************
* 技能位移器
***********************************************/
type SkillMover interface {
	IsSkillMove() bool
	SetSkillMove(bool)

	GetSkillMoveSpeed() float64
	SetSkillMoveSpeed(float64)

	GetSkillMoveTime() int64
	SetSkillMoveTime(int64)

	GetSkillMoveHeading() *geom.Vector2d
	SetSkillMoveHeading(*geom.Vector2d)

	UpdateSkillMove(obj GameObject, dt int64)

	// GetSkillMoveDist() int
	// SetSkillMoveDist(int)

	// GetSkillMovedDist() int
	// SetSkillMovedDist(int)
}

type BaseSkillMover struct {
	bMove   bool
	speed   float64
	time    int64
	heading *geom.Vector2d
	// dist      int
	// movedDist int
}

func NewSkillMover() SkillMover {
	return &BaseSkillMover{
		bMove:   false,
		speed:   0,
		time:    0,
		heading: nil,
		// dist:      0,
		// movedDist: 0,
	}
}

func (m *BaseSkillMover) IsSkillMove() bool                          { return m.bMove }
func (m *BaseSkillMover) SetSkillMove(move bool)                     { m.bMove = move }
func (m *BaseSkillMover) GetSkillMoveSpeed() float64                 { return m.speed }
func (m *BaseSkillMover) SetSkillMoveSpeed(speed float64)            { m.speed = speed }
func (m *BaseSkillMover) GetSkillMoveTime() int64                    { return m.time }
func (m *BaseSkillMover) SetSkillMoveTime(time int64)                { m.time = time }
func (m *BaseSkillMover) GetSkillMoveHeading() *geom.Vector2d        { return m.heading }
func (m *BaseSkillMover) SetSkillMoveHeading(heading *geom.Vector2d) { m.heading = heading }

func (m *BaseSkillMover) UpdateSkillMove(obj GameObject, dt int64) {
	if !m.IsSkillMove() {
		return
	}

	if m.GetSkillMoveTime() <= 0 {
		m.SetSkillMove(false)
		return
	}

	speed := m.GetSkillMoveSpeed()
	heading := m.GetSkillMoveHeading()
	moveVec := heading.MulN(speed * float64(dt) / 1000)
	obj.GetPosition().Add(moveVec)
}

// func (m *BaseSkillMover) GetSkillMoveDist() int       { return m.dist }
// func (m *BaseSkillMover) SetSkillMoveDist(dist int)   { m.dist = dist }
// func (m *BaseSkillMover) GetSkillMovedDist() int      { return m.movedDist }
// func (m *BaseSkillMover) SetSkillMovedDist(dist int)  { m.movedDist = dist }

// =========================================================================================================
// 边界
// =========================================================================================================
type Boundary interface {
	MinX() float64
	MinY() float64
	MinZ() float64
	MaxX() float64
	MaxY() float64
	MaxZ() float64
}

type BaseBoundary struct {
	minX, maxX float64
	minY, maxY float64
	minZ, maxZ float64
}

func NewBoundary(minX, minY, minZ float64, maxX, maxY, maxZ float64) Boundary {
	return &BaseBoundary{
		minX: minX,
		minY: minY,
		minZ: minZ,
		maxX: maxX,
		maxY: maxY,
		maxZ: maxZ,
	}
}

func (b *BaseBoundary) MinX() float64 { return b.minX }
func (b *BaseBoundary) MinY() float64 { return b.minY }
func (b *BaseBoundary) MinZ() float64 { return b.minZ }
func (b *BaseBoundary) MaxX() float64 { return b.maxX }
func (b *BaseBoundary) MaxY() float64 { return b.maxY }
func (b *BaseBoundary) MaxZ() float64 { return b.maxZ }

/* 技能关联器 */
type SkillRelative interface {
	GetUnit() BattleUnit
	SetUnit(BattleUnit)

	GetSkill() Skill
	SetSkill(Skill)
}

type BaseSkillRelative struct {
	unit  BattleUnit
	skill Skill
}

func NewSkillRelative() SkillRelative {
	return &BaseSkillRelative{}
}

func (sr *BaseSkillRelative) GetUnit() BattleUnit     { return sr.unit }
func (sr *BaseSkillRelative) SetUnit(unit BattleUnit) { sr.unit = unit }
func (sr *BaseSkillRelative) GetSkill() Skill         { return sr.skill }
func (sr *BaseSkillRelative) SetSkill(skill Skill)    { sr.skill = skill }

/* 变换 */
type Transform interface {
	GetPosition() *geom.Vector2d
	SetPosition(*geom.Vector2d)

	GetHeading() *geom.Vector2d
	SetHeading(*geom.Vector2d)

	GetBoundingShape(transform bool) geom.Shape
	SetBoundingShape(geom.Shape)
}

type BaseTransform struct {
	pos           *geom.Vector2d
	heading       *geom.Vector2d
	boundingShape geom.Shape
}

func NewTransform() Transform {
	return &BaseTransform{
		pos:     geom.NewVector2d(0, 0),
		heading: geom.NewVector2d(1, 0),
	}
}

func (t *BaseTransform) GetPosition() *geom.Vector2d {
	return t.pos
}
func (t *BaseTransform) SetPosition(pos *geom.Vector2d) {
	t.pos = pos
}

func (t *BaseTransform) GetHeading() *geom.Vector2d {
	return t.heading
}

func (t *BaseTransform) SetHeading(norm *geom.Vector2d) {
	t.heading = norm
}

func (t *BaseTransform) GetBoundingShape(transform bool) geom.Shape {
	if transform {
		// mat := geom.MatrixTranslateV(t.GetPosition())
		// t.boundingShape.Transform(mat)
		angle := t.heading.Angle()
		mat := geom.MatrixRotate(angle)
		mat = mat.Mul(geom.MatrixTranslateV(t.pos))
		t.boundingShape.Transform(mat)
	}
	return t.boundingShape
}

func (t *BaseTransform) SetBoundingShape(shape geom.Shape) {
	t.boundingShape = shape
}
