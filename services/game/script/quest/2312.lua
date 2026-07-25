local quest_id = 2312

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "We need your help, noble explorer. Our kingdom is currently facing a big threat, and we are in desperate need of a courageous explorer willing to fight for us, and that's how you ended up here. Please understand, though, that since we need place our faith in you, we'll have to test your skills first before we can stand firmly behind you. Will it be okay for you to do this for us?") then
			me:dialog(npc, "Hmmm... you must be unsure of your combat skills. We'll be here waiting for you, so come see us when you're ready.", false, false)
			return
		end

		q:start(npc, true)
		me:dialog(npc, "Keep moving forward, and you'll see #bRenegade Spores#k, the Spores that turned their backs on the Kingdom of Mushroom. We'd appreciate it if you can teach them a lesson or two, and bring back #b50 Mutated Spores#k in return.", false, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "Did you teach those Renegade Spores a lesson?", false, true)

		local code = me:exchange({ item = { [4000499] = 50 } }, { exp = 11500 })
		if code == ExchangeResult.LackCapacity then
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
		me:dialog(npc, "That was amazing. I apologize for doubting your abilities. Please save our Kingdom of Mushroom from this crisis!", true, true)
	end
}
