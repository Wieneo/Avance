ALTER TABLE public."Tasks"
    ADD COLUMN "Interval" integer;

ALTER TABLE public."Tasks"
    ADD COLUMN "LastRun" timestamp with time zone;

CREATE OR REPLACE FUNCTION "GetTask"() RETURNS bigint
    AS '
    DECLARE
    taskid bigint := 0;
    BEGIN
    LOCK TABLE "Tasks" IN ROW EXCLUSIVE MODE;
    taskid := (SELECT "ID" FROM "Tasks" WHERE "Status" = 0 AND (EXTRACT(EPOCH FROM (NOW() - "LastRun")) >= "Interval" OR ("Interval" IS NOT NULL AND "LastRun" IS NULL) OR ("Interval" IS NULL AND "LastRun" IS NULL)) ORDER BY "ID" LIMIT 1 FOR UPDATE);
    UPDATE "Tasks" SET "Status" = 1 WHERE "ID" = taskid;
    RETURN taskid;
    END;'
    LANGUAGE 'plpgsql';