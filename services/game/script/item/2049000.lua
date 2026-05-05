-- Item name (String.wz/Consume.img.xml): 백의 주문서 1%

function on_scroll(me, scroll_slot, target_slot)
    local target = nil
    if target_slot < 0 then
        target = me:equipped(target_slot)
    else
        target = me:item(InventoryType.Equipment, target_slot)
    end
    if target == nil then
        return false
    end
    local max = target:wz():enhance_chance()
    local used = target:enhance_count()
    local remain = target:enhance_chance()
    if max == nil or used == nil or remain == nil or used + remain >= max then
        return false
    end

    return me:enhance(scroll_slot, target_slot, nil, nil, function(target, scroll)
        local chance = target:enhance_chance()
        if chance == nil then
            return
        end
        local recover = scroll:wz():recover()
        if recover == nil then
            recover = 0
        end
        target:enhance_chance(chance + recover)
    end)
end
