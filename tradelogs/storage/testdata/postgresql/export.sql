--
-- PostgreSQL database dump
--

-- Dumped from database version 10.11
-- Dumped by pg_dump version 10.11

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

--
-- Data for Name: reserve; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

INSERT INTO public.reserve (id, address, name) VALUES (1, '0x63825c174ab367968EC60f061753D3bbD36A0D8F', '');
INSERT INTO public.reserve (id, address, name) VALUES (2, '0x0000000000000000000000000000000000000000', '');
INSERT INTO public.reserve (id, address, name) VALUES (3, '0x21433Dec9Cb634A23c6A4BbcCe08c83f5aC2EC18', '');
INSERT INTO public.reserve (id, address, name) VALUES (4, '0x56e37b6b79d4E895618B8Bb287748702848Ae8c0', '');


--
-- Data for Name: token; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

INSERT INTO public.token (id, address, symbol) VALUES (1, '0x595832F8FC6BF59c85C527fEC3740A1b7a361269', '');
INSERT INTO public.token (id, address, symbol) VALUES (2, '0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE', '');
INSERT INTO public.token (id, address, symbol) VALUES (3, '0xfe5F141Bf94fE84bC28deD0AB966c16B17490657', '');
INSERT INTO public.token (id, address, symbol) VALUES (6, '0xdd974D5C2e2928deA5F71b9825b8b646686BD200', '');
INSERT INTO public.token (id, address, symbol) VALUES (7, '0xF433089366899D83a9f26A773D59ec7eCF30355e', '');
INSERT INTO public.token (id, address, symbol) VALUES (12, '0x89d24A6b4CcB1B6fAA2625fE562bDD9a23260359', '');
INSERT INTO public.token (id, address, symbol) VALUES (13, '0x0F5D2fB29fb7d3CFeE444a200298f468908cC942', '');
INSERT INTO public.token (id, address, symbol) VALUES (16, '0x4156D3342D5c385a87D264F90653733592000581', '');
INSERT INTO public.token (id, address, symbol) VALUES (18, '0x514910771AF9Ca656af840dff83E8264EcF986CA', '');
INSERT INTO public.token (id, address, symbol) VALUES (19, '0xC5bBaE50781Be1669306b9e001EFF57a2957b09d', '');
INSERT INTO public.token (id, address, symbol) VALUES (22, '0x23Ccc43365D9dD3882eab88F43d515208f832430', '');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

INSERT INTO public.users (id, address, "timestamp") VALUES (1, '0x64C0372Ebc1F812398Bf3e475117Fd48D5EfB035', '2018-10-11 08:45:11+00');
INSERT INTO public.users (id, address, "timestamp") VALUES (2, '0x96159285B88d578BEb3e241E5F8a1bBe31e4C683', '2018-10-11 08:48:02+00');
INSERT INTO public.users (id, address, "timestamp") VALUES (3, '0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F', '2018-10-11 08:48:26+00');
INSERT INTO public.users (id, address, "timestamp") VALUES (8, '0x17D79F467243c5DB655282Ce6187127c42986413', '2018-10-11 08:59:32+00');
INSERT INTO public.users (id, address, "timestamp") VALUES (9, '0xA41983E9baa92bA284A75dB1dB2bCbAfb763B033', '2018-10-11 08:59:32+00');
INSERT INTO public.users (id, address, "timestamp") VALUES (11, '0x0826601F28B691CEEa2Be05EC1c922Ea0eC2d82D', '2018-10-11 09:04:41+00');


--
-- Data for Name: wallet; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

INSERT INTO public.wallet (id, address, name) VALUES (1, '0x0000000000000000000000000000000000000000', '');
INSERT INTO public.wallet (id, address, name) VALUES (3, '0x0000000000000000000000000000000000631762', '');
INSERT INTO public.wallet (id, address, name) VALUES (5, '0x0000000000000000000000000000000000631776', '');
INSERT INTO public.wallet (id, address, name) VALUES (6, '0x000000000000000000000000000000000063177F', '');
INSERT INTO public.wallet (id, address, name) VALUES (8, '0xDECAF9CD2367cdbb726E904cD6397eDFcAe6068D', 'Myetherwallet');
INSERT INTO public.wallet (id, address, name) VALUES (9, '0x0000000000000000000000000000000000631793', '');


