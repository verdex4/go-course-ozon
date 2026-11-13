SELECT
    DISTINCT(name)
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE o.item = 'наушники'
ORDER BY u.name;