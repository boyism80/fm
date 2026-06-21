-- Skill name (String.wz/Skill.img.xml): 암살

local combat = require("script/lib/combat")

function on_attack_4221001(me, skill, damages)
	combat.apply_venom(me, damages, Skill.Venom4220005)
end
