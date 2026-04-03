-- Skill name (String.wz/Skill.img.xml): 드레인

function on_attack_4101005(me, skill, damages)
	apply_skill_drain_on_attack(me, skill, damages)
end

function on_activated_4101005(me, skill, params)
end
