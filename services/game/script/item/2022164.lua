-- Item name (String.wz/Consume.img.xml): 작은 어둠의 큐브

return {
	on_item_gain = function(me, item_id)
		local team = me:carnival_team()
		if team == nil then
			return
		end
		local enemy = team:enemy()
		if enemy == nil then
			return
		end
		enemy:debuff(me:map(), 8)
	end,
}
