/*************************
* THEME-BASED VALIDATION *
*************************/

"use strict";

window.Survana = window.Survana || {};

(function (Survana) {

    Survana.Validation = Survana.Validation || {};

    /* THEME-BASED VALIDATION */

    var message_els = {},
        count = 0;


    function onSkipQuestion(e) {
        var btn = e.currentTarget,
            q_id = btn.getAttribute('data-question'),
            q_msg = message_els[q_id];

        if (!q_msg) {
            return false;
        }

        //hide the error message
        hide_validation_message(document.getElementById(q_id));

        //mark this question as having no answer
        Survana.NoAnswer(q_id);

        return false;
    }

    function onHideMessage(e) {
        var btn = e.currentTarget,
            q_id = btn.getAttribute('data-question'),
            q_msg = message_els[q_id];

        if (!q_msg) {
            return false;
        }

        hide_validation_message(document.getElementById(q_id));

        return false;
    }

    function newErrorElement(question) {
        //temporary error message. this should be implemented by the current theme.
        var errdiv = document.createElement('div'),
            errmsg = document.createElement('div'),
            q_id = question.getAttribute('id'),
            skipbtn,
            hidebtn;


        errdiv.setAttribute('class','s-error alert alert-warning alert-dismissible');
        errdiv.setAttribute('role', 'alert');
        errdiv.setAttribute('id', 'survana-message-' + count);
        count++;

        if (!question.classList.contains('no-skip')) {
            skipbtn = document.createElement('button');
            skipbtn.setAttribute('type', 'button');
            skipbtn.setAttribute('class', 'btn btn-xs btn-default');
            skipbtn.setAttribute('data-question', q_id);
            skipbtn.innerHTML = 'Prefer Not to Answer';
            skipbtn.addEventListener('click', onSkipQuestion);
        }

        hidebtn = document.createElement('button');
        hidebtn.setAttribute('type', 'button')
        hidebtn.setAttribute('data-dismiss', 'alert');
        hidebtn.setAttribute('class', 'close pull-left');
        hidebtn.setAttribute('data-question', q_id);
        hidebtn.addEventListener('click', onHideMessage);
        hidebtn.innerHTML = '&times;'

        errdiv.appendChild(errmsg);
        if (skipbtn) {
            errdiv.appendChild(skipbtn);
        }

        errdiv.appendChild(hidebtn);

        question.insertBefore(errdiv, question.firstChild);

        return errdiv;
    }

    function show_validation_message(question, message) {
        var q_id = question.getAttribute('id'),
            err_el;

        //reuse a previous error message
        if (message_els[q_id]) {
            err_el = document.getElementById(message_els[q_id]);
        }

        //if no previous error message was found, create a new one
        if (!err_el) {
            //create a new error message element
            err_el = newErrorElement(question);
            //cache the error elementme
            message_els[q_id] = err_el.id;
        }

        err_el.firstChild.innerHTML = message;
        err_el.classList.remove('hidden');


        //assume the 'form-group' is the last child of the <question>
        //add .has-error to it
        question.lastChild.classList.add('has-error');
    }

    function hide_validation_message(question) {
        var q_id = question.getAttribute('id'),
            err_el,
            i;

        if (message_els[q_id]) {
            err_el = document.getElementById(message_els[q_id]);

            if (err_el) {
                err_el.classList.add('hidden');
            }
        }

        question.lastChild.classList.remove('has-error');
    }

    Survana.Validation.ShowMessage = show_validation_message;
    Survana.Validation.HideMessage = hide_validation_message;
}(window.Survana));
