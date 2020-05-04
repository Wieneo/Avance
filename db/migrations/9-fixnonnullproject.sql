BEGIN;

UPDATE public."Projects" SET "Description" = 'Placeholder';

ALTER TABLE public."Projects"
    ALTER COLUMN "Description" SET NOT NULL;

END;