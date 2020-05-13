BEGIN;
CREATE TABLE public."Statuses"
(
    "ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    "Enabled" boolean NOT NULL DEFAULT true,
    "Name" text COLLATE pg_catalog."default" NOT NULL,
    "DisplayColor" text COLLATE pg_catalog."default" NOT NULL,
    "TicketsVisible" boolean NOT NULL DEFAULT true,
    "Project" bigint NOT NULL,
    CONSTRAINT "Statuses_pkey" PRIMARY KEY ("ID"),
    CONSTRAINT "FK_Project" FOREIGN KEY ("Project")
        REFERENCES public."Projects" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;
END;