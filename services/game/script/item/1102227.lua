-- Item (1102227)

return {
	on_item_gain = function(me, item_id)
		if me == nil then
			return
		end
		me:run_quest_hook(1303, "on_quest_start")
		local quest = me:quest(1303)
		if quest == nil then
			return
		end
		quest:record_ex("have0", "1")
	end
}
