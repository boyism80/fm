-- Skill name (String.wz/Skill.img.xml): 홀리 차지 : 검

local util = require("script/lib/skill")

function on_activated_1221003(me, skill, params)
	util.apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

