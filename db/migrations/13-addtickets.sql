BEGIN;

CREATE TABLE public."Tickets"
(
    "ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    "Title" text COLLATE pg_catalog."default" NOT NULL,
    "Description" text COLLATE pg_catalog."default" NOT NULL,
    "Queue" bigint NOT NULL,
    "Owner" bigint,
    "Severity" bigint NOT NULL,
    "Status" bigint NOT NULL,
    "CreatedAt" timestamp with time zone NOT NULL,
    "LastModified" timestamp with time zone NOT NULL,
    "StalledUntil" timestamp with time zone,
    "Meta" json NOT NULL,
    CONSTRAINT "Tickets_pkey" PRIMARY KEY ("ID"),
    CONSTRAINT "FK_Owner" FOREIGN KEY ("Owner")
        REFERENCES public."Users" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT,
    CONSTRAINT "FK_Queue" FOREIGN KEY ("Queue")
        REFERENCES public."Queue" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT,
    CONSTRAINT "FK_Severity" FOREIGN KEY ("Severity")
        REFERENCES public."Severities" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT "FK_Status" FOREIGN KEY ("Status")
        REFERENCES public."Statuses" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public."Tickets"
    OWNER to tixter;

END;