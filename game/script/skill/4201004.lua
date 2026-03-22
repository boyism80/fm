-- Skill name (String.wz/Skill.img.xml): 스틸

function on_attack_4201004(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Stun)
end

function on_activated_4201004(me, skill, params)
end
