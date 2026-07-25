local quest_id = 2314

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_yes_no(npc, "In order to rescue the princess, you must first navigate the Mushroom Forest. King Pepe set up a powerful barrier forbidding anyone from entering the castle. Please investigate this matter for us.") then
			me:dialog(npc, "Please do not lose faith in our Kingdom of Mushroom.", false, false)
			return
		end

		me:dialog(npc, "You'll run into the barrier at the Mushroom Forest by heading east of where you are standing right now. Please be careful. I hear that the area is infested with crazy, fear-inducing monsters.", false, true)
		q:start(npc, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "I see that you have thoroughly investigated the barrier at the Mushroom Forest. What was it like?", false, true)

		local code = me:exchange({}, { exp = 8300 })
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "I see, so it was indeed not a regular barrier by any means. Great work there. If not for you help, we wouldn't have had a clue as to what that was all about.", false, true)
		q:force_complete(npc)
	end
}
