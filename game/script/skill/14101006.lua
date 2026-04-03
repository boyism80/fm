-- Skill name (String.wz/Skill.img.xml): 뱀파이어

function on_attack_14101006(me, skill, damages)
	apply_skill_drain_on_attack(me, skill, damages)
end

function on_activated_14101006(me, skill, params)
end
