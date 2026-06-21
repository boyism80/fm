-- Skill name (String.wz/Skill.img.xml): 트리플 스로우

local combat = require("script/lib/combat")

function on_attack_4121007(me, skill, damages)
	combat.apply_venom(me, damages, Skill.Venom)
end
