-- Skill name (String.wz/Skill.img.xml): 파워 가드

local util = require("script/lib/skill")

function on_activated_1101007(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Powerguard, "x")
end
