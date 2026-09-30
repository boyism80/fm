-- Item name (String.wz/Consume.img.xml): 카니발 포인트 2

return {
	on_item_gain = function(me, item_id)
		local team = me:carnival_team()
		if team == nil then
			return
		end
		team:add_cp(me, 2)
	end,
}
