--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.14
-- Dumped by pg_dump version 10.7

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.addresses DROP CONSTRAINT IF EXISTS addresses_pkey;
ALTER TABLE IF EXISTS ONLY public.addresses DROP CONSTRAINT IF EXISTS addresses_address_key;
ALTER TABLE IF EXISTS public.addresses ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.addresses_id_seq;
DROP TABLE IF EXISTS public.addresses;
DROP EXTENSION IF EXISTS plpgsql;
DROP SCHEMA IF EXISTS public;
--
-- Name: public; Type: SCHEMA; Schema: -; Owner: reserve_stats
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO reserve_stats;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: reserve_stats
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: reserve_stats
--

CREATE TABLE public.addresses (
    id integer NOT NULL,
    address text NOT NULL,
    type text NOT NULL,
    description text,
    "timestamp" timestamp without time zone,
    last_updated timestamp without time zone NOT NULL
);


ALTER TABLE public.addresses OWNER TO reserve_stats;

--
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: reserve_stats
--

CREATE SEQUENCE public.addresses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.addresses_id_seq OWNER TO reserve_stats;

--
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reserve_stats
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

COPY public.addresses (id, address, type, description, "timestamp", last_updated) FROM stdin;
1	0x63825c174ab367968EC60f061753D3bbD36A0D8F	reserve	main reserve	2018-02-07 21:15:57	2019-03-28 06:27:09.061943
2	0x8bC3da587DeF887B5C822105729ee1D6aF05A5ca	pricing_operator	pricing operator 1	\N	2019-03-28 06:29:35.274457
3	0x9224016462B204C57Eb70e1D69652f60bcAF53A8	pricing_operator	pricing operator 2	2018-02-07 10:54:42	2019-03-28 06:29:55.826882
4	0x0A3D5C8894bBE1E9113e4eD6f0c3B0D4Fa6b131E	sanity_operator	sanity operator 1	\N	2019-03-28 06:30:44.283876
5	0xd0643BC0D0C879F175556509dbcEe9373379D5C3	sanity_operator	sanity operator 2	\N	2019-03-28 06:30:56.400864
6	0x44d34A119BA21A42167FF8B77a88F0Fc7BB2Db90	cex_deposit_address	binance deposit	\N	2019-03-28 06:31:46.325932
7	0x0C8fd73Eaf6089eF1B91231D0A07D0d2cA2b9d66	cex_deposit_address	huobi deposit	\N	2019-03-28 06:32:02.254027
\.


--
-- Name: addresses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.addresses_id_seq', 7, true);


--
-- Name: addresses addresses_address_key; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_address_key UNIQUE (address);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: reserve_stats
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM reserve_stats;
GRANT ALL ON SCHEMA public TO reserve_stats;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

