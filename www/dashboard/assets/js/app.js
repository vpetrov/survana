"use strict";

var app = angular.module('dashboardApp', [
    'ngRoute',
    'dashboard'
]);

app.config(['$routeProvider', '$controllerProvider', function ($routeProvider, $controllerProvider) {
    $routeProvider.
        when("/", {
            templateUrl: 'home',
            controller: 'HomeCtrl'
        }).
        when("/forms", {
            templateUrl: 'forms',
            controller: 'FormListCtrl'
        }).
        when("/forms/create", {
            templateUrl: 'forms/create',
            controller: 'FormEditCtrl'
        }).
        when("/forms/:id", {
            templateUrl: 'forms/view',
            controller: 'FormViewCtrl'
        }).
        when("/forms/edit/:id", {
            templateUrl: 'forms/edit',
            controller: 'FormEditCtrl'
        }).
        when("/studies", {
            templateUrl: 'studies',
            controller: 'StudyListCtrl'
        }).
        when("/studies/create", {
            templateUrl: 'studies/create',
            controller: 'StudyEditCtrl'
        }).
        when("/studies/:id", {
            templateUrl: 'studies/view',
            controller: 'StudyViewCtrl'
        }).
        when("/studies/edit/:id", {
            templateUrl: 'studies/edit',
            controller: 'StudyEditCtrl'
        }).
        when("/studies/publish/:id", {
            templateUrl: 'studies/publish',
            controller: 'StudyPublishCtrl'
        }).
        when("/studies/subjects/:id", {
            templateUrl: 'studies/subjects',
            controller: 'StudySubjectsCtrl'
        }).
        otherwise({
            redirectTo: "/"
        });

    app.controller = $controllerProvider.register;
}]);

// register the http interceptor which controls the spinner
// based on httpInterceptor code from http://docs.angularjs.org/api/ng.$http and
// http://stackoverflow.com/questions/18238227/delay-an-angular-js-http-service
app.config(['$httpProvider', function ($httpProvider) {
    $httpProvider.interceptors.push(['$q', '$timeout', function ($q, $timeout) {

        var show = 0,
            waitBeforeShow = 1000, // ms
            spinner = angular.element('.navbar-spinner');

        //shows the spinning element, but only if a request is still outstanding
        function showSpinner() {
            if (show > 0) {
                spinner.removeClass('invisible');
            }
        }

        // hides the spinning element
        function hideSpinner() {
            if (show > 0) {
                //decrement the number of requests that need spinning
                show--;

                if (show === 0) {
                    spinner.addClass('invisible');
                }
            }
        }

        return {
            //turns on the spinner when a request is about to be made
            "request": function (config) {
                //increment the number of requests made
                show++;
                $timeout(showSpinner, waitBeforeShow, false);

                return config || $q.when(config);
            },

            "requestError": function (rejection) {
                hideSpinner();
                return $q.reject(rejection);
            },

            // turns off the spinner when a response has been received
            "response": function (response) {
                hideSpinner();
                return response || $q.when(response);
            },

            "responseError": function (rejection) {
                hideSpinner();
                return $q.reject(rejection);
            }
        }
    }]);
}]);

//register an http interceptor that adds 'X-Requested-With': 'XMLHttpRequest' to all XHR requests
app.config(['$httpProvider', function ($httpProvider) {
    $httpProvider.interceptors.push(['$q', function ($q) {
        return {
            "request": function (config) {
                console.log('Setting X-Requested-With to', config.url);
                config.headers['X-Requested-With'] = 'XMLHttpRequest';

                return config || $q.when(config);
            }
        }
    }]);
}]);

app.run(function($rootScope /*, $location */) {
    console.log('app run');
    $rootScope.$on('$routeChangeError', function (nge, current, previous, rejection) {
        //when unauthorized, redirect to /login
        /*if (rejection !== undefined && rejection.status === 401) {
            $location.path("/login");
        }*/
        console.log('Route change error', rejection);
    });
});
