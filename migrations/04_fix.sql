ALTER TABLE nodes ALTER COLUMN geom TYPE geometry(Point, 4326)
USING ST_Transform(ST_SetSRID(geom, 3857), 4326);

UPDATE roads r
SET weight_road = ST_Distance(
    ST_SetSRID(ST_MakePoint(n1.x, n1.y), 4326)::geography,
    ST_SetSRID(ST_MakePoint(n2.x, n2.y), 4326)::geography
)
FROM
    (SELECT id, ST_X(geom) as x, ST_Y(geom) as y FROM nodes) n1,
    (SELECT id, ST_X(geom) as x, ST_Y(geom) as y FROM nodes) n2
WHERE r.from_node = n1.id AND r.to_node = n2.id;