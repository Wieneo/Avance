BEGIN;
UPDATE public."Users" SET "Permissions" = '{}';

ALTER TABLE public."Users"
    ALTER COLUMN "Permissions" SET NOT NULL;

END;