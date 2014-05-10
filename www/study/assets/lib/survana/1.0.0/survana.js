/* survana.js

Defines an API for all survana-enabled client-side modules. Auto-detects the path to survana.js, stores a list of all
available questionnaire engines, and defines methods for generating HTML from JSON definition of forms, as well as
utility methods for loading scripts and changing Survana themes.
*/

"use strict";

(function(window, document, undefined) {

    /** Generates HTML from a questionnaire's description
     * @param form The JSON description of a questionnaire
     * @returns {DocumentFragment} An HTML rendering in a DocumentFragment object
     */
    function questionnaire(form) {

        if (!form || !form['fields'] || !Survana.Theme) {
            return null
        }

        var Q = new window.Survana.Engine[window.Survana.Theme](document);

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
    }

    /** Loads a js file using <script> tags appended to <body>
     * @param path      {String}    Path to the Javascript file
     * @param success   {Function}  The success callback, the first argument is the onload event object
     * @param error     {Function}  The error callback, the first rgument is the onerror event object
     */
    function load_script(path, success, error) {
        var script = document.createElement('script');
        script.setAttribute('src', path);
        script.setAttribute('type', 'text/javascript');

        script.onload = success;
        script.onerror = error;

        //append <script> to this document's <body>
        document.body.appendChild(script);
    }

    /** Loads and updates the current theme. Questionnaires should be re-rendered
     * @param theme_id  {String}    ID of the theme to load
     * @param success   {Function}  The success callback, called when the theme file was loaded
     * @param error     {Function}  The error callback in case of a network failure
     */
    function set_theme(theme_id, success, error) {

        //load the theme dynamically
        if (!window.Survana.Engine[theme_id]) {
            console.log('Loading theme', theme_id);

            load_script(window.Survana.ScriptPath + window.Survana.ThemePath + theme_id + '/survana-' + theme_id + '.js',
                function () {
                    Survana.Theme = theme_id;
                    success && success(theme_id);
                },
                error);
        } else {
            //if the theme is already available, set it as the current theme and call the success function
            Survana.Theme = theme_id;
            success && success(theme_id);
        }
    }

    /** Searches for the path to a script specified by 'scriptName'
     * @param scriptName {String} The name of the script to search for.
     * @returns {String|null} Path to the 'scriptName', or null if not found
     */
    function detect_script_path(scriptName) {
        //detect script path
        var scripts = document.getElementsByTagName('script'),
            s,
            src;

        for (s in scripts) {
            if (!scripts.hasOwnProperty(s)) {
                continue;
            }

            src = scripts[s].src;
            if (!src) {
                continue;
            }

            //
            if (src.indexOf(scriptName) >= 0) {
                return src.substr(0, src.length - scriptName.length);
            }
        }

        return null;
    }

    //API
    window.Survana = {
        //paths
        ThemePath:  "theme/",
        ScriptPath: detect_script_path('survana.js'),

        //properties
        Engine :    {},
        Theme:      null,
        Version:    "1.0.0",

        //methods
        Questionnaire:  questionnaire,
        SetTheme:       set_theme,
        LoadScript:     load_script,

        //Switches
        DesignerMode: true
    };
})(window, document);
