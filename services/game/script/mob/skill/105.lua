-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 105

local function heal_mob(target, effect)
    if effect.x ~= 0 then
        target:add_hp(effect.x)
    end
    if effect.y ~= 0 then
        target:add_mp(effect.y)
    end
end

return {
    on_mob_skill_choose = function(mob, skill)
        return true
    end,

    on_mob_skill = function(mob, controller, skill)
        local effect = skill:effect()
        if mob:skill_area(skill) == nil then
            heal_mob(mob, effect)
            return true
        end

        for _, other in ipairs(mob:skill_mobs(skill)) do
            if other:oid() ~= mob:oid() then
                mob:map():remove_mob(other:oid(), MobDieAnimation.FadeOut)
                heal_mob(mob, effect)
                break
            end
        end
        return true
    end
}
