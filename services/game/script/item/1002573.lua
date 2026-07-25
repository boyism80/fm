-- Item (1002573)

return {
	on_item_gain = function(me, item_id)
		if me == nil then
			return
		end
		me:run_quest_hook(1204, "on_quest_start")
		local quest = me:quest(1204)
		if quest == nil then
			return
		end
		quest:record_ex("have2", "1")
	end
}
