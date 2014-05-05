/* survana-storage.js

Survana.Storage is an abstraction of persistent storage available in browsers. It can make use of localStorage,
IndexedDB or FileWriter, and will automatically select the best adapter. To override the automatic selection, set
Survana.Storage = {} before loading this file. The API is fully asynchronous (even for localStorage), and all methods
take a success and error callback.

    Survana.Storage.Name                            - (string)  Name of the storage adapter
    Survana.Storage.IsAvailable                     - (boolean) Whether any storage is available
    Survana.Storage.Get(key, success, error)        - Returns an object or value identified by 'key'.
                                                      Returns multiple values if 'key' is an object with keys/defaults
    Survana.Storage.Set(key, value, success, error) - Saves an object 'value' identified by 'key'.
                                                      Can save multiple values if 'key' is an object with keys/values.
    Survana.Storage.Save(object, success, error)    - Saves multiple keys and values stored in 'object'.
                                                      This is a shortcut for Set(object, null, success, error)
    Survana.Storage.Remove(key, success, error)     - Removes the object identified by 'key'.
                                                      To remove multiple objects, pass an object with keys/values.
    Survana.Storage.SetScope(scope)                 - Sets a scope for all 'key' parameters of Get/Set/Remove

Note: currently, only the LocalStorage adapter is implemented.
Dependencies: none

@author Victor Petrov <victor_petrov@harvard.edu>
@license BSD
@date 05/01/2014
*/

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

