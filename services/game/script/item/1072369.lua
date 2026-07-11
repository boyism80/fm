-- Item (1072369)

function on_item_gain_1072369(me, item_id)
	if me == nil then
		return
	end
	me:run_quest_hook(1201, "on_quest_start_1201")
	local quest = me:quest(1201)
	if quest == nil then
		return
	end
	quest:record_ex("have", "1")
end
