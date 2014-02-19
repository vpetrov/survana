function NextPage(btn) {
    var study_id = btn.getAttribute('study_id'),
        workflow = JSON.parse(localStorage[study_id + "-workflow"]),
        current = localStorage[study_id + "-current"] | 0;

    current++;

    console.log('changing page to', workflow[current]);

    localStorage[study_id + "-current"] = current;

    window.location.href = workflow[current];
}

function FinishSurvey(btn) {

}

window.survana = {
    NextPage: NextPage,
    FinishSurvey: FinishSurvey
};
