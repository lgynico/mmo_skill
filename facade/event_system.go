package facade

type (
	GameEvent struct {
		Target GameObject
		Params []any
	}

	EventLinstener interface {
		OnEventTrigger(eventType EventType, event *GameEvent)
	}

	EventSystem struct {
		listeners map[EventType][]EventLinstener
	}

	EventHandler interface {
		OnImmediately(*GameEvent) bool
		OnEnterScene(*GameEvent) bool
		OnLeaveScene(*GameEvent) bool
		OnTakeDamage(*GameEvent) bool
		OnGiveDamage(*GameEvent) bool
		OnTakeHealing(*GameEvent) bool
		OnGiveHealing(*GameEvent) bool
		OnDead(*GameEvent) bool
		OnReborn(*GameEvent) bool
		OnAddBuff(*GameEvent) bool
		OnKillTarget(*GameEvent) bool
		OnOverHeal(*GameEvent) bool
		OnGiveHealFront(*GameEvent) bool
	}
)

func NewGameEvent(target GameObject, params ...any) *GameEvent {
	return &GameEvent{
		Target: target,
		Params: params,
	}
}

func NewEventSystem() *EventSystem {
	return &EventSystem{
		listeners: map[EventType][]EventLinstener{},
	}
}

func (p *EventSystem) SendEvent(eventType EventType, evt *GameEvent) {
	if listeners, ok := p.listeners[eventType]; ok {
		for _, listener := range listeners {
			listener.OnEventTrigger(eventType, evt)
		}
	}
}

func (p *EventSystem) AddListener(eventType EventType, listener EventLinstener) {
	listeners, ok := p.listeners[eventType]
	if !ok {
		listeners = make([]EventLinstener, 0, 8)
		p.listeners[eventType] = listeners
	}

	for _, l := range listeners {
		if l == listener {
			return
		}
	}

	listeners = append(listeners, listener)
	p.listeners[eventType] = listeners
}

func (p *EventSystem) RemoveListener(eventType EventType, listener EventLinstener) {
	listeners, ok := p.listeners[eventType]
	if !ok {
		return
	}

	for i, l := range listeners {
		if l == listener {
			p.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
		}
	}
}

type EventType int32

const (
	EVENT_IMMEDIATELY     EventType = iota // 立即生效0
	EVENT_ENTER_SCENE                      // 进入场景1
	EVENT_LEAVE_SCENE                      // 离开场景2
	EVENT_GIVE_DAMAGE                      // 造成伤害3
	EVENT_TAKE_DAMAGE                      // 受到伤害
	EVENT_GIVE_HEAL                        // 造成治疗
	EVENT_TAKE_HEAL                        // 受到治疗
	EVENT_DEAD                             // 死亡
	EVENT_REBORN                           // 复活
	EVENT_ADD_BUFF                         // 加buff
	EVENT_KILL_TARGET                      // 击杀目标
	EVENT_OVER_HEAL                        // 过量治疗
	EVENT_GIVE_HEAL_FRONT                  // 造成治疗前
)

var mapId2EventName = map[EventType]string{
	EVENT_IMMEDIATELY:     "immediately",
	EVENT_ENTER_SCENE:     "enterScene",
	EVENT_LEAVE_SCENE:     "leaveScene",
	EVENT_GIVE_DAMAGE:     "giveDamage",
	EVENT_TAKE_DAMAGE:     "takeDamage",
	EVENT_GIVE_HEAL:       "giveHeal",
	EVENT_TAKE_HEAL:       "takeHeal",
	EVENT_DEAD:            "dead",
	EVENT_REBORN:          "reborn",
	EVENT_ADD_BUFF:        "addBuff",
	EVENT_KILL_TARGET:     "killTarget",
	EVENT_OVER_HEAL:       "overHeal",
	EVENT_GIVE_HEAL_FRONT: "giveHealfront",
}

var mapEventName2Id = map[string]EventType{
	"immediately":   EVENT_IMMEDIATELY,
	"enterScene":    EVENT_ENTER_SCENE,
	"leaveScene":    EVENT_LEAVE_SCENE,
	"giveDamage":    EVENT_GIVE_DAMAGE,
	"takeDamage":    EVENT_TAKE_DAMAGE,
	"giveHeal":      EVENT_GIVE_HEAL,
	"takeHeal":      EVENT_TAKE_HEAL,
	"dead":          EVENT_DEAD,
	"reborn":        EVENT_REBORN,
	"addBuff":       EVENT_ADD_BUFF,
	"killTarget":    EVENT_KILL_TARGET,
	"overHeal":      EVENT_OVER_HEAL,
	"giveHealfront": EVENT_GIVE_HEAL_FRONT,
}

func GetEventNameMap() map[EventType]string {
	return mapId2EventName
}

func GetEventName(eventType EventType) (string, bool) {
	name, ok := mapId2EventName[eventType]
	return name, ok
}

func GetTypeByEventName(name string) (EventType, bool) {
	id, ok := mapEventName2Id[name]
	return id, ok
}

func SendEvent(eventType EventType, target GameObject, params ...any) {
	if eventSys := target.GetEventSystem(); eventSys != nil {
		event := NewGameEvent(target, params...)
		eventSys.SendEvent(eventType, event)
	}
}
