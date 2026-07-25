-- Skill name (String.wz/Skill.img.xml): 건 부스터

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
	end
}
