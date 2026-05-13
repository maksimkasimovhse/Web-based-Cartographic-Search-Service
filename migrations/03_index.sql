CREATE INDEX ON places USING GIST(geom);
CREATE INDEX ON nodes USING GIST(geom);
CREATE INDEX ON roads (from_node);
CREATE INDEX ON roads (to_node);