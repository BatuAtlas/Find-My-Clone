-- Table: public.GPS

-- DROP TABLE IF EXISTS public."GPS";

CREATE TABLE IF NOT EXISTS public."GPS"
(
    "user" bigint NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    lat double precision NOT NULL,
    lon double precision NOT NULL,
    CONSTRAINT "GPS_pkey" PRIMARY KEY ("user", "timestamp"),
    CONSTRAINT "GPS_user_fkey" FOREIGN KEY ("user")
        REFERENCES public."User" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."GPS"
    OWNER to postgres;
-- Index: idx_gps_brin

-- DROP INDEX IF EXISTS public.idx_gps_brin;

CREATE INDEX IF NOT EXISTS idx_gps_brin
    ON public."GPS" USING btree
    ("user" ASC NULLS LAST, "timestamp" DESC NULLS FIRST)
    TABLESPACE pg_default;