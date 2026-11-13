SELECT
    name,
    age
FROM users
WHERE age = (
    SELECT
        MAX(age) as m
    FROM users
);

-- WHERE age = MAX(age) писать нельзя