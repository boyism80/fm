-- Skill name (String.wz/Skill.img.xml): 더블 스탭

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_venom(me, damages, Skill.Venom4220005)
	end
}
