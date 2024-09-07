package buff

import "github.com/lgynico/mmo_skill/facade"

type HaloBuff struct {
	facade.BuffHandler
}

func (b *HaloBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	haloId := buff.GetConfig().Params[0]
	halo := unit.GetContext().CreateHalo(haloId, buff.GetSkill())
	if halo != nil {
		unit.GetHaloMgr().AddHalo(buff.GetId(), halo)
	}
}

func (b *HaloBuff) OnRemoveBuff(buff facade.Buff, unit facade.BattleUnit) {
	b.BuffHandler.OnRemoveBuff(buff, unit)
	unit.GetHaloMgr().RemoveHalo(buff.GetId())
}
