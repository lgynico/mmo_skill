package effect

import (
	"log"

	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/proto"
	"github.com/lgynico/mmo_skill/utils"
)

// 创建子弹
type CreateProjectile struct{}

func NewCreateProjectile() facade.SkillEffect {
	return &CreateProjectile{}
}

func (e *CreateProjectile) OnEffect(skill facade.Skill, config *conf.SkillEffectRow, target facade.BattleUnit) {
	if projRow, ok := conf.ConfMgr.SkillConfAdapter.ProjectileRows[config.Param1]; ok {
		unit := skill.GetUnit()
		ctx := unit.GetContext()
		projectile := ctx.CreateProjectile(skill, projRow)
		if config.Param2 == 1 {
			enemies := common.FindAllEnemyHeros(unit, false)
			if len(enemies) > 0 {
				idx := utils.RandIntByCrypto(0, len(enemies))
				target := enemies[idx]
				projectile.SetObjTarget(target)
				projectile.SetPosTarget(target.GetPosition())
			} else {
				// 没有目标
				return
			}
		} else {
			projectile.SetObjTarget(target)
			projectile.SetPosTarget(target.GetPosition())
		}
		ctx.AddObj(projectile)

		// 去掉普攻创建的子弹
		if !skill.GetConfig().IsAtk {
			action := &proto.BattleAction{
				Time:      ctx.GetTimeMillis(),
				Unit:      unit.GetId(),
				Key:       int32(consts.ACTION_CREATE_PROJ),
				Value:     int64(projectile.GetCfgId()),
				Target:    projectile.GetObjTarget().GetId(),
				UnitPos:   facade.ConvertVector2d(projectile.GetPosition()),
				TargetPos: facade.ConvertVector2d(projectile.GetPosTarget()),
			}
			ctx.RecordAction(action)
		}

		return
	}

	log.Printf("config not exists: name=%s, row=%d\n", conf.CONF_SKILL_PROJECTILE, config.Param1)
}
