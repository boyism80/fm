-- Item (1032061)

return {
	on_item_gain = function(me, item_id)
		if me == nil then
			return
		end
		me:run_quest_hook(1206, "on_quest_start")
		local quest = me:quest(1206)
		if quest == nil then
			return
		end
		quest:record_ex("have1", "1")
	end
}
