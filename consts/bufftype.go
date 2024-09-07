package consts

type BuffType int32

const (
	BUFF_ATTRIBUTE BuffType = 1 // 属性
	BUFF_DAMAGE    BuffType = 2 // 伤害
	BUFF_HEALING   BuffType = 3 // 治疗
	BUFF_STATE     BuffType = 4 // 状态
	BUFF_MOVE      BuffType = 5 // 位移
	BUFF_SHIELD    BuffType = 6 // 护盾
	BUFF_HOLA      BuffType = 7 // 光环
	BUFF_ENERGY    BuffType = 8 // 能量
)
