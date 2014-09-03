/* survana.js

 Defines an API for all Survana-enabled client-side modules.
 - Auto-detects the path to survana.js
 - Provides method for loading scripts dynamically
 - Defines debug methods (Assert, Error, Log, Warn)
 */

"use strict";

(function (window, document) {

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

    function assert() {
        console.assert.apply(console, arguments);
    }

    function log() {
        console && console.log.apply(console, arguments);
    }

    function warn() {
        console && console.warn.apply(console, arguments);
    }

    function error() {
        console && console.error.apply(console, arguments);
    }

    //Element Helpers

    //returns the value of an element based on its declared field type
    function field_value(element, field_type) {
        //prefer to use field_type as the type of the field represented by 'element'
        switch (field_type) {
            default:
                break; //todo: implement custom components
        }

        switch (element.tagName.toLowerCase()) {
            case 'input':
                switch (element.getAttribute('type')) {
                    case 'radio':
                    case 'checkbox':
                        if (element.checked) {
                            if (element.value === Survana.NO_ANSWER) {
                                return null;
                            }

                            return element.value;
                        }
                        return undefined;
                    default:
                        return element.value;
                }
                break;
            case 'button':
            case 'select':
                return element.value;
            default:
                return undefined;
        }
    }

    /**
     * Group HTMLForm elements by name attribute
     * @param form {HTMLFormElement} The HTML form to scan for elements
     * @returns {Object} key = name of element, value = array of elements with same name
     * @constructor
     */
    function group_elements(form) {
        var result = {},
            el,
            name;

        if (!form.elements) {
            return result;
        }

        for (var i = 0; i < form.elements.length; ++i) {
            el = form.elements[i];
            name = el.getAttribute('name');

            if (!name) {
                continue;
            }

            //skip input element with type 'hidden'
            if (el.getAttribute('type') === 'hidden') {
                continue;
            }

            if (result[name] === undefined) {
                result[name] = [el];
            } else {
                result[name].push(el);
            }
        }

        return result;
    }

    function values_from_group(group, field_type) {
        var result = [],
            value;

        if (!group) {
            return result;
        }

        for (var i = 0; i < group.length; ++i) {
            value = field_value(group[i], field_type);
            if (value === null) {
                result = null;
                break;
            }

            if (value !== undefined) {
                result.push(value);
            }
        }

        return result;
    }

    function no_answer(question_id) {
        var question = document.getElementById(question_id),
            children;

        if (!question) {
            return;
        }

        question.classList.add(Survana.NO_ANSWER);

        children = question.querySelectorAll('input,select');

        //mark all inputs for no validation
        for (var i = 0; i < children.length; ++i) {
            children[i].value = children[i].defaultValue;
            children[i].classList.add(Survana.NO_ANSWER);
        }
    }

    function get_matrix_group(field, groups) {
        var result = {},
            row,
            row_id,
            i;

        for (i = 0; i < field.rows.length; ++i) {
            row = field.rows[i];
            row_id = row.id || field.id + ":" + (i + 1);

            if (groups[row_id]) {
                result[row_id] = groups[row_id];
            } else {
                result[row_id] = null;
            }
        }

        return result;
    }

    /**
     * Returns all form fields grouped by name and their values as arrays
     * @param form_el {HTMLFormElement} The <form> element to parse
     */
    function get_form_fields(form_el, schemata) {
        var fields = {},
            field,
            groups = group_elements(form_el),
            group,
            row_id,
            i,
            j;

        //if no schema was provided, attempt to fetch it from Survana.Schema
        schemata = schemata || Survana.Schema[form.id];
        if (!schemata) {
            Survana.Error('No Schema found for form ' + form_el.id);
            return false;
        }

        for (i = 0; i < schemata.fields.length; ++i) {
            field = schemata.fields[i];

            //special case: matrix containers. This loops effectively unwraps the matrix
            //and treats each question as if it were top-level
            if (field.type === 'matrix') {
                //iterate over each matrix row
                for (j = 0; j < field.rows.length; ++j) {
                    row_id = field.rows[j].id;
                    fields[row_id] = values_from_group(groups[row_id]);
                }
                continue;
            }

            //all other fields
            group = groups[field.id];

            if (!group) {
                console.log('Skipping answers for question', field.id);
                continue;
            }
            if (group[0].classList.contains(Survana.NO_ANSWER)) {
                fields[field.id] = null;
            } else {
                fields[field.id] = values_from_group(group);
            }
        }

        return fields;
    }

    //API
    var Survana = {

        //'constants'
        NO_ANSWER: 'no-answer',

        //Properties
        ScriptPath: detect_script_path('survana.js'),
        Version: "1.0.0",

        //Methods
        LoadScript: load_script,
        FieldValue: field_value,
        GroupElements: group_elements,
        ValuesFromGroup: values_from_group,
        NoAnswer: no_answer,
        FormFields: get_form_fields,

        //Dev methods
        Assert: assert,
        Error: error,
        Log: log,
        Warn: warn,

        //Switches

        //Schema
        Schema: {}

    };

    window.Survana = window.Survana || {};

    for (var p in Survana) {
        if (!Survana.hasOwnProperty(p)) {
            continue;
        }

        if (window.Survana[p] === undefined || window.Survana[p] === null) {
            window.Survana[p] = Survana[p];
        }
    }


    Survana.OnFormLoad = function () {
        //read any baked-in form information
        var script_elements = document.querySelectorAll('script.schema');

        if (!script_elements.length) {
            return;
        }

        for (var i = 0; i < script_elements.length; ++i) {
            var script = script_elements[i],
                json_string = script.innerHTML,
                schemata;

            if (!json_string.length) {
                continue;
            }

            try {
                schemata = JSON.parse(json_string);
            } catch (e) {
                Survana.Error(e);
                continue;
            }

            Survana.Schema[schemata.id] = schemata;
        }
    };

    function on_dom_content_loaded() {
        document.removeEventListener('DOMContentLoaded', on_dom_content_loaded);
        Survana.OnFormLoad();
    }

    //register an onReady handler, i.e. $(document).ready(). Caveat: does not support older versions of IE
    document.addEventListener("DOMContentLoaded", on_dom_content_loaded);
})(window, document);
