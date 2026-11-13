SELECT
    u.name,
    COALESCE(SUM(o.price), 0) as total_cost -- количества нет
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id
ORDER BY total_cost DESC;