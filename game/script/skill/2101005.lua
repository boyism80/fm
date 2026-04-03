-- Skill name (String.wz/Skill.img.xml): 포이즌 브레스

function on_attack_2101005(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Poison)
end

function on_activated_2101005(me, skill, params)
end

-- Skill name (String.wz/Skill.img.xml): 포이즌 브레스

function on_activated_2101005(me, skill, params)
end
