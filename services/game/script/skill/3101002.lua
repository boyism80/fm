-- Skill name (String.wz/Skill.img.xml): 보우 부스터

local util = require("script/lib/skill")

function on_activated_3101002(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end
