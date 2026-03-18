SELECT
    sp.plan_id,
    sp.url,
    sp.created_at,
    p.title,
    p.description,
    p.plan_table
FROM scraped sp
JOIN plans p ON sp.plan_id = p.plan_id
WHERE EXISTS (
    SELECT 1
    FROM jsonb_array_elements(p.plan_table) AS row_elem
    WHERE jsonb_array_length(row_elem->'SubRows') > 0
);
