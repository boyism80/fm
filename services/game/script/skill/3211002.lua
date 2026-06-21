-- Skill name (String.wz/Skill.img.xml): 퍼펫

local util = require("script/lib/skill")

function on_activated_3211002(me, skill, params)
	util.mark_archer_puppet(me, skill, params)
end

function on_buff_3211002(me, skill)
	util.apply_archer_puppet_buff(me, skill)
end

function on_unbuff_3211002(me, skill)
	util.apply_archer_puppet_unbuff(me, skill)
end
