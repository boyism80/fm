-- Skill name (String.wz/Skill.img.xml): 에너지 드레인

function on_attack_15111001(me, skill, damages)
	apply_skill_drain_on_attack(me, skill, damages)
end

function on_activated_15111001(me, skill, params)
end
