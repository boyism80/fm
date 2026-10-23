-- NPC name (String.wz/Npc.img.xml): 크리스

local maker = require("script/lib/maker")

local CATEGORIES = {
    {
        text = "광석 제련",
        prompt = "어떤 광석을 제련하고 싶으신가요?#b",
        equip = false,
        recipes = {
            { item = 4011000, text = "청동", mats = { { 4010000, 10 } }, cost = 300 },
            { item = 4011001, text = "강철", mats = { { 4010001, 10 } }, cost = 300 },
            { item = 4011002, text = "미스릴", mats = { { 4010002, 10 } }, cost = 300 },
            { item = 4011003, text = "아다만티움", mats = { { 4010003, 10 } }, cost = 500 },
            { item = 4011004, text = "은", mats = { { 4010004, 10 } }, cost = 500 },
            { item = 4011005, text = "오리할콘", mats = { { 4010005, 10 } }, cost = 500 },
            { item = 4011006, text = "금", mats = { { 4010006, 10 } }, cost = 800 },
        },
    },
    {
        text = "보석 제련",
        prompt = "어떤 보석을 제련하고 싶으신가요?#b",
        equip = false,
        recipes = {
            { item = 4021000, text = "가넷", mats = { { 4020000, 10 } }, cost = 500 },
            { item = 4021001, text = "자수정", mats = { { 4020001, 10 } }, cost = 500 },
            { item = 4021002, text = "아쿠아마린", mats = { { 4020002, 10 } }, cost = 500 },
            { item = 4021003, text = "에메랄드", mats = { { 4020003, 10 } }, cost = 500 },
            { item = 4021004, text = "오팔", mats = { { 4020004, 10 } }, cost = 500 },
            { item = 4021005, text = "사파이어", mats = { { 4020005, 10 } }, cost = 500 },
            { item = 4021006, text = "토파즈", mats = { { 4020006, 10 } }, cost = 500 },
            { item = 4021007, text = "다이아몬드", mats = { { 4020007, 10 } }, cost = 1000 },
            { item = 4021008, text = "흑수정", mats = { { 4020008, 10 } }, cost = 3000 },
        },
    },
    {
        text = "아이언 호그의 강철 발굽을 가져왔어요..",
        confirm = "아이언 호그의 철발굽에는 슬픈 전설이.. 아니라 숨겨진 잠재된 힘이 있다는걸 아는사람은 극히 드물지요. 그것을 제게 주시지 않겠어요? 그렇다면 보상은 섭섭치 않게 드리지요.",
        equip = false,
        recipes = {
            { item = 4011001, mats = { { 4000039, 100 } }, cost = 1000 },
        },
    },
    {
        text = "아대 합성",
        prompt = "아대 합성이요? 어떤걸 원하세요?#b",
        equip = true,
        recipes = {
            { item = 1472023, text = "블러드 기간틱#k - 도적 레벨. 60#b", mats = { { 1472022, 1 }, { 4011007, 1 }, { 4021000, 8 }, { 2012000, 10 } }, cost = 80000 },
            { item = 1472024, text = "사파이어 기간틱#k - 도적 레벨. 60#b", mats = { { 1472022, 1 }, { 4011007, 1 }, { 4021005, 8 }, { 2012002, 10 } }, cost = 80000 },
            { item = 1472025, text = "다크 기간틱#k - 도적 레벨. 60#b", mats = { { 1472022, 1 }, { 4011007, 1 }, { 4021008, 3 }, { 4000046, 5 } }, cost = 100000 },
        },
    },
}

local TEXT = {
    greet = "제가 바로 이 수리점의 주인이지요! 핫핫핫. 약간의 수수료만 주신다면 원하시는 서비스를 제공해 드리도록 하지요.#b",
    ask_qty = function(item)
        return "#t" .. item .. "# 아이템을 원하시나요? 얼마나 만들어보시고 싶으신가요?"
    end,
    confirm = function(item, qty)
        return "흐음. #t" ..
            item ..
            "# " ..
            qty ..
            "개를 제작하고 싶으신거군요? 좋아요. 약간의 수수료와 재료만 갖고오신다면 충분히 만들어 드릴 수 있답니다. 재료는 아래와 같아요. 아, 그리고 인벤토리 공간이 충분한지도 확인해 주세요.#b"
    end,
    unit = "개",
    no_meso = "서비스를 받을 메소가 부족하시네요.",
    failed = "재료가 부족하시거나 인벤토리 공간이 부족한 건 아니신지 확인해 보시겠어요?",
    done = "휴, 다 되었답니다. 여기 완성품이에요.",
    done_next = true,
}

return {
    on_click = function(me, npc)
        maker.run(me, npc, CATEGORIES, TEXT)
    end
}
