-- Skill name (String.wz/Skill.img.xml): 라이트닝 차지

local util = require("script/lib/skill")

function on_activated_15101006(me, skill, params)
	util.apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

