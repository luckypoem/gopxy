/**
 * Main worker entry point.
 */
addEventListener("fetch", event => {
    // Fail-safe in case of an unhandled exception
    console.log(event.request.url);
    event.passThroughOnException();
    event.respondWith(processRequest(event.request, event));
});

/**
 * Handle all non-proxied requests. Send HTML or CSS on for further processing
 * and pass everything else through unmodified.
 * @param {*} request - Original request
 * @param {*} event - Original worker event
 */
async function processRequest(request, event) {
    // Proxy the request
    let init = {
        method: request.method,
        redirect: "manual",
        headers: [...request.headers]
    };
    const url = new URL(request.url);
    let proxyOrigin = url.origin;
    const proxyUrl = 'https:/' + url.pathname + url.search;
    let originalDomain = url.pathname.substr(1);
    const domainEnd = originalDomain.indexOf('/');
    if (domainEnd >= 0)
        originalDomain = originalDomain.substr(0, domainEnd);
    const response = await fetch(proxyUrl, init);
    if (response) {
        // Process test responses
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.indexOf("text/") !== -1) {
            let content = await response.text();
            let init = {
                method: request.method,
                headers: [...response.headers]
            };
            const newResponse = new Response(content, init);
            newResponse.headers.set('X-Debug-Path', url.pathname)
            newResponse.headers.set('X-Debug-Search', url.search)
            return newResponse;
        }
    }

    return response;
}