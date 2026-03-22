-- Skill name (String.wz/Skill.img.xml): 매직 컴포지션

function on_attack_2211006(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Freeze)
end

function on_activated_2211006(me, skill, params)
end
