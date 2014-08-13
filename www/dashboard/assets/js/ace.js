(function () {
    "use strict";
    var app = angular.module("ace", []);

    app.directive('ace', ['$timeout', function ($timeout) {
        return {
            restrict: 'A',
            require: '?ngModel',
            scope: false,
            link: function (scope, elem, attrs, ngModel) {
                var node = elem[0],
                    editor = ace.edit(node),
                    session = editor.getSession(),
                    model_updated = false;

                session.setMode("ace/mode/json");

                // on model change, update the editor's view
                ngModel.$render = function () {
                    try {
                        var code = JSON.stringify(ngModel.$viewValue, null, 4);
                        model_updated = true;
                        editor.setValue(code);
                    } catch (e) {
                        console.error('bad JSON model', e);
                    }

                    model_updated = false;
                };

                // on edit change, update the model
                editor.on('change', function (changes) {
                    //
                    if (model_updated) {
                        return;
                    }

                    $timeout(function () {
                        scope.$apply(function () {
                            var value = editor.getValue();

                            try {
                                var json_value = JSON.parse(value)
                                ngModel.$setViewValue(json_value);
                            } catch (e) {
                                //invalid json (this will always happen while the user is typing JSON code)
                            }
                        });
                    });
                });
            }
        }
    }]);
})();