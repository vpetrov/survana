"use strict";

var adminControllers = angular.module('adminControllers', []);

adminControllers.controller('HomeCtrl', ['$scope', '$http',
    function HomeCtrl($scope, $http) {
    }
]);

angular.module("admin").controller('LoginCtrl', ['$scope', '$http', '$location',
    function LoginCtrl($scope, $http, $location) {
        $scope.data = { key : "" };
        $scope.doLogin = function (q) {
            $http.post("session", $scope.data).success(function () {
                console.log('SUCCESS, redirection to /');
                $location.path("/");
            }).error(function() {
                    console.log('LOGIN FAILED');
                })
        }
    }
]);
