-- Item (1122007)

function on_item_gain_1122007(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1301, "on_quest_start_1301")
	local quest = me:quest(1301)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
