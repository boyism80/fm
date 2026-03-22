-- Skill name (String.wz/Skill.img.xml): 위협

function on_activated_1201006(me, skill, params)
	local map = me:map()
	if map == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	local x_val = effect.x or 0
	local y_val = effect.y or 0
	if duration_ms <= 0 or (x_val == 0 and y_val == 0) then
		return
	end
	local mob_count = effect.mob_count or 6
	local posX, posY = me:position()
	local lt = effect.lt
	local rb = effect.rb
	if lt == nil or rb == nil then
		return
	end
	local minX = posX + math.min(lt.x, rb.x)
	local maxX = posX + math.max(lt.x, rb.x)
	local minY = posY + math.min(lt.y, rb.y)
	local maxY = posY + math.max(lt.y, rb.y)
	local mobs = map:objects(ObjectType.Mob, { area = { minX = minX, minY = minY, maxX = maxX, maxY = maxY } })
	local n = 0
	for _, mob in ipairs(mobs) do
		if n >= mob_count then
			break
		end
		local status_values = {}
		if x_val ~= 0 then
			status_values[MobStatus.Watk] = x_val
		end
		if y_val ~= 0 then
			status_values[MobStatus.Wdef] = y_val
		end
		if next(status_values) ~= nil then
			mob:set_status(status_values, duration_ms, skill, me)
		end
		n = n + 1
	end
end
