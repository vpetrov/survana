/* survana-workflow.js

Survana.Workflow contains function that are responsible for the control flow during the taking of a survey.

Dependencies: survana-storage.js

@author Victor Petrov <victor_petrov@harvard.edu>
@license BSD
@date 05/01/2014
*/

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana) {

    if (!Survana.Storage || !Survana.Storage.IsAvailable) {
        console.error('Survana Storage is not available.');
        return;
    }

    var context = {
        workflow: {},
        current: 0
    };

    /** Handles errors reported by Survana.Storage
     * @todo Log the error on the server, display notification to user
     * @param {Error} e
     */
    function onStorageError(e) {
        console.error(e);
    }

    function onFormLoaded() {
        Survana.Storage.Get(context, function (result) {
            context = result;
            context.current |= 0; //convert 'current' to a number
        }, onStorageError);
    }

    function onDOMContentLoaded () {
        //remove this handler
        document.removeEventListener("DOMContentLoaded", onDOMContentLoaded, false);

        //call the onLoad function
        onFormLoaded();
    }

    /** Callback function for goign to the next form. This function performs response validation and will load the next
     * form or scroll the page to the first error.
     * @param btn {HTMLButtonElement} The source button
     */
    Survana.NextPage = function (btn) {

        //disable the button
        if (btn) {
            btn.setAttribute('disabled', 'disabled');
        }

        //if validation succeeds, go to the next form
        if (Survana.Validation.Validate(document.forms[0])) {
            context.current++;
            //Store the incremented value of 'current'
            Survana.Storage.Set('current', context.current, function () {
                //load the next form
                window.location.href = context.workflow[context.current];
            }, onStorageError);

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

    /** Terminates the survey by disabling the source button and removing all workflow from storage.
     * @param btn {HTMLButtonElement} The source button
     */
    Survana.FinishSurvey = function (btn) {
        if (btn) {
            btn.setAttribute('disabled', 'disabled');
        }

        //remove the entire workflow from storage
        Survana.Storage.Remove(context, null, onStorageError);
    };


    //register an onReady handler, i.e. $(document).ready(). Caveat: does not support older versions of IE
    document.addEventListener("DOMContentLoaded", onDOMContentLoaded);
}(window.Survana));
