local quest_id = 2322

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_yes_no(npc, "Like I told you, just breaking the barrier cannot be a cause for celebration. That's because our castle for the Kingdom of Mushroom completely denies entry of anyone outside our kingdom, so it'll be hard for you to do that. Hmmm... to figure out a way to enter, can you...investigate the outer walls of the castle first?") then
		me:dialog(npc, "Really? Is there another way you can penetrate the castle? If you don't know of one, then just come see me.", false, true)
		return
	end

	me:dialog(npc, "Walk past the Mushroom Forest and when you reach the #bSplit Road of Choice#k, just walk towards the castle. Good luck.", false, true)
	q:start(npc, true)
end

function on_end(me, npc)
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
