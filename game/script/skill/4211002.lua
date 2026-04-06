-- Skill name (String.wz/Skill.img.xml): 어썰터

function on_attack_4211002(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
	apply_venom(me, damages, Skill.Venom4220005)
end

function on_activated_4211002(me, skill, params)
end
