-- Skill name (String.wz/Skill.img.xml): 위협

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
		if effect.time <= 0 or (effect.x == 0 and effect.y == 0) then
			return
		end
		local mob_count = effect.mob_count or 6
		local pos_x, pos_y = me:position()
		local lt = effect.lt
		local rb = effect.rb
		if lt == nil or rb == nil then
			return
		end
		local min_x = pos_x + math.min(lt.x, rb.x)
		local max_x = pos_x + math.max(lt.x, rb.x)
		local min_y = pos_y + math.min(lt.y, rb.y)
		local max_y = pos_y + math.max(lt.y, rb.y)
		local mobs = map:objects(ObjectType.Mob, { area = { minX = min_x, minY = min_y, maxX = max_x, maxY = max_y } })
		local n = 0
		for _, mob in ipairs(mobs) do
			if n >= mob_count then
				break
			end
			local buff_values = {}
			if effect.x ~= 0 then
				buff_values[MobBuff.Watk] = effect.x
			end
			if effect.y ~= 0 then
				buff_values[MobBuff.Wdef] = effect.y
			end
			if next(buff_values) ~= nil then
				mob:buff(buff_values, effect.time, skill, me)
			end
			n = n + 1
		end
	end
}
