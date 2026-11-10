-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 129

local ex = require("script/lib/expedition")

return {
    -- An expedition battle keeps its members, so the banish is left out there.
    on_mob_skill_choose = function(mob, skill)
        return not ex.in_battle(mob:map())
    end,

    on_mob_skill = function(mob, controller, skill)
        for _, target in ipairs(mob:skill_characters(skill, controller)) do
            mob:banish(target)
        end

        return true
    end
}
