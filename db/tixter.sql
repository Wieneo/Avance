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
-- Name: Users; Type: TABLE; Schema: public; Owner: tixter
--

CREATE TABLE public."Users" (
    "ID" bigint NOT NULL,
    "Username" text NOT NULL,
    "Password" text NOT NULL,
    "Mail" text NOT NULL
);


ALTER TABLE public."Users" OWNER TO tixter;

--
-- Name: Users_ID_seq; Type: SEQUENCE; Schema: public; Owner: tixter
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
-- Name: Version; Type: TABLE; Schema: public; Owner: tixter
--

CREATE TABLE public."Version" (
    "Schema" integer NOT NULL
);


ALTER TABLE public."Version" OWNER TO tixter;

--
-- Data for Name: Users; Type: TABLE DATA; Schema: public; Owner: tixter
--

COPY public."Users" ("ID", "Username", "Password", "Mail") FROM stdin;
1	Admin	$2y$12$ghd7daO/gSufBzmFzJZSYuv7HRplIia1kktycoELfEfGhiIromR1u	johann@gnaucke.com
\.


--
-- Data for Name: Version; Type: TABLE DATA; Schema: public; Owner: tixter
--

COPY public."Version" ("Schema") FROM stdin;
2
\.


--
-- Name: Users_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: tixter
--

SELECT pg_catalog.setval('public."Users_ID_seq"', 1, true);


--
-- Name: Users Users_pkey; Type: CONSTRAINT; Schema: public; Owner: tixter
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "Users_pkey" PRIMARY KEY ("ID");


--
-- PostgreSQL database dump complete
--

