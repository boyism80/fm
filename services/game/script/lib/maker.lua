local M = {}

local function recipe_text(recipe, text)
    if text.option ~= nil then
        return text.option(recipe)
    end
    if recipe.text ~= nil then
        return recipe.text
    end
    return "#z" .. recipe.item .. "##k" .. recipe.label
end

local function choose_recipe(me, npc, category, text)
    if category.confirm ~= nil then
        if not me:dialog_yes_no(npc, category.confirm) then
            return nil
        end
        return category.recipes[1]
    end

    local options = {}
    for i, recipe in ipairs(category.recipes) do
        options[i] = recipe_text(recipe, text)
    end
    local sel = me:dialog_list(npc, category.prompt, options)
    if sel == nil then
        return nil
    end
    return category.recipes[sel]
end

function M.craft(me, npc, category, recipe, text)
    local qty = 1
    if not category.equip then
        qty = me:dialog_number(npc, text.ask_qty(recipe.item, recipe), 1, 1, 100)
        if qty == nil then
            return
        end
    end

    local stimulator = nil
    if category.stimulator then
        stimulator = recipe.stimulator or text.stimulator
    end
    local cost_items = {}
    for _, mat in ipairs(recipe.mats) do
        cost_items[mat[1]] = mat[2] * qty
    end
    local meso = recipe.cost * qty
    if text.confirm ~= nil then
        local prompt = text.confirm(recipe.item, qty, recipe)
        if stimulator ~= nil and text.lists_stimulator then
            prompt = prompt .. "\r\n#i" .. stimulator .. "# #t" .. stimulator .. "# 1" .. text.unit
        end
        for _, mat in ipairs(recipe.mats) do
            prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] * qty .. text.unit
        end
        if meso > 0 then
            prompt = prompt .. "\r\n#i4031138# " .. meso .. " 메소"
        end
        if text.ask ~= nil then
            prompt = prompt .. text.ask
        end
        if not me:dialog_yes_no(npc, prompt) then
            if text.declined ~= nil then
                me:dialog(npc, text.declined)
            end
            return
        end
    end

    if me:meso() < meso then
        me:dialog(npc, text.no_meso)
        return
    end
    local cost = { item = cost_items }
    if meso > 0 then
        cost.meso = meso
    end
    local reward = { item = { [recipe.item] = (recipe.gain or 1) * qty } }
    if category.stimulator then
        if stimulator ~= nil then
            cost_items[stimulator] = 1
        end
        -- One stimulated make in ten breaks and takes the materials with it.
        if math.random(0, 9) == 0 then
            reward = nil
        else
            reward.random = true
        end
    end
    local code = me:exchange(cost, reward)
    if code == ExchangeResult.LackCapacity and text.no_space ~= nil then
        me:dialog(npc, text.no_space)
        return
    end
    if code ~= ExchangeResult.OK then
        me:dialog(npc, text.failed, false, text.failed_next == true)
        return
    end
    if reward == nil then
        me:dialog(npc, text.broken)
        return
    end
    me:dialog(npc, text.done, false, text.done_next == true)
end

function M.run(me, npc, categories, text)
    local category_options = {}
    for i, category in ipairs(categories) do
        category_options[i] = category.text
    end
    local category_sel = me:dialog_list(npc, text.greet, category_options)
    if category_sel == nil then
        return
    end
    local category = categories[category_sel]
    if category.info ~= nil then
        me:dialog(npc, category.info, false, true)
        return
    end

    local recipe = choose_recipe(me, npc, category, text)
    if recipe == nil then
        return
    end
    M.craft(me, npc, category, recipe, text)
end

return M
