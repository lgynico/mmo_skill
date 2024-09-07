package proto

type (
	Vector struct {
		X int64 `json:"x,omitempty"`
		Y int64 `json:"y,omitempty"`
		Z int64 `json:"z,omitempty"`
	}

	Attr struct {
		Key   int32 `json:"key,omitempty"`
		Value int64 `json:"value,omitempty"`
	}

	BattleAction struct {
		Time      int64   `json:"time,omitempty"`      // 时间
		Unit      int32   `json:"unit,omitempty"`      // 单位（ 队伍+位置，如：11 表示 红色方1号位， 12 表示 蓝色方2号位 ）
		Key       int32   `json:"key,omitempty"`       // 动作类型
		Value     int64   `json:"value,omitempty"`     // 动作对应的id
		Target    int32   `json:"target,omitempty"`    // 动作目标
		UnitPos   *Vector `json:"unitPos,omitempty"`   // 单位坐标
		TargetPos *Vector `json:"targetPos,omitempty"` // 目标坐标
	}

	BattleTeam struct {
		Units []*BattleUnit `json:"units,omitempty"`
	}

	BattleUnit struct {
		ID    int32   `json:"id,omitempty"`    //
		Attrs []*Attr `json:"attrs,omitempty"` // 属性
	}

	BattleRecord struct {
		TimeElapsed int64           `json:"timeElapsed,omitempty"` // 时长
		WhichWin    int32           `json:"whichWin,omitempty"`    // 哪边赢（ 0=平手，1=红方，2=蓝方 ）
		Teams       []*BattleTeam   `json:"teams,omitempty"`       // 队伍（0=红方, 1=蓝方）
		Stats       []*BattleStat   `json:"stats,omitempty"`       // 战斗统计
		Actions     []*BattleAction `json:"actions,omitempty"`     // 动作
	}

	BattleStat struct {
		ID        int32 `json:"id,omitempty"`        // 单位id
		Score     int32 `json:"score,omitempty"`     // 评分
		Damage    int32 `json:"damage,omitempty"`    // 伤害
		Healing   int32 `json:"healing,omitempty"`   // 治疗量
		Suffering int32 `json:"suffering,omitempty"` // 承伤量
	}
)
