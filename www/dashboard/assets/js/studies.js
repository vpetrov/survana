(function () {
    var app = angular.module("studies", ['preview']);

    app.controller('StudyListCtrl', ['$scope', '$http',
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
                    $scope.selected.splice(index, 1);
                }
            };

            $scope.isSelected = function (study_id) {
                return ($scope.selected.indexOf(study_id) >= 0);
            };

            $scope.deleteStudy = function (study_id) {
                $scope.message = "";

                $http.delete('study', {params: {'id': study_id}}).success(function (response, code, request) {
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

    app.controller('StudyEditCtrl', ['$scope', '$http', '$window', '$location', '$routeParams',
        function StudyEditCtrl($scope, $http, $window, $location, $routeParams) {
            $scope.study = {
                name: "",
                title: "",
                description: "",
                version: "",
                form_ids: []
            };

            $scope.create = ($routeParams.id === undefined);
            $scope.forms = [];
            $scope.loading = false;
            $scope.message = "";

            $scope.changed = false;
            $scope.study_forms = [];

            //get all forms
            $http.get('forms/list').success(function (response, code, request) {
                if (response.success) {
                    $scope.forms = response.message;

                    //update the list of forms in the study, if the study info has been downloaded already
                    if ($scope.study.form_ids.length) {
                        resolveStudyForms();
                    }
                } else {
                    console.log('Error message', response.message);
                }
            }).error(function () {
                console.log("Error fetching", $location.path())
            });

            //if we're editing a form, the 'id' route param will be set
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
                var form, form_id, i;

                if ($scope.study.form_ids == undefined) {
                    $scope.study.form_ids = [];
                }

                //replace form stubs with actual forms, if they're present
                for (i = 0; i < $scope.study.form_ids.length; i++) {
                    form_id = $scope.study.form_ids[i];

                    //skip invalid entries
                    if (!form_id) {
                        continue;
                    }

                    form = findForm(form_id);

                    if (form) {
                        $scope.study_forms[i] = form;
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

                if (!form) {
                    return
                }

                if (!$scope.study.form_ids) {
                    $scope.study.form_ids = [];
                }

                $scope.study.form_ids.push(form_id);
                $scope.study_forms.push(form);
                $scope.changed = true;
            };

            $scope.removeForm = function (index) {
                try {
                    $scope.study.form_ids.splice(index, 1);
                    $scope.study_forms.splice(index, 1);
                } catch (e) {
                    console.log('error', e);
                }

                $scope.changed = true;
            };

            //if the save operation was successful
            function onSaveSuccess(response, code, request) {
                $scope.loading = false;
                $scope.changed = false;

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
                    $http.put('studies/edit', $scope.study, {params: $routeParams}).
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

    app.controller('StudyViewCtrl', ['$scope', '$window', '$location', '$routeParams', '$http', '$templateCache',
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

            $scope.study_forms = [];

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
                        if ($scope.study.form_ids.length) {
                            fetchForms($scope.study.form_ids);
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
                $http.get('forms/list', {params: {fields: true}}).success(function (response, code, request) {
                    if (response.success) {
                        $scope.forms = response.message;

                        //update the list of forms in the study, if the study info has been downloaded already
                        if ($scope.study.form_ids.length) {
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
                var form, form_id, i;

                //replace form stubs with actual forms, if they're present
                for (i = 0; i < $scope.study.form_ids.length; i++) {
                    form_id = $scope.study.form_ids[i];

                    //skip invalid entries
                    if (!form_id) {
                        continue;
                    }

                    form = findForm(form_id);

                    if (form) {
                        $scope.study_forms[i] = form;
                    }
                }

                //as soon as the study forms are resolved, we can render the current form
                //if we update current.form, the watch will trigger the update
                $scope.current.form = $scope.study_forms[$scope.current.index];
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
                if ($scope.study && $scope.study_forms && $scope.study_forms.length) {
                    $scope.current.form = $scope.study_forms[newIndex];
                }
            });
        }
    ]);

    app.controller('StudyPublishCtrl', ['$scope', '$window', '$location', '$routeParams', '$http', '$templateCache',
        function StudyPublishCtrl($scope, $window, $location, $routeParams, $http, $templateCache) {
            $scope.study = null;
            $scope.study_url = null;
            $scope.forms = null;
            $scope.study_forms = [];
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

            $scope.publishStudy = function () {
                $scope.publishing = true;
                $scope.current.index = 0;
                $scope.error = false;
                $scope.message = null;
                $scope.current.form = $scope.study_forms[0];
            };

            $scope.unpublishStudy = function () {
                $scope.message = "";

                $http.put('studies/edit', {published: false}, {params: $routeParams}).
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
                if (($scope.current.index + 1) < $scope.study_forms.length) {
                    $scope.current.index++;
                    $scope.current.form = $scope.study_forms[$scope.current.index];
                } else {
                    //otherwise, we need to mark the study object as 'published' and save it
                    $scope.study.published = true;
                    saveStudy({published: true});
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


            function copyStudy() {
                //create a copy of $scope.study, and replace all forms with stubs
                var study = {
                        name: $scope.study.name,
                        title: $scope.study.title,
                        description: $scope.study.description,
                        published: $scope.study.published,
                        version: $scope.study.version,
                        form_ids: []
                    },
                    i;

                //extract form IDs, since we don't actually want to send copies of the forms, just the form IDs
                for (i in $scope.study.form_ids) {
                    //replace each form with a form-proxy object, which just contains the form id
                    study.form_ids[i] = { 'id': $scope.study.form_ids[i].id };
                }

                return study;
            }

            function saveStudy(changes) {
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
                //fetch the study JSON and store it in $scope.study
                $http.get('study', {params: $routeParams}).success(function (response, code, request) {
                    if (response.success) {
                        $scope.study = response.message;

                        //fetch form definitions for all the forms
                        if ($scope.study.form_ids.length) {
                            fetchForms($scope.study.form_ids);
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
                $http.get('forms/list', {params: {fields: true}}).success(function (response, code, request) {
                    if (response.success) {
                        $scope.forms = response.message;

                        //update the list of forms in the study, if the study info has been downloaded already
                        if ($scope.study.form_ids.length) {
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
                var form, form_id, i;

                //replace form stubs with actual forms, if they're present
                for (i = 0; i < $scope.study.form_ids.length; i++) {
                    form_id = $scope.study.form_ids[i];

                    //skip invalid entries
                    if (!form_id) {
                        continue;
                    }

                    form = findForm(form_id);

                    if (form) {
                        $scope.study_forms[i] = form;
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

    app.controller('StudySubjectsCtrl', ['$scope', '$http', '$window', '$location', '$routeParams',
        function StudySubjectsCtrl($scope, $http, $window, $location, $routeParams) {

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
})();