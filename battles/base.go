package battles

import (
	"log"

	"github.com/lgynico/mmo_skill/buff"
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/gameobj/projectile"
	"github.com/lgynico/mmo_skill/gameobj/trap"
	"github.com/lgynico/mmo_skill/geom"
	"github.com/lgynico/mmo_skill/halo"
	"github.com/lgynico/mmo_skill/proto"
)

type BaseBattle struct {
	facade.Boundary
	// facade.GameObjectMgr
	facade.SkillDerivantCreator
	// scene.EndCondition
	// scene.SceneHandler
	time              int64
	teamRed, teamBlue [5]facade.BattleUnit
	// units          [10]facade.BattleUnit
	projectiles map[facade.GameObject]bool
	traps       map[facade.GameObject]bool
	timeLimit   int64
	endType     int
	record      *proto.BattleRecord
	tasks       []facade.Task
	// mapActionCount map[int]*utils.Pair
}

func NewBaseBattle(teamRed, teamBlud [5]facade.BattleUnit) facade.Battle {
	return &BaseBattle{
		Boundary:    facade.CommonBoundary,
		time:        0,
		timeLimit:   consts.BATTLE_TIME_LIMIT_SEC * 1000,
		teamRed:     teamRed,
		teamBlue:    teamBlud,
		projectiles: make(map[facade.GameObject]bool),
		traps:       make(map[facade.GameObject]bool),
		record: &proto.BattleRecord{
			Teams:   make([]*proto.BattleTeam, 2),
			Stats:   make([]*proto.BattleStat, 0),
			Actions: make([]*proto.BattleAction, 0),
		},
		// mapActionCount: make(map[int]*utils.Pair),
		tasks: make([]facade.Task, 0),
	}
}

func (b *BaseBattle) Init() {
	teamRed := &proto.BattleTeam{
		Units: make([]*proto.BattleUnit, 0),
	}

	for i := 0; i < 5; i++ {
		unit := b.teamRed[i]
		if unit == nil {
			continue
		}

		b.initUnit(unit, int32(11+i), 1, geom.NewVector2d(0, 1))
		teamRed.Units = append(teamRed.Units, unitToProto(unit))
		stat := &proto.BattleStat{ID: unit.GetId()}
		b.record.Stats = append(b.record.Stats, stat)
	}
	b.record.Teams[0] = teamRed

	teamBlue := &proto.BattleTeam{
		Units: make([]*proto.BattleUnit, 0),
	}
	for i := 0; i < 5; i++ {
		unit := b.teamBlue[i]
		if unit == nil {
			continue
		}

		b.initUnit(unit, int32(21+i), 2, geom.NewVector2d(0, -1))
		teamBlue.Units = append(teamBlue.Units, unitToProto(unit))
		stat := &proto.BattleStat{ID: unit.GetId()}
		b.record.Stats = append(b.record.Stats, stat)
	}
	b.record.Teams[1] = teamBlue
}

func (b *BaseBattle) initUnit(unit facade.BattleUnit, id int32, teamId int32, heading *geom.Vector2d) {
	unit.SetContext(b)
	unit.SetId(id)
	unit.SetTeamId(teamId)

	vec := facade.UnitPos[int(id)]
	unit.SetPosition(vec.Copy())
	unit.SetHeading(heading)

	facade.SendEvent(facade.EVENT_IMMEDIATELY, unit, nil)

	ep := unit.GetAttrs().Get(consts.Energy)
	unit.SetEp(ep)
	unit.SetHp(unit.GetAttrs().Get(consts.MaxHp))
	// unit.SetHp(1000000)
}

func (b *BaseBattle) Start() *proto.BattleRecord {
	log.Println("战斗开始！！！！！！")
	for {
		if b.checkOver() {
			break
		}

		b.time += facade.DeltaTime

		// handle tasks
		b.handleTasks()

		// handle buffs
		b.handleHalos()
		b.handleBuffs()
		// handle derivants
		b.handleDerivants()
		// handle unit actions
		b.updateObjects()
	}

	log.Println("战斗结束！！！！！！")
	b.record.TimeElapsed = b.time
	if facade.IsTimeout(b.endType) || facade.IsTeamRedDestory(b.endType) {
		b.record.WhichWin = int32(consts.BATTLERESULT_BLUE_WIN)
		log.Println("蓝方胜！！！！！！")
	} else {
		b.record.WhichWin = int32(consts.BATTLERESULT_RED_WIN)
		log.Println("红方胜！！！！！！")
	}

	// var str string
	// for key, pair := range b.mapActionCount {
	// 	str += fmt.Sprintf("(key:%d, cnt:%d, len:%d)", key, pair.GetKey(), pair.GetVal())
	// }
	// facade.Logger.D(fmt.Sprintf("战斗动作总结: %s", str))

	return b.record
}

func (b *BaseBattle) AddTask(task facade.Task) {
	b.tasks = append(b.tasks, task)
}

func (b *BaseBattle) RecordAction(action *proto.BattleAction) {
	b.record.Actions = append(b.record.Actions, action)
	// pair, ok := b.mapActionCount[int(action.Key)]
	// if !ok {
	// 	pair = utils.NewPair(0, 0)
	// 	b.mapActionCount[int(action.Key)] = pair
	// }

	// bs, _ := proto.Marshal(action)
	// pair.SetKey(pair.GetKey().(int) + 1)
	// pair.SetVal(pair.GetVal().(int) + len(bs))
}

func (b *BaseBattle) RecordDamage(id int32, val int32) {
	for _, stat := range b.record.Stats {
		if stat.ID == id {
			stat.Damage += val
			return
		}
	}
}

