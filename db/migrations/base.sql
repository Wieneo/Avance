--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2
-- Dumped by pg_dump version 12.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: Groups; Type: TABLE; Schema: public; 
--

CREATE TABLE public."Groups" (
    "ID" bigint NOT NULL,
    "Name" text NOT NULL,
    "Permissions" json NOT NULL
);



--
-- Name: Groups_ID_seq; Type: SEQUENCE; Schema: public; 
--

ALTER TABLE public."Groups" ALTER COLUMN "ID" ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Groups_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: Projects; Type: TABLE; Schema: public; 
--

CREATE TABLE public."Projects" (
    "ID" bigint NOT NULL,
    "Name" text NOT NULL,
    "Description" text NOT NULL
);


--
-- Name: Projects_ID_seq; Type: SEQUENCE; Schema: public; 
--

ALTER TABLE public."Projects" ALTER COLUMN "ID" ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Projects_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: Queue; Type: TABLE; Schema: public; 
--

CREATE TABLE public."Queue" (
    "ID" bigint NOT NULL,
    "Name" text NOT NULL,
    "Project" bigint NOT NULL
);


--
-- Name: Queue_ID_seq; Type: SEQUENCE; Schema: public; 
--

ALTER TABLE public."Queue" ALTER COLUMN "ID" ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Queue_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: Users; Type: TABLE; Schema: public; 
--

CREATE TABLE public."Users" (
    "ID" bigint NOT NULL,
    "Username" text NOT NULL,
    "Password" text NOT NULL,
    "Mail" text NOT NULL,
    "Permissions" json NOT NULL,
    "Firstname" text NOT NULL,
    "Lastname" text NOT NULL
);



--
-- Name: Users_ID_seq; Type: SEQUENCE; Schema: public; 
--

ALTER TABLE public."Users" ALTER COLUMN "ID" ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Users_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: Version; Type: TABLE; Schema: public; 
--

CREATE TABLE public."Version" (
    "Schema" integer NOT NULL
);



--
-- Name: map_User_Group; Type: TABLE; Schema: public; 
--

CREATE TABLE public."map_User_Group" (
    "ID" bigint NOT NULL,
    "UserID" bigint NOT NULL,
    "GroupID" bigint NOT NULL
);



--
-- Name: map_User_Group_ID_seq; Type: SEQUENCE; Schema: public; 
--

ALTER TABLE public."map_User_Group" ALTER COLUMN "ID" ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."map_User_Group_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: Groups Groups_pkey; Type: CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."Groups"
    ADD CONSTRAINT "Groups_pkey" PRIMARY KEY ("ID");


--
-- Name: Projects Projects_pkey; Type: CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."Projects"
    ADD CONSTRAINT "Projects_pkey" PRIMARY KEY ("ID");


--
-- Name: Queue Queue_pkey; Type: CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."Queue"
    ADD CONSTRAINT "Queue_pkey" PRIMARY KEY ("ID");


--
-- Name: Users Users_pkey; Type: CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "Users_pkey" PRIMARY KEY ("ID");


--
-- Name: map_User_Group map_User_Group_pkey; Type: CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."map_User_Group"
    ADD CONSTRAINT "map_User_Group_pkey" PRIMARY KEY ("ID");


--
-- Name: None; Type: INDEX; Schema: public; 
--

CREATE INDEX "None" ON public."Queue" USING btree ("Project");


--
-- Name: map_User_Group FK_GroupID_Group; Type: FK CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."map_User_Group"
    ADD CONSTRAINT "FK_GroupID_Group" FOREIGN KEY ("GroupID") REFERENCES public."Groups"("ID") ON DELETE CASCADE;


--
-- Name: Queue FK_Queue_Projekt; Type: FK CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."Queue"
    ADD CONSTRAINT "FK_Queue_Projekt" FOREIGN KEY ("Project") REFERENCES public."Projects"("ID");


--
-- Name: map_User_Group FK_UserID_User; Type: FK CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public."map_User_Group"
    ADD CONSTRAINT "FK_UserID_User" FOREIGN KEY ("UserID") REFERENCES public."Users"("ID") ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

