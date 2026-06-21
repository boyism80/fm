-- Skill name (String.wz/Skill.img.xml): 어벤져

local combat = require("script/lib/combat")

function on_attack_14111002(me, skill, damages)
	combat.apply_venom(me, damages, Skill.VenomCygnus)
end
