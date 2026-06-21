-- Skill name (String.wz/Skill.img.xml): 아이스 스트라이크

local combat = require("script/lib/combat")

function on_attack_2211002(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Freeze)
end
