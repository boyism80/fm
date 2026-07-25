local quest_id = 2315

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "A powerful barrier of magic, huh? Then what should we do...? If we can't find a way to break that barrier, then we can't save the princess. If it's impossible to physically break through, as you mentioned, then how about requesting help from our #bMinister of Magic#k?") then
			me:dialog(npc, "Please do not forget our plea for help.", false, false)
			return
		end

		q:start(npc, true)
		me:dialog(npc, "Please go see him immediately. The #bMinister of Magic#k may seem a bit on the edge, but he's very knowledgeable, and I'm sure he'll know what to do.", false, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "What? You investigated the barrier at the Mushroom Forest?", false, true)

		local code = me:exchange({}, { exp = 4000 })
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "Hmmm...this is interesting. It's a barrier set up by someone with a powerful force of magic, which means there's no way we can manually break through it.", false, true)
		q:force_complete(npc)
	end
}
