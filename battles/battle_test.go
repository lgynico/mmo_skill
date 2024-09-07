package battles

import (
	"testing"

	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/gameobj/hero"
)

func TestMain(m *testing.M) {
	loadConfigs()
	m.Run()
}

func loadConfigs() {
	conf.ConfMgr.Init("../deploy/")
}

func Test_Battle(t *testing.T) {
	teamRed := [5]facade.BattleUnit{}
	teamRed[0] = unitFromMonster(112021240, 1)
	// teamRed[1] = unitFromMonster(112021240, 2)
	teamRed[2] = unitFromMonster(112011240, 3)
	// teamRed[3] = unitFromMonster(1, 4)
	// teamRed[4] = unitFromMonster(112011240, 5)

	teamBlue := [5]facade.BattleUnit{}
	teamBlue[0] = unitFromMonster(112021240, 1)
	// teamBlue[1] = unitFromMonster(112021240, 2)
	teamBlue[2] = unitFromMonster(112011240, 3)
	// teamBlue[3] = unitFromMonster(1, 4)
	// teamBlue[4] = unitFromMonster(112011240, 5)

	b := NewBaseBattle(teamRed, teamBlue)
	b.Init()

	b.Start()
}

func Benchmark_Battle(b *testing.B) {
	teamRed := [5]facade.BattleUnit{}
	teamRed[0] = unitFromMonster(112021240, 1)
	teamRed[1] = unitFromMonster(112021240, 2)
	teamRed[2] = unitFromMonster(112011240, 3)
	// teamRed[3] = unitFromMonster(1, 4)
	teamRed[4] = unitFromMonster(112011240, 5)

	teamBlue := [5]facade.BattleUnit{}
	teamBlue[0] = unitFromMonster(112021240, 1)
	teamBlue[1] = unitFromMonster(112021240, 2)
	teamBlue[2] = unitFromMonster(112011240, 3)
	// teamBlue[3] = unitFromMonster(1, 4)
	teamBlue[4] = unitFromMonster(112011240, 5)

	for i := 0; i < b.N; i++ {
		facade := NewBaseBattle(teamRed, teamBlue)
		facade.Init()
		facade.Start()
	}
}

func unitFromMonster(id int, index int) facade.BattleUnit {
	monsterCfg, _ := conf.ConfMgr.MonsterEntry.GetRowByInt(id)
	// h := entry.NewHeroWithMonster(monsterCfg)
	unit := hero.NewHero(int32(20+index), 0, monsterCfg.HeroId, monsterCfg.Lv)
	attrs := getAttrs(monsterCfg.HeroId)
	unit.SetAttrs(attrs)
	return unit
}

func getAttrs(heroID int) facade.Attrs {
	var (
		hero, _ = conf.ConfMgr.HeroEntry.GetRowByInt(heroID)
		attrs   = map[int]int{}
	)

	for k, v := range hero.Attrs {
		attrs[k] = int(v)
	}

	return facade.NewAttr(attrs)
}

// func NewHeroWithMonster(monsterCfg *conf.MonsterEntryRow) *Hero {
// 	h := &chproto.Hero{
// 		Info: &chproto.HeroInfo{
// 			CfgId:     uint32(monsterCfg.HeroId),
// 			Lvl:       uint32(monsterCfg.Lv),
// 			Rank:      uint32(monsterCfg.Rank),
// 			Attrs:     make(map[uint32]*chproto.SimpleAttr),
// 			Equipping: make(map[uint32]uint32),
// 		},
// 		DetailAttrs: make(map[uint32]*chproto.Attr),
// 		GmAttrs:     make(map[uint32]uint32),
// 	}

// 	hero := NewHero()
// 	hero.Init(h)
// 	return hero
// }
