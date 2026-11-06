-- Converts apertures.actual_width / actual_height from BIGINT to BIGINT[] with one element per unit (length = quantity).
BEGIN;

CREATE TABLE IF NOT EXISTS apertures_actual_dims_backup AS
SELECT id, quantity, actual_width, actual_height FROM apertures;

ALTER TABLE apertures
    ALTER COLUMN actual_width  DROP DEFAULT,
    ALTER COLUMN actual_height DROP DEFAULT;

ALTER TABLE apertures
    ALTER COLUMN actual_width  TYPE bigint[]
        USING array_fill(COALESCE(actual_width, 0)::bigint,  ARRAY[GREATEST(COALESCE(quantity, 0), 0)::int]),
    ALTER COLUMN actual_height TYPE bigint[]
        USING array_fill(COALESCE(actual_height, 0)::bigint, ARRAY[GREATEST(COALESCE(quantity, 0), 0)::int]);

COMMIT;
