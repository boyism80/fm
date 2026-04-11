-- Skill name (String.wz/Skill.img.xml): 코마

function on_attack_11111003(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_11111003(me, skill, params)
	-- Dawn Warrior Coma consumes combo orbs the same way as Hero Coma.
	consume_combo_orbs(me)
end
