/* survana-storage.js

Survana.Storage is an abstraction of persistent storage available in browsers. It can make use of localStorage,
IndexedDB or FileWriter, and will automatically select the best adapter. To override the automatic selection, set
Survana.Storage = {} before loading this file. The API is fully asynchronous (even for localStorage), and all methods
take a success and error callback.

    Survana.Storage.Name                            - (string)  Name of the storage adapter
    Survana.Storage.IsAvailable                     - (boolean) Whether any storage is available
    Survana.Storage.Get(key, success, error)        - Returns an object or value identified by 'key'
    Survana.Storage.Set(key, value, success, error) - Saves an object 'value' identified by 'key'
    Survana.Storage.Remove(key, success, error)     - Removes the object identified by 'key'
    Survana.Storage.SetPrefix(prefix)               - Sets a prefix for all 'key' parameters of Get/Set/Remove

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

    //no prefix by default
    var prefix = '';

    /** Sets a prefix for all storage keys
     * @param p
     */
    function set_prefix(p) {
        if (p) {
            prefix = p + "-";
        }
    }

    /** Checks to see if window.localStorage is available (and works)
     * @param success {Function} A success callback that receives a boolean result: true/false
     * @param error {Function} A callback that handles any errors during testing
     */
    function is_local_storage_available(success, error) {

        if (window === undefined || window.localStorage === undefined || !window.localStorage) {
            //no, localStorage is not available
            return success && success(false);
        }

        var key     = 'survana-test',
            value   = 'test',
            value2;

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

    /** Returns an item by id from localStorage
     * @param key       The id of the object to retrieve from localStorage
     * @param success   The callback function that receives the result on success
     * @param error     The error callback.
     */
    function local_storage_get(key, success, error) {
        var result;

        try {
            result = JSON.parse(localStorage[prefix + key]);
        } catch (e) {
            return error && error(e);
        }

        success && success(result);
    }

    /** Saves a value in localStorage under the given 'key' id.
     * @param key       The id of the object to save
     * @param value     The object (any JSON value) to save
     * @param success   The success function. It receives boolean true on success, false on failure.
     * @param error     The error callback
     */
    function local_storage_set(key, value, success, error) {
        try {
            localStorage[prefix + key] = JSON.stringify(value);
        } catch (e) {
            return error && error(e);
        }

        success && success(true);
    }

    /** Removes an object from localStorage based on the key id
     * @param key       The id of the object to remove
     * @param success   The success callback
     * @param error     The error callback
     */
    function local_storage_remove(key, success, error) {
        try {
            delete localStorage[prefix + key];
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
     * @param key       The id of the value
     * @param success   The success function; receives the value as its parameter
     * @param error     The error callback
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
     * @param key       An identifier associated with the value
     * @param value     The value to save
     * @param success   The success callback
     * @param error     The error callback
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

    /** Calls the adapter's Remove method after performing parameter checking
     * @param key       The id of the value to remove
     * @param success   The success callback
     * @param error     The error callback
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
                'Remove': storage_remove,
                'SetPrefix': set_prefix
            };
        });

        console.log("Survana Storage: using adapter", Survana.Storage.Name);
    }
}(window.Survana));