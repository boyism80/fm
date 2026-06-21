-- Skill name (String.wz/Skill.img.xml): 뱀파이어

local combat = require("script/lib/combat")

function on_attack_14101006(me, skill, damages)
	combat.drain_hp_from_damage(me, skill, damages)
	combat.apply_venom(me, damages, Skill.VenomCygnus)
end

function on_activated_14101006(me, skill, params)
end
