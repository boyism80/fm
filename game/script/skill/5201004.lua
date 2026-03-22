-- Skill name (String.wz/Skill.img.xml): 페이크 샷

function on_attack_5201004(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Stun)
end

function on_activated_5201004(me, skill, params)
end
