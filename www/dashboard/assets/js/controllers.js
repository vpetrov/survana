"use strict";

var dashboard = angular.module('dashboard', []);

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


dashboard.controller('StudyListCtrl', ['$scope', '$http',
    function StudyListCtrl($scope, $http) {
        $scope.studies = [];
        $scope.selected = [];
        $scope.message = '';
        $scope.search = '';
        $scope.loading = true;

        $http.get('studies/list').success(function (response, code, request) {
            if (response.success) {
                $scope.studies = response.message;
                $scope.loading = false;
            } else {
                console.log('Error message', response.message);
            }
        }).error(function () {
                $scope.loading = false;
                console.log("Error fetching studies/list");
            });

        //toggles a selected study on or off
        $scope.toggle = function (study_id) {
            var index = $scope.selected.indexOf(study_id);
            if (index === -1) {
                $scope.selected.push(study_id);
            } else {
                $scope.selected.splice(index,1);
            }
        };

        $scope.isSelected = function (study_id) {
            return ($scope.selected.indexOf(study_id) >= 0);
        };

        $scope.deleteStudy = function (study_id) {
            $scope.message = "";

            $http.delete('study', {params:{'id': study_id}}).success(function (response, code, request) {
                if (code === 204) {
                    $scope.removeStudy(study_id);
                } else {
                    $scope.message = 'Invalid response from server: ' + response;
                }
            }).error(function (response) {
                    $scope.message = response;
                    console.log('error', response);
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
                $scope.deleteStudy($scope.selected[i]);
            }
        };

        $scope.removeStudy = function (study_id) {
            for (var i = 0; i < $scope.studies.length; ++i) {
                if ($scope.studies[i].id === study_id) {
                    $scope.studies.splice(i, 1);
                    break;
                }
            }

            //remove it from selected
            var index = $scope.selected.indexOf(study_id);

            if (index >= 0) {
                $scope.selected.splice(index, 1);
            }
        }
    }
]);

dashboard.controller('StudyEditCtrl', ['$scope', '$http', '$window', '$location', '$routeParams',
    function StudyEditCtrl($scope, $http, $window, $location, $routeParams) {
        $scope.study = {
            name: "",
            title: "",
            description: "",
            version: "",
            forms: []
        };

        $scope.create = ($routeParams.id === undefined);
        $scope.forms = [];
        $scope.loading = false;
        $scope.message = "";

        //get all forms
        $http.get('forms/list').success(function (response, code, request) {
            if (response.success) {
                $scope.forms = response.message;

                //update the list of forms in the study, if the study info has been downloaded already
                if ($scope.study.forms.length) {
                    resolveStudyForms();
                }
            } else {
                console.log('Error message', response.message);
            }
        }).error(function () {
                console.log("Error fetching", $location.path())
            });

        //if we're editing a form, 'id' will be set
        if (!$scope.create) {
            //fetch the form JSON and store it in $scope.form
            $http.get('study', {params: $routeParams}).success(function (response, code, request) {
                if (response.success) {
                    $scope.study = response.message;

                    //update the list of forms in the study, if the form datastore is available
                    if ($scope.forms.length) {
                        resolveStudyForms();
                    }
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                    console.log("Error fetching", $location.path())
                });
        }

        //by default, the backend will store just pointers to the forms ({'id':form_id}).
        //this will dereference all pointers using $scope.forms.
        function resolveStudyForms() {
            var form, form_proxy, i;

            if (!$scope.study || !$scope.study.forms) {
                return;
            }

            //replace form stubs with actual forms, if they're present
            for (i = 0; i < $scope.study.forms.length; i++) {
                form_proxy = $scope.study.forms[i];

                //skip invalid entries
                if (!form_proxy || !form_proxy.id) {
                    continue;
                }

                form = findForm(form_proxy.id);

                if (form) {
                    $scope.study.forms[i] = form;
                }
            }
        }

        function findForm(form_id) {
            var form;

            for (var i = 0; i < $scope.forms.length; ++i) {
                form = $scope.forms[i];

                if (form && form.id === form_id) {
                    return form;
                }
            }
        }

        $scope.addForm = function (form_id) {

            var form = findForm(form_id);

            if (form) {
                if (!$scope.study.forms) {
                    $scope.study.forms = [form];
                } else {
                    $scope.study.forms.push(form);
                }
            }
        };

        $scope.removeForm = function (index) {
            try {
                $scope.study.forms.splice(index,1);
            } catch (e) {
                console.log('error', e);
            }
        };

        //if the save operation was successful
        function onSaveSuccess(response, code, request) {
            $scope.loading = false;

            var id;

            //no content = an 'edit' operation, therefore $scope.study should have a valid id
            if (code == 204 && ($scope.study.id !== undefined)) {
                id = $scope.study.id;
            } else if (response.success && response.message && response.message.id) {
                id = response.message.id;
            }

            //either redirect to view the new form, or show the message from server
            if (id) {
                $location.path('/studies/' + id);
            } else {
                $scope.message = response.message;
            }
        }

        //if the save operation failed
        function onSaveError(response) {
            $scope.loading = false;
            $scope.message = "Failed to save study (" + response + ")";
        }

        $scope.saveStudy = function () {
            //reset state
            $scope.message = "";

            if (!$scope.study.name.length) {
                $scope.message = 'Please enter a name for this study';
                return
            }

            if ($scope.create) {
                $http.post('studies/create', $scope.study).
                    success(onSaveSuccess).
                    error(onSaveError);
            } else {

                //create a copy of $scope.study, and replace all forms with stubs
                var study = {
                        name: $scope.study.name,
                        title: $scope.study.title,
                        description: $scope.study.description,
                        forms: []
                    },
                    i;

                //extract form IDs, since we don't actually want to send copies of the forms, just the form IDs
                for (i in $scope.study.forms) {
                    //replace each form with a form-proxy object, which just contains the form id
                    study.forms[i] = { 'id': $scope.study.forms[i].id };
                }

                $http.put('studies/edit', study, {params: $routeParams}).
                    success(onSaveSuccess).
                    error(onSaveError);
            }
        };

        $scope.discardStudy = function () {
            $window.history.back();
        };

        $scope.stopEvent = stopEvent;
    }
]);

dashboard.controller('StudyViewCtrl', ['$scope', '$window', '$location', '$routeParams', '$http', '$templateCache',
    function ($scope, $window, $location, $routeParams, $http, $templateCache) {
        $scope.study = {};
        $scope.forms = [];
        $scope.current = {
            index: 0,
            form: {}
        };
        $scope.size = 'M';
        $scope.template = null;
        $scope.theme = 'bootstrap';

        function fetchTemplate(theme_id, theme_version) {
            var url = 'theme?id=' + theme_id + '&version=' + theme_version + '&preview=true&study=true',
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

        function fetchStudy() {
            //fetch the form JSON and store it in $scope.form
            $http.get('study', {params: $routeParams}).success(function (response, code, request) {
                if (response.success) {
                    $scope.study = response.message;

                    //fetch form definitions for all the forms
                    if ($scope.study.forms.length) {
                        fetchForms($scope.study.forms);
                    }
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                    console.log("Error fetching", url)
                });
        }

        function fetchForms(form_list) {
            if (!form_list.length) {
                return
            }

            //get forms, including the fields
            $http.get('forms/list', {params:{fields:true}}).success(function (response, code, request) {
                if (response.success) {
                    $scope.forms = response.message;

                    //update the list of forms in the study, if the study info has been downloaded already
                    if ($scope.study.forms.length) {
                        resolveStudyForms();
                    }
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                    console.log("Error fetching", $location.path())
                });
        }

        //by default, the backend will store just pointers to the forms ({'id':form_id}).
        //this will dereference all pointers using $scope.forms.
        function resolveStudyForms() {
            var form, form_proxy, i;

            //replace form stubs with actual forms, if they're present
            for (i = 0; i < $scope.study.forms.length; i++) {
                form_proxy = $scope.study.forms[i];

                //skip invalid entries
                if (!form_proxy || !form_proxy.id) {
                    continue;
                }

                form = findForm(form_proxy.id);

                if (form) {
                    $scope.study.forms[i] = form;
                }
            }

            //as soon as the study forms are resolved, we can render the current form
            //if we update current.form, the watch will trigger the update
            $scope.current.form = $scope.study.forms[$scope.current.index];
        }

        function findForm(form_id) {
            var form;

            for (var i = 0; i < $scope.forms.length; ++i) {
                form = $scope.forms[i];

                if (form && form.id === form_id) {
                    return form;
                }
            }
        }

        $scope.resize = function (size) {
            $scope.size = size;
        };

        $scope.getStudyDate = function () {
                return (new Date($scope.study.created_on)).toLocaleDateString();
        };

        //when 'theme' changes, notify Survana
        $scope.$watch('theme', function (newTheme, oldTheme) {
            Survana.Theme.SetTheme(newTheme,
                function () {
                    fetchTemplate(newTheme, Survana.Version);
                    fetchStudy();
                },
                function () {
                    console.error('Failed to load Survana Themes!');
                });
        });

        //when the current index changes, change the current form
        $scope.$watch('current.index', function (newIndex, oldIndex) {
            if ($scope.study && $scope.study.forms && $scope.study.forms.length) {
                $scope.current.form = $scope.study.forms[newIndex];
            }
        });
    }
]);

dashboard.controller('StudyPublishCtrl', ['$scope', '$window', '$location', '$routeParams', '$http', '$templateCache',
    function StudyPublishCtrl($scope, $window, $location, $routeParams, $http, $templateCache) {
        $scope.study = null;
        $scope.study_url = null;
        $scope.forms = null;
        $scope.current = {
            index: -1,
            form: null,
            rendered: -1
        };
        $scope.rendered = null;

        $scope.template = null;
        $scope.theme = 'bootstrap';

        $scope.message = null;

        $scope.publishing = false;
        $scope.unpublishing = false;
        $scope.error = false;

        $scope.publishStudy = function() {
            $scope.publishing = true;
            $scope.current.index = 0;
            $scope.error = false;
            $scope.message = null;
            $scope.current.form = $scope.study.forms[0];
        };

        $scope.unpublishStudy = function () {
            $scope.message = "";

            $http.put('studies/edit', {published:false}, {params: $routeParams}).
                success(function (response, code, request) {
                    //we're done.
                    if (code === 204) {
                        $scope.study.published = false;
                    } else {
                        $scope.message = response.message;
                    }
                }).error(function (response) {
                    $scope.message = response;
                });
        };

        $scope.errorPublishing = function (message) {
            $scope.error = true;
            $scope.publishing = false;
            $scope.message = message;
        };

        $scope.finishPublishing = function () {
            $scope.publishing = false;
            $scope.current.index = null;
            $scope.current.form = null;
            $scope.message = null;
            $scope.error = false;
        };

        function nextForm() {
            //if there are forms left
            if (($scope.current.index + 1) < $scope.study.forms.length) {
                $scope.current.index++;
                $scope.current.form = $scope.study.forms[$scope.current.index];
            } else {
                //otherwise, we need to mark the study object as 'published' and save it
                $scope.study.published = true;
                saveStudy({published:true});
            }
        }

        //watch the numerical value of 'current.rendered'. The actual HTML data is going to be stored in $scope.rendered
        $scope.$watch('current.rendered', function (newVal, oldVal) {
            if (!$scope.rendered) {
                console.warn('No data was rendered.');
                return;
            }

            var url = "studies/publish?id=" + $scope.study.id + "&f=" + $scope.current.rendered;

            $http.post(url, $scope.rendered).success(function (response, code, request) {
                //go to the next form
                if (code === 204) {
                    nextForm();
                } else {
                    $scope.errorPublishing(response.message);
                }
            }).error(function (response) {
                $scope.errorPublishing(response);
            });
        });

        $scope.isCurrent = function (index) {
            return index === $scope.current.index;
        };

        $scope.selectLink = function (e) {
            var link = angular.element('a.study-link');
            $scope.selectNode(link[0]);
        };

        $scope.selectNode = function (node) {
            var range, selection;

            if (window.getSelection && document.createRange) {
                selection = window.getSelection();
                range = document.createRange();
                range.selectNodeContents(node);
                selection.removeAllRanges();
                selection.addRange(range);
            } else if (document.selection && document.body.createTextRange) {
                range = document.body.createTextRange();
                range.moveToElementText(node);
                range.select();
            }
        }

        function fetchTemplate(theme_id, theme_version) {
            var url = 'theme?id=' + theme_id + '&version=' + theme_version + '&publish=true&study=true',
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


        function copyStudy () {
            //create a copy of $scope.study, and replace all forms with stubs
            var study = {
                    name: $scope.study.name,
                    title: $scope.study.title,
                    description: $scope.study.description,
                    published: $scope.study.published,
                    version: $scope.study.version,
                    forms: []
                },
                i;

            //extract form IDs, since we don't actually want to send copies of the forms, just the form IDs
            for (i in $scope.study.forms) {
                //replace each form with a form-proxy object, which just contains the form id
                study.forms[i] = { 'id': $scope.study.forms[i].id };
            }

            return study;
        }

        function saveStudy (changes) {
            //reset state
            $scope.message = "";

            if (changes === undefined || !changes) {
                changes = $scope.study;
            }

            $http.put('studies/edit', changes, {params: $routeParams}).
                success(function (response, code, request) {
                    //we're done.
                    if (code === 204) {
                        $scope.finishPublishing();
                    } else {
                        $scope.errorPublishing(response.message);
                    }
                }).error(function (response) {
                    $scope.errorPublishing(response);
                });
        }


        function fetchStudy() {
            //fetch the form JSON and store it in $scope.form
            $http.get('study', {params: $routeParams}).success(function (response, code, request) {
                if (response.success) {
                    $scope.study = response.message;

                    //fetch form definitions for all the forms
                    if ($scope.study.forms.length) {
                        fetchForms($scope.study.forms);
                    }

                    //update study_url
                    $scope.study_url = $window.location.protocol + "//" + $window.location.host + "/study/?" + $scope.study.id;
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                    console.log("Error fetching", url)
                });
        }

        function fetchForms(form_list) {
            if (!form_list.length) {
                return
            }

            //get forms, including the fields
            $http.get('forms/list', {params:{fields:true}}).success(function (response, code, request) {
                if (response.success) {
                    $scope.forms = response.message;

                    //update the list of forms in the study, if the study info has been downloaded already
                    if ($scope.study.forms.length) {
                        resolveStudyForms();
                    }
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                    console.log("Error fetching", $location.path())
                });
        }

        //by default, the backend will store just pointers to the forms ({'id':form_id}).
        //this will dereference all pointers using $scope.forms.
        function resolveStudyForms() {
            var form, form_proxy, i;

            //replace form stubs with actual forms, if they're present
            for (i = 0; i < $scope.study.forms.length; i++) {
                form_proxy = $scope.study.forms[i];

                //skip invalid entries
                if (!form_proxy || !form_proxy.id) {
                    continue;
                }

                form = findForm(form_proxy.id);

                if (form) {
                    $scope.study.forms[i] = form;
                }
            }
        }

        function findForm(form_id) {
            var form;

            for (var i = 0; i < $scope.forms.length; ++i) {
                form = $scope.forms[i];

                if (form && form.id === form_id) {
                    return form;
                }
            }
        }

        //when 'theme' changes, notify Survana
        $scope.$watch('theme', function (newTheme, oldTheme) {
            Survana.Theme.SetTheme(newTheme,
                function () {
                    fetchTemplate(newTheme, Survana.Version);
                },
                function () {
                    console.error('Failed to load Survana Theme!');
                });
        });

        fetchStudy();
    }
]);

dashboard.controller('StudySubjectsCtrl', ['$scope', '$http', '$window', '$location', '$routeParams',
    function StudyEditCtrl($scope, $http, $window, $location, $routeParams) {

        //create a file upload field
        var fileUploader = document.createElement('input');

        fileUploader.setAttribute('type', 'file');
        fileUploader.classList.add('hidden');
        fileUploader.addEventListener('change', onFileUpload);

        document.body.appendChild(fileUploader);

        $scope.study = {
            id: $routeParams.id,
            name: "",
            title: "",
            description: "",
            version: "",
            forms: [],
            published: false,
            participants: {}
        };

        $scope.loading = false;
        $scope.message = "";

        $scope.stopEvent = stopEvent;
        $scope.search = "";

        function fetchStudy() {
            //fetch the form JSON and store it in $scope.form
            $http.get('study', {params: $routeParams}).success(function (response, code, request) {
                if (response.success) {
                    $scope.study = response.message;
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                    console.log("Error fetching", url)
                });
        }

        $scope.uploadFileDialog = function () {
            fileUploader.click();
        };

        function showError(msg) {
            console.error('Upload error:', msg);
            $scope.$apply(function () {
                $scope.message = msg;
                $scope.loading = false;
            });
        }

        function onFileUpload() {
            if (fileUploader.files.length == 0) {
                return
            }

            $scope.$apply(function () {
                $scope.loading = true;
            });

            //choose first file
            var file = fileUploader.files[0],
                freader = new FileReader();

            //when the file has been read, extract all the IDs as object keys
            freader.onloadend = function (e) {
                if (!e.total || !freader.result) {
                    showError("The file you selected is either empty, or could not be read.");
                    return;
                }

                //create an array of IDs
                var ids = freader.result.split("\n"),
                    result = [],
                    id,
                    i;

                if (!ids.length) {
                    showError("No participant IDs could be found in the selected file.");
                    return;
                }

                //trim and filter invalid IDs
                for (i = 0; i < ids.length; ++i) {
                    id = ids[i].trim();
                    if (!id.length) {
                        continue;
                    }

                    result.push(id.toUpperCase());
                }

                //send the IDs to the server
                uploadIDs(result);
            };

            freader.onerror = freader.onabort = function (e) {
                showError(freader.error);
            };

            freader.readAsText(file);
        }

        function uploadIDs(ids) {

            function onUploadSuccess(response) {
                console.log(arguments);
                $scope.loading = false;

                if (response.success) {
                    $scope.message = "";
                    //update current study subjects with the copy on the server
                    $scope.study.subjects = response.message;
                } else {
                    $scope.message = response.message;
                }
            }

            function onUploadError() {
                console.log(arguments);
                $scope.message = "Upload failed.";
            }

            $http.put('studies/subjects', ids, {params: $routeParams}).
                success(onUploadSuccess).
                error(onUploadError);
        }

        fetchStudy();
    }
]);

dashboard.controller('FormListCtrl', ['$scope', '$http',
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
                $scope.selected.splice(index,1);
            }
        };

        $scope.isSelected = function (form_id) {
            return ($scope.selected.indexOf(form_id) >= 0);
        };

        $scope.deleteForm = function (form_id) {
            $scope.message = "";

            $http.delete('form', {params:{'id': form_id}}).success(function (response, code, request) {
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

dashboard.controller('FormEditCtrl', ['$scope', '$http', '$window', '$location', '$timeout', '$routeParams',
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

dashboard.controller('FormViewCtrl', ['$scope', '$location', '$routeParams', '$http', '$templateCache',
    function ($scope, $location, $routeParams, $http, $templateCache) {

        $scope.form = {};
        $scope.size = 'M';
        $scope.theme = 'bootstrap';
        $scope.template = null;

        //expect the iframe to replace this function. TODO: figure out a better way of telling the iframe to call
        //window.validateForm()
        $scope.verifyForm = function () {}

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

/* DIRECTIVES */

dashboard.directive('ace', ['$timeout', function ($timeout) {
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
                    var code = JSON.stringify(ngModel.$viewValue,null,4);
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

dashboard.directive("questionnaire", ['$window', '$compile', '$timeout', function ($window, $compile, $timeout) {
    return {
        restrict: 'A',
        require: '?ngModel',
        scope: false,
        link: function (scope, elem, attrs, ngModel) {

            function updateTemplate() {
                $compile(elem.contents())(scope);

                //re-render the model
                ngModel.$render();
            }

            //quick hack to pass this value form the scope to the iframe
            $window.study_id = function () {
                if (scope.study) {
                    return scope.study.id;
                }

                return null;
            };

            //register a NextPage() function that can be called within the questionnaire preview iframe
            $window.NextPage = function () {
                $timeout(function () {
                    if ((scope.current.index + 1) < scope.study.forms.length) {
                        scope.current.index++;
                    }
                });
            };

            $window.FinishSurvey = function () {
                $timeout(function () {
                    scope.current.index = 0;
                });
            };

            scope.verifyForm = function () {
                if (!elem || !elem[0] || !elem[0].contentWindow || !elem[0].contentWindow.validateForm) {
                    return
                }

                elem[0].contentWindow.validateForm();
            };

            //when current_form changes, update the template
            scope.$watch('current.index', updateTemplate);

            scope.$watch('template', function(val) {

                //nothing to do?
                if (!val) {
                    return
                }

                var frame = elem[0],
                    doc = frame.contentDocument || frame.contentWindow.document;

                //document.write() is the fastest way to update the contents.
                doc.open();
                doc.write(scope.template);
                doc.close();

                updateTemplate();
            });

            //update the view
            ngModel.$render = function () {

                var frame = elem[0],
                    doc = frame.contentDocument || frame.contentWindow.document,
                    node = doc.getElementById('content'),
                    validation_node = doc.getElementById('validation'),
                    result,
                    validation;


                //make sure a theme, a template and a rendering node are available
                if (!Survana.Theme || !scope.template || !node) {
                    return;
                }

                result = Survana.Theme.Questionnaire(ngModel.$viewValue);

                if (result) {
                    if (node.hasChildNodes()) {
                        node.removeChild(node.firstChild);
                    }

                    //append the form
                    node.appendChild(result);

                    //validation configuration
                    validation = Survana.Validation.ExtractConfiguration(ngModel.$viewValue);

                    if (validation_node && validation) {
                        validation_node.innerHTML = validation;
                        //rely on the fact that template will include a 'startValidation()' function
                        if (frame.contentWindow.startValidation) {
                            frame.contentWindow.startValidation();
                        }
                    }

                    //if we're supposed to save rendered data
                    if (attrs['render']) {
                        //skip a digest cycle to let the updateTemplate() digest to finish,
                        //otherwise, innerHTML is still the old html, before any watches are updated by the new changes
                        $timeout(function () {
                            //store the HTML data into the variable pointed to by data-render
                            scope[attrs['render']] = "<!DOCTYPE html><html>" + doc.documentElement.innerHTML + "</html>";
                            //update the currently rendered form index
                            scope.current.rendered = scope.current.index;
                        });
                    }
                }
            }
        }
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

                var src_index = e.originalEvent.dataTransfer.getData('text/plain')|0,
                    dest_index = elem.attr('data-list-index')|0;

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
