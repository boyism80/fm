-- Skill name (String.wz/Skill.img.xml): 페럴라이즈

local combat = require("script/lib/combat")

function on_attack_2121006(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Freeze)
end
