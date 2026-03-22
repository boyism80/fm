-- Skill name (String.wz/Skill.img.xml): 차지 블로우

function on_attack_1211002(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Stun)

	local adv_charge = me:skill(Skill.AdvancedCharge)
	local keep_chance = 0
	if adv_charge ~= nil then
		local adv_effect = get_skill_effect(adv_charge)
		if adv_effect ~= nil and adv_effect.x ~= nil then
			keep_chance = adv_effect.x
		end
	end
	if math.random(1, 100) > keep_chance then
		me:unbuff(BuffFlag.WkCharge)
	end
end

function on_activated_1211002(me, skill, params)
end
