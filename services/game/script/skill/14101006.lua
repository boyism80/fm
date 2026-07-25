-- Skill name (String.wz/Skill.img.xml): 뱀파이어

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.drain_hp_from_damage(me, skill, damages)
		combat.apply_venom(me, damages, Skill.VenomCygnus)
	end
}
