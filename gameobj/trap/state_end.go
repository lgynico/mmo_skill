package trap

import (
	"log"

	"github.com/lgynico/mmo_skill/facade"
)

var EndState facade.FSMState = nil

func init() {
	EndState = &endState{facade.NewFSMState()}
}

type endState struct {
	facade.FSMState
}

func (s *endState) OnEnter(obj facade.GameObject, dt int64) {
	// obj.GetContext().Leave(obj)
	log.Printf("[%.2f] 陷阱移除: %s\n", float64(obj.GetContext().GetTimeMillis())/1000, obj.(*Trap).GetConfig().Name)
}
