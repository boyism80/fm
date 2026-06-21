-- Skill name (String.wz/Skill.img.xml): 소울 에로우
local util = require("script/lib/skill")

function on_activated_13101003(me, skill, params)
	util.apply_buff_fixed(me, skill, BuffFlag.SoulArrow, 1)
end
