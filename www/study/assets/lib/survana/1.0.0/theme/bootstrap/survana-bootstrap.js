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

var BootstrapEngine = function (doc) {

    var types = {
        button: radio_button,
        input: input,
        text: text,
        radio: radio,
        checkbox: checkbox,
        option: option,
        optgroup: optgroup,
        select: select,
        instructions: instructions,
        separator: separator,
        group: group,
        matrix: matrix
    };

    var radio_count = 0,
        checkbox_count = 0,
        input_count = 0;

    function _html(elem, field, value) {
        elem.innerHTML = value || field.html || "";
    }

    function _value(elem, value) {
        elem.setAttribute('value', value);
    }

    function _size(elem, field, s) {
        var c = elem.getAttribute('class') || "";

        if (s === undefined || !s) {
            s = getSizes(field).control;
        }

        if (s.l) {
            c += ' col-lg-' + s.l;
        }

        if (s.m) {
            c += ' col-md-' + s.m;
        }

        if (s.s) {
            c += ' col-sm-' + s.s;
        }

        if (s.xs) {
            c += ' col-xs-' + s.xs;
        }

        elem.setAttribute('class', c);
    }

    function _align(elem, field) {

        var c = elem.getAttribute('class') || "";

        console.log('align ', elem.getAttribute('id'), field.align)

        if (field.align) {
            elem.setAttribute('class', c + ' text-' + field.align);
        }
    }

    //the field refers to the control element to which a label should be attached.
    function _label(field, sizes) {

        //field.label must exist
        if (field.label === undefined) {
            return
        }

        var elem = doc.createElement("label");

        if (sizes === undefined) {
            sizes = getSizes(field);
        }

        elem.setAttribute('for', field.id);
        elem.setAttribute('class', 'control-label');

        _align(elem, field.label);
        _size(elem, field, sizes.label);
        _html(elem, field.label);

        return elem;
    }

    function _striped(elem, field) {
        var c = elem.getAttribute('class') || "";

        if (field.striped) {
            c += ' ' + 'striped';

            elem.setAttribute('class', c);
        }
    }

    function _hover(elem, field) {
        var c = elem.getAttribute('class') || "";

        if (field.hover) {
            c += ' ' + 'hover';
            elem.setAttribute('class', c);
        }
    }

    function _noanswer(elem, field) {
        var c = elem.getAttribute('class') || "";

        if (field.noanswer) {
            c += ' ' + 'no-answer text-muted';
            elem.setAttribute('class', c);
        }
    }

    function group(field) {

        var container = doc.createElement('div'),
            elem,
            f,
            i,
            c = "group",
            child_id = 0;

        if (field.fields === undefined) {
            return null;
        }

        for (i in field.fields) {
            f = field.fields[i];

            //a 'radio'-group has 'radio' fields, a 'checkbox'-group has 'checkbox'-fields, etc
            if (f.type === undefined) {
                f.type = field.group;
            }

            //auto-generate an ID if necessary
            if (f.id === undefined) {
                f.id = field.id + ":" + child_id;
                child_id += 1;
            }

            //auto-generate a 'name' property
            if (f.name === undefined) {
                f.name = field.name || field.id;
            }

            //crate the child element
            elem = by_type(f);

            //append child to list
            if (elem) {
                container.appendChild(elem);
            }
        }



        switch (field.group) {
            case 'button': c += ' btn-group btn-group-justified';
                           container.setAttribute('data-toggle', 'buttons');
                           break;
        }

        container.setAttribute('class', c);

        return container;
    }

    //Internally, the matrix will always think it takes up 12 columns. Set the size on a container element.
    function matrix(field) {
        console.log('Creating matrix element');

        var container = doc.createElement('div'),
            header = doc.createElement('div'),
            column_row_wrapper = doc.createElement('div'),
            column_row = doc.createElement('div'),
            column_header_wrapper = doc.createElement('div'),
            matrix_columns = field.columns.length - 1,
            total_col_width,
            col_width,
            fcol_width,
            i,
            el,
            inner_el;

        //first column size
        fcol_width = field.columns[0].size;

        if (!fcol_width) {
            fcol_width = matrix_width;
        }

        total_col_width = {
            'l': ncolumns - fcol_width.l,
            'm': ncolumns - fcol_width.m,
            's': ncolumns - fcol_width.s
        };

        //each column's width is a portion of the total width allocated for the columns.
        col_width = {
            'l': Math.ceil(total_col_width.l / matrix_columns),
            'm': Math.ceil(total_col_width.m / matrix_columns),
            's': Math.ceil(total_col_width.s / matrix_columns)
        };

        container.setAttribute('class', 'matrix');
        _hover(container, field);
        _striped(container, field);

        /* MATRIX HEADER */

        //matrix header row
        header.setAttribute('class', 'row matrix-header-row');

        /*  push the columns to the right the same distance as the width of the first column
            the width of the rest of the columns should be the same as for the controls in each cell
            the labels in the header and on each control must match
        */
        column_row_wrapper.setAttribute('class', ' col-lg-push-' + fcol_width.l +
                                                 ' col-md-push-' + fcol_width.m +
                                                 ' col-sm-push-' + fcol_width.s +
                                                 ' col-lg-' + total_col_width.l +
                                                 ' col-md-' + total_col_width.m +
                                                 ' col-sm-' + total_col_width.s);

        //make a row to contain all the column header cells
        column_header_wrapper.setAttribute('class', 'row');

        //start with the second column, since the first column is config for the labels
        for (i = 1; i < field.columns.length; ++i) {
            el = doc.createElement('div');
            //mark each header cell. it should be hidden on 'xs' devices, and the text in the middle
            el.setAttribute('class', ' matrix-header hidden-xs text-center' +
                                     ' col-lg-' + col_width.l +
                                     ' col-md-' + col_width.m +
                                     ' col-sm-' + col_width.s);

            _noanswer(el, field.columns[i]);

            inner_el = doc.createElement('div');
            inner_el.setAttribute('class', 'matrix-header-text');
            _html(inner_el, field.columns[i]);
            _noanswer(inner_el, field.columns[i]);

            //append the header text to header cell
            el.appendChild(inner_el);
            //append the header cell to the header row
            column_header_wrapper.appendChild(el);
        }

        //append the header row to the header row wrapper
        column_row_wrapper.appendChild(column_header_wrapper);
        //append the header row wrapper to the matrix header
        header.appendChild(column_row_wrapper);

        //append the header to the matrix container
        container.appendChild(header);

        /*  MATRIX ROWS
            Each  header label is repeated here and displayed on 'xs' screens
         */

        var row,
            row_label,
            padded_row_label,
            control_row_wrapper,
            control_row,
            control_wrapper,
            control_el,
            j;

        if (field.equalize) {
            //first, if equalize was enabled, pad all labels with non breakable spaces to make all labels
            //have the same width. We need to determine the length of the longest label and reset the paddedHtml prop
            var max_label = 0;
            for (i = 0; i < field.rows.length; ++i) {
                if (field.rows[i].html.length > max_label) {
                    max_label = field.rows[i].html.length;
                }

                //remove any previous padding settings
                field.rows[i].paddedHtml = '';
            }

            //now, loop over all labels again and pad them until max_label
            var field_length, padding, npad;
            for (i = 0; i < field.rows.length; ++i) {
                field_length = field.rows[i].html.length;
                if (field_length < max_label) {
                    //create a string of repeating 'nbsp;'
                    npad = max_label - field_length;
                    //npad -= ~~(npad/3); //subtract 30%
                    padding = Array(npad).join("&nbsp; ") || "";
                    field.rows[i].paddedHtml = field.rows[i].html + padding;
                } else {
                    field.rows[i].paddedHtml = field.rows[i].html;
                }
            }
        }

        //loop over all user-supplied rows
        for (i = 0; i < field.rows.length; ++i) {
            //create a new matrix row
            row = doc.createElement('div');
            row.setAttribute('class', 'row matrix-row');

            //create a new row label
            row_label = doc.createElement('label');
            var row_label_class = 'matrix-question control-label';
            if (field.equalize) {
                row_label_class += ' equalized';
            }
            row_label.setAttribute('class', row_label_class);
            _size(row_label, null, fcol_width);

            //prepend numbers, if desired
            if (field.numbers) {
                field.rows[i].html = (i + 1) + ". " + (field.rows[i].html || "");
            }
            _html(row_label, field.rows[i]);

            row.appendChild(row_label);

            //append equalized label
            if (field.equalize) {
                padded_row_label = document.createElement('label');
                padded_row_label.setAttribute('class', 'matrix-question control-label equalized-padded');
                _size(padded_row_label, null, fcol_width);
                if (field.numbers) {
                    field.rows[i].paddedHtml = (i + 1) + ". " + (field.rows[i].paddedHtml || "");
                }

                _html(padded_row_label, field.rows[i], field.rows[i].paddedHtml);

                row.appendChild(padded_row_label);
            }

            //control row wrapper
            control_row_wrapper = doc.createElement('div');
            //the wrapper should have the same width as the column headers
            _size(control_row_wrapper, null, total_col_width);

            //create a new row
            control_row = doc.createElement('div');
            control_row.setAttribute('class', 'row');

            var row_id = field.rows[i].id;

            //create a control for each column
            for (j = 1; j < field.columns.length; ++j) {
                control_wrapper = doc.createElement('div');
                //the control should have the same width as a column
                _size(control_wrapper, null, col_width);

                field.columns[j]['-input-label-class'] = 'visible-xs text-left';
                //id= row_id:0,1,2,...
                field.columns[j].id = row_id + ":" + (j - 1);
                field.columns[j].name = row_id;

                //create the element
                control_el = by_type(field.columns[j], field.matrix || default_matrix_type);

                //append element to the wrapper
                control_wrapper.appendChild(control_el);

                //append the wrapper to the current matrix row
                control_row.appendChild(control_wrapper)
            }

            control_row_wrapper.appendChild(control_row);
            row.appendChild(control_row_wrapper);
            container.appendChild(row);
        }

        return container;
    }

    //returns an <input> element
    function input(field) {
        console.log('Creating input element');

        //create <input>
        var elem = doc.createElement('input'),
            id = field.id || 'input_' + (input_count++),
            name = field.name || id;

        elem.setAttribute('id', id);
        elem.setAttribute('name', name);
        elem.setAttribute('type', field.type || "text");
        elem.setAttribute('class', 'form-control');

        if (field.placeholder) {
            elem.setAttribute('placeholder', field.placeholder);
        }

        return elem;
    }

    //syntactic sugar for input()
    function text(field) {
        return input(field);
    }

    function radio(field) {
        var container = doc.createElement('div'),
            elem = doc.createElement('input'),
            label = doc.createElement('label'),
            label_text = doc.createElement('span'),
            id = field.id || 'radio_' + radio_count++,
            name = field.name || id;

        container.setAttribute('class', 'radio');

        elem.setAttribute('id', id);
        elem.setAttribute('name', name);
        elem.setAttribute('type', 'radio');
        _value(elem, field.value);

        if (field['-input-label-class']) {
            label_text.setAttribute('class', field['-input-label-class']);
        }

        _html(label_text, field);


        //create the label without using _label, because this label is special (no need for default 'label' props)
        label.setAttribute('for', elem.id);
        label.setAttribute('class', 'control-label');

        //append the radio button and text to the label
        label.appendChild(elem);
        label.appendChild(label_text);

        //append the label to the container
        container.appendChild(label);

        return container;
    }

    function checkbox(field) {
        var container = doc.createElement('div'),
            elem = doc.createElement('input'),
            label = doc.createElement('label'),
            label_text = doc.createElement('span'),
            id = field.id || 'checkbox_' + checkbox_count++,
            name = field.name || id;

        container.setAttribute('class', 'checkbox');

        elem.setAttribute('id', id);
        elem.setAttribute('name', name);
        elem.setAttribute('type', 'checkbox');

        _value(elem, field.value);

        if (field['-input-label-class']) {
            label_text.setAttribute('class', field['-input-label-class']);
        }

        _html(label_text, field);

        //create the label without using _label, because this label is special (no need for default 'label' props)
        label.setAttribute('for', elem.id);
        label.setAttribute('class', 'control-label');

        //append the radio button and text to the label
        label.appendChild(elem);
        label.appendChild(label_text);

        //append the label to the container
        container.appendChild(label);

        return container;
    }

    /*
         <label class="btn btn-default">
            <input type="radio" name="sex" id="sex:0">Male</input>
         </label>
     */
    function radio_button(field) {
        var elem = doc.createElement('label'),
            child = doc.createElement('input'),
            id = field.id || 'button_' + radio_count++,
            name = field.name || id;

        elem.setAttribute('class', 'btn btn-default');

        child.setAttribute('type', 'radio');
        child.setAttribute('id', id)
        child.setAttribute('name', name);

        _value(child, field.value);

        elem.innerHTML += field.html;

        if (child) {
            elem.appendChild(child);
        }

        return elem;
    }

    function option(field) {
        var elem = doc.createElement('option');

        _html(elem, field);

        if (field.value !== undefined) {
            elem.setAttribute('value', field.value)
        }

        return elem;
    }

    function optgroup(field) {
        var elem = doc.createElement('optgroup'),
            child;

        if (field.html) {
            elem.setAttribute('label', field.html);
        }

        if (field.fields) {
            for (var f in field.fields) {
                if (field.hasOwnProperty(f)) {
                    child = option(field[f]);

                    if (child) {
                        elem.appendChild(child);
                    }
                }
            }
        }

        return elem;
    }

    function select(field) {
        var elem = doc.createElement('select'),
            c;

        elem.setAttribute('class', 'form-control');

        if (field.fields) {
            for (var f in field.fields) {
                if (field.fields.hasOwnProperty(f)) {
                    c = field.fields[f];

                    //an optgroup will have a 'fields' property, simple options won't have it
                    if (c.fields === undefined) {
                        c = option(c);
                    } else {
                        c = optgroup(c);
                    }

                    if (c) {
                        elem.appendChild(c);
                    }
                }
            }
        }

        return elem;
    }

    function instructions(field) {
        var elem = doc.createElement('blockquote');

        _align(elem, field);
        _html(elem, field);

        return elem;
    }

    function separator(field) {
        return doc.createElement('hr');
    }

    //returns an element generated from one of the functions in the 'types' object,
    //or null.
    function by_type(field, t) {

        var type = field.type || t,
            elem;

        //supported field?
        if (type && (types[type] !== undefined)) {

            //set the field type, in case it wasn't already set
            field.type = type;

            //generate the element
            elem = types[type](field);
        }

        return elem;
    }

    function getSizes(field) {

        var result = {
                control: {},
                label: {}
            },
            cw, //computed control width for a given device size
            lw, //computed label width for a given device size
            max_lw,//maximum allowed control width for a given device size, based on the value of 'cw'
            i;

        for (i in control_width) {
            if (control_width.hasOwnProperty(i)) {

                //reset control width and label width
                cw = null;
                lw = null;

                //the control width is the minimum value from the default or user supplied
                if (field.size !== undefined) {
                    cw = field.size[i];
                }

                //fix invalid widths
                if (cw === undefined || cw === null || cw <= 0 || cw > ncolumns) {
                    cw = control_width[i];
                }

                //set control width for the current screen size
                result.control[i] = cw;

                if (field.label !== undefined) {
                    //set the maximum label width for this field
                    max_lw = ncolumns - cw;

                    //no room for a label? set label width to full size to create row
                    if (max_lw == 0) {
                        lw = ncolumns;
                    } else {
                        //did the user specify a label size?
                        if (field.label.size !== undefined) {
                            lw = field.label.size[i];
                        }

                        //make sure the values are sane. if not, expand to max size
                        if (lw === undefined || lw === null || lw < 0) {
                            lw = max_lw;
                        }

                        //if lw + cw > max_lw, bootstrap will automatically place the label and control on separate rows
                    }

                    //setting label width to 0 will prevent it from showing
                    if (lw) {
                        //set label width for the current screen size
                        result.label[i] = lw;
                    }
                }
            }
        }

        return result;
    }

    //generates a textual addon (either prefix or suffix)
    function addon_text(field) {
        var elem;

        elem = document.createElement('span');
        elem.setAttribute('class', 'input-group-addon');

        if (typeof field === "object") {
            //assume .html exists
            _html(elem, field);
        } else {
            //assume string, number or other simple type
            _html(elem, {html:field});
        }

        return elem;
    }

    //generates an addon element (suffix or prefix).
    //in theory, this can be any control supported by bootstrap.
    function addon(field) {
        var elem;

        //TODO: arrays
        if (typeof field === "object") {
            //if a type has been specified, use it
            if (field.type !== undefined) {
                elem = by_type(field);
            } else if (field.html !== undefined) {
                //if there is no type, but there's a label, assume it's a text addon
                elem = addon_text(field);
            }
        } else {
            elem = addon_text(field);
        }

        return elem;
    }

    //returns the control for a question. this can either be an element, or a group of elements (input with suffix)
    function control(field) {
        var container,
            elem;

        //generate the actual control
        elem = by_type(field);

        //generate any suffixes or prefixes
        if (field.suffix || field.prefix) {
            var prefix, suffix;

            //create container group
            //TODO: figure out if <select> works with input group
            container = document.createElement('div');
            container.setAttribute('class', 'input-group');

            if (field.prefix) {
                prefix = addon(field.prefix);
                if (prefix) {
                    container.appendChild(prefix);
                } else {
                    console.warn(field.id,"prefix declared, but no element was generated for",field.prefix);
                }
            }

            //append the child after the prefix
            container.appendChild(elem)

            if (field.suffix) {
                suffix = addon(field.suffix);
                if (suffix) {
                    container.appendChild(suffix);
                } else {
                    console.warn(field.id,"suffix declared, but no element generated for",field.suffix);
                }
            }
        }

        return container || elem;
    }

    // <question id="%QID%"><div class="form-group"> ... </div></question>
    function question(field) {
        var q, form_group, row, label, cwrap, elem, sizes;

        sizes = getSizes(field);

        //<question>
        q = doc.createElement('question');
        q.setAttribute('id', field.id);

        //form group
        form_group = doc.createElement('div');
        form_group.setAttribute('class', 'form-group');

        //row
        row = doc.createElement('div');
        row.setAttribute('class', 'row');

        //control wrapper
        cwrap = doc.createElement('div');
        _size(cwrap, field, sizes.control);

        //control element
        elem = control(field);
        cwrap.appendChild(elem);

        //label
        label = _label(field);

        if (label) {
            row.appendChild(label);
        }


        row.appendChild(cwrap);

        form_group.appendChild(row);
        q.appendChild(form_group);

        return q;
    }

    function form(field) {
        var elem = doc.createElement('form');

        elem.setAttribute('role', 'form');
        //elem.setAttribute('class', 'form-vertical');
        if (field.id) {
            elem.setAttribute('id', field.id);
        } else {
            //autogenerate an id
            elem.setAttribute('id', 'form-'+String((new Date()).valueOf()));
        }

        return elem;
    }

    return {
        form: form,
        question: question
    };
};

Survana.engine[id] = BootstrapEngine;

//set this as the default theme, if no default exists
if (Survana.theme === undefined) {
    Survana.theme = id;
}
