CREATE TABLE places (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    category TEXT,
    osm_id BIGINT,
    geom GEOMETRY(Point, 4326)
);

CREATE TABLE nodes (
    id BIGINT PRIMARY KEY,
    geom GEOMETRY(Point, 4326)
);

CREATE TABLE roads (
    id BIGSERIAL PRIMARY KEY,
    from_node BIGINT,
    to_node BIGINT,
    weight_road FLOAT,
    oneway TEXT,
    highway TEXT,
    FOREIGN KEY (from_node) REFERENCES nodes(id),
    FOREIGN KEY (to_node) REFERENCES nodes(id)
);