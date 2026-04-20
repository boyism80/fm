-- Skill name (String.wz/Skill.img.xml): 빅뱅

function on_attack_2221001(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Freeze)
end

function on_activated_2221001(me, skill, params)
end
