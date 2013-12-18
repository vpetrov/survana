"use strict";

var adminControllers = angular.module('adminControllers', []);

adminControllers.controller('HomeCtrl', ['$scope', '$http',
    function HomeCtrl($scope, $http) {
    }
]);

angular.module("admin").controller('LoginCtrl', ['$scope', '$http', '$location', '$window',
    function LoginCtrl($scope, $http, $location, $window) {
        $scope.data = { email : "" };
        $scope.doLogin = function (q) {
            $http.post("login", $scope.data).success(function (response, status) {
                if (response.success == true) {
                    if ((response.message !== undefined) && (response.message.redirect !== undefined)) {
                        $window.location.href = response.message.redirect;
                    } else {
                        $location.path("/");
                    }
                } else {
                    console.log('SERVER ERROR / LOGIN FAILED')
                }

            }).error(function() {
                    console.log('LOGIN FAILED');
                })
        }
    }
]);
