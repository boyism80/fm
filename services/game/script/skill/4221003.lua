-- Skill name (String.wz/Skill.img.xml): 쇼다운

local combat = require("script/lib/combat")

function on_activated_4221003(me, skill, params)
end

function on_mob_buff_4221003(mob, skill, causer)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.x <= 0 then
		return
	end
	mob:exp_rate(mob:exp_rate() + effect.x)
	mob:drop_rate(mob:drop_rate() + effect.x)
end

function on_mob_unbuff_4221003(mob, skill, causer)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.x <= 0 then
		return
	end
	mob:exp_rate(mob:exp_rate() - effect.x)
	mob:drop_rate(mob:drop_rate() - effect.x)
end

function on_attack_4221003(me, skill, damages)
	combat.apply_showdown(me, skill, damages)
end
