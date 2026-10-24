-- Skill name (String.wz/Skill.img.xml): 프리져

local util = require("script/lib/skill")

return {
    on_attack = function(me, skill, damages)
        util.summon_status_attack(me, skill, damages, MobBuff.Freeze)
    end,

    on_activated = function(me, skill, params)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        if effect.time <= 0 then
            return
        end
        me:buff(skill, BuffFlag.Summon, 1)
    end,

    on_buff = function(me, skill)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        if effect.time <= 0 then
            return
        end
        local wz = skill:wz()
        local level = skill:level()
        if level == nil or level <= 0 then
            return
        end
        me:create_summon(wz:id(), level, effect.time, SummonMovementType.CircleFollow, SummonType.Normal)
    end,

    on_unbuff = function(me, skill)
        local wz = skill:wz()
        me:remove_summon(wz:id())
    end
}
