package survana

import (
        "crypto/ecdsa"
        "crypto/rand"
        "crypto/elliptic"
        "crypto/sha512"
        "math/big"
        "encoding/json"
        "encoding/hex"
        "time"
       )

const (
        EC_P521 = iota
      )

type PrivateKey struct {
    Id string
    Type int
    *ecdsa.PrivateKey
}

type serializableKey struct {
    Id string
    Type int
    D, X, Y *big.Int
}

func NewPrivateKey() *PrivateKey {
    return &PrivateKey{
        Type: EC_P521,
    }
}

//generate a new Elliptical private key using the P521 curve and /dev/urandom
func GeneratePrivateKey() (private_key *PrivateKey, err error) {
    ec_key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
    if err != nil {
        return nil, err
    }

    private_key = &PrivateKey{
        PrivateKey: ec_key,
    }

    private_key.Id, err = GenerateKeyId()
    if err != nil {
        return
    }

    return
}

//generates a new Id by concatenating the current time in nanoseconds with 16 random bytes
func GenerateKeyId() (id string, err error) {
    //get the number of nanoseconds since the Epoch
    t := time.Now().UnixNano();
    b := make([]byte, 1)

    //read random bytes
    _, err = rand.Read(b)
    if err != nil {
        return
    }

    //convert the timestamp to string, then to byte array, and append it to the random bytes
    b = []byte(string(t) + string(b))

    //generate a sha512 hash of the bytes
    hashed := sha512.New().Sum(b)

    //return the hex representation of the hash
    id = hex.EncodeToString(hashed)

    return
}

func (key *PrivateKey) MarshalJSON() (data []byte, err error) {
    skey := &serializableKey{
        Id: key.Id,
        Type: key.Type,
        D: key.D,
        X: key.PublicKey.X,
        Y: key.PublicKey.Y,
    }

    data, err = json.Marshal(skey)
    return
}

func (key *PrivateKey) UnmarshalJSON(data []byte) (err error) {
    skey := &serializableKey{}
    err = json.Unmarshal(data, skey)
    if err != nil {
        return
    }

    key.Id = skey.Id
    key.Type = skey.Type
    key.PrivateKey = &ecdsa.PrivateKey{}

    var curve elliptic.Curve

    switch (key.Type) {
        case EC_P521: curve = elliptic.P521()
    }

    //.PublicKey is the embedded .PrivateKey.PublicKey
    key.PublicKey = ecdsa.PublicKey{
                        Curve: curve,
                        X: skey.X,
                        Y: skey.Y,
                    }

    key.D = skey.D

    return
}
