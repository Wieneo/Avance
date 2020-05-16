-- This script was generated by a beta version of the Schema Diff utility in pgAdmin 4. 
-- This version does not include dependency resolution, and may require manual changes 
-- to the script to ensure changes are applied in the correct order.
-- Please report an issue for any failure with the reproduction steps. 
BEGIN;
-- WARNING:
-- The SQL statement below would normally be used to alter the datatype for the Enabled column, however,
-- the current datatype cannot be cast to the target datatype so this conversion cannot be made automatically.

 -- ALTER TABLE public."Severities"
 --     ALTER COLUMN "Enabled" TYPE boolean;

-- WARNING:
-- The SQL statement below would normally be used to alter the datatype for the Name column, however,
-- the current datatype cannot be cast to the target datatype so this conversion cannot be made automatically.

 -- ALTER TABLE public."Severities"
 --     ALTER COLUMN "Name" TYPE text COLLATE pg_catalog."default";

-- WARNING:
-- The SQL statement below would normally be used to alter the datatype for the DisplayColor column, however,
-- the current datatype cannot be cast to the target datatype so this conversion cannot be made automatically.

 -- ALTER TABLE public."Severities"
 --     ALTER COLUMN "DisplayColor" TYPE text COLLATE pg_catalog."default";

ALTER TABLE public."Severities"
    ADD COLUMN "Project" bigint;
ALTER TABLE public."Severities"
    ADD CONSTRAINT "FK_Severity_Project" FOREIGN KEY ("Project")
    REFERENCES public."Projects" ("ID") MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
    NOT VALID;

END;