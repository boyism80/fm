-- Item (1022073)

function on_item_gain_1022073(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1202, "on_quest_start_1202")
	local quest = me:quest(1202)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
