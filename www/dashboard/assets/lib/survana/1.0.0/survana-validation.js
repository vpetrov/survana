/*************
 * VALIDATION *
 *************/

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana) {
    var validation_config = {};

    Survana.Validation = {};

    Survana.Validation.NO_VALIDATION = 's-no-validation';
    Survana.Validation.INVALID = 's-invalid';

    Survana.Validation.SetConfiguration = function (form, config, messages) {
        validation_config[form.id] = {
            config: config,
            messages: messages || Survana.Validation.Messages
        };
    };

    Survana.Validation.Skip = function (question_id) {
        var question = document.getElementById(question_id),
            children,
            cl;

        if (question === undefined || !question) {
            return
        }

        cl = question.classList;

        //mark the question for no validation
        cl.remove(Survana.Validation.INVALID);

        children = question.querySelectorAll('input,select');

        //mark all inputs for no validation
        for (var i = 0; i < children.length; ++i) {
            children[i].classList.add(Survana.Validation.NO_VALIDATION);
        }
    };

    Survana.Validation.Constraints = {
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

    /**
     * Validate <form> elements, based on a validation configuration pre-built when publishing the form
     * and custom validation messages
     * @param form The HTMLFormElement being validated
     * @param config (optional) Validation configuration
     * @param messages (optional) Custom error messages
     * @constructor
     */
    Survana.Validation.Validate = function (form, config, messages) {

        if (!form) {
            throw new Error('Invalid validation form supplied to Survana.Validate');
        }

        if (!config && validation_config[form.id]) {
            config = validation_config[form.id].config;
        }

        if (!config) {
            throw new Error('Invalid validation configuration supplied to Survana.Validate');
        } else {
            //cache this configuration object
            Survana.Validation.SetConfiguration(form, config, messages);
        }

        function groupElementsByName(elements) {
            var result = {},
                el,
                name;

            for (var i = 0; i < elements.length; ++i) {
                el = elements[i];
                name = el.getAttribute('name');

                if (!name) {
                    continue;
                }

                //skip elements marked as 'do not validate'
                if (el.classList.contains(Survana.Validation.NO_VALIDATION)) {
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

        //returns the value of an element based on its declared field type
        function getValueByType(element, field_type) {
            switch (element.tagName.toLowerCase()) {
                case 'input':
                        switch (element.getAttribute('type')) {
                            case 'radio':
                            case 'checkbox':
                                if (element.checked) {
                                    return element.value;
                                }
                                return undefined;
                            default: return element.value;
                        }
                        break;
                case 'button':
                case 'select':
                    return element.value;
                default:
                    return undefined;
            }
        }

        function getValuesFromGroup(group, field_type) {
            var result = [],
                value;

            if (!group) {
                return result;
            }

            for (var i = 0; i < group.length; ++i) {
                value = getValueByType(group[i], field_type);
                if (value !== undefined) {
                    result.push(value);
                }
            }

            return result;
        }

        function invalid(field, field_config, constraint, src_input) {
            var message = (Survana.Validation.Messages[constraint] &&
                    Survana.Validation.Messages[constraint](field_config[constraint])) ||
                    Survana.Validation.Messages.InvalidField,
                question = document.getElementById(field);

            if (!question) {
                console.error('No such question:', field);
                return;
            }

            question.classList.add(Survana.Validation.INVALID);

            //call the theme's error message handler
            if (Survana.Theme && Survana.Theme.ShowValidationMessage) {
                Survana.Theme.ShowValidationMessage(question, message, src_input);
            }
        }

        function valid(field, field_config) {
            var question = document.getElementById(field);

            if (!question) {
                console.error('No such question:', field);
                return;
            }

            question.classList.remove(Survana.Validation.INVALID);

            if (Survana.Theme && Survana.Theme.HideValidationMessage) {
                Survana.Theme.HideValidationMessage(question);
            }
        }

        var field, field_config, constraint, elements, group, values, is_valid;

        elements = groupElementsByName(form.elements);

        //loop through all known fields
        for (field in config) {
            if (!config.hasOwnProperty(field)) {
                continue;
            }

            field_config = config[field];

            //skip fields with no type information
            if (!field_config.type) {
                continue;
            }

            //find the controls responsible for this field
            group = elements[field];
            if (group === undefined || !group) {
                continue;
            }

            values = getValuesFromGroup(group, field_config.type);

            is_valid = true;
            //check all user-specified constraints
            for (constraint in field_config) {
                if (!field_config.hasOwnProperty(constraint) || !Survana.Validation.Constraints[constraint]) {
                    continue;
                }

                //verify constraint
                if (!Survana.Validation.Constraints[constraint](values, field_config[constraint])) {
                    console.log('Constraint', constraint, '=', field_config[constraint], 'failed validation; values=', values);
                    invalid(field, field_config, constraint);
                    is_valid = false;
                    break;
                }
            }

            //mark this field as valid
            if (is_valid) {
                valid(field, field_config);
            }
        }
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
        'optional': function () {
            return "This field is required";
        },
        'numeric': function () {
            return "This field requires a numeric value";
        },
        'max': function (max) {
            return ["Please enter a value that's less than ", max, "."].join("");
        },
        'min': function (min) {
            return ["Please enter a value greater than ", min, "."].join("");
        }
    };
}(window.Survana));
