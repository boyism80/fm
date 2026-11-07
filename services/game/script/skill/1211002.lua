-- Skill name (String.wz/Skill.img.xml): 차지 블로우

local combat = require("script/lib/combat")

return {
    on_attack = function(me, skill, damages)
        combat.apply_prob_status(me, skill, damages, MobBuff.Stun)

        local adv_charge = me:skill(Skill.AdvancedCharge)
        local keep_chance = 0
        if adv_charge ~= nil then
            local adv_effect = adv_charge:effect()
            if adv_effect ~= nil and adv_effect.x ~= nil then
                keep_chance = adv_effect.x
            end
        end
        -- Each mob the blow damaged rolls for the charge; a blow that damaged none keeps it.
        for _, hits in pairs(damages or {}) do
            if combat.total_damage_to_mob(hits) > 0 and math.random(1, 100) > keep_chance then
                me:unbuff(BuffFlag.WkCharge)
                return
            end
        end
    end
}
