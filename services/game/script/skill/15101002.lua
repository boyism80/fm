-- Skill name (String.wz/Skill.img.xml): 너클 부스터

local util = require("script/lib/skill")

function on_activated_15101002(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end
