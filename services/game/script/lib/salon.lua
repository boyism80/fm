-- Coupon salons: a look off the NPC's list is paid for with one coupon.

local M = {}

local FACES = {
    [0] = { 20000, 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20012, 20014, 20016, 20020, 20017 },
    [1] = { 21000, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21008, 21012, 21014, 21016, 21020, 21017 },
}

-- Returns the picked look, nil off the list, and whether the coupon was taken for it.
function M.pick(me, npc, prompt, styles, coupon)
    local pick = me:dialog_style(npc, prompt, styles)
    if pick == nil or styles[pick] == nil then
        return nil, false
    end
    return styles[pick], me:exchange({ item = { [coupon] = 1 } }) == ExchangeResult.OK
end

-- The faces keep the current lens color.
function M.faces(me)
    local face = me:face()
    local styles = {}
    for _, style in ipairs(me:gender() == 0 and FACES[0] or FACES[1]) do
        styles[#styles + 1] = style + face % 1000 - face % 100
    end
    return styles
end

function M.skin(me, npc, coupon)
    if me:dialog(npc, "아름다운 피부로 가꾸시고 싶으신가요? 제가 도와드리도록 하지요. 서비스를 이용하려면 #b#t" .. coupon ..
            "##k 아이템을 가져오셔야 합니다.", false, true) == false then
        return
    end
    local skin, paid = M.pick(me, npc, "원하시는 피부를 골라보세요~", { 0, 1, 2, 3, 4 }, coupon)
    if skin == nil then
        return
    end
    if not paid then
        me:dialog(npc, "죄송하지만 쿠폰을 가져오시지 않으면 피부관리를 해드릴 수 없답니다.")
        return
    end
    me:skin(skin)
    me:dialog(npc, "자~ 다 되었답니다. 아름다운 피부는 미의 기본.. 앞으로도 아름다운 피부를 가꾸시고 싶으시다면 저희 피부관리실을 찾아주세요.")
end

return M
