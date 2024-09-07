package conf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lgynico/mmo_skill/consts"
)

func init() {
	ConfMgr.addAdapter("skill", &SkillConfAdapter{})
	ConfMgr.addAdapter("hero", &HeroConfAdapter{})
}

/* 效果装饰器类型 */
const (
	EFF_DECO_NONE int = iota
	EFF_DECO_REPEAT
	EFF_DECO_WHILE
	EFF_DECO_WAIT
	EFF_DECO_IF
	EFF_DECO_PROB
	EFF_DECO_REPEAT_WITH_DELAY
)

type EffectDeco struct {
	Type   int
	Params []interface{}
}

func NewEffectDeco(t int, params ...interface{}) *EffectDeco {
	return &EffectDeco{
		Type:   t,
		Params: params,
	}
}

type SkillCondition struct {
	IsOppoTarget bool
	AttrName     string
	Comparator   string
	Value        int
	IsPerc       bool
}

type SkillEffectRowEx struct {
	*SkillEffectRow
	Event int
	Deco  *EffectDeco
}

type SkillLevelRowEx struct {
	*SkillLevelRow
	Effects    []*SkillEffectRowEx
	PasEffects []*SkillEffectRowEx
}

type SkillEntryRowEx struct {
	*SkillEntryRow
	Lv1Row *SkillLevelRowEx
	Lv2Row *SkillLevelRowEx
	Lv3Row *SkillLevelRowEx
	Lv4Row *SkillLevelRowEx
}

type SkillTrapRowEx struct {
	*SkillTrapRow
	BuffRows []*SkillBuffRow
}

type SkillProjectileRowEx struct {
	*SkillProjectileRow
	Effects []*SkillEffectRowEx
}

// type SkillBuffRowEx struct{ *SkillBuffRow }

type SkillHaloRowEx struct {
	*SkillHaloRow
	BuffRows []*SkillBuffRow
}

/********************************
* Adapter
********************************/
type SkillConfAdapter struct {
	BaseConfAdapter
	LevelRows      map[int]*SkillLevelRowEx
	EntryRows      map[int]*SkillEntryRowEx
	TrapRows       map[int]*SkillTrapRowEx
	ProjectileRows map[int]*SkillProjectileRowEx
	HolaRows       map[int]*SkillHaloRowEx
}

// var SkillConfAdapter = &SkillConfAdapter{
// 	EntryRows:      make(map[int]*SkillEntryRowEx),
// 	LevelRows:      make(map[int]*SkillLevelRowEx),
// 	TrapRows:       make(map[int]*SkillTrapRowEx),
// 	ProjectileRows: make(map[int]*SkillProjectileRowEx),
// 	HolaRows:       make(map[int]*SkillHaloRowEx),
// }

func (a *SkillConfAdapter) onLoadComplete() {
	a.adaptTrapRows()
	a.adaptLevelRows()
	a.adaptEntryRows()
	a.adaptProjectileRows()
	a.adaptHaloRows()
}

func (a *SkillConfAdapter) adaptProjectileRows() {
	a.ProjectileRows = make(map[int]*SkillProjectileRowEx, len(ConfMgr.SkillProjectile.Rows))

	for _, row := range ConfMgr.SkillProjectile.Rows {
		rowEx := &SkillProjectileRowEx{
			SkillProjectileRow: row,
			Effects:            make([]*SkillEffectRowEx, 0),
		}

		if eff := a.wrapEffect(row.Effect1, row.EffDeco1, 0); eff != nil {
			rowEx.Effects = append(rowEx.Effects, eff)
		}
		if eff := a.wrapEffect(row.Effect2, row.EffDeco2, 0); eff != nil {
			rowEx.Effects = append(rowEx.Effects, eff)
		}
		if eff := a.wrapEffect(row.Effect3, row.EffDeco3, 0); eff != nil {
			rowEx.Effects = append(rowEx.Effects, eff)
		}

		a.ProjectileRows[row.Id] = rowEx
	}
}

func (a *SkillConfAdapter) adaptTrapRows() {
	a.TrapRows = make(map[int]*SkillTrapRowEx, len(ConfMgr.SkillTrap.Rows))

	for _, trapRow := range ConfMgr.SkillTrap.Rows {
		buffRows := make([]*SkillBuffRow, 0, len(trapRow.Buffs))
		for _, buffId := range trapRow.Buffs {
			buffRow, ok := ConfMgr.SkillBuff.GetRowByInt(buffId)
			if !ok {
				panic(fmt.Sprintf("config not exists: %s %d", CONF_SKILL_BUFF, buffId))
			}
			buffRows = append(buffRows, buffRow)
		}
		rowEx := &SkillTrapRowEx{
			SkillTrapRow: trapRow,
			BuffRows:     buffRows,
		}
		a.TrapRows[trapRow.Id] = rowEx
	}
}

func (a *SkillConfAdapter) adaptLevelRows() {
	a.LevelRows = make(map[int]*SkillLevelRowEx, len(ConfMgr.SkillLevel.Rows))

	for _, row := range ConfMgr.SkillLevel.Rows {
		rowEx := &SkillLevelRowEx{
			SkillLevelRow: row,
			Effects:       make([]*SkillEffectRowEx, 0),
			PasEffects:    make([]*SkillEffectRowEx, 0),
		}

		if eff := a.wrapEffect(row.Effect1, row.EffDeco1, 0); eff != nil {
			rowEx.Effects = append(rowEx.Effects, eff)
		}
		if eff := a.wrapEffect(row.Effect2, row.EffDeco2, 0); eff != nil {
			rowEx.Effects = append(rowEx.Effects, eff)
		}
		if eff := a.wrapEffect(row.Effect3, row.EffDeco3, 0); eff != nil {
			rowEx.Effects = append(rowEx.Effects, eff)
		}

		if eff := a.wrapEffect(row.PasEffect1, row.PasEffDeco1, row.PasEvent1); eff != nil {
			rowEx.PasEffects = append(rowEx.PasEffects, eff)
		}
		if eff := a.wrapEffect(row.PasEffect2, row.PasEffDeco2, row.PasEvent2); eff != nil {
			rowEx.PasEffects = append(rowEx.PasEffects, eff)
		}
		if eff := a.wrapEffect(row.PasEffect3, row.PasEffDeco3, row.PasEvent3); eff != nil {
			rowEx.PasEffects = append(rowEx.PasEffects, eff)
		}

		a.LevelRows[row.Id] = rowEx
	}
}

