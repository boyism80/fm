-- Item (1122058)

return {
	on_item_gain = function(me, item_id)
		if me == nil then
			return
		end
		me:run_quest_hook(1302, "on_quest_start")
		local quest = me:quest(1302)
		if quest == nil then
			return
		end
		quest:record_ex("have", "1")
	end
}
