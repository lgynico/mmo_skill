package common

import (
	"math"
	"math/rand/v2"

	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/geom"
	"github.com/lgynico/mmo_skill/utils"
)

func MoveAhead(pos, heading *geom.Vector2d, speed float64, dt int64) *geom.Vector2d {
	v := heading.NormalizeN()
	v.Mul(speed * float64(dt) / 1000)
	return pos.AddN(v)
}

func Move(pos, targetPos *geom.Vector2d, speed float64, dt int64) *geom.Vector2d {
	heading := targetPos.SubN(pos)
	heading.NormalizeN()
	heading.Mul(speed * float64(dt) / 1000)
	return pos.AddN(heading)
}

// 不可越过终点
func MoveNoCross(pos, targetPos *geom.Vector2d, speed float64, dt int64, stopDist float64) (*geom.Vector2d, int64) {
	vHeading := pos.SubN(targetPos)
	vHeading.Normalize()
	dis := vHeading.MulN(stopDist)
	realTargetPos := targetPos.AddN(dis)

	heading := realTargetPos.SubN(pos)
	heading.Normalize()
	dis = heading.MulN(speed * float64(dt) / 1000)
	nowPos := pos.AddN(dis)

	newHeading := realTargetPos.SubN(nowPos)
	newHeading.Normalize()

	if !newHeading.Eq(heading) {
		toTarget := realTargetPos.SubN(pos)
		dist := toTarget.Length()
		useTime := int64((dist / speed) * 1000)
		return realTargetPos, dt - useTime
	}

	return nowPos, 0
}

func Arrive(pos, targetPos *geom.Vector2d, deceleration, decelerationWeaker float64, maxSpeed float64, dt int64, stopDist float64) *geom.Vector2d {
	toTarget := targetPos.SubN(pos)
	dist := toTarget.Length()
	if utils.IsFloat64Gt(dist, stopDist) {
		speed := dist / (deceleration * decelerationWeaker)
		speed = math.Min(speed, maxSpeed)

		toTarget.Normalize()
		toTarget.Mul(speed * float64(dt) / 1000)
		return pos.AddN(toTarget)

	}

	return pos
}

func EnsureBuffValue(buff facade.Buff, target facade.BattleUnit, isDamage bool) int {
	var val int
	if buff.GetConfig().IsOuterVal { // 外部参数
		val = buff.GetValue()
	} else { // 计算值
		valueType := buff.GetConfig().Params[0]
		value := buff.GetConfig().Params[1]
		switch valueType {
		case VALUE_TYPE_ATK_PERC:
			atk := buff.GetSrc().GetAttrs().Get(consts.Attack)
			valF := float64(atk) * (float64(value) / 10000.0)
			val = int(math.Round(valF))
		case VALUE_TYPE_MAX_HP_PERC:
			hp := target.GetAttrs().Get(consts.MaxHp)
			valF := float64(hp) * (float64(value) / 10000.0)
			val = int(math.Round(valF))
		case VALUE_TYPE_CURR_HP_PERC:
			currHp := target.GetHp()
			valF := float64(currHp) * (float64(value) / 10000.0)
			val = int(math.Round(valF))
		default:
			val = value
		}
	}

	// 伤害数值享受伤害加成
	if isDamage {
		attrs := buff.GetSrc().GetAttrs()
		val = damageAddition(val, attrs)
	}
	return val
}

func damageAddition(dmg int, attrs facade.Attrs) int {
	d := float64(dmg)
	dmi := float64(attrs.GetPerc(consts.DamageIncrease))

	d = d * dmi / 10000
	dmg += int(math.Round(d))
	dmg += attrs.GetExtr(consts.DamageIncrease)
	return dmg
}

// 计算攻击间隔
func CalAttackInterval(asp int) int {
	atkPerSec := float64(asp) / 10.0
	return int(1000 / atkPerSec)
}

func BuffFilter(buffs []facade.Buff, filter func(facade.Buff) bool) []facade.Buff {
	if len(buffs) == 0 {
		return buffs
	}
	// unit := buffs[0].GetSrc().GetId()
	// log.Printf("%d 过滤buff之前，buff列表大小为：%d\n", unit, len(buffs))
	bs := buffs[:0]
	// bs := make([]facade.Buff, 0, len(buffs))
	for _, buff := range buffs {
		if filter(buff) {
			bs = append(bs, buff)
		}
	}
	// log.Printf("%d 过滤buff之后，buff列表大小为：%d\n", unit, len(bs))
	return bs
}

