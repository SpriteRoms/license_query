package main

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type WebLicense struct {
	LicenseId  string
	DeviceId   string
	ExpireTime string
	Products   []string
	Timestamp  string
}

func MarshalLicense(licenses WebLicense) ([]byte, error) { //marshal using msgpack
	data, err := msgpack.Marshal(licenses)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal licenseInfos: %v", err)
	}

	//encrypt using aes
	data, err = AESEncrypt(data, []byte(LicDataAesKey), []byte(LicDataAesVec))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt licenseInfos: %v", err)
	}

	return data, nil
}

func UnmarshalLicense(data []byte) (*WebLicense, error) {
	//decrypt using aes
	data, err := AESDecrypt(data, []byte(LicDataAesKey), []byte(LicDataAesVec))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt licenseInfos: %v", err)
	}

	var licenses WebLicense
	err = msgpack.Unmarshal(data, &licenses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal licenseInfos: %v", err)
	}

	return &licenses, nil
}
