(function () {
    "use strict";

    var dashboard = angular.module('dashboard', ['studies', 'forms']);

    dashboard.controller('HomeCtrl', ['$scope', '$http',
        function HomeCtrl($scope, $http) {
            console.log('HomeCtrl running')
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

    dashboard.directive('loading', [function () {
        return function (scope, elem, attrs) {
            var expr = attrs.loading;

            //add a watch to monitor the value of the expression
            scope.$watch(expr, function (val) {
                if (val) {
                    elem.prop('disabled', true);
                } else {
                    elem.prop('disabled', false);
                }
            });
        }
    }]);

    dashboard.directive("draggable", ['$window', function ($window) {
        return {
            restrict: 'A',
            require: '?ngModel',
            scope: false,
            link: function (scope, elem, attrs, ngModel) {
                var node = elem[0];

                if (ngModel === undefined) {
                    console.error('draggable must have a scope model', elem);
                }

                function onDragStart(e) {
                    var src = angular.element(e.currentTarget);


                    if (e.originalEvent && (e.originalEvent.dataTransfer !== undefined)) {
                        e.originalEvent.dataTransfer.effectAllowed = 'move';
                        e.originalEvent.dataTransfer.setData('text/plain', src.attr('data-list-index'));
                        e.originalEvent.dataTransfer.setDragImage(e.currentTarget, 0, 0);
                    }

                    e.currentTarget.style.opacity = '0.5'; //decrease opacity
                }

                var dd = 0;

                //add 'dragover' class
                function onDragEnter(e) {
                    e.currentTarget.classList.add('dragover');
                }

                function onDragOver(e) {

                    if (!e.currentTarget.classList.contains('dragover')) {
                        e.currentTarget.classList.add('dragover');
                    }

                    stopEvent(e);

                    if (e.originalEvent.dataTransfer) {
                        e.originalEvent.dataTransfer.dropEffect = 'move';
                    }
                }

                //remove 'dragover' class
                function onDragLeave(e) {
                    e.currentTarget.classList.remove('dragover');
                }

                function onDragEnd(e) {
                    e.currentTarget.style.opacity = '1.0';
                }

                function onDragDrop(e) {

                    stopEvent(e);

                    var src_index = e.originalEvent.dataTransfer.getData('text/plain') | 0,
                        dest_index = elem.attr('data-list-index') | 0;

                    e.currentTarget.classList.remove('dragover');

                    scope.$apply(function () {

                        var delta = 0,
                            temp = ngModel.$viewValue[src_index];

                        //if source was moved up, the new element has shifted this index by 1
                        if (src_index < dest_index) {
                            dest_index += 1;
                        } else {
                            src_index += 1;
                        }

                        ngModel.$viewValue.splice(dest_index, 0, temp);
                        ngModel.$viewValue.splice(src_index, 1);
                    });


                    return false;
                }

                //attach events
                elem.on('dragstart', onDragStart);
                elem.on('dragenter', onDragEnter);
                elem.on('dragover', onDragOver);
                elem.on('dragleave', onDragLeave);
                elem.on('dragend', onDragEnd);
                elem.on('drop', onDragDrop);

                elem.parent().addClass('draggable-container');
            }
        }
    }]);


//TODO: provide this as a service
    function stopEvent($e) {
        if ($e.stopPropagation) {
            $e.stopPropagation();
        }

        if ($e.preventDefault) {
            $e.preventDefault();
        }

        $e.cancelBubble = true;
        $e.returnValue = false;
    }

})();