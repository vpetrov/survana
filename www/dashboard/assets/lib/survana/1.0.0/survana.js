(function(window, document, undefined) {
    'use strict';

    if (window.Survana === undefined) {
        window.Survana = {
            engine : {},
            theme: null,
            version: "1.0.0"
        }
    }

    var Survana = window.Survana,
        scriptName = 'survana.js',
        themePath = 'theme/',
        scriptPath;

    /* Generate HTML from a questionnaire description */
    Survana.Questionnaire = function (form) {

        if (form === undefined || form.fields === undefined || !Survana.theme) {
            return null
        }

        var Q = new Survana.engine[Survana.theme](document);

        // parses a list of fields
        var questionnaire = document.createDocumentFragment(),
            form_el = Q.form(form),
            i,
            elem,
            nfields = form.fields.length;

        //loop over all fields
        for (i = 0; i < nfields; ++i) {
            elem = Q.question(form.fields[i]);
            if (elem) {
                form_el.appendChild(elem);
            }
        }

        questionnaire.appendChild(form_el);
        return questionnaire;
    };

    /* Loads and updates the current theme. Questionnaires should be re-rendered */
    Survana.setTheme = function (theme_id, success, error) {

        //load the theme dynamically
        if (Survana.engine[theme_id] === undefined) {
            console.log('Loading theme', theme_id);

            loadScript(scriptPath + themePath + theme_id + '/survana-' + theme_id + '.js',
                function () {
                    Survana.theme = theme_id;
                    success();
                },
                error);
        } else {
            Survana.theme = theme_id;
            if (success) {
                success();
            }
        }
    };

    //Loads a js file using <script> tags appended to <body>
    function loadScript(path, success, error) {
        var script = document.createElement('script');
        script.setAttribute('src', path);
        script.setAttribute('type', 'text/javascript');

        script.onload = success;
        script.onerror = error;

        //append <script> to this document's <body>
        document.body.appendChild(script);
    }

    //detect script path
    var scripts = document.getElementsByTagName('script'),
        s,
        src;

    for (s in scripts) {
        src = scripts[s].src;

        if (src === undefined || !src.length) {
            continue;
        }

        if (src.indexOf(scriptName) >= 0) {
            scriptPath = src.substr(0, src.length - scriptName.length);
        }
    }

    if (!scriptPath) {
        throw Error("Failed to detect Survana's script path.");
    }
})(window, document);
