-- Skill name (String.wz/Skill.img.xml): 패닉

local util = require("script/lib/skill")

function on_activated_11111002(me, skill, params)
	util.consume_combo_orbs(me)
end
