-- Skill name (String.wz/Skill.img.xml): 몬스터 라이딩

local MAX_FATIGUE = 100

return {
	on_activated = function(me, skill, params)
		local taming_mob = me:equipped(EquipmentPart.TamingMob)
		if taming_mob == nil or me:equipped(EquipmentPart.Saddle) == nil then
			return false
		end
		local _, _, fatigue = me:mount()
		if fatigue >= MAX_FATIGUE then
			return false
		end
		return me:ride(skill, taming_mob:wz():id())
	end,

	on_ride = function(me, skill)
		me:start_mount_fatigue()
	end,

	on_rider_equipment_changed = function(me, skill, part)
		if part == EquipmentPart.TamingMob or part == EquipmentPart.Saddle then
			me:dismount()
		end
	end
}
