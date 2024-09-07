package facade

import (
	"errors"
	"math"

	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/geom"
	"github.com/lgynico/mmo_skill/proto"
	"github.com/lgynico/mmo_skill/utils"
)

const RPC_SERV_TAG = "RpcServ"

var (
	ErrSceneNotExist    = errors.New("scene not exist")
	ErrSceneCfgNotExist = errors.New("scene cfg not exist")
	ErrPlayerIdIsNil    = errors.New("playerId is nil")
)

func CalSkillSingTime(skill Skill) int {
	config := skill.GetConfig()
	if !config.IsAtk {
		return config.SingTime
	}

	attrs := skill.GetUnit().GetAttrs()
	attackSpeed := attrs.Get(consts.AttackSpeed)
	if attackSpeed == 0 {
		return config.SingTime
	}

	asp := float64(attackSpeed) / 10000
	return int(float64(config.SingTime) * asp)

	// return calRatioTime(config.SingTime, attrs.GetBase(AttackSpeed), attrs.Get(AttackSpeed))
}

func CalSkillSpellTime(skill Skill) int {
	config := skill.GetConfig()
	return config.ProcessTime
}

func CalSkillShakeTime(skill Skill) int {
	config := skill.GetConfig()
	if !config.IsAtk {
		return config.ShakeTime
	}

	attrs := skill.GetUnit().GetAttrs()
	attackSpeed := attrs.Get(consts.AttackSpeed)
	if attackSpeed == 0 {
		return config.ShakeTime
	}

	asp := float64(attackSpeed) / 10000
	return int(float64(config.ShakeTime) * asp)
	// return calRatioTime(config.ShakeTime, attrs.GetBase(AttackSpeed), attrs.Get(AttackSpeed))
}

func CalSkillCd(skill Skill) int {
	config := skill.GetConfig()
	return config.CoolDown
	// if !config.IsAtk {
	// 	return config.CoolDown
	// }

	// attrs := skill.GetUnit().GetAttrs()
	// return calRatioTime(1000*attrs.GetBase(AttackSpeed)/10, attrs.GetBase(AttackSpeed), attrs.Get(AttackSpeed))
}

func calRatioTime(time int, base int, curr int) int {
	if base >= curr {
		return time
	}

	ratio := 1.0 * float64(base) / float64(curr)
	return int(1.0 * float64(time) * ratio)
}

func ConvertVector2d(vec *geom.Vector2d) *proto.Vector {
	return &proto.Vector{
		X: int64(vec.GetX() * 1000),
		Y: 0,
		Z: int64(vec.GetY() * 1000),
	}
}

func BuffFilter(buffs []Buff, filter func(Buff) bool) []Buff {
	if len(buffs) == 0 {
		return buffs
	}
	// unit := buffs[0].GetSrc().GetId()
	// log.Printf("%d 过滤buff之前，buff列表大小为：%d\n", unit, len(buffs))
	bs := buffs[:0]
	// bs := make([]battle.Buff, 0, len(buffs))
	for _, buff := range buffs {
		if filter(buff) {
			bs = append(bs, buff)
		}
	}
	// log.Printf("%d 过滤buff之后，buff列表大小为：%d\n", unit, len(bs))
	return bs
}

func FindTarget(unit BattleUnit, targetType int, radius int, enemyTarget BattleUnit) BattleUnit {
	// if targetType == 0 || (targetType&TARGET_TYPE_SELF) == TARGET_TYPE_SELF {
	// 	return unit
	// }

	var rad float64
	if radius == 0 {
		rad = math.MaxFloat64
	} else {
		rad = float64(radius) / 100
		rad *= rad
	}

	if targetType&TARGET_TYPE_ENEMY == TARGET_TYPE_ENEMY {
		dist := unit.GetPosition().DistSqTo(enemyTarget.GetPosition())
		// if dist <= rad {
		if utils.IsFloat64El(dist, rad) {
			return enemyTarget
		}
	}

	for _, u := range unit.GetContext().GetUnits() {
		if !IsTarget(u, unit, targetType) {
			continue
		}

		dist := unit.GetPosition().DistSqTo(u.GetPosition())
		// if dist <= rad {
		if utils.IsFloat64El(dist, rad) {
			return u
		}

	}

	return nil
}

func IsTarget(targetUnit BattleUnit, srcUnit BattleUnit, targetType int) bool {
	teamId := checkCharmState(srcUnit)

	bTarget := false
	if (targetType & TARGET_TYPE_SELF) == TARGET_TYPE_SELF {
		bTarget = bTarget || srcUnit.GetId() == targetUnit.GetId()
	}
	if (targetType & TARGET_TYPE_TEAMMATE) == TARGET_TYPE_TEAMMATE {
		bTarget = bTarget || (teamId == targetUnit.GetTeamId() && srcUnit.GetId() != targetUnit.GetId())
	}
	if (targetType & TARGET_TYPE_ENEMY) == TARGET_TYPE_ENEMY {
		bTarget = bTarget || teamId != targetUnit.GetTeamId() && srcUnit.GetId() != targetUnit.GetId()
	}
	if (targetType & TARGET_TYPE_DIE) == TARGET_TYPE_DIE {
		bTarget = bTarget || targetUnit.IsDead()
	} else {
		bTarget = bTarget && !targetUnit.IsDead()
	}

	return bTarget
}

func checkCharmState(srcUnit BattleUnit) int32 {
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

func EnergyRecover(unit BattleUnit, skill Skill, typ int) {
	addEp := unit.GetAttrs().Get(typ)
	ep := addEp + unit.GetEp()
	maxEp := unit.GetAttrs().Get(consts.MaxEnergy)
	change := unit.AddEp(ep, maxEp, false)

	if change > 0 {
		ctx := unit.GetContext()
		action := &proto.BattleAction{
			Time:  ctx.GetTimeMillis(),
			Unit:  unit.GetId(),
			Key:   int32(consts.ACTION_ENERGY_CHG),
			Value: int64(addEp),
		}
		ctx.RecordAction(action)
	}
}

func IsTeamBlueDestory(mark int) bool {
	return (BATTLE_END_BLUE_DESTORY & mark) == BATTLE_END_BLUE_DESTORY
}

func IsTeamRedDestory(mark int) bool {
	return (BATTLE_END_RED_DESTORY & mark) == BATTLE_END_RED_DESTORY
}

func IsTimeout(mark int) bool {
	return (BATTLE_END_TIME_OUT & mark) == BATTLE_END_TIME_OUT
}
