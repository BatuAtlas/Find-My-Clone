-- Table: public.Usersettings

-- DROP TABLE IF EXISTS public."Usersettings";

CREATE TABLE IF NOT EXISTS public."Usersettings"
(
    "user" bigint NOT NULL,
    mail character varying(320) COLLATE pg_catalog."default" NOT NULL,
    password character varying(64) COLLATE pg_catalog."default",
    notifications jsonb DEFAULT '[]'::jsonb,
    CONSTRAINT "Usersettings_pkey" PRIMARY KEY ("user"),
    CONSTRAINT unique_mail UNIQUE (mail),
    CONSTRAINT "Usersettings_user_fkey" FOREIGN KEY ("user")
        REFERENCES public."User" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."Usersettings"
    OWNER to postgres;