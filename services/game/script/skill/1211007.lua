-- Skill name (String.wz/Skill.img.xml): 선더 차지 : 검

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
	end
}
