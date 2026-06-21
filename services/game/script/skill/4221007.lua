-- Skill name (String.wz/Skill.img.xml): 부메랑 스텝

local combat = require("script/lib/combat")

function on_attack_4221007(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
	combat.apply_venom(me, damages, Skill.Venom4220005)
end
