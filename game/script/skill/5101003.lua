-- Skill name (String.wz/Skill.img.xml): 더블 어퍼

function on_attack_5101003(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Stun)
end

function on_activated_5101003(me, skill, params)
end
