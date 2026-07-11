local quest_id = 2318

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_accept(npc, "Hmmm... I looked into the making of the Spores while you were gathering up the Poison Mushroom Caps, and realised that we'll need more materials for it. I want you to gather up one more set of items. Can you do it?") then
		me:dialog(npc, "I understand it's not an easy task, but I can't make #bKiller Mushroom Spores#k without them. Please reconsider.", false, false)
		return
	end

	q:start(npc, true)
	me:dialog(npc, "Okay, I want you to defeat the Regenade Spores and bring back #b50 Mutated Spores#k in return.", false, true)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "Did you gather up all the necessary ingredients for it?", false, true)

	local code = me:exchange({ item = { [4000499] = 50 } }, { exp = 11500 })
	if code == ExchangeResult.LackCapacity then
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	me:dialog(npc, "Okay, these should be enough for me to make the #bKiller Mushroom Spores.#k Please hold on for a bit.", false, true)
	q:force_complete(npc)
	me:dialog(npc, "Okay, here are the Killer Mushroom Spores. Hopefully this will be enough for you to save our princess and help regain our kingdom. Good luck!", true, false)

	local code2 = me:exchange({}, { item = { [2430014] = 1 } })
	if code2 ~= ExchangeResult.OK then
		return
	end
end
