package main

import (
    "io/ioutil"
    "os"
	"neuroinformatics.harvard.edu/survana"
    )

func GetPrivateKey(keypath string) (private_key *survana.PrivateKey, err error) {
    var create_new_key bool = true

    if len(keypath) > 0 {
        _, err := os.Stat(keypath)

        if !os.IsNotExist(err) {
            create_new_key = false
        }
    }

    if create_new_key {
        private_key, err = survana.GeneratePrivateKey()
        if err != nil {
            return
        }

        //save the key to file before returning
        err = SavePrivateKey(private_key, keypath)
    } else {
        //load the private key from an existing file
        private_key, err = ReadPrivateKey(keypath)
    }

    return
}

func ReadPrivateKey(keypath string) (private_key *survana.PrivateKey, err error) {
    keydata, err := ioutil.ReadFile(keypath)
    if err != nil {
        return
    }

    private_key = survana.NewPrivateKey()
    err = private_key.UnmarshalJSON(keydata)

    return
}

func SavePrivateKey(key *survana.PrivateKey, keypath string) (err error) {
    keydata, err := key.MarshalJSON()
    if err != nil {
        return
    }

    err = ioutil.WriteFile(keypath, keydata, 0600)
    return
}
