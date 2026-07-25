-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 131

local function mob_facing_left(mob)
    return mob:stance() % 2 ~= 0
end

local function world_bounds_from_effect(mob, effect)
    local mob_x, mob_y = mob:position()
    local bounds = effect.bounds
    if mob_facing_left(mob) then
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
		local bounds = effect.bounds
		if bounds.left ~= 0 or bounds.top ~= 0 or bounds.right ~= 0 or bounds.bottom ~= 0 then
		    bounds = world_bounds_from_effect(mob, effect)
		end
		mob:create_mist(skill, effect.x * 10, MistType.Poison, bounds)
		return true
	end
}
