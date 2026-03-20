-- Skill name (String.wz/Skill.img.xml): 메디테이션

function on_activated(me, skill, params)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end

	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end

	local mad = effect.mad or 0
	if mad <= 0 then
		return
	end

	me:buff(skill, BuffFlag.MagicAtk, mad)

	-- Party support is not implemented, so apply to self only.
end

function on_unbuff(me, skill)
end
