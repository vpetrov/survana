/* survana-request.js
Implements common request functionality (jsonp, etc)
*/

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana) {

    function post_json(url, data, success, error) {

        var req = new XMLHttpRequest(),
            post_data = JSON.stringify(data),
            result_json;

        function on_request_load() {
            console.log("on_request_load", req.readyState, arguments);
        }

        function on_request_loadend() {
            console.log("on_request_loadend", req.readyState, arguments);
            if (req.readyState === XMLHttpRequest.DONE) {
                //on OK
                if (req.status === 200) {
                    try {
                        result_json = JSON.parse(req.responseText);
                    } catch (e) {
                        console.log("Survana.Request: JSON.parse() failed", e, "on", req.responseText);
                        error && error(e);
                        return
                    }
                    success && success(result_json);
                } else {
                    error && error(new Error(req.statusText));
                }

            }
        }

        function on_request_change() {
            console.log("on_request_change", req.readyState, arguments);
        }

        req.onload = on_request_load;
        req.onloadend = on_request_loadend;
        req.onreadystatechange = on_request_change;
        req.onerror = error;
        req.onabort = error;

        req.open("POST", url, true);
        req.send(post_data);
    }

    //API
    Survana.Request = {
        'PostJSON': post_json
    };
}(window.Survana));