local quest_id = 2323

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog(npc, "Ah! There might be a way... if you can utilize the spine vine that we have grown for the protection of our castle, then you just might be able to enter the premise!", false, true) then
			return
		end
		if not me:dialog_accept(npc, "If you can somehow eliminate the spines from the spine vine, then you'll be able to climb over the castle wall using the vine. Of course, that'll also require a Vine Remover...") then
			me:dialog(npc, "This will be the only way for you to enter the castle. Please think it through", false, false)
			return
		end
		if not me:dialog(npc, "The #bSpine Remover#k is created out of extracts from mysterious herbs at the highlands of El Naths. King Pepe used these herbs to intoxicate the pigs and take over the Mushroom Forest. #bIntoxicated Pig Tail#k is where you'll find the extracts of the herb. Please gather up #b100 Intoxicated Pig Tails#k and take them over to #bMinister of Magic.#k", false, true) then
			return
		end

		q:start(npc, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "Hmmm I see... so they have completely shut off the entrance and everything.", false, true)

		local code = me:exchange({}, { exp = 11000 })
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "Good job navigating through the area.", false, true)
		q:force_complete(npc)
	end
}
