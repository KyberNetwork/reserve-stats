---
title: API Reference

language_tabs: # must be one of https://git.io/vQNgJ
  - shell

toc_footers:
  - <a href='#'>Sign Up for a Developer Key</a>
  - <a href='https://github.com/lord/slate'>Documentation Powered by Slate</a>

includes:
  - tradelogs/trade_logs
  - tradelogs/trade_summary
  - tradelogs/asset_volume
  - tradelogs/reserve_volume
  - tradelogs/wallet_fee
  - users/users
  - users/user_list
  - price-analytics-data/price_analytics_data
  - errors

search: true
---

# Introduction


# Authentication
Authentication follow: https://tools.ietf.org/html/draft-cavage-http-signatures-10

Required headers:

- **Digest**
- **Authorization**
- **Signature**
- **Nonce**

