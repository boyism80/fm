-- Skill name (String.wz/Skill.img.xml): 파워 가드

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_buff_from_effect(me, skill, BuffFlag.Powerguard, "x")
	end
}
