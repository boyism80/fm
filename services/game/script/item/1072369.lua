-- Item (1072369)

return {
	on_item_gain = function(me, item_id)
		if me == nil then
			return
		end
		me:run_quest_hook(1201, "on_quest_start")
		local quest = me:quest(1201)
		if quest == nil then
			return
		end
		quest:record_ex("have", "1")
	end
}
