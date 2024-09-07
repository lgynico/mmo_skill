package consts

type ActionType int32

const (
	ACTION_DO_NOTHING  ActionType = 0  // 不做事
	ACTION_MOVE        ActionType = 1  // 移动
	ACTION_ATTACK      ActionType = 2  // 攻击
	ACTION_SPELL       ActionType = 3  // 放技能
	ACTION_HURT        ActionType = 4  // 伤害
	ACTION_HEAL        ActionType = 5  // 治疗
	ACTION_ADD_BUFF    ActionType = 6  // 加buff
	ACTION_REM_BUFF    ActionType = 7  // 减buff
	ACTION_DAED        ActionType = 8  // 死亡
	ACTION_ENERGY_CHG  ActionType = 9  // 能量变化
	ACTION_CREATE_PROJ ActionType = 10 // 创建子弹
)
