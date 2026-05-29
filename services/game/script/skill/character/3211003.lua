-- Skill name (String.wz/Skill.img.xml): 아이스 샷

function on_attack_3211003(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Freeze)
end

function on_activated_3211003(me, skill, params)
end
