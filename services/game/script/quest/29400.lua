-- 칭호 도전 - 노련한 사냥꾼 (Quest.wz/QuestData/29400.img): 칭호 도전 - 노련한 사냥꾼

local medal = require("script/lib/medal")

local quest_id = 29400
local MEDAL = 1142004
local GOAL = 100000

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if me:dialog(npc, "#v" .. MEDAL .. "# #e#b#t" .. MEDAL .. "##k\r\n\r\n - 제한시간 30일\r\n - 몬스터 100,000마리 사냥\r\n#n※ 자신의 레벨 이상의 몬스터만 인정\r\n   (120레벨 이상의 캐릭터는 120레벨 이상의 몬스터 인정)\r\n\r\n이 훈장의 주인이 될 자격이 있는지 시험해 보시겠소?", false, true) == false then
			return
		end
		if me:dialog(npc, medal.ranking_text(medal.EXPERT_HUNTER), true, true) == false then
			return
		end
		if me:dialog(npc, "자, 30일의 시간을 줄테니 사냥을 마치고 나에게 돌아오시오. 제한시간 내에 나에게 와서 확인을 받아야만 인정받을 수 있다는 것을 꼭 명심하시오. 그리고 이 도전을 완료하거나 포기하지 않는 이상 다른 칭호에 도전할 수는 없다는 것도 알아두시오.", true, true) == false then
			return
		end
		q:start(npc, true)
		q:record_ex("mg", tostring(GOAL))
		q:record_ex("mon", "0")
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local mon = tonumber(q:record_ex("mon")) or 0
		local options = { "현재 랭킹을 보고 싶습니다." }
		if mon >= GOAL then
			options[2] = "도전에 성공했습니다. 훈장을 지급받겠습니다."
		else
			options[2] = "도전을 포기하고 싶습니다."
		end
		local sel = me:dialog_list(npc, "그대가 사냥한 몬스터의 수는 " .. mon .. " 마리고 목표는 " .. GOAL .. " 마리라네. 무엇을 하겠는가?", options)
		if sel == nil then
			return
		end
		if sel == 1 then
			me:dialog(npc, medal.ranking_text(medal.EXPERT_HUNTER))
			return
		end
		if mon < GOAL then
			q:forfeit()
			me:dialog(npc, "다음에 도전하려면 언제든지 또 오시오.")
			return
		end
		if q:expired() then
			q:forfeit()
			me:dialog(npc, "이미 30일의 제한 시간이 지난 것 같소. 다시 도전하시오.")
			return
		end

		-- The score takes the first place medal from its holder, so the room is checked before it is sent.
		if me:empty_slots(InventoryType.Equipment) < 2 then
			me:dialog(npc, "장비 인벤토리 공간이 충분한지 확인해주게.")
			return
		end
		local rank = submit_medal_score(me, medal.EXPERT_HUNTER, mon)
		if rank == nil then
			me:dialog(npc, "지금은 기록을 확인할 수 없구려. 잠시 후 다시 오시오.")
			return
		end
		local items = medal.reward(me, medal.EXPERT_HUNTER, rank)
		if not medal.holds(me, MEDAL) then
			items[MEDAL] = 1
		end
		if me:exchange({}, { item = items }) ~= ExchangeResult.OK then
			me:dialog(npc, "장비 인벤토리 공간이 충분한지 확인해주게.")
			return
		end
		q:force_complete(npc)
		me:dialog(npc, "호오, 대단하오. 훈장을 받을 자격이 충분하구려.")
	end,
}
