-- Item name (String.wz/Consume.img.xml): 망토 체력 주문서 60%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Cape, nil, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end
