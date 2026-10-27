-- NPC name (String.wz/Npc.img.xml): 요정 웡키

local pq = require("script/lib/party_quest")

local GROUP_NAME = "orbis_party_quest"
local MIN_PARTY_SIZE = 6
local MIN_LEVEL = 51
local SCALE_LEVEL = 70
local FEATHER = 4001158
local BRACELET = 1082232
local ENTRY_MAP = 200080101
local GIVEUP_MAP = 920010000

local MISSION_RECORD = "wonky.mission"
local MISSION_SECONDS = 3600
local FEEL_RECORD = "wonky.feel"
local FOODS_RECORD = "wonky.foods"
local FED_AT_RECORD = "wonky.fed_at"
local MIN_FEEL = 100
local MAX_FEEL = 500
local FEEL_DECAY_SECONDS = 5040

local MISSION_CLASSES = {
    { class = Class.Fighter, name = "파이터" },
    { class = Class.Page, name = "페이지" },
    { class = Class.Spearman, name = "스피어맨" },
    { class = Class.FpWizard, name = "위자드(불/독)" },
    { class = Class.IlWizard, name = "위자드(얼음/번개)" },
    { class = Class.Cleric, name = "클레릭" },
    { class = Class.Hunter, name = "헌터" },
    { class = Class.Crossbowman, name = "사수" },
    { class = Class.Assassin, name = "어쌔신" },
    { class = Class.Bandit, name = "시프" },
}

local BLESSINGS = {
    { item = 2022090, text = "공격력, 마력 상승 마법" },
    { item = 2022091, text = "물리방어력, 마법방어력 상승 마법" },
    { item = 2022093, text = "이동속도, 점프력 상승 마법" },
    { item = 2022092, text = "명중율, 회피율 상승 마법" },
}

local FOODS = {
    {
        item = 2020001,
        name = "통닭",
        gain = 1,
        chance = { four = 2, three = 1, two = 2, one = 2, fresh = 4, earlier = 3 },
        text = "이건 바깥은 바삭하고 속은 부드러운 #b통닭#k 이잖아? 맛있겠다~"
    },
    {
        item = 2022001,
        name = "팥죽",
        gain = 1,
        chance = { four = 2, three = 3, two = 4, one = 5, fresh = 8, earlier = 6 },
        text = "와~ 이건 따뜻한 #b#t2022001##k 이잖아? 보기만 해도 군침이 돋는걸~"
    },
    {
        item = 2020004,
        name = "햄버거",
        gain = 1,
        chance = { four = 2, three = 3, two = 4, one = 5, fresh = 8, earlier = 6 },
        text = "맛있는 #b햄버거#k 다! 많이 먹으면 살이 늘어날지도 모르겠지만.. "
    },
    {
        item = 2022003,
        name = "장어 구이",
        gain = 2,
        chance = { four = 2, three = 7, two = 10, one = 13, fresh = 19, earlier = 16 },
        text = "이건 #b장어 구이#k 잖아? 맛있게 구워져서 굉장히 좋아~"
    },
    {
        item = 2001001,
        name = "쭈쭈바",
        gain = 2,
        chance = { four = 2, three = 20, two = 27, one = 30, fresh = 39, earlier = 36 },
        text = "시원하고 달콤한 #b쭈쭈바#k 잖아~ 더운 여름날엔 이게 최고지~"
    },
    {
        item = 2020028,
        name = "초콜릿",
        gain = 2,
        chance = { four = 2, three = 26, two = 32, one = 40, fresh = 50, earlier = 47 },
        text = "달콤한 #b초콜릿#k이잖아~? 음~ 맛있겠다!"
    },
    {
        item = 2001002,
        name = "팥빙수",
        gain = 3,
        chance = { four = 2, three = 30, two = 40, one = 50, fresh = 69, earlier = 65 },
        text = "아이 차가워! #b팥빙수#k 잖아? 팥과 연유, 떡.. 그리고 빙수! 시원하고 달콤한 이 맛이란.. 음~ "
    },
}

local GIFTS = {
    {
        { chance = 20, item = 2020013, count = 1 },
        { chance = 20, item = 2020014, count = 1 },
        { chance = 2, item = 2040708, count = 1 },
        { chance = 18, item = 2020015, count = 1 },
        { chance = 20, item = 2000002, count = 5 },
        { chance = 20, item = 2000002, count = 10 },
    },
    {
        { chance = 20, item = 2020013, count = 1 },
        { chance = 20, item = 2020014, count = 1 },
        { chance = 2, item = 2040708, count = 1 },
        { chance = 2, item = 2040707, count = 1 },
        { chance = 16, item = 2020015, count = 1 },
        { chance = 40, item = 2000002, count = 10 },
    },
    {
        { chance = 20, item = 2020013, count = 1 },
        { chance = 20, item = 2020014, count = 1 },
        { chance = 2, item = 2040707, count = 1 },
        { chance = 2, item = 2040708, count = 1 },
        { chance = 6, item = 2020015, count = 15 },
        { chance = 50, item = 2000002, count = 10 },
    },
}

