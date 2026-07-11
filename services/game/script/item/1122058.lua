-- Item (1122058)

function on_item_gain_1122058(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1302, "on_quest_start_1302")
	local quest = me:quest(1302)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
