-- Skill name (String.wz/Skill.img.xml): 새비지 블로우

local combat = require("script/lib/combat")

function on_attack_4201005(me, skill, damages)
	combat.apply_venom(me, damages, Skill.Venom4220005)
end

function on_activated_4201005(me, skill, params)
end
