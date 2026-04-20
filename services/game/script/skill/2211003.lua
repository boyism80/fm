-- Skill name (String.wz/Skill.img.xml): 썬더 스피어

function on_attack_2211003(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_2211003(me, skill, params)
end
