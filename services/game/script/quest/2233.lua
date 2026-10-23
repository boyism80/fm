local quest_id = 2233

return {
    on_end = function(me, npc)
        local q = me:quest(quest_id)
        if q == nil then
            return
        end
        if not q:started() then
            q:start(npc, true)
            return
        end

        if me:family().total_reputation < 1000 then
            me:dialog(npc, "흐음, 아직 명성도 1000은 모으지 못한 모양이지? 조금 어려울 수 있지만 주니어들과 함께 사냥을 하다보면 금방 이룩할 수 있을걸세!")
            return
        end
        me:dialog(npc, "오~! 명성도가 1,000이 넘다니.. 자네도 이제 유명한 모험가가 되었어. 정말 잘 했네!", false, true)
        me:dialog(npc,
            "자네의 활약을 보면 나도 피가 다시 끓는군. 빨리 젊고 유능한 주니어를 다시 찾아봐야 겠어. 어떤가? 명성도를 올리는 것은 그리 어렵지 않지? 핵심은 주니어가 빠르게 성장하도록 지원과 관심을 아끼지 않는 것이라네..",
            false, true)
        me:dialog(npc, "그럼, 앞으로도 좋은 엘더가 되어 주게나.", false, true)

        if me:exchange(nil, { exp = 2400 }) ~= ExchangeResult.OK then
            return
        end
        q:force_complete(npc)
        me:show_effect(EffectType.QuestCompletion)
    end
}