--
-- Data for Name: tradelogs; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (1, '2018-10-11 08:45:11+00', 6494045, '0xce56f715862b458bfe9a2fc7059707efee5827f36f5be01edddc013277aa99fa', 1.25605787600702801, 1.25605787600702801, 1, 1, 2, 1, 2, 1424.55751099999998, 1.25605787600702801, 1, 1.74278030295975128, 0, 0, 0, 'KyberSwap', NULL, NULL, 225.662016177461652, 'coingecko', 91, false, true, '0x8177573B5557e3a2213d4aEc44abe7BaEF6D737D', '0x64C0372Ebc1F812398Bf3e475117Fd48D5EfB035', 537709, 8.0000000000000005e-09, 0.00430167200000000031);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (2, '2018-10-11 08:48:02+00', 6494053, '0xce85c30aaffd7d13722a6e3f3b4014c575a4c6b5f9973d66a0547e76ce4aba41', 0.594180157001091391, 0.594180157001091391, 2, 3, 2, 1, 2, 3170.0215575518755, 0.594180157001091391, 1, 0.8244249678390142, 0, 0, 0, 'KyberSwap', NULL, NULL, 225.662016177461652, 'coingecko', 56, false, true, '0x96159285B88d578BEb3e241E5F8a1bBe31e4C683', '0x96159285B88d578BEb3e241E5F8a1bBe31e4C683', 159886, 1.00000000000000002e-08, 0.00159885999999999991);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (3, '2018-10-11 08:48:26+00', 6494055, '0xeabbfa4e070cd23a9be7b8bbec47ee8e07f037440d72faa983178642b9d491fc', 0.00100000000000000002, 0.00100000000000000002, 3, 2, 6, 2, 1, 0.00100000000000000002, 0.532114987602218537, 3, 0, 0.00138750000000000006, 0, 0, 'KyberSwap', '117.4.192.22', 'VN', 225.662016177461652, 'coingecko', 18, false, true, '0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F', '0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F', 217132, 1.22999999999999993e-08, 0.00267072360000000005);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (4, '2018-10-11 08:54:03+00', 6494072, '0x963fe56aaa2add83bcd8ed1326c0b150c738afb462f4711a25b8ffa83d464d6a', 1.00454822773370611, 1.00454822773370611, 1, 7, 2, 3, 2, 341.230954139999994, 1.00454822773370611, 1, 0.669029119670648376, 0, 0, 0, 'KyberSwap', NULL, NULL, 225.662016177461652, 'coingecko', 126, false, false, '0x91BbA529a4e469758CD3832305586fD1e8161eDd', '0x64C0372Ebc1F812398Bf3e475117Fd48D5EfB035', 539072, 8.0000000000000005e-09, 0.00431257600000000026);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (5, '2018-10-11 08:54:39+00', 6494074, '0x190eb1f5432253005e95f5291d9acf8c8c13b752c1c80cd161271d07ba5a7b84', 0.00111867653591824802, 0.00111867653591824802, 3, 6, 2, 1, 2, 0.599999999999999978, 0.00111867653591824802, 5, 0.0015521636935865691, 0, 0, 0, 'KyberSwap', '117.4.192.22', 'VN', 225.662016177461652, 'coingecko', 22, false, false, '0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F', '0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F', 224625, 1.22999999999999993e-08, 0.00276288750000000007);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (6, '2018-10-11 08:56:37+00', 6494088, '0x3ef81b263846d9a8674938ad23c346685e34bea48cd5cb3110727498b4f298f8', 0.00223495740203999991, 0.00111747870101999995, 3, 6, 12, 3, 3, 0.599999999999999978, 0.220337540440739582, 6, 0.000744240814879320017, 0.000744240814879320017, 0, 0, 'KyberSwap', '117.4.192.22', 'VN', 225.662016177461652, 'coingecko', 17, false, false, '0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F', '0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F', 352384, 1.22999999999999993e-08, 0.00433432320000000024);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (7, '2018-10-11 08:56:37+00', 6494088, '0x00b0343a905007b8a1b89c73b19a85b3a7307286841570c14643c81644ad8b39', 1.00519751875436181, 1.00519751875436181, 1, 13, 2, 1, 2, 3061.07267668513168, 1.00519751875436181, 1, 1.39471155727167684, 0, 0, 0, 'KyberSwap', NULL, NULL, 225.662016177461652, 'coingecko', 46, false, false, '0xEF50EeD70D0Ff96354368749877BF51BFF73E4eb', '0x64C0372Ebc1F812398Bf3e475117Fd48D5EfB035', 629667, 8.0000000000000005e-09, 0.0050373359999999999);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (8, '2018-10-11 08:59:32+00', 6494106, '0x2c7c31a01431e6d9abb1769b87b99f8c0d6f660e7fcada13276a2f7cce52abfa', 12, 12, 8, 2, 16, 2, 1, 12, 4397.85565915000006, 8, 0, 9.99000000000000021, 0, 6.66000000000000014, 'ThirdParty', NULL, NULL, 225.662016177461652, 'coingecko', 5, false, true, '0x17D79F467243c5DB655282Ce6187127c42986413', '0x17D79F467243c5DB655282Ce6187127c42986413', 183969, 4.10000000000000032e-08, 0.00754272899999999975);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (9, '2018-10-11 08:59:32+00', 6494106, '0xb05789dd7a2f5431fb8ff70e15ba89732b6f9eb8382debb2f00f5fd42f27d960', 0.471510913161030443, 0.471510913161030443, 9, 2, 18, 2, 1, 0.471510913161030443, 295.712956760479472, 9, 0, 0.65422139201092977, 0, 0, 'KyberSwap', '93.241.201.194', 'DE', 225.662016177461652, 'coingecko', 48, false, true, '0xA41983E9baa92bA284A75dB1dB2bCbAfb763B033', '0xA41983E9baa92bA284A75dB1dB2bCbAfb763B033', 189854, 4.6999999999999999e-09, 0.000892313799999999995);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address, gas_used, gas_price, transaction_fee) VALUES (10, '2018-10-11 09:04:26+00', 6494129, '0x31b7bdc67dcd63a6175086eb34ce37559f7bca7788d6ebf08f193fdddeb92e45', 1.00617464264818923, 1.00617464264818923, 1, 19, 2, 1, 2, 3125.2175400000001, 1.00617464264818923, 1, 1.39606731667436268, 0, 0, 0, 'KyberSwap', NULL, NULL, 225.662016177461652, 'coingecko', 104, false, false, '0x8177573B5557e3a2213d4aEc44abe7BaEF6D737D', '0x64C0372Ebc1F812398Bf3e475117Fd48D5EfB035', 542417, 8.0000000000000005e-09, 0.00433933600000000006);
INSERT INTO public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, original_eth_amount, user_address_id, src_address_id, dst_address_id, src_reserve_address_id, dst_reserve_address_id, src_amount, dst_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender, receiver_address) VALUES (11, '2018-10-11 09:04:41+00', 6494131, '0xd847a7247826203560a393af2e0a64e732e9c8bdc817bd92758d2ab26804c659', 0.0500000000000000028, 0.0500000000000000028, 11, 2, 22, 2, 4, 0.0500000000000000028, 489.501869250098025, 1, 0, 0.0693750000000000061, 0, 0, 'KyberSwap', NULL, NULL, 225.662016177461652, 'coingecko', 5, false, true, '0x0826601F28B691CEEa2Be05EC1c922Ea0eC2d82D', '0x0826601F28B691CEEa2Be05EC1c922Ea0eC2d82D');


--
-- Data for Name: big_tradelogs; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--



--
-- Name: big_tradelogs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.big_tradelogs_id_seq', 1, false);


--
-- Name: reserve_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.reserve_id_seq', 5, true);


--
-- Name: token_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.token_id_seq', 23, true);


--
-- Name: tradelogs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.tradelogs_id_seq', 11, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.users_id_seq', 11, true);


--
-- Name: wallet_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.wallet_id_seq', 11, true);


--
-- PostgreSQL database dump complete
--

