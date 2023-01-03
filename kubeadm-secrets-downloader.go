// You can edit this code!
// Click here and start typing.
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "flag"
    "os/exec"
    "os"
    "github.com/pkg/errors"
)

// function from https://github.com/kubernetes/kubernetes/blob/5f0e7932f73d958be62cd1b034e9e2a2d9525f08/cmd/kubeadm/app/phases/copycerts/copycerts.go#L276
// DecryptBytes takes a byte slice of encrypted data and an encryption key and returns a decrypted byte slice of data.
// The key must be an AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256
func DecryptBytes(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, errors.New("size of data is less than the nonce")
    }

    nonce, out := data[:nonceSize], data[nonceSize:]
    out, err = gcm.Open(nil, nonce, out, nil)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func DecryptString(dataBase64 string, keyEncoded string) (string, error) {
    // secrets in kubernetes are base64 encoded:
    data, _ := base64.StdEncoding.DecodeString(dataBase64)

    // certificate-key in kubeadm is encoded in hexadecimal:
    key, _ := hex.DecodeString(keyEncoded)

    result, err := DecryptBytes(data, key)
    return string(result[:]), err
}

func GetCertsFromKubernetesSecret() (map[string]any) {

    // parse json:
    cmd := exec.Command("kubectl", "get", "secrets", *secret, "-n", *secretnamespace, "-o", "json")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("%v", err)
        os.Exit(1)
    }
    var resultExec map[string]any
    if err := json.Unmarshal([]byte(output), &resultExec); err != nil {
         fmt.Printf("\n%v %v\n", string(output), err)
         os.Exit(1)
    }
    return resultExec["data"].(map[string]any)
}

func Test() {
    result, _ := DecryptString("ZHynzOio2A+jk2sSikStHH3SDEZ/zgUxfZr8KXOhSPeMUEp3Io2x9YjtftltDLEwaQhukVrkl7GlZExf+PCSXWHzLjWEKTxloYbLLXKnHX2OjS4yyiTSSL5AM2Xtsc83S5EtRZJIVavtHTCg/uuoVluj0043cMcK+mxnQ2znEd35bBNse2j/bB8UoAlTTfMZv471QqxDQEM7nie3s0epkWblN+oW/nAO3GwBPye4FXxUagcT5hVfybh4sttoA781iC1B1xRiwGuIeNw7/vPl523yGf0A2tdqQ3OfhsHEnNHIfKe33Qvw9NF7dohznPVkxVbZXpNBKMRR7bZkpvo9gu4j1/cEI8t23AhaiOGVaQ2Yzm9PGIx7pSO9kh87fszpGKvCuK28Ddnnvq3PAKdFOh9TJWvSPi8LQmcszCeEDs5osMogLyD3m4CVvN5wnwk8ZMh5MxJQn0SQz0xhfE5q7kZ3jzpyU1YmJMzcIlDBOq8FLSYWP3v3q98nIRh7Vb0MTZw20xbkf467A7gpLxNzETrW7pj6InwJfHF8qypCpnO3EJ+CKF6fD1zzCG3Iduq0pFta6ZOv49dEtyUzcNFG/2VYsigEhJTTcoNJPPpV5mrJwWPOIkHGKjgixMBMeb3oyWVNah9AfArD4oIgnORQ4mCiOXBYZAZvI2jyyeIxZSju9zdb0MCJqEmMFkMMzvwixE1Wi6404MKKoxt8rcWorpd6wkTWMs3bPnw2sTQ3cazOEXn4lewli6Md7qwwGMFWevkvZptbXNh8fLJXN4zLzGQmJuuHeObfkEjK3jhD3VjFq7Ekm5zPZENonbC6Gdy4WPOk+jxiFKuZ0py95KNxmD0x0yUg+yMsz6a2CUw1ZgtaCv9Q+oi4VT2puODf0yDhU3jdZJg/zvmYePNyzJa1bnVSR1M3fAc9Auia3i9tsxvU3/w7Z9zfMILWdEE1eAICMhwumKPN2cLhmXhK6riMymFRPHXj3Dk3Gm3GJmXq7Ut8ZYZoBVH363gzsfU1YHL/tg5ERLHiTYRuC6NIKQgMLzR6jpYyGADnRQADyqNtQQv0NN7vIpph7+VahdOuVM6XsqNT578a3wWja+Kz/4xjKh/VWAPG0ysDSJY9i8g+8la+XMiEZVCHO9QtYdJdYe9bQxtpjHfl7TaNeMrmYqU15CVERsAHE+saCfkzk1zu+zffDshJZvuHX33LSk7BfvBBwnGIDglC/9L4In08cLc/m2OZKLhdmdBHJXHJcq4eOyf74cPrXbNl7Gd3SimxoXgTswPaNkARFAmURellGgfT4oZnzfB9Z9ee+m0TGeBS2fWsGKTWMPQfANKgACSmrH+nm7R3T6mm2PnY4g1RJrd+CXFsOZ5XXqZ66ukuVviHQKKvKXMK94AUgD16Q9Px4ACVrDcRbbArU58RbIqD9zisMvJptm5pqVTnxxLbgDCP2oWpjH3dq6q1R8Sm0SyGY6QwXVUh4Sk2F21jYBSwqL9yCfKJgO4js3w=",
        "a6c1aa299c575f057288aafdb5eff22acefbbbe359a1adee749b8d0f6d22d7d1")
    fmt.Println(result)
}

var (
    key  *string
    secret *string
    secretnamespace *string
)

func init() {
   key = flag.String("key", "", "kubeadm --certificate-key generated during upload-certs phase")
   secret = flag.String("secret", "kubeadm-certs", "kubernetes secrets")
   secretnamespace = flag.String("secretnamespace", "kube-system", "kubernetes secrets namespace")
}


func main() {
    flag.Parse()
	
    certs := GetCertsFromKubernetesSecret()

    // decrypt each cert in the secret:    
    for cert, certEncoded := range certs {
      fmt.Println(cert, ":")
      certDecoded, _ := DecryptString(certEncoded.(string), *key)
      fmt.Println(certDecoded)
      fmt.Println("")
    }
}