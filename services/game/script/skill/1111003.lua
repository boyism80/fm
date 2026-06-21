-- Skill name (String.wz/Skill.img.xml): 패닉 : 검

local util = require("script/lib/skill")

function on_activated_1111003(me, skill, params)
	-- Panic consumes all combo orbs (leaving at least 1) after the attack is fully resolved.
	util.consume_combo_orbs(me)
end
