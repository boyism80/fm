-- Skill name (String.wz/Skill.img.xml): 드레인

local combat = require("script/lib/combat")

function on_attack_4101005(me, skill, damages)
	combat.drain_hp_from_damage(me, skill, damages)
end

function on_activated_4101005(me, skill, params)
end
