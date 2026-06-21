-- Skill name (String.wz/Skill.img.xml): 어벤져

local combat = require("script/lib/combat")

function on_attack_4111005(me, skill, damages)
	combat.apply_venom(me, damages, Skill.Venom)
end

function on_activated_4111005(me, skill, params)
end
