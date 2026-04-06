-- Skill name (String.wz/Skill.img.xml): 백스핀 블로우

function on_attack_5101002(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_5101002(me, skill, params)
end
