package survana

import (
	"testing"
    "crypto/elliptic"
    "crypto/ecdsa"
    "crypto/rand"
    _ "log"
)

type fake_rand struct {
}

func (r *fake_rand) Read(p []byte) (n int, err error) {
    np := len(p)
    i := 0

    for i = 0; i < np; i++ {
        p[i] = byte(i);
    }

    return i, nil
}

func TestNewPrivateKey(t *testing.T) {
    key := NewPrivateKey()

    //only EC_P521 is supported
    if key.Type != EC_P521 {
		t.Errorf("key.Type = %v, expected %v (EC_P521)", key.Type, EC_P521)
    }
}

func BenchmarkNewPrivateKey(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = NewPrivateKey()
    }
}

func TestGeneratePrivateKey(t *testing.T) {
    key, err := GeneratePrivateKey()

    if err != nil {
        t.Errorf("err = %v", err)
    }

    if len(key.Id) == 0 {
        t.Errorf("len(key.Id) == %v, expected non-zero", len(key.Id))
    }

    //only EC_P521 is supported
    if key.Type != EC_P521 {
		t.Errorf("key.Type = %v, expected %v (EC_P521)", key.Type, EC_P521)
    }

    //test that the private key is not nil
    if key.PrivateKey == nil {
        t.Errorf("key.PrivateKey = %v, expected non-nil", key.PrivateKey)
    }
}

func BenchmarkGeneratePrivateKey(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = GeneratePrivateKey()
    }
}

func TestGenerateKeyId(t *testing.T) {
    id, err := GenerateKeyId()

    if err != nil {
        t.Errorf("err = %v", err)
    }

    if len(id) == 0 {
        t.Errorf("len(id) == %v, expected non-zero", len(id))
    }

    id2, err := GenerateKeyId()

    if err != nil {
        t.Errorf("err2 = %v", err)
    }

    if len(id2) == 0 {
        t.Errorf("len(id2) == %v, expected non-zero", len(id2))
    }

    if id == id2 {
        t.Errorf("id1 == id2, expected them to be different\n\tid1=%v\n\tid2=%v\n", id, id2)
    }
}

func BenchmarkGenerateKeyId(b *testing.B) {
    //fake_rng := &fake_rand{}
    curve := elliptic.P521()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = ecdsa.GenerateKey(curve, rand.Reader)
    }
}
