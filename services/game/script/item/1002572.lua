-- Item (1002572)

function on_item_gain_1002572(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1204, "on_quest_start_1204")
	local quest = me:quest(1204)
	if quest == nil then
		return
	end
	quest:record_ex("have1", "1")
end
