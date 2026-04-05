-- Skill name (String.wz/Skill.img.xml): 포이즌 봄

function on_attack_14111006(me, skill, damages)
	apply_venom(me, damages, Skill.VenomCygnus)
end

function on_activated_14111006(me, skill, params)
end
