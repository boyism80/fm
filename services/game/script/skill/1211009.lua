-- Skill name (String.wz/Skill.img.xml): 매직 크래쉬

local combat = require("script/lib/combat")

function on_activated_1211009(me, skill, params)
	combat.for_each_mob_in_skill_area(me, skill, function(mob)
		mob:clear_buffs(MobBuff.MagicDefenseUp)
	end)
end
