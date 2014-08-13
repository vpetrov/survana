(function () {
    "use strict";

    var app = angular.module("preview", []);

    app.directive("questionnaire", ['$window', '$compile', '$timeout', function ($window, $compile, $timeout) {
        return {
            restrict: 'A',
            require: '?ngModel',
            scope: false,
            link: function (scope, elem, attrs, ngModel) {

                function updateTemplate() {
                    $compile(elem.contents())(scope);

                    //re-render the model
                    ngModel.$render();
                }

                function extract_schemata(form_fields) {
                    var schemata = {},
                        field,
                        id,
                        i;

                    if (!form_fields) {
                        return schemata;
                    }

                    for (i = 0; i < form_fields.length; ++i) {
                        id = form_fields[i].id;
                        if (!id) {
                            //user forgot to supply id
                            //todo: show error to the user
                            console.warn("No ID set for question #" + id);
                        }
                        schemata[id] = {
                            "id": id,
                            "type": form_fields[i].type,
                            "index": i
                        }
                    }

                    return schemata;
                }

                //quick hack to pass this value form the scope to the iframe
                $window.study_id = function () {
                    if (scope.study) {
                        return scope.study.id;
                    }

                    return null;
                };

                //register a NextPage() function that can be called within the questionnaire preview iframe
                $window.NextPage = function () {
                    $timeout(function () {
                        if ((scope.current.index + 1) < scope.study.form_ids.length) {
                            scope.current.index++;
                        } else {
                            scope.current.index = 0;
                        }
                    });
                };

                scope.verifyForm = function () {
                    if (!elem || !elem[0] || !elem[0].contentWindow || !elem[0].contentWindow.validateForm) {
                        return
                    }

                    elem[0].contentWindow.validateForm();
                };

                //when current_form changes, update the template
                scope.$watch('current.index', updateTemplate);

                scope.$watch('template', function (val) {

                    //nothing to do?
                    if (!val) {
                        return
                    }

                    var frame = elem[0],
                        doc = frame.contentDocument || frame.contentWindow.document;

                    //document.write() is the fastest way to update the contents.
                    doc.open();
                    doc.write(scope.template);
                    doc.close();

                    updateTemplate();
                });

                //update the view
                ngModel.$render = function () {

                    var frame = elem[0],
                        doc = frame.contentDocument || frame.contentWindow.document,
                        node = doc.getElementById('content'),
                        schemata_node = doc.getElementById('schemata'),
                        validation_node = doc.getElementById('validation'),
                        result,
                        schemata,
                        validation;


                    //make sure a theme, a template and a rendering node are available
                    if (!Survana.Theme || !scope.template || !node) {
                        return;
                    }

                    result = Survana.Theme.Questionnaire(ngModel.$viewValue);

                    if (result) {
                        if (node.hasChildNodes()) {
                            node.removeChild(node.firstChild);
                        }

                        //append the form
                        node.appendChild(result);

                        //form schemata
                        schemata = extract_schemata(ngModel.$viewValue.fields);

                        if (schemata_node && schemata) {
                            schemata_node.innerHTML = JSON.stringify(schemata);
                        }

                        //validation configuration
                        validation = Survana.Validation.ExtractConfiguration(ngModel.$viewValue);

                        if (validation_node && validation) {
                            validation_node.innerHTML = validation;
                            //rely on the fact that template will include a 'startValidation()' function
                            if (frame.contentWindow.startValidation) {
                                frame.contentWindow.startValidation();
                            }
                        }

                        //if we're supposed to save rendered data
                        if (attrs['render']) {
                            //skip a digest cycle to let the updateTemplate() digest to finish,
                            //otherwise, innerHTML is still the old html, before any watches are updated by the new changes
                            $timeout(function () {
                                //store the HTML data into the variable pointed to by data-render
                                scope[attrs['render']] = "<!DOCTYPE html><html>" + doc.documentElement.innerHTML + "</html>";
                                //update the currently rendered form index
                                scope.current.rendered = scope.current.index;
                            });
                        }
                    }
                }
            }
        }
    }]);
})();