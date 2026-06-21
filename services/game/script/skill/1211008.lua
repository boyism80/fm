-- Skill name (String.wz/Skill.img.xml): 라이트닝 차지 : 둔기

local util = require("script/lib/skill")

function on_activated_1211008(me, skill, params)
	util.apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

