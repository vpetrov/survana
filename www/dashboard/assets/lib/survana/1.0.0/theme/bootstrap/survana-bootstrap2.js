if (typeof Survana === undefined) {
    throw Error("Survana-Bootstrap requires Survana");
}

var id = "bootstrap",
    engine_version = "1.0.0",
    ncolumns = 12,
    //default controls take up all columns (an entire row)
    control_width = {
        'l': ncolumns,
        'm': ncolumns,
        's': ncolumns,
        'xs': ncolumns
    },
    //default width for the first column
    //xs doesn't make sense, because the matrix will be resized to vertical on xs screens
    matrix_width = {
        'l': 4,
        'm': 4,
        's': 4
    },
    default_matrix_type = 'radio';

var templates = {};

templates.html = "<span>{{ html }}</span>";

(function compileTemplates() {
    for (var tpl_name in templates) {
        if (templates.hasOwnProperty(tpl_name)) {
            //replace the source with compiled JS functions
            templates[tpl_name] = Handlebars.compile(templates[tpl_name]);
        }
    }
}());

var BootstrapEngine = function (doc) {

    var DP = new DOMParser();

    /**
     * Creates a container element for the form
     * @param fields
     * @returns {HTMLElement|*}
     */
    function form(fields) {
        return doc.createElement('form');
    }

    /**
     * Returns a single question
     * @param field
     * @returns {String}
     */
    function question(field) {
        return templates.question(field);
    }

    /**
     * Parses an entire JSON form and returns an HTMLDocument element
     * @param fields
     * @returns {HTMLDocument}
     */
    function questionnaire(fields) {
        var form = "<form>",
            f, r, result;

        for (var field in fields) {
            if (fields.hasOwnProperty(field)) {
                f = fields[field];
                if (templates[f.type]) {
                    r = question(f);
                    if (r) {
                        form += r;
                    }
                }
            }
        }

        form += "</form>";
        return DP.parseFromString(form, "text/html");
    }

    return {
        form: form,
        question: question,
        questionnaire: questionnaire
    };
};

Survana.engine[id] = BootstrapEngine;

//set this as the default theme, if no default exists
if (Survana.theme === undefined) {
    Survana.theme = id;
}
