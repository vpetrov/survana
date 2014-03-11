/***********
* WORKFLOW *
***********/

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana) {

    var study_id = sessionStorage['study_id'],
        localWorkflow = localStorage[study_id + "-workflow"],
        workflow,
        current;

    function onFormLoaded() {
        if (localWorkflow) {
            workflow = JSON.parse(localWorkflow);
            current = localStorage[study_id + "-current"] | 0;
        }
    }

    function onDOMContentLoaded () {
        //remove this handler
        document.removeEventListener("DOMContentLoaded", onDOMContentLoaded, false);

        //call the onLoad function
        onFormLoaded();
    }

    Survana.NextPage = function (btn) {

        //disable the button
        if (btn) {
            btn.setAttribute('disabled', 'disabled');
        }

        //if validation succeeds, go to the next form
        if (Survana.Validation.Validate(document.forms[0])) {
            current++;
            localStorage[study_id + "-current"] = current;
            window.location.href = workflow[current];
        } else if (btn) {
            btn.removeAttribute('disabled');
            //scroll to first error
            var error_el = document.forms[0].querySelector('.s-error');
            if (error_el) {
                var y = error_el.offsetTop - 100;
                window.scrollBy(0, y);
            }
        }
    };

    Survana.FinishSurvey = function (btn) {
        if (btn) {
            btn.setAttribute('disabled', 'disabled');
        }

        delete localStorage[study_id + "-current"];
        delete localStorage[study_id + "-workflow"];
    };


    //register an onReady handler, i.e. $(document).ready(). Caveat: does not support older versions of IE
    document.addEventListener("DOMContentLoaded", onDOMContentLoaded);
}(window.Survana));
