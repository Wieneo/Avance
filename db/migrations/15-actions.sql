BEGIN;
CREATE TABLE public."Actions"
(
    "ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    "Type" integer NOT NULL,
    "Title" text NOT NULL,
    "Content" text NOT NULL,
    "Ticket" bigint NOT NULL,
    "IssuedAt" timestamp with time zone NOT NULL,
    "IssuedBy" bigint NOT NULL,
    PRIMARY KEY ("ID"),
    CONSTRAINT "FK_Ticket" FOREIGN KEY ("Ticket")
        REFERENCES public."Tickets" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
        NOT VALID,
    CONSTRAINT "FK_User" FOREIGN KEY ("IssuedBy")
        REFERENCES public."Users" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
        NOT VALID
);
END;