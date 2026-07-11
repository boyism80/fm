-- Item (1122010)

function on_item_gain_1122010(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1205, "on_quest_start_1205")
	local quest = me:quest(1205)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