local function feel(me)
    local value = me:records():get(FEEL_RECORD)
    if value < MIN_FEEL then
        return MIN_FEEL
    end
    return value
end

local function mission(me)
    local code = me:records():get(MISSION_RECORD)
    if code == 0 then
        code = math.random(0, #MISSION_CLASSES * 4 - 1)
        me:records():set(MISSION_RECORD, code + 1, { every = MISSION_SECONDS })
        return code
    end
    return code - 1
end

local function mission_target(code)
    local tier = code % 4
    return MISSION_CLASSES[math.floor(code / 4) + 1], 51 + tier * 5, 55 + tier * 5
end

local function mission_cleared(me, party, code)
    local map = me:map()
    local target, low, high = mission_target(code)
    local found = false
    for _, mem in ipairs(party:members()) do
        if mem ~= nil then
            local ch = map:characters()[mem:id()]
            if ch == nil then
                return false
            end
            if ch:class() == target.class and ch:level() >= low and ch:level() <= high then
                found = true
            end
        end
    end
    return found
end

local function bless(me, party, blessing)
    local map = me:map()
    for _, mem in ipairs(party:members()) do
        if mem ~= nil then
            local ch = map:characters()[mem:id()]
            if ch ~= nil then
                ch:message("요정 웡키의 미션을 수행하여 30분간 " .. blessing.text .. "이 적용됩니다.")
                ch:use_item(blessing.item)
            end
        end
    end
    me:records():remove(MISSION_RECORD)
end

local function offer_mission(me, npc, party)
    local code = mission(me)
    local target, low, high = mission_target(code)
    local level = string.format("%d~%d", low, high)
    if not me:dialog_yes_no(npc, "직업은 #b" .. target.name .. "#k, 레벨은 #b" .. level .. "#k 인 친구를 파티에 참가시켜보겠어? 뭔가 좋은 일이 생길지도 몰라.\r\n하지만 제안일 뿐이니까 꼭 할 필요는 없어. 만약 그냥 입장하고 싶다면 아니오를 선택하면 돼.") then
        me:dialog(npc, "알았어, 여신의 탑으로 이동시켜 줄게. 꼭 임무를 완수하고 여신 미네르바의 행방을 찾아줘~")
        return nil, true
    end
    if not mission_cleared(me, party, code) then
        me:dialog(npc, "미션을 만족하지 못했거나, 파티원이 모두 이곳에 있지 않아. 직업은 #b" .. target.name .. "#k,  레벨은 #b" .. level .. "#k 라구.")
        return nil, false
    end
    local blessing = BLESSINGS[math.random(1, #BLESSINGS)]
    me:dialog(npc, "와우~ 대단한데? 내가 말한 대로 딱 맞췄어. #b입장할 때 30분간 " .. blessing.text .. "#k을 걸어줄게. 여신의 흔적을 찾는데 도움이 될거야.")
    return blessing, true
end

local function gift(me, mood)
    local tier = 1
    if mood >= 300 then
        tier = 3
    elseif mood >= 200 then
        tier = 2
    end
    local roll = math.random(1, 100)
    for _, entry in ipairs(GIFTS[tier]) do
        roll = roll - entry.chance
        if roll <= 0 then
            if me:exchange(nil, { item = { [entry.item] = entry.count } }) ~= ExchangeResult.OK then
                me:message("인벤토리 공간이 부족합니다.", Msg.PinkText)
                return
            end
            me:message("웡키가 음식을 먹고 감사의 표시로 특별한 선물을 주었습니다!", Msg.PinkText)
            return
        end
    end
end

local function taste(me, npc, index, food, history)
    local streak = 0
    for i = #history, 1, -1 do
        if tonumber(history:sub(i, i)) ~= index then
            break
        end
        streak = streak + 1
    end
    if streak == 4 then
        me:dialog(npc, "#b#t" .. food.item .. "##k... 좋아하긴 해도 반복해서 먹으면 질린다구..")
        return food.chance.four
    end
    if streak == 3 then
        me:dialog(npc, "#b#t" .. food.item .. "##k... 다른걸 먹으면 안될까?")
        return food.chance.three
    end
    if streak == 2 then
        me:dialog(npc, "와~ #b#t" .. food.item .. "##k! 고마워~ 잘 먹을게!")
        return food.chance.two
    end
    if streak == 1 then
        me:dialog(npc, "와~ #b#t" .. food.item .. "##k! 고마워~ 잘 먹을게!")
        return food.chance.one
    end
    if history:find(tostring(index), 1, true) == nil then
        me:dialog(npc, "와! 이건 #b#t" .. food.item .. "##k! 정말 고마워~ 잘 먹을게!")
        return food.chance.fresh
    end
    me:dialog(npc, food.text)
    return food.chance.earlier
end

local function feed(me, npc)
    local options = {}
    for _, food in ipairs(FOODS) do
        table.insert(options, "#t" .. food.item .. "#")
    end
    local index = me:dialog_list(npc, "와~ 맛있는 음식을 준다고? 고마워~! 어떤 음식을 줄거야?\r\n#b", options)
    local food = FOODS[index]
    if food == nil then
        return
    end

    local mood = feel(me)
    local fed_at = now()
    if me:records():get(FED_AT_RECORD) + FEEL_DECAY_SECONDS < fed_at then
        mood = math.max(MIN_FEEL, mood - 30)
    end
    me:records():set(FED_AT_RECORD, fed_at)

    local history = me:records():text(FOODS_RECORD)
    if history == "" then
        history = "0000"
    end
    if me:exchange({ item = { [food.item] = 1 } }, nil) ~= ExchangeResult.OK then
        me:message("웡키의 기분이 나빠졌습니다.", Msg.PinkText)
        me:records():set(FEEL_RECORD, math.max(MIN_FEEL, mood - 3))
        me:dialog(npc, "#b#t" .. food.item .. "##k .. 분명 제대로 갖고 있는거야? 장난치지 말라구!")
        return
    end
    me:play_sound("Party1/Clear")
    local chance = taste(me, npc, index, food, history)
    me:records():set_text(FOODS_RECORD, history:sub(2) .. index)
    if math.random(1, 100) <= chance then
        me:message(food.name .. "을(를) 먹고 웡키의 기분이 무척 좋아졌습니다.", Msg.PinkText)
        if math.random(1, 10) == 1 then
            gift(me, mood)
        end
        mood = mood + 5
    else
        me:message(food.name .. "을(를) 먹고 웡키의 기분이 좋아졌습니다.", Msg.PinkText)
        mood = mood + food.gain
    end
    me:records():set(FEEL_RECORD, math.min(MAX_FEEL, mood))
end

local function greeting(me)
    local mood = feel(me)
    if mood >= 400 then
        return "#e<파티퀘스트 : 여신의 흔적>#n\r\n기다렸다구! 왜 이제야 오는거야~?! 얼른 여신의 탑으로 모험을 떠나자구!"
    end
    if mood >= 300 then
        return "#e<파티퀘스트 : 여신의 흔적>#n\r\n#b#h0##k! 돌아와줬구나~ 반가워! 헤헤. 언제든지 여신의 탑을 모험하고 싶으면 나에게 이야기 하라고~"
    end
    return "#e<파티퀘스트 : 여신의 흔적>#n\r\n이야아~ 반가워! 여신의 탑을 모험하고 싶어?"
end

local function try_start(me, npc)
    local party = me:party()
    if party == nil then
        me:dialog(npc, "파티원을 모으거나 파티에 참가한 후에 여신의 흔적을 찾으러 갈 수 있어.")
        return
    end
    local blessing, proceed = offer_mission(me, npc, party)
    if not proceed then
        return
    end
    if party:leader_id() ~= me:id() then
        me:dialog(npc, "흐음. 너는 파티장이 아닌 것 같은데? 퀘스트를 시작하려면 파티장이 내게 말을 걸어야 한다구.")
        return
    end
    local map = me:map()
    if map == nil then
        return
    end
    local ok = true
    local in_map = 0
    local has_gm = false
    for _, mem in ipairs(party:members()) do
        if mem ~= nil then
            if mem:level() < MIN_LEVEL then
                ok = false
            end
            local ch = map:characters()[mem:id()]
            if ch ~= nil then
                if pq.is_gm(ch) then
                    has_gm = true
                end
                in_map = in_map + 1
            end
        end
    end
    if not ok or (in_map ~= MIN_PARTY_SIZE and not has_gm) then
        me:dialog(npc, "파티 요구 조건이 맞지 않는 것 같은데? 파티는 #r레벨 51 이상의 파티원으로만 이루어진 6명의 파티#k만 입장이 가능해. 다시 한번 확인해 봐.")
        return
    end
    local group = state_machine(GROUP_NAME)
    if group == nil then
        me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
        return
    end
    if group:full() then
        me:dialog(npc, "이미 다른 파티가 이 안에서 퀘스트에 도전하는 중이야. 잠시 후 다시 시도해 주거나 채널을 변경해 봐.")
        return
    end
    if blessing ~= nil then
        bless(me, party, blessing)
    end
    local sm, err = group:start_party(me, party, SCALE_LEVEL)
    if sm == nil then
        if err == StateMachineError.Full then
            me:dialog(npc, "이미 다른 파티가 이 안에서 퀘스트에 도전하는 중이야. 잠시 후 다시 시도해 주거나 채널을 변경해 봐.")
        else
            me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
            log("orbis_party_quest start_party:", err)
        end
    end
end

local function claim_bracelet(me, npc)
    if not pq.has_item(me, FEATHER, 10) then
        me:dialog(npc, "#b여신의 깃털 10개#k는 가지고 있는거야? 여신의 아이템을 받고 싶으면 #b여신의 깃털 10개#k를 가지구 오라구~ 또 만들어 줄게!")
        return
    end
    local code = me:exchange(
        { item = { [FEATHER] = 10 } },
        { item = { [BRACELET] = 1 } }
    )
    if code == ExchangeResult.LackCapacity then
        me:dialog(npc, "인벤토리 공간을 확보한 뒤 다시 말을 걸어줘.")
        return
    end
    if code ~= ExchangeResult.OK then
        me:dialog(npc, "#b여신의 깃털 10개#k는 가지고 있는거야? 여신의 아이템을 받고 싶으면 #b여신의 깃털 10개#k를 가지구 오라구~ 또 만들어 줄게!")
        return
    end
    me:dialog(npc, "고마워~~! 여기 #b여신의 팔찌#k를 받아~ 혹시라도 또 필요하면 #b여신의 깃털 10개#k를 가지구 오라구~ 또 만들어 줄게!")
end

return {
    on_click = function(me, npc)
        local map = me:map()
        if map == nil or map:wz() == nil then
            return
        end
        local map_id = map:wz():id()
        if map_id == GIVEUP_MAP then
            if not pq.is_leader(me) then
                me:dialog(npc, "아직 포기하긴 이르다구? 조금 더 열심히 노력해 봐~")
                return
            end
            if me:dialog_yes_no(npc, "음.. 정말 이곳에서 나가고 싶어? 여기서 나간다면 처음부터 다시 시작해야만 해.") then
                me:map(920011200)
            else
                me:dialog(npc, "아직 포기하긴 이르다구? 조금 더 열심히 노력해 봐~")
            end
            return
        end
        if map_id ~= ENTRY_MAP then
            return
        end
        local sel = me:dialog_list(npc, greeting(me) .. " 참, 파티에 내가 말해주는 특정한 캐릭터가 1명이상 있으면 웡키의 축복을 걸어 줄게.\r\n\r\n#b", {
            "입장을 신청한다.",
            "여신의 탑에 대해 묻는다.",
            "웡키에게 먹을 것을 준다.",
            "여신의 깃털을 다른 아이템과 바꾼다.",
        })
        if sel == 1 then
            try_start(me, npc)
        elseif sel == 2 then
            me:dialog(npc, "며칠 전에 갑자기 오르비스 구름 위에 가장 아름다운 여신, 미네르바와 신비한 탑이 나타났어.")
            me:dialog(npc,
                "그런데 어느 순간 갑자기 여신은 사라지고 탑은 제어 불능에 빠져버렸어. 아무래도 여신이 탑 어딘가에 갇혀 있는 것 같아. 너는 어떻게 생각해? 그 탑을 탐험하고 여신의 흔적을 찾아 떠나보고 싶지 않아?")
        elseif sel == 3 then
            feed(me, npc)
        elseif sel == 4 then
            local craft = me:dialog_list(npc,
                "나의 배고픔을 달래 주던 #b#h ##k님이잖아! 그동안 고마웠다는 표시로 #b여신의 깃털을 여신의 아이템으로#k 만들어 줄게. 뭘 만들어 줄까?#b\r\n", {
                    "여신의 팔찌 #r(여신의 깃털 10개 필요)",
                })
            if craft == 1 then
                claim_bracelet(me, npc)
            end
        end
    end
}
