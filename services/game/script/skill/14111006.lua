-- Skill name (String.wz/Skill.img.xml): 포이즌 봄

local combat = require("script/lib/combat")

return {
    -- The bomb leaves a poison mist where it lands, whether or not it hit a mob.
    on_attack = function(me, skill, damages, info)
        combat.apply_venom(me, damages, Skill.VenomCygnus)
        local effect = skill:effect()
        if effect == nil or effect.time <= 0 then
            return
        end
        local multiplier = combat.compute_poison_tick_multiplier(me, skill)
        me:create_mist(skill, effect.time, MistType.Poison, me:skill_area(skill, info.x, info.y), 2000, multiplier)
    end
}
