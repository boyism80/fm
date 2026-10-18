-- 칭호 도전 - 기부왕 (Quest.wz/QuestData/29503.img): 칭호 도전 - 기부왕

local medal = require("script/lib/medal")

local quest_id = 29503
local MIN_DONATION = 100000
local TOWNS = {
	"헤네시스 기부왕",
	"엘리니아 기부왕",
	"페리온 기부왕",
	"커닝시티 기부왕",
	"슬리피우드 기부왕",
	"노틸러스 기부왕",
	"리스항구 기부왕",
}

local function ranking_text(kind)
	return medal.ranking_text(kind, " 메소", true) .. "\r\n\r\n현 기부왕의 기부액수는 공개할 수 없으니 이해해 주시오. \r\n그리고 매월 초에 기록이 초기화된다는 것도 잊지 마시오."
end

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local sel = me:dialog_list(npc, "어떤 마을의 기부왕에 도전하시겠습니까?", TOWNS)
		if sel == nil then
			return
		end
		local donor = medal.DONORS[sel]
		local item = medal.item(donor.kind)
		if me:dialog(npc, "#v" .. item .. "# #e#b#t" .. item .. "##k\r\n\r\n - 제한시간 1시간\r\n - 마을에 최고액을 기부\r\n\r\n#n이 훈장의 주인이 될 자격이 있는지 시험해 보시겠소?", false, true) == false then
			return
		end
		if me:dialog(npc, ranking_text(donor.kind), true, true) == false then
			return
		end
		q:start(npc, true)
		q:record_ex("money", "0")
		q:record_ex("trymap", tostring(donor.map))
		me:dialog(npc, "그럼, 행운을 빌겠소. 1시간 이내에 나에게 돌아와서 심사를 받으시오. 그리고 이 도전을 완료하거나 포기하지 않는 이상 다른 칭호에 도전할 수는 없다는 것을 알아두시오.")
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local money = tonumber(q:record_ex("money")) or 0
		local sel = me:dialog_list(npc, "그대는 " .. money .. " 메소를 기부하였군. 무엇을 하시겠소?", {
			"현재 랭킹을 보고 싶습니다.",
			"스페셜 훈장에 도전하겠습니다.",
			"도전을 포기하고 싶습니다.",
		})
		if sel == nil then
			return
		end
		if sel == 1 then
			local town = me:dialog_list(npc, "어느 마을의 랭킹을 보고 싶소?", TOWNS)
			if town ~= nil then
				me:dialog(npc, ranking_text(medal.DONORS[town].kind))
			end
			return
		end
		if sel == 3 then
			if me:dialog_accept(npc, "더 이상 도전하지 않을 테요? 다른 칭호 퀘스트에 도전하기 위해서는 이 퀘스트를 포기해야 하오만...") then
				q:forfeit()
			end
			return
		end

		local try_map = tonumber(q:record_ex("trymap")) or 0
		local donor = medal.donor_by_map(try_map)
		if donor == nil or me:map():wz():id() ~= try_map then
			me:dialog(npc, "그대가 기부를 한 마을은 #b#m" .. try_map .. "##k인 걸로 아오만. 이 마을에 기부하려면 먼저 다른 마을의 기부를 포기하시오. 물론 한번 기부한 돈은 돌려받을 수 없겠지만...")
			return
		end
		if money < MIN_DONATION then
			me:dialog(npc, "명예의 신을 속일 수는 없소. 도전을 했으면 제대로 기부를 하고 오시오.")
			return
		end
		if q:expired() then
			q:forfeit()
			me:dialog(npc, "이미 1시간의 제한 시간이 지난 것 같소. 다시 도전하시오.")
			return
		end

		local rank = submit_medal_score(me, donor.kind, money)
		if rank == nil then
			me:dialog(npc, "지금은 기록을 확인할 수 없구려. 잠시 후 다시 오시오.")
			return
		end
		if me:exchange({}, { item = medal.reward(me, donor.kind, rank) }) ~= ExchangeResult.OK then
			me:dialog(npc, "장비 인벤토리 공간이 충분한지 확인해주시오.")
			return
		end
		q:force_complete(npc)
		if rank == 1 then
			me:dialog(npc, "호오, 대단하오. 훈장을 받을 자격이 충분하구려.")
		else
			me:dialog(npc, "그대는 " .. money .. " 메소를 기부하였네만... 기부왕이 되기에는 모자라는군...\r\n 비록 최고는 아니지만 그대도 대단하다는 것은 틀림 없소.")
		end
	end,
}
