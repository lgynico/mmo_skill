package facade

import (
	"github.com/lgynico/mmo_skill/proto"
)

type Battle interface {
	SkillDerivantCreator

	Init()
	Start() *proto.BattleRecord
	AddTask(Task)
	RecordAction(*proto.BattleAction)
	RecordDamage(id int32, val int32)
	RecordHealing(id int32, val int32)
	RecordSuffering(id int32, val int32)

	GetTimeMillis() int64
	GetDeltaTime() int64

	GetObjs() []GameObject
	RemoveObj(obj GameObject)
	AddObj(obj GameObject)

	GetUnits() []BattleUnit
}

// /* Base */
// type BaseBattle struct {
// 	SkillDerivantCreator
// 	*Boundary
// 	time              int64
// 	teamRed, teamBlue [5]BattleUnit
// 	projectiles       map[GameObject]bool
// 	traps             map[GameObject]bool
// 	timeLimit         int64
// 	endType           int
// 	record            *proto.BattleRecord
// }

// func NewBaseBattle(teamRed, teamBlud [5]BattleUnit) Battle {
// 	return &BaseBattle{
// 		SkillDerivantCreator: NewSkillDerivantCreator(),
// 		time:                 0,
// 		timeLimit:            int64(conf.ConfMgr.ConstParamAdapter.BattleTimeLimit * 1000),
// 		teamRed:              teamRed,
// 		teamBlue:             teamBlud,
// 		projectiles:          make(map[GameObject]bool),
// 		traps:                make(map[GameObject]bool),
// 	}
// }

// func (b *BaseBattle) Init() {
// 	b.Boundary = boundary
// 	for i := 0; i < 5; i++ {
// 		unit := b.teamRed[i]
// 		if unit == nil {
// 			continue
// 		}

// 		unit.SetContext(b)

// 		id := 10 + i + 1
// 		unit.SetId(int32(id))

// 		vec := unitPos[id]
// 		unit.SetPosition(geom.NewVector2d(vec.X, vec.Z))
// 		unit.SetHeading(geom.NewVector2d(1, 0))

// 		unit.SetTeamId(1)
// 	}

// 	for i := 0; i < 5; i++ {
// 		unit := b.teamBlue[i]
// 		if unit == nil {
// 			continue
// 		}

// 		unit.SetContext(b)

// 		id := 20 + i + 1
// 		unit.SetId(int32(id))

// 		vec := unitPos[id]
// 		unit.SetPosition(geom.NewVector2d(vec.X, vec.Z))
// 		unit.SetHeading(geom.NewVector2d(-1, 0))

// 		unit.SetTeamId(2)
// 	}
// }

// func (b *BaseBattle) Start() *proto.BattleRecord {
// 	b.record = &proto.BattleRecord{
// 		Actions: make([]*proto.BattleAction, 0),
// 	}

// 	for {
// 		if b.checkOver() {
// 			break
// 		}

// 		b.time += deltaTime

// 		// handle buffs
// 		b.handleBuffs()
// 		// handle derivants
// 		b.handleDerivants()
// 		// handle unit actions
// 		b.updateObjects()
// 	}

// 	b.record.TimeElapsed = int32(b.time)
// 	if IsTimeout(b.endType) || IsTeamRedDestory(b.endType) {
// 		b.record.WhichWin = int32(proto.BattleResult_BLUE_WIN)
// 	} else {
// 		b.record.WhichWin = int32(proto.BattleResult_RED_WIN)
// 	}
// 	return b.record
// }

// func (b *BaseBattle) AddTask(task Task) {}

// func (b *BaseBattle) RecordAction(action *proto.BattleAction) {
// 	b.record.Actions = append(b.record.Actions, action)
// }

// func (b *BaseBattle) GetTimeMillis() int64 {
// 	return b.time
// }

// func (b *BaseBattle) GetObjs() []GameObject {
// 	objs := make([]GameObject, 0)
// 	for _, unit := range b.teamRed {
// 		if unit != nil {
// 			objs = append(objs, unit)
// 		}
// 	}
// 	for _, unit := range b.teamBlue {
// 		if unit != nil {
// 			objs = append(objs, unit)
// 		}
// 	}
// 	for p := range b.projectiles {
// 		objs = append(objs, p)
// 	}
// 	for t := range b.traps {
// 		objs = append(objs, t)
// 	}
// 	return objs
// }

// func (b *BaseBattle) AddObj(obj GameObject) {
// 	switch obj.GetObjType() {
// 	case proto.GameObjectType_OBJ_PROJECTILE:
// 		b.projectiles[obj] = true
// 	case proto.GameObjectType_OBJ_TRAP:
// 		b.traps[obj] = true
// 	}
// }

// func (b *BaseBattle) RemoveObj(obj GameObject) {
// 	if p, ok := obj.(Projectile); ok {
// 		delete(b.projectiles, p)
// 		return
// 	}

// 	if t, ok := obj.(Trap); ok {
// 		delete(b.traps, t)
// 		return
// 	}
// }

// func (b *BaseBattle) GetUnits() []BattleUnit {
// 	units := make([]BattleUnit, 0, 10)
// 	for _, unit := range b.teamRed {
// 		if unit != nil {
// 			units = append(units, unit)
// 		}
// 	}
// 	for _, unit := range b.teamBlue {
// 		if unit != nil {
// 			units = append(units, unit)
// 		}
// 	}

// 	return units
// }

// func (b *BaseBattle) checkOver() bool {
// 	if b.time >= b.timeLimit {
// 		b.endType |= BATTLE_END_TIME_OUT
// 		return true
// 	}

// 	teamRedDestory := true
// 	for _, unit := range b.teamRed {
// 		if unit != nil && !unit.IsDead() {
// 			teamRedDestory = false
// 			break
// 		}
// 	}

// 	if teamRedDestory {
// 		b.endType |= BATTLE_END_RED_DESTORY
// 		return true
// 	}

// 	teamBlueDestory := true
// 	for _, unit := range b.teamBlue {
// 		if unit != nil && !unit.IsDead() {
// 			teamBlueDestory = false
// 			break
// 		}
// 	}

// 	if teamBlueDestory {
// 		b.endType |= BATTLE_END_BLUE_DESTORY
// 		return true
// 	}

// 	return false
// }

// func (b *BaseBattle) handleBuffs() {
// 	for _, unit := range b.teamRed {
// 		if unit != nil && !unit.IsDead() {
// 			unit.GetBuffMgr().UpdateBuffs(unit, deltaTime)
// 		}
// 	}

// 	for _, unit := range b.teamBlue {
// 		if unit != nil && !unit.IsDead() {
// 			unit.GetBuffMgr().UpdateBuffs(unit, deltaTime)
// 		}
// 	}
// }

// func (b *BaseBattle) handleDerivants() {
// 	for p := range b.projectiles {
// 		p.Update(deltaTime)
// 	}

// 	for t := range b.traps {
// 		t.Update(deltaTime)
// 	}
// }

// func (b *BaseBattle) updateObjects() {
// 	for _, unit := range b.teamRed {
// 		if unit != nil && !unit.IsDead() {
// 			unit.Update(deltaTime)
// 		}
// 	}

// 	for _, unit := range b.teamBlue {
// 		if unit != nil && !unit.IsDead() {
// 			unit.Update(deltaTime)
// 		}
// 	}
// }
