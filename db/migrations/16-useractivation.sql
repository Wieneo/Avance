BEGIN;
ALTER TABLE public."Users"
    ADD COLUMN "Active" boolean NOT NULL DEFAULT true;
END;