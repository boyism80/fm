-- Skill name (String.wz/Skill.img.xml): 헤이스트

local combat = require("script/lib/combat")

return {
    on_activated = function(me, skill, params)
        combat.for_each_near_party_member(me, skill, function(ch)
            ch:buff(skill)
            if ch ~= me then
                ch:show_skill_effect(skill, SkillEffectType.Affected)
            end
        end)
    end
}
