-- Item (1082232)

function on_item_gain_1082232(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1203, "on_quest_start_1203")
	local quest = me:quest(1203)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
