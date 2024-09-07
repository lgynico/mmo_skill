package halo

import (
	"log"

	"github.com/lgynico/mmo_skill/facade"
)

type Halo struct {
	facade.Halo
	facade.HaloHandler
}

func NewHalo(id int) facade.Halo {
	halo := &Halo{
		Halo: facade.NewHalo(id),
	}

	halo.SetHandler(halo)
	return halo
}

func (h *Halo) OnEffectHalo() {
	log.Printf("[%.2f] 光环 %d(%s) 生效\n", float64(h.GetUnit().GetContext().GetTimeMillis())/1000, h.GetConfig().Id, h.GetConfig().Name)
	unit := h.GetUnit()
	ctx := unit.GetContext()
	units := ctx.GetUnits()
	config := h.GetConfig()
	radiusSq := float64(config.Radius * config.Radius)
	for _, buffCfg := range config.BuffRows {
		val := 0
		if buffCfg.IsOuterVal {
			val = h.GetSkill().GetValue()
		}

		for _, target := range units {
			if facade.IsTarget(target, unit, config.TargetType) &&
				unit.GetPosition().DistSqTo(target.GetPosition()) <= radiusSq {

				buff := ctx.CreateBuff(buffCfg.Id, h.GetSkill(), val)
				buff.SetValue(buff.CalValue(unit))
				target.GetBuffMgr().AddBuff(buff)
			}
		}
	}
}
