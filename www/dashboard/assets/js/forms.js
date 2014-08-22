(function () {
    "use strict";

    var app = angular.module('forms', ['preview', 'ace']);

    app.controller('FormListCtrl', ['$scope', '$http',
        function FormListCtrl($scope, $http) {
            $scope.forms = [];
            $scope.selected = [];
            $scope.message = '';
            $scope.search = '';
            $scope.max_selected = 10;
            $scope.loading = true;

            $http.get('forms/list').success(function (response, code, request) {
                if (response.success) {
                    $scope.loading = false;
                    $scope.forms = response.message;
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                $scope.loading = false;
                console.log("Error fetching", $location.path())
            });

            $scope.toggle = function (form_id) {
                var index = $scope.selected.indexOf(form_id);
                if (index === -1) {
                    $scope.selected.push(form_id);
                } else {
                    $scope.selected.splice(index, 1);
                }
            };

            $scope.isSelected = function (form_id) {
                return ($scope.selected.indexOf(form_id) >= 0);
            };

            $scope.deleteForm = function (form_id) {
                $scope.message = "";

                $http.delete('form', {params: {'id': form_id}}).success(function (response, code, request) {
                    if (code === 204) {
                        $scope.removeForm(form_id);
                    } else {
                        $scope.message = 'Invalid response from server: ' + response;
                    }
                }).error(function (response) {
                    $scope.message = response.message;
                });
            };

            $scope.deleteSelected = function () {
                $scope.message = "";

                if (!$scope.selected.length) {
                    return;
                }

                if ($scope.selected.length > $scope.max_selected) {
                    return
                }

                for (var i = 0; i < $scope.selected.length; ++i) {
                    $scope.deleteForm($scope.selected[i]);
                }
            };

            $scope.removeForm = function (form_id) {
                for (var i = 0; i < $scope.forms.length; ++i) {
                    if ($scope.forms[i].id === form_id) {
                        $scope.forms.splice(i, 1);
                        break;
                    }
                }

                //remove it from selected
                var index = $scope.selected.indexOf(form_id);

                if (index >= 0) {
                    $scope.selected.splice(index, 1);
                }
            }

        }
    ]);

    app.controller('FormEditCtrl', ['$scope', '$http', '$window', '$location', '$timeout', '$routeParams',
        function AceCtrl($scope, $http, $window, $location, $timeout, $routeParams) {

            $scope.loading = false;
            $scope.error = false;
            $scope.message = "";
            //whether we're in 'create' or 'edit' mode
            $scope.create = ($routeParams.id === undefined);

            $scope.form = {
                name: "MyForm",
                title: "My First Form",
                fields: [
                    {
                        "type": "text",
                        "id": "subject_id",
                        "label": {
                            "html": "Subject ID:"
                        }
                    }
                ]
            };

            //if we're editing a form, 'id' will be set
            if (!$scope.create) {
                //fetch the form JSON and store it in $scope.form
                $http.get('form', {params: $routeParams}).success(function (response, code, request) {
                    if (response.success) {
                        $scope.form = response.message;
                    } else {
                        console.log('Error message', response.message);
                    }
                }).error(function () {
                    console.log("Error fetching", $location.path())
                });
            }

            //if the save operation was successful
            function onSaveSuccess(response, code, request) {
                $scope.loading = false;

                var id;

                //no content = an 'edit' operation, therefore $scope.form should have a valid id
                if (code == 204 && ($scope.form.id !== undefined)) {
                    id = $scope.form.id;
                } else if (response.success && response.message && response.message.id) {
                    id = response.message.id;
                }

                //either redirect to view the new form, or show the message from server
                if (id) {
                    $location.path('/forms/' + id);
                } else {
                    $scope.message = response.message;
                }
            }

            //if the save operation failed
            function onSaveError(response) {
                $scope.loading = false;
                $scope.error = true;
                $scope.message = "Failed to save form (" + response + ")";
            }

            //on Save click
            $scope.saveCode = function () {
                //reset state
                $scope.message = "";
                $scope.error = false;
                $scope.loading = true;

                //the server url is the same, except for the leading slash
                if ($scope.create) {
                    $http.post('forms/create', $scope.form).
                        success(onSaveSuccess).
                        error(onSaveError);
                } else {
                    $http.put('forms/edit', $scope.form, {params: $routeParams}).
                        success(onSaveSuccess).
                        error(onSaveError);
                }
            };

            //on Discard
            $scope.discardCode = function ($event) {

                var button = $($event.target);

                if ($scope.loading) {
                    button.popover({
                        animation: true,
                        placement: 'auto',
                        html: true,
                        content: '<i class="glyphicon glyphicon-exclamation-sign"></i> Please wait for the Save request to complete.',
                        trigger: 'manual',
                        delay: {
                            show: 0,
                            hide: 5
                        }
                    }).popover('show');

                    $timeout(function () {
                        $(button).popover('hide');
                    }, 5000);
                } else {
                    //go back to wherever we came from
                    $window.history.back();
                }
            };

        }
    ]);

    app.controller('FormViewCtrl', ['$scope', '$location', '$routeParams', '$http', '$templateCache',
        function ($scope, $location, $routeParams, $http, $templateCache) {

            $scope.form = {};
            $scope.size = 'M';
            $scope.theme = 'bootstrap';
            $scope.template = null;

            $scope.resize = function (size) {
                $scope.size = size;
            };

            function fetchTemplate(theme_id, theme_version) {
                var url = 'theme?id=' + theme_id + '&version=' + theme_version + '&preview=true',
                    cachedTemplate = $templateCache.get(url);

                if (cachedTemplate) {
                    $scope.template = cachedTemplate;
                    return;
                }

                //fetch the theme template and cache it
                $http.get(url).success(function (response, code, request) {
                    $templateCache.put(url, response);
                    $scope.template = response;
                }).error(function () {
                    console.log("Error fetching", $location.path())
                });
            }

            function fetchForm() {
                //fetch the form JSON and store it in $scope.form
                $http.get('form', {params: $routeParams}).success(function (response, code, request) {
                    if (response.success) {
                        $scope.form = response.message;
                    } else {
                        console.log('Error message', response.message);
                    }
                }).error(function () {
                    console.log("Error fetching", url)
                });
            }

            $scope.getFormDate = function () {
                return (new Date($scope.form.created_on)).toLocaleDateString();
            };

            $scope.validateForm = function () {
                $scope.$broadcast('validate');
            };

            //when 'theme' changes, notify Survana
            $scope.$watch('theme', function (newTheme, oldTheme) {
                Survana.Theme.SetTheme(newTheme,
                    function () {
                        fetchTemplate(newTheme, Survana.Version);
                        fetchForm();
                    },
                    function () {
                        console.error('Failed to load Survana Themes!');
                    });
            });
        }
    ]);
})();