-- Skill name (String.wz/Skill.img.xml): 위협

function on_activated_1201006(me, skill, params)
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
