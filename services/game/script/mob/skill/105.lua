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
        local is_range = mob:skill_area(skill) ~= nil
        local targets = {}
        if is_range then
            local other = mob:skill_targets(skill, ObjectType.Mob)[1]
            if other ~= nil then
                table.insert(targets, other)
            end
        else
            table.insert(targets, mob)
        end

        local map = mob:map()
        for _, target in ipairs(targets) do
            if is_range then
                map:remove_mob(target:oid(), MobDieAnimation.FadeOut)
                heal_mob(mob, effect)
            else
                heal_mob(target, effect)
            end
        end

        return true
    end
}
