-- Skill name (String.wz/Skill.img.xml): MP 리커버리
-- MP recovery requires HP >= 10% of max HP.

function on_preactive(me, skill)
	local max_hp = me:get_max_hp()
	local min_hp = math.max(1, math.floor(max_hp * 10 / 100))
	return me:get_hp() >= min_hp
end

function on_active(me, skill)
end
