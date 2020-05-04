BEGIN;
-- WARNING:
-- The SQL statement below would normally be used to alter the datatype for the ID column, however,
-- the current datatype cannot be cast to the target datatype so this conversion cannot be made automatically.

 -- ALTER TABLE public."Projects"
 --     ALTER COLUMN "ID" TYPE bigint;

ALTER TABLE public."Projects"
    ADD COLUMN "Description" text COLLATE pg_catalog."default";

CREATE INDEX "None"
    ON public."Queue"("Project");

END;