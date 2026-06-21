-- Skill name (String.wz/Skill.img.xml): 닌자 앰부쉬

local combat = require("script/lib/combat")

function on_activated_4221004(me, skill, params)
	combat.apply_ninja_ambush_skill(me, skill)
end
