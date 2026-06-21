-- Skill name (String.wz/Skill.img.xml): 차지 블로우

local combat = require("script/lib/combat")

function on_attack_1211002(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)

	local adv_charge = me:skill(Skill.AdvancedCharge)
	local keep_chance = 0
	if adv_charge ~= nil then
		local adv_effect = adv_charge:effect()
		if adv_effect ~= nil and adv_effect.x ~= nil then
			keep_chance = adv_effect.x
		end
	end
	if math.random(1, 100) > keep_chance then
		me:unbuff(BuffFlag.WkCharge)
	end
end