func (b *BaseBattle) RecordHealing(id int32, val int32) {
	for _, stat := range b.record.Stats {
		if stat.ID == id {
			stat.Healing += val
			return
		}
	}
}

func (b *BaseBattle) RecordSuffering(id int32, val int32) {
	for _, stat := range b.record.Stats {
		if stat.ID == id {
			stat.Suffering += val
			return
		}
	}
}

func (b *BaseBattle) GetTimeMillis() int64 {
	return b.time
}

func (b *BaseBattle) GetDeltaTime() int64 {
	return facade.DeltaTime
}

func (b *BaseBattle) GetObjs() []facade.GameObject {
	objs := make([]facade.GameObject, 0)
	for _, unit := range b.teamRed {
		if unit != nil {
			objs = append(objs, unit)
		}
	}
	for _, unit := range b.teamBlue {
		if unit != nil {
			objs = append(objs, unit)
		}
	}
	for p := range b.projectiles {
		objs = append(objs, p)
	}
	for t := range b.traps {
		objs = append(objs, t)
	}
	return objs
}

func (b *BaseBattle) AddObj(obj facade.GameObject) {
	switch obj.GetObjType() {
	case consts.GAMEOBJECT_PROJECTILE:
		b.projectiles[obj] = true
	case consts.GAMEOBJECT_TRAP:
		b.traps[obj] = true
	}
}

func (b *BaseBattle) RemoveObj(obj facade.GameObject) {
	if obj.GetObjType() == consts.GAMEOBJECT_PROJECTILE {
		delete(b.projectiles, obj)
		return
	}
	if obj.GetObjType() == consts.GAMEOBJECT_TRAP {
		delete(b.traps, obj)
		return
	}
}

func (b *BaseBattle) GetUnits() []facade.BattleUnit {
	units := make([]facade.BattleUnit, 0, 10)
	for _, unit := range b.teamRed {
		if unit != nil {
			units = append(units, unit)
		}
	}
	for _, unit := range b.teamBlue {
		if unit != nil {
			units = append(units, unit)
		}
	}

	return units
}

func (b *BaseBattle) checkOver() bool {
	if b.time >= b.timeLimit {
		b.endType |= facade.BATTLE_END_TIME_OUT
		log.Println("战斗超时......")
		return true
	}

	teamRedDestory := true
	for _, unit := range b.teamRed {
		if unit != nil && !unit.IsDead() {
			teamRedDestory = false
			break
		}
	}

	if teamRedDestory {
		b.endType |= facade.BATTLE_END_RED_DESTORY
		log.Println("红方全灭......")
		return true
	}

	teamBlueDestory := true
	for _, unit := range b.teamBlue {
		if unit != nil && !unit.IsDead() {
			teamBlueDestory = false
			break
		}
	}

	if teamBlueDestory {
		b.endType |= facade.BATTLE_END_BLUE_DESTORY
		log.Println("蓝方全灭......")
		return true
	}

	return false
}

func (b *BaseBattle) handleTasks() {
	for i := len(b.tasks) - 1; i >= 0; i++ {
		task := b.tasks[i]
		if task.Exec(b) {
			b.tasks = append(b.tasks[:i], b.tasks[i+1:]...)
		}
	}
}

func (b *BaseBattle) handleHalos() {
	for _, unit := range b.teamRed {
		if unit != nil && !unit.IsDead() {
			unit.GetHaloMgr().Update(facade.DeltaTime)
		}
	}

	for _, unit := range b.teamBlue {
		if unit != nil && !unit.IsDead() {
			unit.GetHaloMgr().Update(facade.DeltaTime)
		}
	}
}

func (b *BaseBattle) handleBuffs() {
	for _, unit := range b.teamRed {
		if unit != nil && !unit.IsDead() {
			unit.GetBuffMgr().UpdateBuffs(unit, facade.DeltaTime)
		}
	}

	for _, unit := range b.teamBlue {
		if unit != nil && !unit.IsDead() {
			unit.GetBuffMgr().UpdateBuffs(unit, facade.DeltaTime)
		}
	}
}

func (b *BaseBattle) handleDerivants() {
	for p := range b.projectiles {
		p.Update(facade.DeltaTime)
	}

	for t := range b.traps {
		t.Update(facade.DeltaTime)
	}
}

func (b *BaseBattle) updateObjects() {
	for _, unit := range b.teamRed {
		if unit != nil && !unit.IsDead() {
			unit.Update(facade.DeltaTime)
		}
	}

	for _, unit := range b.teamBlue {
		if unit != nil && !unit.IsDead() {
			unit.Update(facade.DeltaTime)
		}
	}
}

func (b *BaseBattle) CreateTrap(config *conf.SkillTrapRowEx, unit facade.BattleUnit) facade.Trap {
	id := 1
	return trap.NewTrap(int32(id), config, unit)
}

func (b *BaseBattle) CreateProjectile(skill facade.Skill, config *conf.SkillProjectileRowEx) facade.Projectile {
	id := 1
	return projectile.NewProjectile(int32(id), skill.GetUnit().GetRoleId(), skill, config)
}

func (b *BaseBattle) CreateBuff(id int, src facade.Skill, val int) facade.Buff {
	return buff.NewBuff(id, src, val)
}

func (b *BaseBattle) CreateHalo(id int, skill facade.Skill) facade.Halo {
	h := halo.NewHalo(id)
	if h != nil {
		h.SetSkill(skill)
		h.SetUnit(skill.GetUnit())
	}

	return h
}
