ALTER TABLE public."Tasks"
    ADD COLUMN "Recipient" text;
ALTER TABLE public."Tasks"
    ADD COLUMN "Ticket" bigint;
ALTER TABLE public."Tasks"
    ADD CONSTRAINT "FK_Ticket" FOREIGN KEY ("Ticket")
    REFERENCES public."Tickets" ("ID") MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;