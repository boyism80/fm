-- Skill name (String.wz/Skill.img.xml): 코마 : 검

function on_attack_1111005(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Stun)
end

function on_activated_1111005(me, skill, params)
	-- Coma (Sword) consumes all combo orbs (leaving at least 1) after the finisher attack.
	consume_combo_orbs(me)
end
