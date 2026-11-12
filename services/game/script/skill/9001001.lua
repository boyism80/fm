-- Skill name (String.wz/Skill.img.xml): 헤이스트

local combat = require("script/lib/combat")

return {
    on_activated = function(me, skill, params)
        combat.for_each_character_in_skill_area(me, skill, function(ch)
            if ch == nil or not ch:is_alive() then
                return
            end
            ch:buff(skill)
            if ch ~= me then
                ch:show_skill_effect(skill, SkillEffectType.Affected)
            end
        end)
    end
}
