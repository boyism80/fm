local M = {}

local function with_commas(fare)
    local text = tostring(fare)
    if #text <= 3 then
        return text
    end
    return text:sub(1, #text - 3) .. "," .. text:sub(-3)
end

-- cfg.routes rows are { map, fare, beginner fare }.
function M.run(me, npc, cfg)
    local greeting = "안녕하세요~! " .. cfg.town
        .. " 중형택시입니다. 다른 마을로 안전하고 빠르게 이동하고 싶으신가요? 그렇다면 저희 택시를 이용해 보세요. 싼 가격으로 원하시는 곳까지 친절하게 모셔다 드리고 있습니다."
    if not me:dialog(npc, greeting, false, true) then
        return
    end

    local class = me:class()
    local beginner = class == Class.Beginner or class == Class.Noblesse
    local fare_column = beginner and 3 or 2
    local options = {}
    for i, route in ipairs(cfg.routes) do
        options[i] = "#m" .. route[1] .. "# (" .. with_commas(route[fare_column]) .. " 메소)"
    end
    local sel = me:dialog_list(npc, beginner and cfg.prompt_beginner or cfg.prompt, options)
    if sel == nil then
        return
    end

    local route = cfg.routes[sel]
    local fare = route[fare_column]
    if not me:dialog_yes_no(npc, "이곳에서 더 이상 볼일이 없으신 모양이로군요. 정말로 #b#m" .. route[1] .. "##k 마을로 이동하시겠습니까? 가격은 #b" .. with_commas(fare) .. " 메소#k" .. cfg.confirm_tail) then
        me:dialog(npc, "이 마을에도 볼거리가 가득하답니다. 다른 마을로 이동하고 싶어지면 언제든지 저희 택시를 이용해 주세요~", false, true)
        return
    end
    if me:exchange({ meso = fare }, nil) ~= ExchangeResult.OK then
        me:dialog(npc, "메소가 부족하시군요. 죄송하지만 요금을 지불하지 않으면 저희 택시를 이용하실 수 없습니다.", false, true)
        return
    end
    me:map(route[1])
end

return M
