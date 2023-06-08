ALTER TABLE sources 
ADD COLUMN new_vids TEXT[],
ADD COLUMN todays_digest TEXT;

---- create above / drop below ----

ALTER TABLE sources 
DROP COLUMN new_vids,
DROP COLUMN todays_digest;
