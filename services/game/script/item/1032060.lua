-- Item (1032060)

function on_item_gain_1032060(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1206, "on_quest_start_1206")
	local quest = me:quest(1206)
	if quest == nil then
		return
	end
	quest:record_ex("have0", "1")
end
