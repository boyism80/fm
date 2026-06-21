-- Skill name (String.wz/Skill.img.xml): 럭키 세븐

local combat = require("script/lib/combat")

function on_attack_4001344(me, skill, damages)
	combat.apply_venom(me, damages, Skill.Venom)
end