(function (Survana) {
    /*
     * LOCAL STORAGE adapter - generally limited to 5MB or so.
     * The API is using callbacks for consistency with other adapters.
     */
    var local_storage = {
        'Name': 'LocalStorage',
        'Get': local_storage_get,
        'Set': local_storage_set,
        'IsAvailable': is_local_storage_available,
        'Remove': local_storage_remove
    };

    //no scope by default
    var scope = '';

    /** Sets a scope for all storage keys
     * @param p
     */
    function set_scope(p) {
        if (p) {
            scope = p + "-";
            console.log('Survana.Storage: new scope:', scope);
        }
    }

    /** Checks to see if window.localStorage is available (and works)
     * @param success   {Function} A success callback that receives a boolean result: true/false
     * @param error     {Function} A callback that handles any errors during testing
     */
    function is_local_storage_available(success, error) {

        if (window === undefined || window.localStorage === undefined || !window.localStorage) {
            //no, localStorage is not available
            return success && success(false);
        }

        var key     = 'survana-test',
            value   = 'test',
            value2;

        //at least in Safari 7 in Private Browsing mode, localStorage is defined, but throws an error when accessed.
        //this will perform a short test to see if we can actually read from and write to localStorage
        try {
            //save value
            window.localStorage[key] = value;

            //read value
            value2 = window.localStorage[key];

            //remove key
            delete window.localStorage[key];

            //check to make sure we read back what we wrote
            if (value2 !== value) {
                //no, localStorage is not available
                return success && success(false);
            }
        } catch (e) {
            //no, localStorage is not available
            return error && error(e);
        }

        //yes, localStorage is available
        success && success(true);
    }

    /** Returns multiple items by id from localStorage. e.g. {'key1':0, 'key2': "abc"} would attempt to fetch 2 keys and
     * would set their default values to 0 and "abc", respectively.
     * @param keys      {Object}    An object specifying the keys to fetch and their default values.
     * @param success   {Function}  The callback function that receives the result on success
     * @param error     {Function}  The error callback.
     */
    function local_storage_get_multi(keys, success, error) {
        var result;

        try {
            for (var id in keys) {
                if (!keys.hasOwnProperty(id)) {
                    continue;
                }

                result = localStorage[scope + id];

                //skip values that aren't available
                if (result === undefined) {
                    continue;
                }

                keys[id] = JSON.parse(result);
            }

        } catch (e) {
            return error && error(e);
        }

        success && success(keys);
    }

    /** Returns one or more items by id(s) from localStorage
     * @param key       {String|Object} The id of the object to retrieve, or an object with keys and default values
     * @param success   {Function}      The callback function that receives the result on success
     * @param error     {Function}      The error callback.
     */
    function local_storage_get(key, success, error) {
        var result;

        //handle requests for multiple items separately
        if (typeof key === "object") {
            return local_storage_get_multi(key, success, error);
        }

        try {
            result = localStorage[scope + key];
            //always return null if the value doesn't exist
            if (result === undefined) {
                result = null;
            } else {
                //otherwise, treat it as JSON
                result = JSON.parse(result);
            }
        } catch (e) {
            return error && error(e);
        }

        success && success(result);
    }

    /** Saves multiple values in localStorage based on the keys object
     * @param obj       {Object}    An object whose keys and values should be saved
     * @param success   {Function}  The success function. It receives boolean true on success, false on failure.
     * @param error     {Function}  The error callback
     */
    function local_storage_set_multi(obj, success, error) {
        var key,
            value;

        try {
            for (key in obj) {
                if (!obj.hasOwnProperty(key)) {
                    continue;
                }
                value = obj[key];
                localStorage[scope + key] = JSON.stringify(value);
            }
        } catch (e) {
            return error && error(e);
        }

        success && success(true);
    }

    /** Saves a value in localStorage under the given 'key' id.
     * @param key       {String|Object} The id of the object to save, or an object with keys and values to save.
     * @param value     {*}             The object (any JSON value) to save. Optional if 'key' is an Object.
     * @param success   {Function}      The success function. It receives boolean true on success, false on failure.
     * @param error     {Function}      The error callback
     */
    function local_storage_set(key, value, success, error) {

        if (typeof key === "object") {
            return local_storage_set_multi(key, success, error);
        }

        try {
            localStorage[scope + key] = JSON.stringify(value);
        } catch (e) {
            return error && error(e);
        }

        success && success(true);
    }

    /** Removes multiple objects from localStorage based on the keys of the 'obj' parameter
     * @param obj       {Object}    An object which holds the keys to remove
     * @param success   {Function}  The success callback
     * @param error     {Function}  The error callback
     */
    function local_storage_remove_multi(obj, success, error) {

        try {
            for (var key in obj) {
                if (!obj.hasOwnProperty(key)) {
                    continue;
                }

                delete localStorage[scope + key];
            }
        } catch (e) {
            return error && error(e);
        }

        success && success(true);
    }

    /** Removes an object from localStorage based on the key id, multiple keys can be removed if 'key' is an object
     * @param key       {String|Object} The id of the object to remove, or an object with keys to remove.
     * @param success   {Function}      The success callback
     * @param error     {Function}      The error callback
     */
    function local_storage_remove(key, success, error) {

        if (typeof key === "object") {
            return local_storage_remove_multi(key, success, error);
        }

        try {
            delete localStorage[scope + key];
        } catch (e) {
            return error && error(e);
        }

        success && success(true);
    }

    //order of priority for storage adapters
    var adapters = [
            local_storage/*,
             indexed_db,
             file_writer*/
        ],
        adapter;

    /** Loops over all adapters (recursively) until it finds one whose IsAvailable() method returns true.
     * @param adapter_list  {Array}     A list of adapters to iterate over, sorted by decreasing order of preference
     * @param index         {Number}    The index to start searching from (or current index to look at, if recursing)
     * @param success       {Function}  The success callback
     * @param error         {Function}  The error callback
     */
    function find_best_adapter(adapter_list, index, success, error) {
        if (!adapter_list) {
            return;
        }

        //if we've exhausted the adapter list, simply call the success function with no result
        if (index >= adapter_list.length) {
            return success && success();
        }

        var a = adapter_list[index];

        if (!a) {
            return;
        }

        //check to see if this adapter is available
        a.IsAvailable(
            //success
            function (result) {
                if (result) {
                    //found a suitable adapter
                    success && success (a);
                } else {
                    //continue searching
                    find_best_adapter(adapter_list, index + 1, success, error);
                }
            },
            error
        );
    }

    /** Calls the adapter's Get method after performing parameter checking
     * @param key       {String|Object} The id of the value, or multiple keys grouped as an Object
     * @param success   {Function}      The success function; receives the value as its parameter
     * @param error     {Function}      The error callback
     */
    function storage_get(key, success, error) {
        //check that a valid key has been specified
        if (!key) {
            return error && error(new Error("Storage.Get: a key is required to retrieve a storage value."));
        }

        //check that an adapter exists
        if (!adapter) {
            return error && error(new Error("No Storage adapter available."));
        }

        //call adapter.Get
        return adapter.Get.apply(this, arguments);
    }

    /** Calls the adapter's Set method after performing parameter checking
     * @param key       {String|Object} An identifier associated with the value
     * @param value     {*}             The value to save. Optional if key is an Object containing keys/values.
     * @param success   {Function}      The success callback
     * @param error     {Function}      The error callback
     */
    function storage_set(key, value, success, error) {
        //check that a valid key has been specified
        if (!key) {
            return error && error(new Error("Storage.Set: a valid key is required."));
        }

        //check that the value is not undefined
        if (value === undefined) {
            return error && error(new Error("Storage.Set: a valid value is required."))
        }

        //check that an adapter exists
        if (!adapter) {
            return error && error(new Error("No Storage adapter available."));
        }

        //call adapter.Set
        return adapter.Set.apply(this, arguments);
    }

    /** A convenient shortcut to Set, for storing objects with keys and values. This method calls Set, but sets the
     * 'value' parameter to null, so that callers don't have to.
     * @param obj       {Object}    The object to save.
     * @param success   {Function}  The success callback
     * @param error     {Function}  The error callback
     */
    function storage_save(obj, success, error) {
        //check that a valid key has been specified
        if (!obj) {
            return error && error(new Error("Storage.Save: a valid object is required."));
        }

        //check that an adapter exists
        if (!adapter) {
            return error && error(new Error("No Storage adapter available."));
        }

        //call adapter.Set with a null value, since 'obj' will contain all keys/values.
        return adapter.Set(obj, null, success, error);
    }

    /** Calls the adapter's Remove method after performing parameter checking
     * @param key       {String|Object} The id of the value to remove
     * @param success   {Function}      The success callback
     * @param error     {Function}      The error callback
     */
    function storage_remove(key, success, error) {
        //check that a valid key has been specified
        if (!key) {
            return error && error(new Error("Storage.Get: a key is required to retrieve a storage value."));
        }

        //check that an adapter exists
        if (!adapter) {
            return error && error(new Error("No Storage adapter available."));
        }

        //call adapter.Get
        return adapter.Remove.apply(this, arguments);
    }

    //only look for an adapter if no manual override has been set
    if (!Survana.Storage) {
        find_best_adapter(adapters, 0, function (a) {
            adapter = a;

            Survana.Storage = {
                'Name': adapter ? adapter.Name : "None",
                'IsAvailable': Boolean(adapter),
                'Get': storage_get,
                'Set': storage_set,
                'Save': storage_save,
                'Remove': storage_remove,
                'SetScope': set_scope
            };

            //attempt to auto-detect the storage scope
            if (document && document.body) {
                var scope = document.body.getAttribute('data-storage-scope');
                if (scope) {
                    console.log('setting the scope');
                    Survana.Storage.SetScope(scope);
                } else {
                    console.log('not setting the scope');
                }
            }
        });

        console.log("Survana Storage: using adapter", Survana.Storage.Name);
    }
}(window.Survana));