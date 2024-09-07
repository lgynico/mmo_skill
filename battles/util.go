package battles

import (
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/proto"
)

func unitToProto(unit facade.BattleUnit) *proto.BattleUnit {
	bu := &proto.BattleUnit{
		ID:    unit.GetId(),
		Attrs: make([]*proto.Attr, 0),
	}

	unit.GetAttrs().Range(func(key, val int) {
		if val > 0 {
			pair := &proto.Attr{Key: int32(key), Value: int64(val)}
			bu.Attrs = append(bu.Attrs, pair)
		}
	})

	return bu
}
