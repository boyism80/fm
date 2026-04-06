-- Skill name (String.wz/Skill.img.xml): 인레이지
function on_activated_1121010(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local combo = me:buff_value(BuffFlag.Combo)
	if combo == nil or combo < 10 then
		return
	end
	consume_combo_orbs(me, 10)
	local pad = effect.pad or 0
	me:buff(skill, BuffFlag.WeaponAtk, pad)
end
