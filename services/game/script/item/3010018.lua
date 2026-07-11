-- Item (3010018)

function on_item_gain_3010018(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1300, "on_quest_start_1300")
	local quest = me:quest(1300)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
