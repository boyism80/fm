-- Skill name (String.wz/Skill.img.xml): 럭키 세븐

function on_activated_4001344(me, skill, params)
end

function on_attack_4001344(me, skill, damages)
	apply_venom(me, damages, Skill.Venom)
end
