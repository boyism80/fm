-- Skill name (String.wz/Skill.img.xml): 트리플 스로우

local combat = require("script/lib/combat")

function on_attack_14111005(me, skill, damages)
	combat.apply_venom(me, damages, Skill.VenomCygnus)
end

function on_activated_14111005(me, skill, params)
end
