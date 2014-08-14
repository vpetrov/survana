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

    dashboard.directive('draggable', [function () {
        return {
            restrict: 'A',
            link: function ($scope, element) {
                $scope.$emit('draggable-element', element);
            }
        }
    }]);

    dashboard.directive('draggableList', ['$window', function (window) {
        return {
            restrict: 'A',
            require: '?ngModel',
            link: function (scope, elem, attrs, ngModel) {

                var inside = null,
                    window_handler = null,
                    drag_item = null;

                if (ngModel === undefined) {
                    console.error('draggable must have a scope model', elem);
                }

                function onDragOut(e) {
                    if (inside) {
                        drag_item.classList.add('dragremove');
                        inside = null;
                    }
                }

                function onDragRemove(e) {
                    e.preventDefault();
                    e.stopPropagation();

                    var src_index = drag_item.getAttribute('data-list-index') | 0;

                    scope.$apply(function () {
                        ngModel.$viewValue.splice(src_index, 1);
                    });

                    return false;
                }

                function onWindowDragOver(e) {
                    e.preventDefault();
                    e.stopPropagation();
                    return false;
                }

                function onWindowDragEnd(e) {
                    e.preventDefault();
                    e.stopPropagation();
                    return false;
                }

                function onDragStart(e) {
                    if (!window_handler) {
                        window.addEventListener('dragenter', onDragOut);
                        window.addEventListener('drop', onDragRemove);
                        window.addEventListener('dragover', onWindowDragOver);
                        window.addEventListener('dragend', onWindowDragEnd);
                        window_handler = true;
                    }

                    var src = angular.element(e.currentTarget);

                    if (e.originalEvent && e.originalEvent.dataTransfer) {
                        e.originalEvent.dataTransfer.effectAllowed = 'move';
                        e.originalEvent.dataTransfer.setData('text/plain', src.attr('data-list-index'));
                        e.originalEvent.dataTransfer.setDragImage(e.currentTarget, 0, 0);
                    }

                    drag_item = e.currentTarget;
                    drag_item.classList.add('dragstart');
                }

                function onDragEnter(e) {
                    e.stopPropagation();
                    e.preventDefault();

                    if (e.currentTarget !== e.target || inside === e.currentTarget) {
                        return false;
                    }

                    if (e.originalEvent && e.originalEvent.dataTransfer) {
                        e.originalEvent.dataTransfer.dropEffect = 'move';
                    }

                    inside = e.currentTarget;
                    inside.classList.add('dragover');
                    drag_item.classList.remove('dragremove');

                    return false;
                }

                function onDragOver(e) {
                    event.preventDefault();
                    event.stopPropagation();
                    return false;
                }

                function onDragLeave(e) {
                    if (inside === e.currentTarget) {
                        return;
                    }

                    e.currentTarget.classList.remove('dragover');
                }

                function onDragDrop(e) {
                    var src_index = e.originalEvent.dataTransfer.getData('text/plain') | 0,
                        dest_index = inside.getAttribute('data-list-index') | 0;

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

                    e.preventDefault();
                    return false;
                }

                function onDragEnd(e) {
                    window.removeEventListener('dragenter', onDragOut);
                    window.removeEventListener('dragover', onWindowDragOver);
                    window.removeEventListener('dragend', onWindowDragEnd);
                    window.removeEventListener('drop', onDragRemove);

                    if (inside) {
                        inside.classList.remove('dragover', 'dragstart', 'dragremove');
                    }

                    drag_item.classList.remove('dragover', 'dragstart', 'dragremove');
                    drag_item = null;
                    window_handler = false;
                }

                scope.$on('draggable-element', function (e, child) {
                    child.on('dragstart', onDragStart);
                    child.on('dragenter', onDragEnter);
                    child.on('dragover', onDragOver);
                    child.on('dragleave', onDragLeave);
                    child.on('drop', onDragDrop);
                    child.on('dragend', onDragEnd);
                });
            }
        }
    }]);

})();