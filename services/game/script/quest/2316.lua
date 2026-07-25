local quest_id = 2316

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "I think i've heard of a potion that breaks these kinds of barriers. I think it's called #bKiller Mushroom Spores#k? Hmmm... outside, you'll find the Mushroom Scholar #bScarrs#k waiting outside. #bScarrs#k is an expert on mushrooms, so go talk to him.") then
			me:dialog(npc, "Why did you even ask if you were going to say no to this?#", false, false)
			return
		end

		q:start(npc, true)
		me:dialog(npc, "I am confident #kScarrs#k will do everything to help you.", false, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "Ah, so you're the explorer people were talking about. I'm #bScarrs, the Royal Mushroom Scholar#k representing the Kingdom of Mushroom. So you need some #kKiller Mushroom Spores#k?", false, true)

		local code = me:exchange({}, { exp = 4200 })
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "#kKiller Mushroom Spores#k... I think i've heard of them before...", false, true)
		q:force_complete(npc)
	end
}
