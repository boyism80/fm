-- Item (1002798)

function on_item_gain_1002798(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1200, "on_quest_start_1200")
	local quest = me:quest(1200)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
