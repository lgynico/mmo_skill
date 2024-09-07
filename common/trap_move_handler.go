package common

import (
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/geom"
)

type TrapMoveHandler interface {
	HandleTrapMove(trap facade.Trap, ctx facade.Battle, dt int64)
}

var mapId2TrapMovehandler = map[int]TrapMoveHandler{
	TRAP_TRACK_UNMOVABLE:     &UnmovableTrapMoveHandler{},
	TRAP_TRACK_FOLLOW_UNIT:   &FollowUnitTrapMoveHandler{},
	TRAP_TRACK_FOLLOW_TARGET: &FollowTargetTrapMoveHandler{},
	TRAP_TRACK_MOVE_FORWARD:  &MoveForwardTrapMoveHandler{},
}

func GetTrapMoveHandler(id int) (TrapMoveHandler, bool) {
	handler, ok := mapId2TrapMovehandler[id]
	return handler, ok
}

/* 不移动 */
type UnmovableTrapMoveHandler struct {
}

func (h *UnmovableTrapMoveHandler) HandleTrapMove(trap facade.Trap, ctx facade.Battle, dt int64) {

}

/* 跟随英雄 */
type FollowUnitTrapMoveHandler struct {
}

func (h *FollowUnitTrapMoveHandler) HandleTrapMove(trap facade.Trap, ctx facade.Battle, dt int64) {
	pos := trap.GetUnit().GetPosition()
	trap.SetPosition(pos)
}

/* 跟随目标 */
type FollowTargetTrapMoveHandler struct {
}

func (h *FollowTargetTrapMoveHandler) HandleTrapMove(trap facade.Trap, ctx facade.Battle, dt int64) {

}

/* 向前移动 */
type MoveForwardTrapMoveHandler struct {
}

func (h *MoveForwardTrapMoveHandler) HandleTrapMove(trap facade.Trap, ctx facade.Battle, dt int64) {
	speed := float64(trap.GetConfig().MoveParams[0])
	dist := speed * float64(dt) / 1000
	v := trap.GetHeading().MulN(dist)

	mat := geom.MatrixTranslateV(v)
	pos := mat.MulV(trap.GetPosition())
	trap.SetPosition(pos)
}
