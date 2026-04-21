-- Skill name (String.wz/Skill.img.xml): 아이언 월

function on_activated_1301006(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.pdd <= 0 and effect.mdd <= 0 then
		return
	end
	local vals = {
		[BuffFlag.WeaponDef] = effect.pdd,
		[BuffFlag.MagicDef] = effect.mdd,
	}
	for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, vals)
	end)
end
