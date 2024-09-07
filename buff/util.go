package buff

import (
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
)

func NewBuff(id int, skill facade.Skill, val int) facade.Buff {
	buff := facade.NewBuff(id, skill, val)
	setBuffHandler(buff)

	if len(buff.GetConfig().SubBuff) > 0 {
		subBuffs := make([]facade.Buff, 0, len(buff.GetConfig().SubBuff))
		for _, buffId := range buff.GetConfig().SubBuff {
			subBuff := facade.NewBuff(buffId, skill, val)
			setBuffHandler(subBuff)
			subBuffs = append(subBuffs, subBuff)
		}

		buff.SetSubBuffs(subBuffs)
	}
	// buff.SetValue(buff.CalValue(src))
	return buff
}

func setBuffHandler(buff facade.Buff) {
	switch consts.BuffType(buff.GetConfig().Type) {
	case consts.BUFF_ATTRIBUTE:
		buff.SetBuffHandler(&AttrBuff{facade.NewBuffHandler()})
	case consts.BUFF_DAMAGE:
		buff.SetBuffHandler(&DamageBuff{facade.NewBuffHandler()})
	case consts.BUFF_HEALING:
		buff.SetBuffHandler(&HealingBuff{facade.NewBuffHandler()})
	case consts.BUFF_STATE:
		buff.SetBuffHandler(NewStatusBuff(facade.NewBuffHandler()))
	case consts.BUFF_MOVE:
		buff.SetBuffHandler(&MoveBuff{facade.NewBuffHandler()})
	case consts.BUFF_SHIELD:
		buff.SetBuffHandler(&ShieldBuff{facade.NewBuffHandler()})
	case consts.BUFF_HOLA:
		buff.SetBuffHandler(&HaloBuff{facade.NewBuffHandler()})
	case consts.BUFF_ENERGY:
		buff.SetBuffHandler(&EnergyBuff{facade.NewBuffHandler()})
	default:
		buff.SetBuffHandler(&DefaultBuff{facade.NewBuffHandler()})
	}
}
