<script lang="ts">
  const domainUrl = "https://app.getcloudfriendly.com";
  const defaultOgImage = `${domainUrl}/assets/images/og-image.png`;
  export let ogImage = defaultOgImage;
  export let description: string;
  export let title: string;
  export let url: string;
  export let ogType = "website";
  export let keywords: string;
</script>

<head>
  <meta charset="utf-8" />
  <link
    rel="icon"
    type="image/svg+xml"
    href="/assets/images/cloudfriendly-logo.svg"
  />
  <link rel="icon" type="image/png" href="/favicon-16x16.png" />
  <link rel="icon" type="image/x-icon" href="/favicon.ico" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link
    rel="sitemap"
    type="application/xml"
    title="Sitemap"
    href="/sitemap.xml"
  />
  <meta name="description" content={description} />
  <link rel="canonical" href={`${domainUrl}${url}`} />
  <meta property="og:url" content={`${domainUrl}${url}`} />
  <meta property="og:title" content={title} />
  <meta property="og:description" content={description} />
  <meta property="og:type" content={ogType || "website"} />
  <meta property="twitter:title" content={title} />
  <meta property="twitter:description" content={description} />
  <meta property="twitter:card" content="summary" />
  <meta property="twitter:creator" content="@hoaxy777" />
  <meta property="og:locale" content="en_US" />
  <meta property="keywords" content={keywords} />
  <meta property="og:image" content={ogImage} />
  <meta property="image" content={ogImage} />
  <title>{title}</title>
</head>
<slot />