func (a *SkillConfAdapter) adaptEntryRows() {
	a.EntryRows = make(map[int]*SkillEntryRowEx, len(ConfMgr.SkillEntry.Rows))

	for _, row := range ConfMgr.SkillEntry.Rows {
		rowEx := &SkillEntryRowEx{
			SkillEntryRow: row,
		}

		rowEx.Lv1Row = a.getLvRowEx(rowEx.Lv1)
		rowEx.Lv2Row = a.getLvRowEx(rowEx.Lv2)
		rowEx.Lv3Row = a.getLvRowEx(rowEx.Lv3)
		rowEx.Lv4Row = a.getLvRowEx(rowEx.Lv4)

		a.EntryRows[row.Id] = rowEx
	}
}

func (a *SkillConfAdapter) wrapEffect(effId int, effDeco string, event int) *SkillEffectRowEx {
	if effId == 0 {
		return nil
	}

	effRow, ok := ConfMgr.SkillEffect.GetRowByInt(effId)
	if !ok {
		panic(fmt.Sprintf("config not exists: name=%s, row=%d", CONF_SKILL_EFFECT, effId))
	}

	deco := a.parseDeco(effDeco)

	return &SkillEffectRowEx{
		SkillEffectRow: effRow,
		Event:          event,
		Deco:           deco,
	}
}

func (a *SkillConfAdapter) parseDeco(text string) *EffectDeco {
	if len(text) == 0 {
		return NewEffectDeco(EFF_DECO_NONE)
	}

	strs := strings.Fields(text)
	switch strs[0] {
	case "Repeat":
		n, err := strconv.Atoi(strs[1])
		if err != nil {
			panic(err)
		}
		if len(strs) > 2 {
			if strs[2] != "WithDelay" {
				panic(fmt.Sprintf("expect 'WithDelay' but met '%s'", strs[2]))
			}
			delay, err := strconv.Atoi(strs[3])
			if err != nil {
				panic(err)
			}
			return NewEffectDeco(EFF_DECO_REPEAT_WITH_DELAY, n, delay)
		}
		return NewEffectDeco(EFF_DECO_REPEAT, n)

	case "While":
		if len(strs) < 4 {
			panic("while condition parameter not enough")
		}
		cond := a.parseCondition(strs[1], strs[2], strs[3])
		return NewEffectDeco(EFF_DECO_WHILE, cond)
	case "Wait":
		n, err := strconv.Atoi(strs[1])
		if err != nil {
			panic(err)
		}
		return NewEffectDeco(EFF_DECO_WAIT, n)
	case "If":
		if len(strs) < 4 {
			panic("if condition parameter not enough")
		}
		cond := a.parseCondition(strs[1], strs[2], strs[3])
		return NewEffectDeco(EFF_DECO_IF, cond)
	case "Prob":
		f, err := strconv.ParseFloat(strs[1], 32)
		if err != nil {
			panic(err)
		}
		return NewEffectDeco(EFF_DECO_PROB, float32(f))
	}

	panic("unexpected deco")
}

func (a *SkillConfAdapter) parseCondition(attr, comp, valStr string) *SkillCondition {
	condition := &SkillCondition{}

	if strings.HasPrefix(attr, "^") {
		attr = strings.ToLower(attr[1:])
		condition.IsOppoTarget = true
	}
	condition.AttrName = attr

	if comp != ">" && comp != "<" && comp != "=" && comp != ">=" && comp != "<=" && comp != "!=" {
		panic("unkown comparator: " + comp)
	}
	condition.Comparator = comp

	if strings.HasSuffix(valStr, "%") {
		valStr = valStr[:len(valStr)-1]
		condition.IsPerc = true
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		panic(err)
	}
	condition.Value = val

	return condition
}

func (a *SkillConfAdapter) getLvRowEx(id int) *SkillLevelRowEx {
	if id <= 0 {
		return nil
	}

	lvRow, ok := a.LevelRows[id]
	if !ok {
		panic(fmt.Errorf("config not exists: name = %s, row = %d", CONF_SKILL_LEVEL, id))
	}

	return lvRow
}

func (a *SkillConfAdapter) adaptHaloRows() {
	a.HolaRows = make(map[int]*SkillHaloRowEx, len(ConfMgr.SkillHalo.Rows))

	for _, row := range ConfMgr.SkillHalo.Rows {
		rowEx := &SkillHaloRowEx{
			SkillHaloRow: row,
			BuffRows:     make([]*SkillBuffRow, 0),
		}

		buffRow, ok := ConfMgr.SkillBuff.GetRowByInt(row.Buff)
		if !ok {
			panic(fmt.Errorf("config not exists: name = %s, row = %d", CONF_SKILL_BUFF, row.Buff))
		}

		if buffRow.Type == int(consts.BUFF_HOLA) {
			panic(fmt.Errorf("halo contains halo buff: id = %d, buff = %d", row.Id, row.Buff))
		}

		rowEx.BuffRows = append(rowEx.BuffRows, buffRow)
		a.HolaRows[row.Id] = rowEx
	}
}
