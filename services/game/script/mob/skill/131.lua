-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 131

return {
    on_mob_skill_choose = function(mob, skill)
        return true
    end,

    on_mob_skill = function(mob, controller, skill)
        local effect = skill:effect()
        mob:create_mist(skill, effect.x * 10, MistType.Poison)
        return true
    end
}
