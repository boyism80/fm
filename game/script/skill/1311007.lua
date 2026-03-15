-- Skill name (String.wz/Skill.img.xml): 파워 크래쉬

function on_activated(me, skill, params)
	local map = me:map()
	if map == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local prop = effect.prop or 0
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
		if math.random(1, 100) <= prop then
			mob:clear_status(MobStatus.WeaponAttackUp)
		end
		n = n + 1
	end
end
