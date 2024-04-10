local function bfs(start, target)
    local queue = {start}
    local visited = {}
    while #queue > 0 do
        local current = table.remove(queue, 1)
        if not visited[current] then
            visited[current] = true
            local members = redis.call('SMEMBERS', current)
            for _, member in pairs(members) do
                if member == target then
                    return 1
                else
                    table.insert(queue, member)
                end
            end
        end
    end
    return 0
end
return bfs(KEYS[1], KEYS[2])