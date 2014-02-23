package survana

import (
	"testing"
    _ "log"
)

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

    if key == nil {
        t.Errorf("key is nil, expected non-nil")
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
    for i := 0; i < b.N; i++ {
        _, _ = GenerateKeyId()
    }
}

func TestMarshalJSON(t *testing.T) {
    key, err := GeneratePrivateKey()

    if err != nil {
        t.Errorf("err = %v", err)
    }

    if key == nil {
        t.Errorf("key is nil, expected non-nil")
    }

    data, err := key.MarshalJSON()

    if err != nil {
        t.Errorf("err = %v", err)
    }

    if len(data) == 0 {
        t.Errorf("len(data) == %v, expected %v", len(data), 0)
    }
}

func BenchmarkMarshalJSON(b *testing.B) {
    key, err := GeneratePrivateKey()

    if err != nil {
        b.Errorf("err = %v", err)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = key.MarshalJSON()
    }
}

func TestUnmarshalJSON(t *testing.T) {
    key, err := GeneratePrivateKey()
    if err != nil {
        t.Errorf("err = %v", err)
    }

    keydata, err := key.MarshalJSON()
    if err != nil {
        t.Errorf("err = %v", err)
    }

    key2 := NewPrivateKey()
    err = key2.UnmarshalJSON(keydata)

    if key2.Id != key.Id {
        t.Errorf("key2.Id = %v, expected %v", key2.Id, key.Id)
    }

    if key2.Type != key.Type {
        t.Errorf("key2.Type = %v, expected %v", key2.Type, key.Type)
    }

    if key2.D.Cmp(key.D) != 0 {
        t.Errorf("key2.PrivateKey.D = %v, expected %v", key2.D, key.D)
    }

    if key2.X.Cmp(key.X) != 0 {
        t.Errorf("key2.PrivateKey.X = %v, expected %v", key2.X, key.X)
    }

    if key2.Y.Cmp(key.Y) != 0 {
        t.Errorf("key2.PrivateKey.Y = %v, expected %v", key2.Y, key.Y)
    }
}

func BenchmarkUnmarshalJSON(b *testing.B) {
    key, err := GeneratePrivateKey()
    if err != nil {
        b.Errorf("err = %v", err)
    }

    keydata, err := key.MarshalJSON()
    if err != nil {
        b.Errorf("err = %v", err)
    }

    key2 := NewPrivateKey()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = key2.UnmarshalJSON(keydata)
    }
}
