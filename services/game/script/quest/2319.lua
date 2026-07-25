local quest_id = 2319

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "Oh, I almost forgot! What was I thinking? I need you to hand this #bSample of Killer Mushroom Spores#k to #bMinister of Magic#k and report the results.") then
			me:dialog(npc, "I know it's not a tough task, so come back to me if you're ready.", false, false)
			return
		end

		q:start(npc, true)
		local code = me:exchange({}, { item = { [4032389] = 1 } })
		if code == ExchangeResult.LackCapacity then
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "The #bMinister of Magic#k told me once the #bKiller Mushroom Spores#k is complete, that he'll want a sample of it as well. I'll give you the sample; now go please hand it in to our #bMinister of Magic.#k", false, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "Are the #bKiller Mushroom Spores#k finally completed?", false, true)

		local code = me:exchange({ item = { [4032389] = 1 } }, { exp = 4200 })
		if code == ExchangeResult.LackCapacity then
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "Okay, so this is the #bKiller Mushroom Spores.#k Thank you, thank you, and please tell #bScarrs#k the same.", false, true)
		q:force_complete(npc)
	end
}
