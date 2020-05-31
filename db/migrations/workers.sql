CREATE TABLE public."Workers"
(
    "ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    "Name" text NOT NULL,
    "LastSeen" timestamp with time zone NOT NULL,
    "Active" boolean NOT NULL,
    PRIMARY KEY ("ID")
);

ALTER TABLE public."Workers"
    ADD CONSTRAINT "Hostname" UNIQUE ("Name");