-- Skill name (String.wz/Skill.img.xml): 포이즌 브레스

local combat = require("script/lib/combat")

function on_attack_2101005(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Poison)
end
