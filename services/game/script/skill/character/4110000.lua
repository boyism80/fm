-- Skill name (String.wz/Skill.img.xml): 알케미스트

function on_passive_4110000(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.x <= 0 then
		return
	end
	me:potion_heal_rate(effect.x)
	me:potion_duration_rate(effect.x)
end

function on_unpassive_4110000(me, skill)
	me:potion_heal_rate(0)
	me:potion_duration_rate(0)
end
