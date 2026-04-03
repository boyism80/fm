-- Skill name (String.wz/Skill.img.xml): 아이스 데몬

function on_attack_2221003(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Poison)
end

function on_activated_2221003(me, skill, params)
end
