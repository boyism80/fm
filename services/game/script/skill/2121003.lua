-- Skill name (String.wz/Skill.img.xml): 파이어 데몬

local combat = require("script/lib/combat")

return {
    on_attack = function(me, skill, damages)
        combat.weaken_damaged(skill, damages, "f")
        combat.apply_prob_status(me, skill, damages, MobBuff.Poison)
    end
}
