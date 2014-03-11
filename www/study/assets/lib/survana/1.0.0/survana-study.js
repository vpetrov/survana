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

function NextPage(btn) {

    //disable the button
    if (btn) {
        btn.setAttribute('disabled', 'disabled');
    }

    current++;
    localStorage[study_id + "-current"] = current;
    window.location.href = workflow[current];
}

function FinishSurvey(btn) {
    if (btn) {
        btn.setAttribute('disabled', 'disabled');
    }

    delete localStorage[study_id + "-current"];
    delete localStorage[study_id + "-workflow"];
}

window.survana = {
    onFormLoaded: onFormLoaded,
    NextPage: NextPage,
    FinishSurvey: FinishSurvey,
};

//register an onReady handler, i.e. $(document).ready(). Caveat: does not support older versions of IE
document.addEventListener("DOMContentLoaded", function () {
    //remove this handler
    document.removeEventListener("DOMContentLoaded", arguments.callee, false);

    //call the onLoad function
    onFormLoaded();
});
