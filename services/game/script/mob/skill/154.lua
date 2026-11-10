-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 154

return {
    on_mob_skill_choose = function(mob, skill)
        return true
    end,

    on_mob_skill = function(mob, controller, skill)
        local effect = skill:effect()
        local targets = mob:skill_mobs(skill)
        for _, target in ipairs(targets) do
            target:buff(MobBuff.Acc, effect.x, effect.time, skill, controller)
        end
        return true
    end
}
