-- 칭호 도전 - 메이플 아이돌스타 (Quest.wz/QuestData/29500.img): 칭호 도전 - 메이플 아이돌스타

local medal = require("script/lib/medal")

local quest_id = 29500
local COOLDOWN_RECORD = "idol_star_challenge"
local COOLDOWN_SEC = 30 * 24 * 60 * 60

return {
    on_start = function(me, npc)
        local q = me:quest(quest_id)
        if q == nil then
            return
        end

        local item = medal.item(medal.IDOL_STAR)
        if me:dialog(npc, "#v" .. item .. "# #e#b#t" .. item .. "##k\r\n\r\n - 인기도 랭킹 1위\r\n\r\n#n이 훈장의 주인이 될 자격이 있는지 시험해 보시겠소?", false, true) == false then
            return
        end
        if me:dialog(npc, medal.ranking_text(medal.IDOL_STAR), true, true) == false then
            return
        end
        if not me:dialog_accept(npc, "자, 업적을 달성하면 나에게 돌아오시오. 업적을 달성하고 나서 나에게 와서 확인을 받아야만 인정받을 수 있다는 것을 꼭 명심하시오. 그리고 이 도전을 완료하거나 포기하지 않는 이상 다른 칭호에 도전할 수는 없다는 것도 알아두시오. 그리고 여기서 도전을 수락하면, 1달 동안 아이돌 스타에 도전 할 수 없으니 신중하게 선택하시오.") then
            return
        end
        local last = me:records():get(COOLDOWN_RECORD)
        if last > 0 and now() - last < COOLDOWN_SEC then
            me:dialog(npc, "아직 도전한지 1달이 지나지 않았소.")
            return
        end
        q:start(npc, true)
        me:records():set(COOLDOWN_RECORD, now())
    end,

    on_end = function(me, npc)
        local q = me:quest(quest_id)
        if q == nil then
            return
        end

        local sel = me:dialog_list(npc, "무엇을 하겠는가?", {
            "현재 랭킹을 보고 싶습니다.",
            "도전을 포기하고 싶습니다.",
            "아이돌스타 랭킹에 도전하겠습니다.",
        })
        if sel == nil then
            return
        end
        if sel == 1 then
            me:dialog(npc, medal.ranking_text(medal.IDOL_STAR))
            return
        end
        if sel == 2 then
            q:forfeit()
            me:dialog(npc, "다음에 도전하려면 언제든지 또 오시오.")
            return
        end

        if me:empty_slots(InventoryType.Equipment) < 1 then
            me:dialog(npc, "장비 인벤토리 공간이 충분한지 확인해주게.")
            return
        end
        local rank = submit_medal_score(me, medal.IDOL_STAR, me:population())
        if rank == nil then
            me:dialog(npc, "지금은 기록을 확인할 수 없구려. 잠시 후 다시 오시오.")
            return
        end
        if rank ~= 1 then
            me:dialog(npc, "흠, 자네는 아이돌스타에 도전하기엔 인기도가 너무 낮소.")
            return
        end
        if me:exchange({}, { item = medal.reward(me, medal.IDOL_STAR, rank) }) ~= ExchangeResult.OK then
            me:dialog(npc, "장비 인벤토리 공간이 충분한지 확인해주게.")
            return
        end
        q:force_complete(npc)
        me:dialog(npc, "호오, 대단하오. 훈장을 받을 자격이 충분하구려.")
    end,
}
