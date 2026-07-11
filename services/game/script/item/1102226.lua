-- Item (1102226)

function on_item_gain_1102226(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1303, "on_quest_start_1303")
	local quest = me:quest(1303)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
