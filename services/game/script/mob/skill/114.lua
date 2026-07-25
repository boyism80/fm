-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 114

local function random_range_heal_hp(effect)
    local hp_scale = effect.x / 1000
    local hp_roll = 950 + 1050 * math.random()
    return math.floor(hp_scale * hp_roll)
end

return {
	on_mob_skill_choose = function(mob, skill)
		return true
	end,

	on_mob_skill = function(mob, controller, skill)
		local effect = skill:effect()
		local bounds = effect.bounds
		local is_range = bounds.left ~= 0 or bounds.top ~= 0 or bounds.right ~= 0 or bounds.bottom ~= 0
		local targets = {}

		if is_range then
		    for _, target in ipairs(mob:objects_in(bounds, ObjectType.Mob)) do
		        table.insert(targets, target)
		    end
		    table.insert(targets, mob)
		else
		    table.insert(targets, mob)
		end

		local hp_amount = effect.x
		if is_range then
		    hp_amount = random_range_heal_hp(effect)
		end

		for _, target in ipairs(targets) do
		    if hp_amount ~= 0 then
		        target:add_hp(hp_amount)
		    end
		    if effect.y ~= 0 then
		        target:add_mp(effect.y)
		    end
		end

		return true
	end
}