// func IsTarget(targetUnit facade.BattleUnit, srcUnit facade.BattleUnit, targetType int) bool {
// 	bTarget := false
// 	if (targetType & facade.TARGET_TYPE_SELF) == facade.TARGET_TYPE_SELF {
// 		bTarget = bTarget || srcUnit.GetId() == targetUnit.GetId()
// 	}
// 	if (targetType & facade.TARGET_TYPE_TEAMMATE) == facade.TARGET_TYPE_TEAMMATE {
// 		bTarget = bTarget || srcUnit.GetTeamId() == targetUnit.GetTeamId()
// 	}
// 	if (targetType & facade.TARGET_TYPE_ENEMY) == facade.TARGET_TYPE_ENEMY {
// 		bTarget = bTarget || srcUnit.GetTeamId() != targetUnit.GetTeamId()
// 	}
// 	return bTarget
// }

func FindUnits(ctx facade.Battle, pos *geom.Vector2d, radius float64, count uint, filter func(facade.BattleUnit) bool) []facade.BattleUnit {
	units := make([]facade.BattleUnit, 0, count)
	if count == 0 {
		return units
	}

	if filter == nil {
		filter = func(unit facade.BattleUnit) bool { return true }
	}

	for _, unit := range ctx.GetUnits() {
		if !filter(unit) {
			continue
		}

		distSq := pos.DistSqTo(unit.GetPosition())
		r := radius + unit.GetBoundingShape(false).Radius()

		if utils.IsFloat64Gt(distSq, r*r) {
			continue
		}

		units = append(units, unit)
		count--
		if count == 0 {
			break
		}
	}

	return units
}

func FindEnemyHeros(src facade.BattleUnit, radius float64, count uint) []facade.BattleUnit {
	filter := func(target facade.BattleUnit) bool {
		return target.GetId() != src.GetId() &&
			target.GetObjType() == consts.GAMEOBJECT_UNIT &&
			target.GetTeamId() != src.GetTeamId()
	}
	return FindUnits(src.GetContext(), src.GetPosition(), radius, count, filter)
}

func FindAllEnemyHeros(src facade.BattleUnit, includeDead bool) []facade.BattleUnit {
	units := make([]facade.BattleUnit, 0, 5)
	for _, unit := range src.GetContext().GetUnits() {
		if unit == src || unit.GetTeamId() == src.GetTeamId() {
			continue
		}
		if !includeDead && unit.IsDead() {
			continue
		}

		units = append(units, unit)
	}
	return units
}

func FindNearestEnemy(src facade.BattleUnit, units []facade.BattleUnit) facade.BattleUnit {
	var findUnit facade.BattleUnit
	var dist float64

	teamId := CheckCharmState(src)
	for _, unit := range units {
		if src == unit {
			continue
		}

		if teamId == unit.GetTeamId() {
			continue
		}

		if unit.IsDead() {
			continue
		}

		if !unit.CanSelect() {
			continue
		}

		if findUnit == nil {
			findUnit = unit
			dist = findUnit.GetPosition().DistSqTo(src.GetPosition())
			continue
		}

		d := findUnit.GetPosition().DistSqTo(src.GetPosition())
		if d < dist {
			dist = d
			findUnit = unit
		}
	}

	return findUnit
}

func CheckCharmState(srcUnit facade.BattleUnit) int32 {
	teamId := srcUnit.GetTeamId()
	if srcUnit.IsCharmState() {
		if teamId == 1 {
			teamId = 2
		} else {
			teamId = 1
		}
	}

	return teamId
}

func FixedPos(pos *geom.Vector2d, unit facade.BattleUnit, ctx facade.Battle) *geom.Vector2d {

	r1 := unit.GetBoundingShape(false).Radius()
	for _, target := range ctx.GetUnits() {
		if target.IsDead() {
			continue
		}

		vShift := pos.SubN(target.GetPosition())
		len := vShift.Length()
		vShift.Normalize()
		r2 := target.GetBoundingShape(false).Radius()

		len = r1 + r2 - len
		if len < 0 {
			continue
		}

		vShift.Mul(len)
		pos.Add(vShift)
	}

	return pos
}

func TrimToInside(point *geom.Vector2d, r *geom.Rect) *geom.Vector2d {
	if r == nil {
		return point
	}

	minX := r.LeftBottom().GetX()
	minY := r.LeftBottom().GetY()
	maxX := r.RightTop().GetX()
	maxY := r.RightTop().GetY()

	point.Trim(minX, minY, maxX, maxY)

	return point
}

