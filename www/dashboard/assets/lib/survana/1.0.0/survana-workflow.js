/****************
* WORKFLOW stub *
****************/

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana) {
    function next_page() {
        Survana.Validation.Validate(document.forms[0]);
        e.stopPropagation();
        return false;
    }

    function finish_survey() {
        //stub
    }

    Survana.Workflow = {
        NextPage: next_page,
        FinishSurvey: finish_survey
    };
}(window.Survana));
