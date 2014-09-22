/* survana-storage.js

 Survana.Storage is an abstraction of persistent storage available in browsers. It can make use of localStorage,
 IndexedDB or FileWriter, and will automatically select the best adapter. To override the automatic selection, set
 Survana.Storage = {} before loading this file (you might need to call Storage.SetScope() manually).
 The API is fully asynchronous (even for localStorage), and all relevant methods take a success and error callback.

 Survana.Storage.Name                            - (string)  Name of the storage adapter
 Survana.Storage.IsAvailable                     - (boolean) Whether any storage is available
 Survana.Storage.Get(key, success, error)        - Returns an object or value identified by 'key'.
 Returns multiple values if 'key' is an object with keys/defaults
 Survana.Storage.All(filter, success, error)     - Returns an object with all keys matching 'filter'
 Survana.Storage.Set(key, value, success, error) - Saves an object 'value' identified by 'key'.
 Can save multiple values if 'key' is an object with keys/values.
 Survana.Storage.Save(object, success, error)    - Saves multiple keys and values stored in 'object'.
 This is a shortcut for Set(object, null, success, error)
 Survana.Storage.Remove(key, success, error)     - Removes the object identified by 'key'.
 To remove multiple objects, pass an object with keys/values.
 Survana.Storage.SetScope(scope)                 - Sets a scope for all 'key' parameters of Get/Set/Remove

 Note: currently, only the LocalStorage adapter is implemented.
 Dependencies: SJCL

 @author Victor Petrov <victor_petrov@harvard.edu>
 @license BSD
 @date 09/22/2014
 */

"use strict";

if (!window.Survana) {
    window.Survana = {};
}

if (!window.sjcl) {
    throw new Error('Survana.Crypto requires SJCL - The Standford Javascript Crypto Library');
}

(function(Survana, sjcl) {
    Survana.Crypto = Survana.Crypto || {};

    Survana.Crypto.Encrypt = function(sensitive, public_data) {
        //generate a random key derived from a random password
        var randomBits = sjcl.random.randomWords(4),
            randomSalt = sjcl.random.randomWords(2),
            //derive a cryptographically secure key from the random password bits and salt
            randomKey = sjcl.misc.pbkdf2(randomBits, randomSalt, 1000, 16 * 8, sjcl.sha256);

        //prepare encryption objects
        var prp = new sjcl.cipher["aes"](randomKey),
            plaintext = sjcl.codec.utf8String.toBits(JSON.stringify(sensitive)),
            iv = sjcl.random.randomWords(3,0),
            auth_data = [],
            tag_length = 128;

        if (public_data) {
            auth_data = sjcl.codec.utf8String.toBits(JSON.stringify(public_data));
        }

        //perform AES encryption with GCM block cipher mode
        var ciphertext = sjcl.mode.gcm.encrypt(prp, plaintext, iv, auth_data, tag_length);

        return {
            "metadata": public_data,
            "password": {
                "data": sjcl.codec.base64.fromBits(randomKey)
            },
            "payload": {
                "cipher":{
                    "name": "aes",
                    "bits": 128
                },
                "blockmode": "gcm",
                "iv": sjcl.codec.base64.fromBits(iv),
                "tag_length": tag_length,
                "ciphertext": sjcl.codec.base64.fromBits(ciphertext)
            }
        }
    }
}(window.Survana, window.sjcl));