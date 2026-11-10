-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 128

return {
    on_mob_skill_choose = function(mob, skill)
        return true
    end,

    on_mob_skill = function(mob, controller, skill)
        local effect = skill:effect()
        local time = effect.time
        if time <= 0 then
            return true
        end

        local targets = mob:skill_characters(skill, controller)

        local x = effect.x
        if x <= 0 then
            x = 1
        end
        for _, ch in ipairs(targets) do
            ch:debuff(DebuffFlag.Seduce, time, x, skill:id(), skill:level())
        end
        return true
    end
}
