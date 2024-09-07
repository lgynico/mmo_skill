package consts

type GameObjectType int32

const (
	GAMEOBJECT_UNKNOWN    GameObjectType = 0
	GAMEOBJECT_UNIT       GameObjectType = 1 // 战斗单位
	GAMEOBJECT_PROJECTILE GameObjectType = 2 // 炮弹
	GAMEOBJECT_TRAP       GameObjectType = 3 // 陷阱
)
