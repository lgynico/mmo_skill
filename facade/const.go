package facade

import "github.com/lgynico/mmo_skill/geom"

const DeltaTime = 1000 / 30

var CommonBoundary = NewBoundary(0, 0, 0, 6, 0, 12)

var UnitPos = map[int]*geom.Vector2d{
	11: geom.NewVector2d(2, 2.5),
	12: geom.NewVector2d(4, 2.5),
	13: geom.NewVector2d(1, 1.5),
	14: geom.NewVector2d(3, 1.5),
	15: geom.NewVector2d(5, 1.5),
	21: geom.NewVector2d(2, 9.5),
	22: geom.NewVector2d(4, 9.5),
	23: geom.NewVector2d(1, 10.5),
	24: geom.NewVector2d(3, 10.5),
	25: geom.NewVector2d(5, 10.5),
}

const (
	SKILL_BURNING_RAGE = 113031
	SKILL_OVER_HEAL    = 123013
	SKILL_HEAL_FRONT   = 123011
	SKILL_NO_DIE       = 133014
)

// 技能类型
const (
	SKILL_TYPE_ACTIVE  = 0x1
	SKILL_TYPE_PASSIVE = 0x2
	SKILL_TYPE_TOGGLE  = 0x4
)

// 技能阶段
const (
	SKILL_PHASE_READY = iota // 准备
	SKILL_PHASE_SING         // 吟唱
	SKILL_PHASE_SPELL        // 施法
	// SKILL_PHASE_GUIDE        // 引导
	SKILL_PHASE_SHAKE // 后摇
	SKILL_PHASE_END   // 结束
)

/* buff类型 */
type BuffType int

const (
	BUFF_STAT_MOVE_BAN        = 1 << iota // 移动禁止
	BUFF_STAT_ATTACK_BAN                  // 攻击禁止
	BUFF_STAT_SKILL_BAN                   // 技能禁止
	BUFF_STATE_INCREDIBLE                 // 无敌
	BUFF_STATE_UNSTOPPABLE                // 霸体
	BUFF_STAT_INVISIBLE                   // 不可见
	BUFF_STAT_CHARM                       // 魅惑
	BUFF_STAT_FINAL_SKILL_BAN             // 大招禁止
)

const (
	TARGET_TYPE_ANYONE int = 0
	TARGET_TYPE_ENEMY  int = 1 << (iota - 1)
	TARGET_TYPE_TEAMMATE
	TARGET_TYPE_SELF
	TARGET_TYPE_DIE
)

const (
	BATTLE_END_RED_DESTORY = 1 << iota
	BATTLE_END_BLUE_DESTORY
	BATTLE_END_TIME_OUT
)
