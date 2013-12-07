"use strict";

var admin = angular.module('admin', [
    'ngRoute',
    'adminControllers'
]);

admin.config(['$routeProvider', '$controllerProvider', function ($routeProvider, $controllerProvider) {
    $routeProvider.
        when("/", {
            templateUrl: 'home',
            controller: 'HomeCtrl'
        }).
        when("/login", {
            templateUrl: 'login',
            controller: 'LoginCtrl'
        }).
        otherwise({
            redirectTo: "/"
        });

    admin.controller = $controllerProvider.register;
}]);

admin.run(function($rootScope, $location) {
    $rootScope.$on('$routeChangeError', function (nge, current, previous, rejection) {
        //when unauthorized, redirect to /login
        if (rejection !== undefined && rejection.status === 401) {
            $location.path("/login");
        }
    });
});
