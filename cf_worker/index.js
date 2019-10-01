/**
 * Main worker entry point.
 */

//refer: https://github.com/pmeenan/cf-workers/blob/master/proxy/proxy.js

const CHECK_CODE = "USE YOUR CODE HERE, EXP: abcdefg";

/**
 * @param {any} body
 * @param {number} status
 * @param {Object<string, string>} headers
 */
function makeRes(body, status = 200, headers = {}) {
    return new Response(body, {status, headers})
}

addEventListener("fetch", event => {
    console.log(event.request.url);
    const rs = processRequest(event.request, event).catch(err => makeRes('cfworker error:\n' + err.stack, 502));
    event.respondWith(rs);
});

function filterKey(key) {
    const lkey = key.toLowerCase();
    if (lkey.startsWith("cf-") || lkey.startsWith("x-")) {
        return true
    }
    return false
}

/**
 * Handle all non-proxied requests. Send HTML or CSS on for further processing
 * and pass everything else through unmodified.
 * @param {*} request - Original request
 * @param {*} event - Original worker event
 */
async function processRequest(request, event) {
    // Proxy the request
    var proxyHeaders = new Headers();
    var rawHeaders = new Headers();
    var kvAll = new Headers(request.headers);
    for(const [k, v] of kvAll.entries()) {
        if(filterKey(k)) {
            continue
        }

        if(k.startsWith('__m_proxy_')) {
            proxyHeaders.set(k, v)
        } else {
            rawHeaders.set(k, v)
        }
    }
    rawHeaders.set('host', proxyHeaders.get('__m_proxy_host'));
    rawHeaders.set('referer', proxyHeaders.get('__m_proxy_referer'));

    if(!proxyHeaders.get('__m_proxy_check_code') || proxyHeaders.get('__m_proxy_check_code') != CHECK_CODE) {
        return makeRes("not found", 404)
    }

    let init = {
        method: request.method,
        redirect: "manual",
        headers: [...rawHeaders]
    };

    schema = proxyHeaders.get('__m_proxy_schema');
    if(!schema) {
        console.log("not found schema? url:" + request.url);
        schema = 'http'
    }

    const url = new URL(request.url);
    const proxyUrl = schema + ':/' + url.pathname + url.search;
    response = await fetch(proxyUrl, init);

    var newRspHeader = new Headers(response.headers);
    for(const [k, v] of kvAll.entries()) {
        newRspHeader.set('X-Debug-' + k, v)
    }
    response.headers = newRspHeader;
    return makeRes(response.body, response.status, newRspHeader);
}