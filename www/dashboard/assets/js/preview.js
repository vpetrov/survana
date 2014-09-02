(function () {
    "use strict";

    var app = angular.module("preview", []);

    //in: size, template
    //out: verify(), render()
    app.directive("formPreview", ['$compile', function ($compile) {
        return {
            restrict: 'E',
            template: '<iframe class="questionnaire size-{{ size }}""></iframe>',
            require: '?ngModel',
            scope: {
                size: '@',
                template: '=',
                study: '='
            },
            link: function (scope, elem, attrs, ngModel) {

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

                scope.$on('validate', function () {
                    console.log('VALIDATE!');
                    if (!elem || !elem[0] || !elem[0].firstChild) {
                        return;
                    }

                    var iframe = elem[0].firstChild;
                    if (!elem[0].firstChild.contentWindow || !elem[0].firstChild.contentWindow.Survana) {
                        return;
                    }

                    var previewSurvana = iframe.contentWindow.Survana,
                        doc = iframe.contentDocument,
                        schemata = extract_schemata(ngModel.$viewValue.fields),
                        validation_config = previewSurvana.Validation.ExtractConfiguration(ngModel.$viewValue);

                    previewSurvana.Validation.Validate(doc.forms[0], schemata, validation_config, undefined);
                });

                scope.$watch('template', function (val) {
                    //nothing to do?
                    if (!val) {
                        return
                    }

                    var iframe = elem[0].firstChild,
                        doc = iframe.contentDocument || iframe.contentWindow.document;

                    //document.write() is the fastest way to update the contents.
                    doc.open();
                    doc.write(scope.template);
                    doc.close();

                    //update template bindings
                    $compile(doc)(scope);
                });

                //update the view
                ngModel.$render = function () {

                    var iframe = elem[0].firstChild,
                        doc = iframe.contentDocument || iframe.contentWindow.document,
                        node = doc.getElementById('content'),
                        result;

                    //make sure a theme, a template and a rendering node are available
                    if (!Survana.Theme || !node) {
                        return;
                    }

                    result = Survana.Theme.Questionnaire(ngModel.$viewValue);

                    if (result) {
                        //remove existing elements
                        while (node.firstChild) {
                            node.removeChild(node.firstChild);
                        }

                        //extract validation and form schema
                        var previewSurvana = iframe.contentWindow.Survana,
                            form = ngModel.$viewValue;

                        if (previewSurvana.Validation) {
                            var schemata = extract_schemata(form.fields),
                                validation_config = previewSurvana.Validation.ExtractConfiguration(form),
                                form_info = JSON.stringify({
                                    id: form.id,
                                    schemata: schemata,
                                    config: validation_config
                            });

                            var script = doc.createElement('script');
                            script.setAttribute('type', 'text/x-survana-schema');
                            script.setAttribute('class', 'schema');
                            script.innerHTML = form_info;

                            //bake validation info into the HTML (so that the rendered HTML contains information about its schema)
                            node.appendChild(script);

                            //update the live validation
                            previewSurvana.Schema[form.id] = schemata;
                            previewSurvana.Validation.SetConfiguration(form.id, validation_config, undefined);
                        }

                        //append the form
                        node.appendChild(result);

                        //send rendered HTML to parent scopes
                        scope.$emit('form:rendered', "<!DOCTYPE html><html>" + doc.documentElement.innerHTML + "</html>")
                    }
                };

                var iframe = elem[0].firstChild,
                    iframeWindow = iframe.contentWindow;

                iframeWindow.Survana = {
                    Workflow: {
                        OnPageLoad: function () {
                            scope.$emit('form:load', iframeWindow);
                        },
                        NextPage: function () {
                            //read the validation configuration again, because the render() function is sometimes
                            //called prior to the validation code initialized
                            var form = ngModel.$viewValue,
                                form_el = iframe.contentDocument.getElementById(form.id),
                                previewSurvana = iframeWindow.Survana,
                                response;

                            response = previewSurvana.Validation.Validate(form_el, form, undefined);

                            console.log('response = ', response);

                            if (response)
                            {
                                scope.$emit('form:next', iframeWindow);
                            }
                        }
                    }
                }
            }
        }
    }]);
})();