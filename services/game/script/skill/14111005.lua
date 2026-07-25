-- Skill name (String.wz/Skill.img.xml): 트리플 스로우

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_venom(me, damages, Skill.VenomCygnus)
	end
}
