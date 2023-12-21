## Deploying the landing page to cloudflare
- deploy `api` service by running the command `wrangler deploy`.
- deploy `www` service by running the command `wrangler pages publish dist`.
- Head to `https://dash.cloudflare.com/3a670acb150d4c26439c8039e539813f/workers-and-pages` -> create application -> create webpage.

# To add domain names
- First make sure to register the domain in the `websites` tab which requires you to change nameservers to cloudflare.
- Head back to `workers-and-pages` click on the pages application deployed and select `custom domains`.
- To add subdomain to `api-worker`, go into the appln and select `triggers` -> `custom domains`.