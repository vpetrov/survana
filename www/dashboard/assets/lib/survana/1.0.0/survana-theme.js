/* survana.js

 Defines an API for all survana-enabled client-side modules. Auto-detects the path to survana.js, stores a list of all
 available questionnaire engines, and defines methods for generating HTML from JSON definition of forms, as well as
 utility methods for loading scripts and changing Survana themes.
 */

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana, document) {

    /** Generates HTML from a questionnaire's description
     * @param form The JSON description of a questionnaire
     * @returns {DocumentFragment} An HTML rendering in a DocumentFragment object
     */
    function questionnaire(form) {

        if (!form || !form['fields'] || !Survana.Theme) {
            return null
        }

        var Q = new Survana.Theme.Engine[Survana.Theme.Id](document);

        // parses a list of fields
        var questionnaire = document.createDocumentFragment(),
            form_el = Q.form(form),
            i,
            elem,
            nfields = form.fields.length;

        //loop over all fields
        for (i = 0; i < nfields; ++i) {
            elem = Q.question(form.fields[i], form_el);
            if (elem) {
                form_el.appendChild(elem);
            }
        }

        questionnaire.appendChild(form_el);

        return questionnaire;
    }

    /** Loads and updates the current theme. Questionnaires should be re-rendered
     * @param theme_id  {String}    ID of the theme to load
     * @param success   {Function}  The success callback, called when the theme file was loaded
     * @param error     {Function}  The error callback in case of a network failure
     */
    function set_theme(theme_id, success, error) {

        //load the theme dynamically
        if (!Survana.Theme.Engine[theme_id]) {
            console.log('Loading theme', theme_id);

            load_script(Survana.Theme.Path + theme_id + '/survana-' + theme_id + '.js',
                function () {
                    Survana.Theme.Id = theme_id;
                    Survana.Theme.Name = Survana.Theme.Engine[theme_id].Name;
                    Survana.Theme.Version = Survana.Theme.Engine[theme_id].Version;
                    Survana.Theme.Current = Survana.Theme.Engine[theme_id];
                    success && success(theme_id);
                },
                error);
        } else {
            //if the theme is already available, set it as the current theme and call the success function
            Survana.Theme.Id = theme_id;
            Survana.Theme.Name = Survana.Theme.Engine[theme_id].Name;
            Survana.Theme.Version = Survana.Theme.Engine[theme_id].Version;
            Survana.Theme.Current = Survana.Theme.Engine[theme_id];
            success && success(theme_id);
        }
    }

    //API
    Survana.Theme = {
        //paths
        Path: Survana.ScriptPath + "theme/",
        //properties
        Engine: {},
        Id: null,
        Name: null,
        Version: null,

        //methods
        Questionnaire: questionnaire,
        SetTheme: set_theme,

        //Switches
        DesignerMode: true
    };
})(window.Survana, document);
