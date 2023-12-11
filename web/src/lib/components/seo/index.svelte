<script lang="ts">
    // export let url: string;
    export let title: string;
    export let description: string;
    export let ogType: string;
    export let ogImage: string;
    export let keywords: string;
</script>


<svelte:head>
    <meta charset="utf-8" />
    <link rel="icon" type="image/svg+xml" href="/favicon.svg">
    <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
    <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
    <link rel="manifest" href="/site.webmanifest">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <!-- <link rel="sitemap" type="application/xml" title="Sitemap" href="/sitemap.xml"> -->
    <meta name="description" content={description} />
    <!-- <link rel="canonical" href={`https://genprose.com${url}`} /> -->
    <!-- <meta property="og:url" content={`https://genprose.com${url}`} /> -->
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
    <!-- Google tag (gtag.js) -->
    <!-- <script async src="https://www.googletagmanager.com/gtag/js?id=G-2XKMY8V0DX"></script>
    <script>
        window.dataLayer = window.dataLayer || [];
        function gtag() { dataLayer.push(arguments); }
        gtag('js', new Date());

        gtag('config', 'G-2XKMY8V0DX');
    </script> -->
</svelte:head>