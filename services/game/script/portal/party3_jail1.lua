local function map_count(group, map_id)
	local map = group:map(map_id)
	if map == nil then
		return 0
	end
	local n = 0
	for _, _ in pairs(map:characters()) do
		n = n + 1
	end
	return n
end

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local group = sm:group()
		if group == nil then
			return
		end
		if map_count(group, 920010910) > 0
			or map_count(group, 920010911) > 0
			or map_count(group, 920010912) > 0 then
			me:notice("이미 감옥에 누군가가 들어가 있습니다.")
			return
		end
		me:map(920010910)
	end
}
