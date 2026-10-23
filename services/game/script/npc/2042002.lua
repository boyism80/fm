-- NPC name (String.wz/Npc.img.xml): 슈피겔만 - 몬스터 카니발

local HUB_MAP = 980000000
local MIN_LEVEL = 30
local MAX_LEVEL = 50
local RANKING_QUEST = 1301
local COIN_ID = 4001129

local function weapons(list)
    local entries = {}
    for i, row in ipairs(list) do
        entries[i] = { item = row[1], coin = row[2], random = true }
    end
    return entries
end

local COIN_EXCHANGE = {
    { label = "#t1122007#(50 코인)", entries = { { item = 1122007, coin = 50, random = true } } },
    { label = "#t2041211#(40 코인)", entries = { { item = 2041211, coin = 40 } } },
    {
        label = "전사용 무기 (한손검·두손검·도끼)",
        entries = weapons({
            { 1302004, 7 }, { 1402006, 7 }, { 1302009, 10 }, { 1402007, 10 }, { 1302010, 20 }, { 1402003, 20 },
            { 1312006, 7 }, { 1412004, 7 }, { 1312007, 10 }, { 1412005, 10 }, { 1312008, 20 }, { 1412003, 20 },
        })
    },
    {
        label = "전사용 무기 (둔기·창·폴암)",
        entries = weapons({
            { 1322015, 7 }, { 1422008, 7 }, { 1322016, 10 }, { 1422007, 10 }, { 1322017, 20 }, { 1422005, 20 },
            { 1432003, 7 }, { 1442003, 7 }, { 1432005, 10 }, { 1442009, 10 }, { 1442005, 20 }, { 1432004, 20 },
        })
    },
    {
        label = "마법사용 무기",
        entries = weapons({
            { 1372001, 7 }, { 1382018, 7 }, { 1372012, 10 }, { 1382019, 10 }, { 1382001, 20 }, { 1372007, 20 },
        })
    },
    {
        label = "궁수용 무기",
        entries = weapons({
            { 1452006, 7 }, { 1452007, 10 }, { 1452008, 20 }, { 1462005, 7 }, { 1462006, 10 }, { 1462007, 20 },
        })
    },
    {
        label = "도적용 무기",
        entries = weapons({
            { 1472013, 7 }, { 1472017, 10 }, { 1472021, 20 }, { 1332014, 7 }, { 1332031, 10 }, { 1332011, 10 },
            { 1332016, 20 }, { 1332003, 20 },
        })
    },
}

local function reward_for_cp(total_cp, winner)
    if total_cp >= 501 then
        if winner then
            return 30000, "A"
        end
        return 25500, "A"
    end
    if total_cp >= 251 then
        if winner then
            return 25500, "B"
        end
        return 20500, "B"
    end
    if total_cp >= 50 then
        if winner then
            return 21000, "C"
        end
        return 17000, "C"
    end
    if winner then
        return 3000, "D"
    end
    return 15000, "D"
end

local function on_reward(me, npc)
    local team = me:carnival_team()
    if team == nil then
        me:map(HUB_MAP, 0)
        return
    end
    local winner = team:is_winner()
    local exp, rank = reward_for_cp(team:total_cp(), winner)
    local message = "아쉽게도 비기거나 지고 말았군. 승리를 위해 좀 더 노력해주게! \r\n#b랭크 : " .. rank
    if winner then
        message = "축하하네. 카니발에서 승리를 거두었군. 자네들의 활약은 잘 지켜보았다네. 다음 번에도 기대하겠네! \r\n#b랭크 : " .. rank
    end
    if not me:dialog(npc, message) then
        return
    end
    me:end_party_quest(RANKING_QUEST)
    me:exchange({}, { exp = exp })
    team:remove_member(me)
    me:map(HUB_MAP, 0)
end

return {
    on_click = function(me, npc)
        local map_id = me:map():wz():id()
        if map_id == 980000010 then
            me:map(HUB_MAP, 0)
            return
        end
        if map_id >= 980000000 and (map_id % 10 == 3 or map_id % 10 == 4) then
            on_reward(me, npc)
            return
        end
        local sel = me:dialog_list(npc, "무얼 하겠는가?", {
            "몬스터 카니발 필드로 이동한다.",
            "몬스터 카니발에 대한 설명을 듣는다.",
            "메이플 코인을 교환한다.",
        })
        if sel == nil then
            return
        end
        if sel == 1 then
            if me:level() < MIN_LEVEL or me:level() > MAX_LEVEL then
                me:dialog(npc, "레벨 30 이상, 50 이하의 캐릭터만 몬스터 카니발을 즐길 수 있다네.")
                return
            end
            me:save_location("MONSTERCARNIVAL")
            me:map(HUB_MAP, 0)
        elseif sel == 2 then
            me:dialog(npc, "직접 그 전율을 느껴보기 전에는 이것이 무엇인지 알 수 없지.")
        elseif sel == 3 then
            local labels = {}
            for i, category in ipairs(COIN_EXCHANGE) do
                labels[i] = category.label
            end
            local pick = me:dialog_list(npc, "교환하고 싶은 아이템을 선택하게나.", labels)
            if pick == nil then
                return
            end
            local entries = COIN_EXCHANGE[pick].entries
            local entry = entries[1]
            if #entries > 1 then
                local choices = {}
                for i, e in ipairs(entries) do
                    choices[i] = "#z" .. e.item .. "#(" .. e.coin .. " 코인)"
                end
                local index = me:dialog_list(npc, "교환하고 싶은 무기를 선택해보라구. 내가 주는 무기는 기본보다 좋은 경우가 많지.", choices)
                if index == nil then
                    return
                end
                entry = entries[index]
            end

            local code = me:exchange({ item = { [COIN_ID] = entry.coin } },
                { item = { [entry.item] = 1 }, random = entry.random })
            if code == ExchangeResult.LackCapacity then
                me:dialog(npc, "인벤토리에 빈 칸이 없는 것 같군. 확인해 보게나.")
                return
            end
            if code ~= ExchangeResult.OK then
                me:dialog(npc, "#b#t" .. COIN_ID .. "##k이 부족한 것 같군.")
                return
            end
            me:dialog(npc, "여기 #b#t" .. entry.item .. "##k일세. 잘 사용하게나.")
        end
    end
}
