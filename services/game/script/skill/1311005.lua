-- Skill name (String.wz/Skill.img.xml): 새크리파이스

local combat = require("script/lib/combat")

return {
    on_attack = function(me, skill, damages)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        local total = 0
        local target = nil
        for mob, hits in pairs(damages) do
            target = mob
            total = combat.total_damage_to_mob(hits)
            break
        end
        if target == nil or total == 0 then
            return
        end
        local hp_loss = math.floor(total * effect.x / 100)
        if hp_loss <= 0 then
            return
        end
        local cur = me:hp()
        local new_hp = math.max(cur - hp_loss, 1)
        me:hp(new_hp, true)
    end
}
