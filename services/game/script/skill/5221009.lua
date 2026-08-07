-- Skill name (String.wz/Skill.img.xml): 마인드 컨트롤

local function mind_control_bounds(me, effect)
	local rng = tonumber(effect.range) or 350
	local w = 200 + rng
	local h = 100 + rng
	local px, py = me:position()
	return px - w, px + w, py - h, py
end

return {
	on_activated = function(me, skill, params)
		local map = me:map()
		if map == nil then
			return
		end
		local effect = skill:effect()
		if effect == nil then
			return
		end
		local duration_ms = effect.time or 0
		if duration_ms <= 0 then
			return
		end
		local prop = tonumber(effect.prop) or 0
		if prop <= 0 then
			return
		end
		local mob_max = tonumber(effect.mob_count) or 0
		if mob_max <= 0 then
			mob_max = 1
		end
		local min_x, max_x, min_y, max_y = mind_control_bounds(me, effect)
		local mobs = map:objects(ObjectType.Mob, {
			area = { minX = min_x, minY = min_y, maxX = max_x, maxY = max_y },
		})
		local checked = 0
		for _, mob in ipairs(mobs) do
			if checked >= mob_max then
				break
			end
			checked = checked + 1
			if math.random(1, 100) <= prop then
				local mwz = mob:wz()
				if mwz == nil or not mwz:boss() then
					mob:buff(MobBuff.Hypnotize, 1, duration_ms, skill, me)
				end
			end
		end
	end
}
