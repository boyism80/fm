-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 131

return {
            left = mob_x + bounds.left,
            top = mob_y + bounds.top,
            right = mob_x + bounds.right,
            bottom = mob_y + bounds.bottom,
        }
    else
        return {
            left = mob_x - bounds.right,
            top = mob_y + bounds.top,
            right = mob_x - bounds.left,
            bottom = mob_y + bounds.bottom,
        }
    end
end

return {
    on_mob_skill_choose = function(mob, skill)
        return true
    end,

    on_mob_skill = function(mob, controller, skill)
        local effect = skill:effect()
        mob:create_mist(skill, effect.x * 10, MistType.Poison, mob:skill_area(skill) or effect.bounds)
        return true
    end
}
