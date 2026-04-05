-- Skill name (String.wz/Skill.img.xml): 부메랑 스텝

function on_attack_4221007(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Stun)
	apply_venom(me, damages, Skill.Venom4220005)
end

function on_activated_4221007(me, skill, params)
end
