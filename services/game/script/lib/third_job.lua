local TRIAL_QUEST = 195000
local STRENGTH_NECKLACE = 4031057
local WISDOM_NECKLACE = 4031058
local MIN_LEVEL = 70

local M = {}

local function has_item(me, item_id)
	return next(me:item(item_id)) ~= nil
end

local function first_trial(me, npc, cfg, q)
	if has_item(me, STRENGTH_NECKLACE) == false then
		me:dialog(npc, cfg.text.trial1_hint)
		return
	end
	if me:dialog(npc, cfg.text.trial1_done, false, true) == false then
		return
	end
	if me:exchange({ item = { [STRENGTH_NECKLACE] = 1 } }, nil) ~= ExchangeResult.OK then
		return
	end
	q:record("job3_trial2_1")
	me:dialog(npc, cfg.text.trial2)
end

local function second_trial(me, npc, cfg, q)
	if has_item(me, WISDOM_NECKLACE) == false then
		me:dialog(npc, cfg.text.trial2_hint)
		return
	end
	if me:dialog(npc, cfg.text.trial2_done, false, true) == false then
		return
	end
	if me:dialog(npc, cfg.text.trials_done, false, true) == false then
		return
	end
	if cfg.text.confirm ~= nil and me:dialog_yes_no(npc, cfg.text.confirm) == false then
		return
	end
	if me:skill_point() > (me:level() - MIN_LEVEL) * 3 then
		me:dialog(npc, cfg.text.sp_left)
		return
	end
	local job = cfg.jobs[me:class()]
	if job == nil then
		me:dialog(npc, cfg.text.not_ready)
		return
	end
	if me:exchange({ item = { [WISDOM_NECKLACE] = 1 } }, nil) ~= ExchangeResult.OK then
		return
	end
	me:class(job.class)
	q:record("job3_clear")
	me:dialog(npc, cfg.text.done .. job.text)
end

function M.advance(me, npc, cfg)
	if me:class_of(cfg.family) == false then
		me:dialog(npc, cfg.text.other_class)
		return
	end
	local q = me:quest(TRIAL_QUEST)
	if q:started() == false then
		q:start("0")
	end
	local value = q:record()

	if string.sub(value, 1, 11) == "job3_trial1" then
		first_trial(me, npc, cfg, q)
		return
	end
	if string.sub(value, 1, 11) == "job3_trial2" then
		second_trial(me, npc, cfg, q)
		return
	end

	if cfg.jobs[me:class()] == nil or me:level() < MIN_LEVEL then
		me:dialog(npc, cfg.greetings[me:class()] or cfg.text.not_ready)
		return
	end
	if me:dialog_yes_no(npc, cfg.text.offer) == false then
		me:dialog(npc, cfg.text.decline)
		return
	end
	q:record("job3_trial1_1")
	me:dialog(npc, cfg.text.trial1)
end

return M
