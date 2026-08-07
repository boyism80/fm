-- NPC name (String.wz/Npc.img.xml): 분수 조각상

local EXIT_MAP = 120000102

local function shuffled_symbols()
	local symbols = { "1", "2", "3", "4", "5", "6" }
	for i = #symbols, 2, -1 do
		local j = math.random(i)
		symbols[i], symbols[j] = symbols[j], symbols[i]
	end
	return table.concat(symbols)
end

local function complete_trial(me, sm)
	local quest = me:quest(6401)
	if quest ~= nil then
		if quest:started() then
			quest:record("q3")
		else
			quest:start("q3")
		end
	end
	sm:finish(0)
	me:map(EXIT_MAP, 1)
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil or sm:group():name() ~= "air_strike_2" then
			me:map(EXIT_MAP, 1)
			return
		end
		if sm:get_property("ready") == nil then
			me:dialog(npc, "... 조나단에게 이야기는 전해 들었다. 내가 자네를 시험해 주도록 하겠다.")
			me:dialog(npc, "내가 말하는 수만큼 동상이 임의로 반짝일 것이다. 반짝인 동상을 같은 순서로 타격한 후 다시 말을 걸라.")
			sm:set_property("ready", "1")
			sm:set_property("stage", "1")
			sm:set_property("progress", "ready")
			return
		end
		local stage = tonumber(sm:get_property("stage")) or 1
		local progress = sm:get_property("progress")
		if progress == "ready" then
			local count = stage + 3
			me:dialog(npc, count .. "개의 동상이 반짝일 것이다. 잘 기억하여 순서대로 타격한 후 다시 말을 걸라.")
			local combo = string.sub(shuffled_symbols(), 1, count)
			sm:set_property("combo", combo)
			sm:set_property("guess", "")
			sm:set_property("progress", "guess")
			local map = me:map()
			for i = 1, count do
				sleep(3500)
				local reactor = map:reactor_by_name(string.sub(combo, i, i))
				if reactor ~= nil then
					reactor:hit(me)
				end
			end
			return
		end
		if progress ~= "guess" then
			return
		end
		local combo = sm:get_property("combo") or ""
		local guess = sm:get_property("guess") or ""
		if guess ~= combo .. combo then
			sm:set_property("stage", "1")
			sm:set_property("progress", "ready")
			me:dialog(npc, ".... 틀렸다. 처음부터 다시 시험하겠다. 준비가 되면 다시 말을 걸라.")
			return
		end
		if stage >= 3 then
			me:dialog(npc, "그대는 나의 시험을 통과했다. 조나단에게 자네의 지혜에 대해 이야기해 두겠다.")
			complete_trial(me, sm)
			return
		end
		sm:set_property("stage", tostring(stage + 1))
		sm:set_property("progress", "ready")
		me:dialog(npc, "올바르게 정답을 맞추었다. 다음 시험을 준비하겠다.")
	end
}
