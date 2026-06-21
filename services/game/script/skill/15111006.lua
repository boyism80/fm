-- Skill name (String.wz/Skill.img.xml): 스파크

local util = require("script/lib/skill")

function on_activated_15111006(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Spark, "x")
end
