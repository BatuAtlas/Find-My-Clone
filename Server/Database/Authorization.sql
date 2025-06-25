-- Table: public.Authorization

-- DROP TABLE IF EXISTS public."Authorization";

CREATE TABLE IF NOT EXISTS public."Authorization"
(
    token character varying(64) COLLATE pg_catalog."default" NOT NULL,
    "user" bigint NOT NULL,
    expires timestamp with time zone,
    CONSTRAINT "Authorization_pkey" PRIMARY KEY (token)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."Authorization"
    OWNER to postgres;