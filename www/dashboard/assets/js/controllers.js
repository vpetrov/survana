"use strict";

var dashboard = angular.module('dashboard', []);

dashboard.controller('HomeCtrl', ['$scope', '$http',
    function HomeCtrl($scope, $http) {
        console.log('HomeCtrl running')
    }
]);

dashboard.controller('FormCtrl', ['$scope',
    function FormCtrl($scope) {
        console.log('FormCtrl running')
    }
]);

dashboard.controller('StudyCtrl', ['$scope', '$http', '$location', '$window',
    function StudyCtrl($scope, $http, $location, $window) {
        console.log('StudyCtrl was created/invoked.');
    }
]);

dashboard.controller('NavigationCtrl', ['$scope', '$location',
    function NavigationCtrl($scope, $location) {
        //glyphicons
        $scope.icons = {
            "dashboard": "home",
            "studies": "th-large",
            "study": "th-large",
            "forms": "list-alt",
            "form": "list-alt",
            "users": "user",
            "user": "user",
            "logs": "align-center"
        };

        $scope.isActive = function (pageUrl) {
            var path = $location.path();

            // we need to use equality for "/", because it's a prefix of all paths
            if (pageUrl === "/") {
                return path === pageUrl
            }

            // for all other page urls, we need to see if the url is a valid prefix
            // this is so the path "/foo/bar" will match the page url "/foo"
            return $location.path().indexOf(pageUrl) === 0;
        }
    }
]);

dashboard.controller('CreateFormCtrl', ['$scope', '$http',
    function AceCtrl($scope, $http) {
        $scope.code = [
            "{",
            "\t\"name\":\"MyForm\",",
            "\t\"title\":\"My First Form\",",
            "\t\"fields\": [",
            "\t\t{",
            "\t\t\t\"s-type\":\"text\",",
            "\t\t\t\"s-id\":\"subject_id\"",
            "\t\t}",
            "\t]",
            "}"
        ].join("\n");

        $scope.loading = false;

        $scope.saveCode = function () {
            var ok = false;
            try {
                JSON.parse($scope.code);
                ok = true
            } catch (e) {
                //there's a JSON error in the user's code. ignore for now
            }

            if (ok) {
                $scope.loading = true;
                $http.post('form/create', $scope.code).
                    //POST succeeded
                    success(function () {
                        $scope.loading = false;
                        console.log('post succeeded', arguments);
                    }).
                    //POST failed
                    error(function () {
                        $scope.loading = false;
                        console.log('post failed', arguments)
                    });
            }
        };

        $scope.discardCode = function () {
            console.log('Discarding code');
            $scope.code = "";
        };

        console.log('Ace ctrl invoked!');
    }
]);

/* DIRECTIVES */

dashboard.directive('ace', ['$timeout', function ($timeout) {
    return {
        restrict: 'A',
        require: '?ngModel',
        scope: false,
        link: function (scope, elem, attrs, ngModel) {
            var node = elem[0],
                editor = ace.edit(node),
                session = editor.getSession();

            session.setMode("ace/mode/json");

            // on model change, update the editor's view
            ngModel.$render = function () {
                editor.setValue(ngModel.$viewValue);
            };

            // on edit change, update the model
            editor.on('change', function () {
                $timeout(function () {
                    scope.$apply(function () {
                        var value = editor.getValue();
                        ngModel.$setViewValue(value);
                    });
                });
            });
        }
    }
}]);
