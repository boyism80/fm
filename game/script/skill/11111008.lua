-- Skill name (String.wz/Skill.img.xml): 샤우트

function on_attack_11111008(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_11111008(me, skill, params)
end
