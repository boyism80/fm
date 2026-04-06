-- Skill name (String.wz/Skill.img.xml): 에로우 봄 : 활

function on_attack_3101005(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_3101005(me, skill, params)
end
