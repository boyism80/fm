-- Skill name (String.wz/Skill.img.xml): 알케미스트

function on_passive_14110003(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.x <= 0 then
		return
	end
end

function on_unpassive_14110003(me, skill)
	me:potion_heal_rate(0)
	me:potion_duration_rate(0)
end

function on_activated_14110003(me, skill, params)
end
