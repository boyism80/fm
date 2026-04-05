-- Skill name (String.wz/Skill.img.xml): 어벤져

function on_attack_14111002(me, skill, damages)
	apply_venom(me, damages, Skill.VenomCygnus)
end

function on_activated_14111002(me, skill, params)
end
