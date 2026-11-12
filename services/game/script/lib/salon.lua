-- Coupon salons: the NPC picks the look from its own list, the coupon pays for it, and the look changes.

local M = {}

local FACES = {
    [0] = { 20000, 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20012, 20014, 20016, 20020, 20017 },
    [1] = { 21000, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21008, 21012, 21014, 21016, 21020, 21017 },
}

-- A pick off the list leaves the coupon untouched.
local function apply(me, npc, coupon, cfg, prompt, styles, set)
    local pick = me:dialog_style(npc, prompt, styles)
    if pick == nil or styles[pick] == nil then
        return
    end
    if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
        me:dialog(npc, cfg.no_coupon)
        return
    end
    set(me, styles[pick])
    me:dialog(npc, cfg.done)
end

local function set_hair(me, style)
    me:hair(style)
end

function M.hair_greeting(cfg)
    return cfg.intro .. "#b#i" .. cfg.hair_coupon .. "# #t" .. cfg.hair_coupon .. "##k 또는 #b#i" .. cfg.color_coupon ..
        "# #t" .. cfg.color_coupon .. "##k " .. cfg.outro
end

-- The styles keep the current hair color.
function M.hair_style(me, npc, cfg)
    local color = me:hair() % 10
    local styles = {}
    for _, style in ipairs(me:gender() == 0 and cfg.male or cfg.female) do
        styles[#styles + 1] = style + color
    end
    apply(me, npc, cfg.hair_coupon, cfg, cfg.style_prompt, styles, set_hair)
end

-- The colors keep the current hair style.
function M.hair_color(me, npc, cfg)
    local base = math.floor(me:hair() / 10) * 10
    local styles = {}
    for _, offset in ipairs(cfg.colors) do
        styles[#styles + 1] = base + offset
    end
    apply(me, npc, cfg.color_coupon, cfg, cfg.color_prompt, styles, set_hair)
end

function M.hair(me, npc, cfg)
    local sel = me:dialog_list(npc, M.hair_greeting(cfg), {
        "머리 스타일 바꾸기",
        "머리 색깔 염색하기",
    })
    if sel == 1 then
        M.hair_style(me, npc, cfg)
    elseif sel == 2 then
        M.hair_color(me, npc, cfg)
    end
end

-- The faces keep the current lens color.
function M.face(me, npc, cfg)
    if me:dialog(npc, cfg.intro .. "#b#t" .. cfg.coupon .. "##k " .. cfg.outro, false, true) == false then
        return
    end
    local face = me:face()
    local styles = {}
    for _, style in ipairs(me:gender() == 0 and FACES[0] or FACES[1]) do
        styles[#styles + 1] = style + face % 1000 - face % 100
    end
    apply(me, npc, cfg.coupon, {
        no_coupon = cfg.no_coupon or "미안하지만 쿠폰 없이는 성형수술을 해 줄수 없다네.",
        done = cfg.done,
    }, cfg.prompt or "어떤 얼굴로 바꿔보고 싶은가?", styles, function(m, style)
        m:face(style)
    end)
end

function M.skin(me, npc, cfg)
    if me:dialog(npc, "아름다운 피부로 가꾸시고 싶으신가요? 제가 도와드리도록 하지요. 서비스를 이용하려면 #b#t" .. cfg.coupon ..
            "##k 아이템을 가져오셔야 합니다.", false, true) == false then
        return
    end
    apply(me, npc, cfg.coupon, {
        no_coupon = "죄송하지만 쿠폰을 가져오시지 않으면 피부관리를 해드릴 수 없답니다.",
        done = "자~ 다 되었답니다. 아름다운 피부는 미의 기본.. 앞으로도 아름다운 피부를 가꾸시고 싶으시다면 저희 피부관리실을 찾아주세요.",
    }, "원하시는 피부를 골라보세요~", { 0, 1, 2, 3, 4 }, function(m, skin)
        m:skin(skin)
    end)
end

return M
