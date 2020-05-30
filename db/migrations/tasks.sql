CREATE TABLE public."Tasks"
(
    "ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    "Task" json NOT NULL,
    "QueuedAt" timestamp with time zone NOT NULL,
    "Status" integer NOT NULL DEFAULT 0,
    "Type" integer NOT NULL,
    CONSTRAINT "Tasks_pkey" PRIMARY KEY ("ID")
);

CREATE FUNCTION "GetTask"() RETURNS bigint
    AS '
     DECLARE
    taskid bigint := 0;
    BEGIN
    LOCK TABLE "Tasks" IN ROW EXCLUSIVE MODE;
    taskid := (SELECT "ID" FROM "Tasks" WHERE "Status" = 0 ORDER BY "ID" LIMIT 1 FOR UPDATE);
    UPDATE "Tasks" SET "Status" = 1 WHERE "ID" = taskid;
    RETURN taskid;
    END;'
    LANGUAGE 'plpgsql';