func TrimToOutside(point *geom.Vector2d, r *geom.Rect) *geom.Vector2d {
	if r == nil {
		return point
	}

	x := point.GetX()
	y := point.GetY()

	minX := r.LeftBottom().GetX()
	minY := r.LeftBottom().GetY()
	maxX := r.RightTop().GetX()
	maxY := r.RightTop().GetY()

	if !(x > minX && x < maxX && y > minY && y < maxY) {
		return point
	}

	dLeft := x - minX
	dRight := maxX - x
	dTop := maxY - y
	dBot := y - minY

	diffs := []float64{dLeft, dRight, dBot, dTop}

	dMin := diffs[0]
	index := 0
	for i, d := range diffs {
		if dMin < d {
			continue
		}

		dMin = d
		index = i
	}

	switch index {
	case 0:
		x = minX
	case 1:
		x = maxX
	case 2:
		y = minY
	case 3:
		y = maxY
	}

	point.SetX(x)
	point.SetY(y)

	return point
}

func RandInRange(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// func HandleTrap(unit facade.BattleUnit, config *conf.SkillTrapRowEx, dir int, dist int, value int) {
// 	ctx := unit.GetContext()
// 	t := ctx.CreateTrap(config, unit)
// 	t.SetValue(value)

// 	// posChooser := GetPosChooser(dir)
// 	// posChooser.Choose(ctx, t, unit, posParams, unit.GetTarget())
// 	ensureTrapPos(t, unit, dir, dir)

// 	// 延迟时间为0则立刻执行
// 	if config.Delay == 0 {
// 		t.Update(ctx.GetDeltaTime())
// 		if config.Times > 1 {
// 			ctx.Enter(t, t.GetPosition())
// 		}
// 	} else {
// 		ctx.Enter(t, t.GetPosition())
// 	}
// }

// func ensureTrapPos(trap facade.Trap, unit facade.GameObject, dir int, dist int) {
// 	heading, pos := chooseTrapPos(unit, dir, dist)
// 	trap.SetHeading(heading)
// 	trap.SetPosition(pos)

// 	angle := heading.Angle()
// 	mat := geom.MatrixRotate(angle)
// 	mat = mat.Mul(geom.MatrixTranslateV(pos))
// 	shape := trap.GetBoundingShape(true)
// 	shape.Transform(mat)
// }

// func chooseTrapPos(unit facade.GameObject, dir int, dist int) (*geom.Vector2d, *geom.Vector2d) {
// 	heading := unit.GetHeading().Copy()
// 	pos := unit.GetPosition().Copy()

// 	radian := float64(0)
// 	if dir >= 0 && dir <= 12 {
// 		if dir == 12 {
// 			dir = 0
// 		}
// 		radian = math.Pi / 6 * float64(dir)
// 	} else if dir <= 16 {
// 		radian = math.Pi / 4 * float64(dir-1)
// 	}
// 	heading.Rotate(radian)

// 	if dist > 0 {
// 		v := heading.NormalizeN()
// 		v.Mul(float64(dist))
// 		pos.Translate(v.GetX(), v.GetY())
// 	}

// 	return heading, pos
// }

func OnDead(ctx facade.Battle, obj facade.GameObject, remove bool) {
	//ctx.PushObjAction(obj.GetId(), proto.ActionId_ACT_DEAD)

	// if remove {
	// 	ctx.Leave(obj)
	// }
}

func UpdateSkillMove(mover facade.SkillMover, dt int64) {
	// 技能位移
	if !mover.IsSkillMove() {
		return
	}

	if mover.GetSkillMoveTime() <= 0 {
		mover.SetSkillMove(false)
		return
	}

	speed := mover.GetSkillMoveSpeed()
	heading := mover.GetSkillMoveHeading()
	moveVec := heading.MulN(speed * float64(dt) / 1000)

	// dist := mover.GetSkillMoveDist()
	// distMoved := mover.GetSkillMovedDist()

	// if distMoved >= dist {
	// 	mover.SetSkillMove(false)
	// 	return
	// }

	// toMovDist := speed * int(dt) / 100
	// if (toMovDist + distMoved) > dist {
	// 	toMovDist = dist - distMoved
	// 	mover.SetSkillMove(false)
	// }

	// distMoved += toMovDist
	// mover.SetSkillMovedDist(distMoved)

	obj := mover.(facade.GameObject)
	// v := obj.GetHeading().MulN(float64(toMovDist))
	obj.GetPosition().Add(moveVec)
}
