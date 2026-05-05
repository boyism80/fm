-- Item name (String.wz/Consume.img.xml): 두손둔기 명중률 주문서 10%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Weapon, WeaponType.Blunt2H, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end
