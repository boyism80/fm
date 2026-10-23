local quest_id = 2232

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

        local juniors = me:family().juniors
        if juniors == 0 then
            me:dialog(npc, "흐음, 아직 자네가 후원해 줄 주니어를 찾지 못한 모양이구만?")
            return
        end
        me:dialog(npc, "오! 자네가 후원해 줄 주니어를 찾았나보군. 정말 잘 했네!", false, true)
        if juniors >= 2 then
            me:dialog(npc, "주니어를 두 명 모두 등록했군? 역시~ 자네는 내가 본 눈이 틀리지 않은 것 같군!", false, true)
        else
            me:dialog(npc, "주니어를 한 명만 등록했다면 다른 한 명도 어서 찾아보게. 당연히 한 명보다는 두 명이 낫지 않겠나?", false, true)
        end
        me:dialog(npc, "그럼, 좋은 엘더가 되어 주게나.", false, true)

        if me:exchange(nil, { exp = 2000 }) ~= ExchangeResult.OK then
            return
        end
        q:force_complete(npc)
        me:show_effect(EffectType.QuestCompletion)
    end
}
