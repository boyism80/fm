-- Skill name (String.wz/Skill.img.xml): 파이어 커튼

local combat = require("script/lib/combat")

return {
    on_activated = function(me, skill, params)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        if effect.time <= 0 then
            return
        end
        local multiplier = combat.compute_poison_tick_multiplier(me, skill)
        me:create_mist(skill, effect.time, MistType.Poison, me:skill_area(skill), 2000, multiplier)
    end
}
