CREATE TABLE public."Recipients"
(
    "ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    "Ticket" bigint NOT NULL,
    "Type" int NOT NULL,
    "User" bigint,
    "Mail" text,
    PRIMARY KEY ("ID"),
    CONSTRAINT "FK_Ticket" FOREIGN KEY ("Ticket")
        REFERENCES public."Tickets" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT "FK_User" FOREIGN KEY ("User")
        REFERENCES public."Users" ("ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);