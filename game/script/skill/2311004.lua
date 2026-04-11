-- Skill name (String.wz/Skill.img.xml): 샤이닝 레이

function on_attack_2311004(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_2311004(me, skill, params)
end
