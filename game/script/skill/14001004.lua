-- Skill name (String.wz/Skill.img.xml): 럭키 세븐

function on_attack_14001004(me, skill, damages)
	apply_venom(me, damages, Skill.VenomCygnus)
end

function on_activated_14001004(me, skill, params)
end
