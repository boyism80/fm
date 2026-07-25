local quest_id = 2313

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "I have told our #bMinister of Home Affairs#k of your abilities. Please go pay a visit to him immediately.") then
			me:dialog(npc, "There's not much time. Please hurry.", false, false)
			return
		end

		q:start(npc, true)
		me:dialog(npc, "Save our kingdom! We believe in you!", false, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local code = me:exchange({}, { exp = 4000 })
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
	end
}
