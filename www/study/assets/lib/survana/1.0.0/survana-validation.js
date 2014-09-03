/*************
 * VALIDATION *
 *************/

"use strict";

window.Survana = window.Survana || {};

(function (document, Survana) {

    Survana.Validation = Survana.Validation || {};

    Survana.Validation.INVALID = 's-invalid';

    /* constraint: function (values, target, groups) { }
     @param values {Array} List of values to validate
     @param constraint value {*} The constraint value specified in the validation schema
     @param fields {Object} All known values in the form, grouped by their question ID
     */
    Survana.Validation.Constraints = {
        'equal': function (response, target, values) {
            var target_values,
                i;

            if (response === undefined || response === null) {
                return false;
            }

            target_values = values[target];

            //return true on skipped responses because no validation should be performed on skipped fields
            if (target_values === null) {
                return true;
            }

            if (!target_values) {
                Survana.Error('Field "' + target + '" does not exist.');
                return false;
            }

            if (response.length !== target_values.length) {
                return false;
            }

            //stringify all the values from 'target_values'
            for (i = 0; i < target_values.length; ++i) {
                target_values[i] = String(target_values[i]);
            }

            //search for each item from 'values' in 'target_values'
            //alternate implementation: sort both arrays and compare items with 1 loop
            for (i = 0; i < response.length; ++i) {
                //if even 1 item could not be found, the constrained failed
                if (target_values.indexOf(String(response[i])) < 0) {
                    return false;
                }
            }

            return true;
        },
        optional: function (values, is_optional) {
            if (is_optional) {
                return true;
            }

            if (!values || !values.length) {
                return false;
            }

            //check all values: if they're undefined, null or empty strings, then
            //the constraint wasn't satisfied
            for (var i = 0; i < values.length; ++i) {
                if (values[i] === undefined || values[i] === null || values[i] === "") {
                    return false;
                }
            }

            return true;
        },
        max: function (values, max_value) {
            if (values === undefined || !values.length) {
                return false;
            }

            max_value = parseFloat(max_value);
            if (isNaN(max_value) || (max_value === Infinity)) {
                return false;
            }

            var v;

            for (var i = 0; i < values.length; ++i) {
                v = parseFloat(values[i]);
                if (isNaN(v) || (v === Infinity) || (v > max_value)) {
                    return false;
                }
            }

            return true;
        },

        min: function (values, min_value) {
            if (values === undefined || !values.length) {
                return false;
            }

            min_value = parseFloat(min_value);
            if (isNaN(min_value) || (min_value === Infinity)) {
                return false;
            }

            var v;

            for (var i = 0; i < values.length; ++i) {
                //convert each value to float and compare it with min_value
                v = parseFloat(values[i]);
                if (isNaN(v) || (v === Infinity) || (v < min_value)) {
                    return false;
                }
            }

            return true;
        }
    };

    function invalid(field, constraint_name, constraint_value, src_input) {
        var message = (Survana.Validation.Messages[constraint_name] &&
                Survana.Validation.Messages[constraint_name](constraint_value)) ||
                Survana.Validation.Messages.invalid,
            question = document.getElementById(field);

        if (!question) {
            console.error('No such question:', field);
            return;
        }

        //call the theme's error message handler
        if (Survana.Validation.ShowMessage) {
            Survana.Validation.ShowMessage(question, message, src_input);
        }
    }

    function valid(field) {
        var question = document.getElementById(field);

        if (!question) {
            console.error('No such question:', field);
            return;
        }

        if (Survana.Validation.HideMessage) {
            Survana.Validation.HideMessage(question);
        }
    }


    /**
     * Validates all the constraints of a single field. Calls invalid() if validation fails.
     * @param field {Object} The Schema field to validate
     * @param values {Object} All form fields with values
     * @returns {Boolean} False if validation fails, Group values otherwise
     */
    Survana.Validation.ValidateField = function (field, values) {
        var response = values[field.id],
            constraint_name,
            constraint_value;

        //skipped fields
        if (response === null) {
            return true;
        }

        //all other fields
        if (!response) {
            return false;
        }

        //check all user-specified constraints
        for (constraint_name in field.validation) {
            if (!field.validation.hasOwnProperty(constraint_name)) {
                continue;
            }

            if (!Survana.Validation.Constraints[constraint_name]) {
                Survana.Warn("Unknown validation constraint:", constraint_name);
                //ignore unknown constraints (TODO: report constraints to author)
                continue;
            }

            constraint_value = field.validation[constraint_name];

            //verify constraint
            if (!Survana.Validation.Constraints[constraint_name](response, constraint_value, values)) {
                console.log('Constraint', constraint_name, '=', constraint_value, 'failed validation; response =', response, 'values =', values);
                invalid(field.id, constraint_name, constraint_value);
                return false;
            }
        }

        //mark the field as valid
        valid(field.id);

        return true;
    };

    /**
     * Validate <form> elements, based on a validation configuration pre-built when publishing the form
     * and custom validation messages. Returns all validated responses.
     * @param form_id {String} The HTMLFormElement being validated
     * @param values {Object} All form fields with values
     * @param schemata {Object} (optional) The form schemata
     * @return {Boolean} Returns all validated responses as an Object, or false if validation failed
     */
    Survana.Validation.Validate = function (form_id, values, schemata) {
        var field,
            i,
            j,
            result;

        //if no schema was provided, attempt to fetch it from Survana.Schema
        schemata = schemata || Survana.Schema[form_id];
        if (!schemata) {
            Survana.Error('No Schema found for form ' + form_id);
            return false;
        }

        //assume the form is valid
        result = true;

        //loop through all known fields
        for (i = 0; i < schemata.fields.length; ++i) {
            field = schemata.fields[i];

            //special case for the matrix element: treat each row as a separate field
            if (field.type === 'matrix') {
                for (j = 0; j < field.rows.length; ++j) {
                    if (!field.validation) {
                        continue;
                    }

                    if (!Survana.Validation.ValidateField(field.rows[j], values)) {
                        result = false;
                    }
                }
            }

            //skip any fields that should not be validated
            if (!field.validation) {
                continue;
            }

            if (!Survana.Validation.ValidateField(field, values)) {
                //validation failed, but keep scanning to display all the fields with errors
                result = false;
            }
        }

        return result;
    };

    /** on* event handler
     * @param el {HTMLElement} The Blur event object
     */
    Survana.Validation.OnEvent = function (el) {
        console.log('onevent', el);

        var form_el,
            form_id,
            field_name = el.getAttribute("name"),
            container_id = el.getAttribute('data-container'),
            field,
            schemata,
            values,
            search_id,
            i;

        if (el.form) {
            form_el = el.form;
            form_id = form_el.id;
        } else {
            form_id = el.getAttribute('data-form');
            form_el = document.forms[form_id];
        }

        Survana.Assert(field_name, el, "Element must have a name attribute");
        Survana.Assert(form_id, el, "Element must have a data-form attribute");

        //if no schema was provided, attempt to fetch it from Survana.Schema
        schemata = Survana.Schema[form_id];
        if (!schemata) {
            Survana.Error('No Schema found for form ' + form_id);
            return;
        }

        //either search for the field name or the container_id
        if (container_id) {
            search_id = container_id;
        } else {
            search_id = field_name;
        }

        //locate the field in the schemata
        for (i = 0; i < schemata.fields.length; ++i) {
            if (schemata.fields[i].id === search_id) {
                field = schemata.fields[i];
                break;
            }
        }

        //now search in the container rows
        if (container_id) {
            //locate the field in the schemata
            for (i = 0; i < field.rows.length; ++i) {
                if (field.rows[i].id === field_name) {
                    field = field.rows[i];
                    break;
                }
            }
        }

        if (!field) {
            Survana.Error("Could not find field " + field_name + " in the schemata.");
            return;
        }

        if (!field.validation) {
            Survana.Warn("No validation configuration for field", field_name);
            return;
        }

        values = Survana.FormFields(form_el, schemata);

        return Survana.Validation.ValidateField(field, values);
    };


    /**
     * Default validation messages. The keys of this object match the names of validation constraints, except for
     * 'invalid', which is a catch-all function.
     * @type Object
     */
    Survana.Validation.Messages = {
        'invalid': function () {
            return "Please enter a valid value for this field.";
        },
        'equal': function () {
            return "This field must be equal to the previous field.";
        },
        'optional': function () {
            return "This field is required";
        },
        'numeric': function () {
            return "This field requires a numeric value";
        },
        'max': function (max) {
            return ["Please enter a value that's less than or equal to ", max, "."].join("");
        },
        'min': function (min) {
            return ["Please enter a value greater than or equal to ", min, "."].join("");
        }
    };
}(document, window.Survana));
