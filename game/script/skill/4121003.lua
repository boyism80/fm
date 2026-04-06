-- Skill name (String.wz/Skill.img.xml): 쇼다운

function on_activated_4121003(me, skill, params)
end

function on_mob_buff_4121003(mob, skill, causer)
	if mob == nil or skill == nil then
		return
	end
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local n = math.floor(tonumber(effect.x) or 0)
	if n <= 0 then
		return
	end
	mob:exp_rate(mob:exp_rate() + n)
	mob:drop_rate(mob:drop_rate() + n)
end

function on_mob_unbuff_4121003(mob, skill, causer)
	if mob == nil or skill == nil then
		return
	end
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local n = math.floor(tonumber(effect.x) or 0)
	if n <= 0 then
		return
	end
	mob:exp_rate(mob:exp_rate() - n)
	mob:drop_rate(mob:drop_rate() - n)
end

function on_attack_4121003(me, skill, damages)
	apply_showdown_on_attack(me, skill, damages)
end
