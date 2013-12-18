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
