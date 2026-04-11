-- Skill name (String.wz/Skill.img.xml): 암살

function on_attack_4221001(me, skill, damages)
	apply_venom(me, damages, Skill.Venom4220005)
end

function on_activated_4221001(me, skill, params)
end
