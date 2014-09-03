/*************************
* THEME-BASED VALIDATION *
*************************/

"use strict";

window.Survana = window.Survana || {};

(function (Survana) {

    Survana.Validation = Survana.Validation || {};

    /* THEME-BASED VALIDATION */

    var messages = {};

    function onSkipQuestion(e) {
        var btn = e.currentTarget,
            q_id = btn.getAttribute('data-question'),
            q_msg = messages[q_id];

        if (!q_msg) {
            return false;
        }

        //hide the error message
        Survana.Theme.Current.HideValidationMessage(document.getElementById(q_id));

        Survana.Validation.Skip(q_id);

        return false;
    }

    function newErrorMessage(question, message) {
        //temporary error message. this should be implemented by the current theme.
        var errdiv = document.createElement('div'),
            errmsg = document.createElement('span'),
            q_id = question.getAttribute('id');

        errdiv.setAttribute('class','s-error alert alert-warning');
        errmsg.innerHTML = message;

        var skipbtn = document.createElement('button');
        skipbtn.setAttribute('type', 'button');
        skipbtn.setAttribute('class', 'btn btn-sm btn-default');
        skipbtn.setAttribute('data-question', q_id);
        skipbtn.innerHTML = 'Prefer Not to Answer';
        skipbtn.addEventListener('click', onSkipQuestion);

        errdiv.appendChild(errmsg);
        errdiv.appendChild(skipbtn);

        question.insertBefore(errdiv, question.firstChild);

        return errdiv;
    }

    function show_validation_message(question, message) {
        var q_id = question.getAttribute('id'),
            err_el;

        //reuse a previous error message
        if (messages[q_id]) {
            err_el = messages[q_id];
            err_el.firstChild.innerHTML = message;
            err_el.classList.remove('hidden');
        } else {
            //create a new error message
            err_el = newErrorMessage(question, message);

            //cache the error element
            messages[q_id] = err_el;
        }

        //assume the 'form-group' is the last child of the <question>
        //add .has-error to it
        question.lastChild.classList.add('has-error');
    };

    function hide_validation_message(question) {
        var q_id = question.getAttribute('id'),
            err_el,
            i;

        if (messages[q_id]) {
            messages[q_id].classList.add('hidden');
        }

        question.lastChild.classList.remove('has-error');
    }
}(window.Survana));
