-- Skill name (String.wz/Skill.img.xml): 메소 업

function on_activated_4111001(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local percent = effect.x or 100
	for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, BuffFlag.MesoUp, percent)
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end

function on_buff_4111001(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local percent = effect.x or 100
	local current = me:bonus_meso_multiplier()
	if current <= 0 then
		current = 100
	end
	me:bonus_meso_multiplier(current + (percent - 100))
end

function on_unbuff_4111001(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local percent = effect.x or 100
	local current = me:bonus_meso_multiplier()
	me:bonus_meso_multiplier(current - (percent - 100))
end
