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
                name: "Untitled Study",
                title: "Untitled Study",
                description: "",
                version: "1",
                form_ids: []
            };

            $scope.create = ($routeParams.id === undefined);
            $scope.forms = [];
            $scope.loading = false;
            $scope.message = "";

            $scope.changed = false;
            $scope.study_forms = [];

            $scope.ready = false;

            var downloading = true;

            //get all forms
            $http.get('forms/list').success(function (response, code, request) {
                if (response.success) {
                    $scope.forms = response.message;

                    //update the list of forms in the study, if the study info has been downloaded already
                    if ($scope.study.form_ids.length) {
                        resolveStudyForms();
                    }

                    if ($scope.create) {
                        $scope.ready = true;
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

                $scope.study_forms = [];

                //replace form stubs with actual forms, if they're present
                for (i = 0; i < $scope.study.form_ids.length; i++) {
                    form_id = $scope.study.form_ids[i];

                    //skip invalid entries
                    if (!form_id) {
                        continue;
                    }

                    form = findForm(form_id);

                    if (form) {
                        $scope.study_forms.push(form);
                    }
                }

                $scope.ready = true;
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

            var remove_ready_watch = $scope.$watch('ready', function (val) {
                if (val) {
                    remove_ready_watch();

                    $scope.$watchCollection('study.form_ids', function (new_ids, old_ids) {
                        if (downloading) {
                            downloading = false;
                            return;
                        }
                        console.log('form_ids have changed! new=', new_ids, 'old=', old_ids);
                        resolveStudyForms();
                        $scope.changed = true;
                    });
                }
            });
        }
    ]);

    app.controller('StudyViewCtrl', ['$scope', '$window', '$location', '$routeParams', '$http', '$templateCache',
        function ($scope, $window, $location, $routeParams, $http, $templateCache) {
            $scope.study = {};
            $scope.forms = [];
            $scope.size = 'M';
            $scope.template = null;
            $scope.theme = null;
            $scope.current = {
                index: 0,
                form: null,
                html: ""
            };

            $scope.study_forms = [];

            var previewWindow = null;

            function fetchTemplate(theme_id, theme_version) {
                console.log('fetch template', theme_id, theme_version);
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
                console.log('fetchStudy');
                //fetch the form JSON and store it in $scope.form
                $http.get('study', {params: $routeParams}).success(function (response, code, request) {
                    if (response.success) {
                        $scope.study = response.message;

                        //fetch form definitions for all the forms
                        if ($scope.study.form_ids.length) {
                            fetchForms($scope.study.form_ids);
                        }

                        $scope.theme = 'bootstrap';
                    } else {
                        console.log('Error message', response.message);
                    }
                }).error(function (err) {
                    console.log("Error fetching", url, err)
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

                //update the current form and the max forms in the preview
                if ($scope.current.index !== null) {
                    $scope.current.form = $scope.study_forms[$scope.current.index];
                    updatePreviewForms();
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

            function updatePreviewForms() {
                if (!previewWindow || !$scope.study.form_ids) {
                    return;
                }
                //update the preview's study progress
                previewWindow.Survana.Workflow.SetProgress($scope.current.index|0, $scope.study.form_ids.length);
            }

            $scope.resize = function (size) {
                $scope.size = size;
            };

            $scope.getStudyDate = function () {
                if (!$scope.study.created_on) {
                    return null;
                }

                return (new Date($scope.study.created_on)).toLocaleDateString();
            };

            //when 'theme' changes, notify Survana
            $scope.$watch('theme', function (newTheme, oldTheme) {
                if (!newTheme) {
                    return;
                }

                console.log('set new theme to', newTheme);
                Survana.Theme.SetTheme(newTheme,
                    function () {
                        fetchTemplate(newTheme, Survana.Version);
                        fetchStudy();
                    },
                    function (err) {
                        console.error('Failed to load Survana theme "' + newTheme + '": ', err);
                    });
            });

            //@note callback parameters: e, formWindow
            $scope.$on('form:load', function (e, newFormWindow) {
                previewWindow = newFormWindow;
                updatePreviewForms();
            });

            //@note callback parameters: e, formWindow
            $scope.$on('form:next', function (e) {
                $scope.$apply(function () {
                    var new_index = ($scope.current.index|0) + 1;
                    if (new_index >= $scope.study.form_ids.length) {
                        new_index = 0;
                    }
                    //go to the next form (or to the first form)
                    $scope.current.index = new_index;
                });
            });

            //when the current index changes, change the current form
            $scope.$watch('current.index', function (newIndex, oldIndex) {
                if ($scope.study && $scope.study_forms && $scope.study_forms.length) {
                    $scope.current.form = $scope.study_forms[newIndex];
                    updatePreviewForms();
                }
            });

            fetchStudy();
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

            $scope.template = null;
            $scope.theme = 'bootstrap';

            $scope.message = null;

            $scope.publishing = false;
            $scope.unpublishing = false;
            $scope.error = false;

            $scope.startPublishingStudy = function () {
                $scope.publishing = true;
                $scope.current.index = 0;
                $scope.error = false;
                $scope.message = null;
                $scope.current.form = $scope.study_forms[0];
            };

            $scope.publishStudy = function () {
                $scope.message = "";

                $http.post('study/publish', null, {params: $routeParams}).
                    success(function (response, code, request) {
                        //we're done.
                        if (code === 200) {
                            $scope.study.published = response.message.published;
                            $scope.study.revision = response.message.revision;
                            $scope.finishPublishing();
                        } else {
                            $scope.message = response.message;
                        }
                    }).error(function (response) {
                        $scope.message = response;
                    });
            };

            $scope.unpublishStudy = function () {
                $scope.message = "";

                $http.post('study/unpublish', null, {params: $routeParams}).
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

            var htmlCache = {};

            function nextForm() {
                //if there are forms left
                if (($scope.current.index + 1) < $scope.study_forms.length) {
                    $scope.current.index++;
                    //check the html cache
                    var current_form_id = $scope.study_forms[$scope.current.index].id;

                    if (htmlCache[current_form_id] !== undefined) {
                        saveForm(htmlCache[current_form_id]);
                    } else {
                        //change current form so that the preview can update itself
                        $scope.current.form = $scope.study_forms[$scope.current.index];
                    }
                } else {
                    $scope.publishStudy();
                }
            }

            function saveForm(html) {
                var url = "studies/publish?id=" + $scope.study.id + "&f=" + $scope.current.index;

                $http.post(url, html).success(function (response, code, request) {
                    //go to the next form
                    if (code === 204) {
                        nextForm();
                    } else {
                        $scope.errorPublishing(response.message);
                    }
                }).error(function (response) {
                    $scope.errorPublishing(response);
                });
            }

            //watch the numerical value of 'current.rendered'. The actual HTML data is going to be stored in $scope.rendered
            $scope.$on('form:rendered', function (e, html) {
                htmlCache[$scope.current.form.id] = html;
                saveForm(html);
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
            };

            function fetchTemplate(theme_id, theme_version) {
                var url = 'theme?id=' + theme_id + '&version=' + theme_version + '&study=true',
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
                //fetch the study JSON and store it in $scope.study
                $http.get('study', {params: $routeParams}).success(function (response, code, request) {
                    if (response.success) {
                        $scope.study = response.message;

                        //fetch form definitions for all the forms
                        if ($scope.study.form_ids.length) {
                            fetchForms($scope.study.form_ids);
                        }

                        //update study_url
                        $scope.study_url = $window.location.protocol + "//" + $window.location.host + "/study/" + $scope.study.id;
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