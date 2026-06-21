-- Skill name (String.wz/Skill.img.xml): 스탠스
local util = require("script/lib/skill")

function on_activated_1121002(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Stance, "prop")
end
