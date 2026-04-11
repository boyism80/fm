-- Skill name (String.wz/Skill.img.xml): 파이어 커튼

function on_activated_12111005(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local lt = effect.lt
	local rb = effect.rb
	if lt == nil or rb == nil then
		return
	end
	local x, y = me:position()
	local left = x + math.min(lt.x, rb.x)
	local right = x + math.max(lt.x, rb.x)
	local top = y + math.min(lt.y, rb.y)
	local bottom = y + math.max(lt.y, rb.y)
	if effect.time <= 0 then
		return
	end
	local multiplier = compute_poison_tick_multiplier(me, skill)
	me:create_mist(skill, effect.time, MistType.Poison, { left = left, top = top, right = right, bottom = bottom }, 2000, multiplier)
end
