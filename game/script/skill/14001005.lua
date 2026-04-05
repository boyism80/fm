-- Skill name (String.wz/Skill.img.xml): 다크니스

function on_attack_14001005(me, skill, damages)
	apply_venom(me, damages, Skill.Venom)
end

function on_activated_14001005(me, skill, params)
end
