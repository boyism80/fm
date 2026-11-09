-- Skill name (String.wz/Skill.img.xml): 디스오더

local combat = require("script/lib/combat")

return {
    on_attack = function(me, skill, damages)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        if effect.time <= 0 then
            return
        end
        for mob, hits in pairs(damages) do
            if combat.damaged(hits) then
                mob:buff({ [MobBuff.Watk] = effect.x, [MobBuff.Wdef] = effect.y }, effect.time, skill, me)
            end
        end
    end
}
