/***********
* WORKFLOW *
***********/

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana) {
    Survana.NextPage = function (e) {
        Survana.Validation.Validate(document.forms[0]);
        e.stopPropagation();
        return false;
    }
    Survana.FinishSurvey = function () {}
}(window.Survana));
