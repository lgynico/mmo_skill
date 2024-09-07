package facade

import "github.com/lgynico/mmo_skill/conf"

/* 光环处理器 */
type HaloHandler interface {
	OnEffectHalo()
}

/* 光环 */
type Halo interface {
	HaloHandler
	SkillRelative
	Update(dt int64)
	GetConfig() *conf.SkillHaloRowEx
	SetHandler(HaloHandler)
}

type BaseHalo struct {
	HaloHandler
	SkillRelative
	timer  int64
	config *conf.SkillHaloRowEx
}

func NewHalo(id int) Halo {
	config, ok := conf.ConfMgr.SkillConfAdapter.HolaRows[id]
	if !ok {
		return nil
	}

	halo := &BaseHalo{
		SkillRelative: NewSkillRelative(),
		config:        config,
		timer:         0,
	}

	halo.SetHandler(halo)
	return halo
}

func (h *BaseHalo) Update(dt int64) {
	h.timer += dt
	if h.timer >= int64(h.config.Cycle) {
		h.HaloHandler.OnEffectHalo()
		h.timer -= int64(h.config.Cycle)
	}
}

func (h *BaseHalo) GetConfig() *conf.SkillHaloRowEx {
	return h.config
}

func (h *BaseHalo) SetHandler(handler HaloHandler) {
	h.HaloHandler = handler
}

func (h *BaseHalo) OnEffectHalo() {
}

/* 光环管理器 */
type HaloMgr interface {
	AddHalo(id int, halo Halo)
	RemoveHalo(id int)
	Update(dt int64)
}

type BaseHaloMgr struct {
	mapId2Halo map[int]Halo
}

func NewHaloMgr() HaloMgr {
	return &BaseHaloMgr{
		mapId2Halo: make(map[int]Halo),
	}
}

func (m *BaseHaloMgr) AddHalo(id int, halo Halo) {
	m.mapId2Halo[id] = halo
}

func (m *BaseHaloMgr) RemoveHalo(id int) {
	delete(m.mapId2Halo, id)
}

func (m *BaseHaloMgr) Update(dt int64) {
	for _, halo := range m.mapId2Halo {
		halo.Update(dt)
	}
}
