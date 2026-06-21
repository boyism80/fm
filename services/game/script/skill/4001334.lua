-- Skill name (String.wz/Skill.img.xml): 더블 스탭

local combat = require("script/lib/combat")

function on_activated_4001334(me, skill, params)
end

function on_attack_4001334(me, skill, damages)
	combat.apply_venom(me, damages, Skill.Venom4220005)
end
