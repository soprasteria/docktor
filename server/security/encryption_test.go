package security

import (
	"testing"
)

func TestGoodEncryption(t *testing.T) {
	data := []struct {
		input string
		key   string
	}{
		{"Foo", "Boo"},
		{"Bar", "Car"},
		{"Bar", ""},
		{"", "Car"},
		{"Long input with more than 16 characters", "Car"},
		{`-----BEGIN RSA FAKE PRIVATE KEY-----
MIIEpQIBAAKCAQEAoK3co8JWl71/ZBI2EGUwIoAdL/9VuHilmJO1KH4R3d+WP9i6
KxBe356g3mhogJUEB82APICJF8gdO9rou0UU+y+N+EQFNkTB+vrsflEe0Jl79CDK
HtXEtrZLQfCfgvRwtmr4qd5hKGsG/fKqJe/EUYg/P/8qCgAUNLe5dsZRRntKHUGJ
uXGehWpCsqd1Xwn+3Q0ScGidvC7hHL1Brt0NtR2NrrMY/KJmoTSgZW3PGUp1GazS
Mj9o1BIrW3f8kFqZOgHgf5NEwZulFY1lY36wLgs0IfoOTH5zlkabo652NX7Dwh1H
JevGrx6hwgO2TCDJKrJbfIkn3LaHQtB02mys5wIDAQABAoIBAH/83ZKH24f1BwNE
SlJW97eTiyxPXA2b9HQLvDKr8Tllvv3OecSRvrCrX4KpcgxyJoq8k5gd3pXV7gte
fSGmjmbi41hNfkVTEZ0EwNqBbWVFzOmSMM3NA2ty25GDxNuAMEPuA3Lz0zppvmEM
KaIcUamMOK2WZ/sdQvsXGaFVqSnVCeRpVX+9NMLo2hdVf9FgJJEnGM3qovZM3ieE
O8InDu3h0zapUAvlqxvpoRyr7u7yK1MZNXzzKlWrZgDdtINbR2C5I35jLrD5gWx4
3e0mfCUvi7OC3D6y+Dz57hWOZXfI27PyPwC/7x/ezX1UvzSqR3kF/Rrs9vuIb//M
pUj0nlkCgYEA0FBWiiXiFKdmJbkKHF6A8P3Y6ziWg8IQuQTa1n4ItI/nT/NjsMCs
Tg24unSk7antey8ktwNAP8F8m9cckZokvGqIo9zH8H39oX4jCJVal46KlFtgc9bH
OCVMI6pkO+Dn/r/9m3fYH7LMu2HJZHkJiXsjESruHpJpBozZLdS6NesCgYEAxXYG
KWMRMPBTCDxEGGh5cz1YnnoKj6CpPV6cqfCHr0S9j05kKXuDRTf6VhAHBTCvfk5p
j8YX7wfMADLDEgTR+wTG3DKPPaqBdNuS5X9T8AAgSrT8u48CWddQFgt2KX1Lfeb4
jBvF7PIy1OwRM9H+ry4bvQHpVTroQOoU7lxxefUCgYEAqh7B3cY8UKO43tuzryFa
afTU/pvTB70nzQFy+jIpR9QxknBxHHrs/D1mfBcgTds1TyFb+X3VLXwFGHvfH+Vj
VOAnwLJgMj2iMQ2C7NKUDithbvEE3vUq8uY6vPG9M81jiP8fzKRdwt1RJ0Ifp0bV
jAocxDtsBVmKHchO4IfWnosCgYEAn25rDTGa6MuyDv0x8g8wuHY4vkRFkLAk1ZM1
pRS7SS2UbEfcIY5DcDkBXEm2kV884xuGqfpEys+dzC8wR7UyoZ26voHoG982hVbg
ZYKIEEjZydgWE44lVMq/M/1vG5K5yF8cIWwvQ+BOYJJ2VUPhgioVZWdMsW9NpVQb
MFXdnZ0CgYEAiXIaeMLOqMnF4l0bWQukvXKob/YxGmnNRyG7b44kduhYpdwXrf5s
uSj4FjHlD2Vo5E5X6pv21VnGWUYnb7h6W+DSx+7BnE1I8Wvdx7veyzaiXY2Kj4Uq
I9bv8M6DJeW+uiUn3lgeWVHzjgRpqVHJIjU47QAgdOdrWeJI3OJ2COM=
-----END RSA FAKE PRIVATE KEY-----`, "A secret"},
	}
	for _, d := range data {
		enc, err := EncryptString(d.input, d.key)
		if err != nil {
			t.Errorf("Unable to encrypt '%v' with key '%v': %v", d.input, d.key, err)
			continue
		}

		dec, err := DecryptString(enc, d.key)
		if err != nil {
			t.Errorf("Unable to decrypt '%v' with key '%v': %v", enc, d.key, err)
			continue
		}
		if dec != d.input {
			t.Errorf("Decrypt Key %v\n  Input: %v\n  Expect: %v\n  Actual: %v", d.key, enc, d.input, enc)
		}
	}
}

func TestDecryptionWithWrongKey(t *testing.T) {

	text := "Un given text"
	firstKey := "The first key"
	secondKey := "Another key"
	enc, err := EncryptString(text, firstKey)
	if err != nil {
		t.Fatalf("Unable to encrypt '%v' with key '%v': %v", text, firstKey, err)
	}

	// Trying to decrypt with wrong key
	dec, err := DecryptString(enc, secondKey)
	if err != nil {
		t.Fatalf("Unable to decrypt '%v' with key '%v': %v", enc, firstKey, err)
	}

	if dec == text {
		t.Fatalf("We should not be able to decrypt '%v' to '%v' with key '%v' because it was encrypted by key '%v'", enc, dec, secondKey, firstKey)
	}
}
