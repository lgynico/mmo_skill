package hero

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/skill"
)

type Hero struct {
	facade.BattleUnit
	config        *conf.HeroEntryRowEx
	stateModifier facade.StateModifier
}

func NewHero(id int32, playerId int64, cfgId int, lvl int) *Hero {
	hero := &Hero{
		BattleUnit: facade.NewBattleUnit(id, playerId, consts.GAMEOBJECT_UNIT, cfgId, lvl),
		config:     conf.ConfMgr.HeroConfAdapter.HeroRows[cfgId],
	}

	hero.SetFSM(facade.NewFSM(hero, InitState))
	hero.SetBuffMgr(facade.NewBuffMgr(hero))

	// attrs := facade.NewAttr(hero.config.Attrs)
	// hero.SetAttrs(attrs)

	skillConfig := hero.config.GetSkillCfgs(lvl)
	skillMgr := skill.NewSkillMgr(skillConfig, hero)
	//skillMgr.SetUnit(hero)
	hero.SetSkillMgr(skillMgr)

	// hero.SetHp(hero.GetAttrs().Get(facade.MaxHp))
	//hero.SetHeading(geom.NewVector2d(1, 0))
	// hero.SetBoundingShape(geom.NewCircle(float64(hero.config.Radius)))

	hero.SetImpl(hero)
	return hero
}
