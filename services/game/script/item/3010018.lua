-- Item (3010018)

return {
	on_item_gain = function(me, item_id)
		if me == nil then
			return
		end
		me:run_quest_hook(1300, "on_quest_start")
		local quest = me:quest(1300)
		if quest == nil then
			return
		end
		quest:record_ex("have", "1")
	end
}
