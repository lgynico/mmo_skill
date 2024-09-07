package common

import (
	"math"

	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/geom"
	"github.com/lgynico/mmo_skill/utils"
)

type PosChooser interface {
	Choose(ctx facade.Battle, chooser facade.GameObject, caster facade.GameObject, params []int, target ...facade.GameObject)
}

type PosChooseTemplate interface {
	doChoose(ctx facade.Battle, caster facade.GameObject, params []int, target ...facade.GameObject) (heading *geom.Vector2d, pos *geom.Vector2d)
}

var mapType2Chooser = map[int]PosChooser{
	POS_TYPE_SELF_POS:  newTrapPosChooser(&SelfPosChooser{}),
	POS_TYPE_DIRECTION: newTrapPosChooser(&DirectionPosChooser{}),
	POS_TYPE_CIRCLE:    newTrapPosChooser(&CircleRangePosChooser{}),
}
var baseChooser = &BasePosChooser{}

func newTrapPosChooser(t PosChooseTemplate) PosChooser {
	base := &BasePosChooser{
		template: t,
	}
	if t == nil {
		base.template = base
	}
	return base
}

func GetPosChooser(t int) PosChooser {
	chooser, ok := mapType2Chooser[t]
	if ok {
		return chooser
	}
	return baseChooser
}

// ==========================================================================
// 基类
// ==========================================================================
type BasePosChooser struct {
	template PosChooseTemplate
}

func (c *BasePosChooser) Choose(ctx facade.Battle, chooser facade.GameObject, caster facade.GameObject, params []int, target ...facade.GameObject) {
	heading, pos := c.template.doChoose(ctx, caster, params, target...)
	chooser.SetHeading(heading)
	chooser.SetPosition(pos)

	angle := heading.Angle()
	mat := geom.MatrixRotate(angle)
	mat = mat.Mul(geom.MatrixTranslateV(pos))
	shape := chooser.GetBoundingShape(true)
	shape.Transform(mat)
}

func (c *BasePosChooser) doChoose(ctx facade.Battle, caster facade.GameObject, params []int, target ...facade.GameObject) (heading *geom.Vector2d, pos *geom.Vector2d) {
	heading = geom.NewVector2d(1, 0)
	pos = geom.NewVector2d(0, 0)
	return
}

// ==========================================================================
// 自身中心点
// ==========================================================================
type SelfPosChooser struct {
	*BasePosChooser
}

func (c *SelfPosChooser) doChoose(ctx facade.Battle, caster facade.GameObject, params []int, target ...facade.GameObject) (heading *geom.Vector2d, pos *geom.Vector2d) {
	heading = caster.GetHeading().Copy()
	pos = caster.GetPosition().Copy()
	return
}

// ==========================================================================
// 固定方向
// ==========================================================================
type DirectionPosChooser struct {
	*BasePosChooser
}

func (c *DirectionPosChooser) doChoose(ctx facade.Battle, caster facade.GameObject, params []int, target ...facade.GameObject) (heading *geom.Vector2d, pos *geom.Vector2d) {
	heading = caster.GetHeading().Copy()
	pos = caster.GetPosition().Copy()

	radian := float64(0)
	posType := params[1]
	if posType >= 0 && posType <= 12 {
		if posType == 12 {
			posType = 0
		}
		radian = math.Pi / 6 * float64(posType)
	} else if posType <= 16 {
		radian = math.Pi / 4 * float64(posType-1)
	}
	heading.Rotate(radian)

	dist := params[2]
	if dist > 0 {
		v := heading.NormalizeN()
		v.Mul(float64(dist))
		pos.Translate(v.GetX(), v.GetY())
	}

	return
}

// ==========================================================================
// 圆内范围
// ==========================================================================
type CircleRangePosChooser struct {
	*BasePosChooser
}

func (c *CircleRangePosChooser) doChoose(ctx facade.Battle, caster facade.GameObject, params []int, target ...facade.GameObject) (heading *geom.Vector2d, pos *geom.Vector2d) {
	radius := params[1]
	if len(target) != 0 && target[0] != nil {
		hisPos := target[0].GetPosition()
		myPos := caster.GetPosition()

		heading = hisPos.SubN(myPos)
		heading.Normalize()

		distSq := hisPos.DistSqTo(myPos)
		radiusSq := radius * radius
		if utils.IsFloat64Gt(distSq, float64(radiusSq)) {
			// 圆上距离目标最近的点
			v := heading.MulN(float64(radius))
			pos = myPos.TranslateN(v.GetX(), v.GetY())
		} else {
			pos = hisPos
		}
	} else {
		// 随机找个点
		tx := utils.RandFloat64Positive(float64(radius))
		ty := utils.RandFloat64Positive(float64(radius))
		pos = caster.GetPosition().TranslateN(tx, ty)
		// 边界trim
		// TODO
		// pos.Trim(0, 0, ctx.GetMaxX(), ctx.GetMaxY())
		heading = pos.SubN(caster.GetPosition())
	}
	return
}
