package password

import (
	"crypto"
	_ "crypto/md5"
	"crypto/rand"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/tredoe/crypt"
	_ "github.com/tredoe/crypt/md5_crypt"
	_ "github.com/tredoe/crypt/sha256_crypt"
	_ "github.com/tredoe/crypt/sha512_crypt"
	"log"
	random "math/rand"
)

const (
	defaultDatasourceCrypt  = "sha512"
	defaultDatasourceHash   = "md5"
	defaultDatasourceLength = 32
	minimumPasswordLength   = 8
	maximumPasswordLength   = 128
	strongPasswordChars     = "!@#$%^&*()_+-?./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// InstructPasswordCrypt Interface used to provide instructions for constructing password crypt
type InstructPasswordCrypt struct {
	cryptAlgorithm crypt.Crypter
	passwordInput  []byte
	passwordLength int
}

// InstructPasswordHash Interface used to provide instructions for constructing password hash
type InstructPasswordHash struct {
	hashAlgorithm crypto.Hash
	hashSumSize   int
}

// Customized encoding to provide strong alphanumeric and special characters for generated passwords
var generatedPasswordEncoding = base64.NewEncoding(generatePasswordEncoding()).WithPadding(base64.NoPadding)

func generatePasswordEncoding() string {
	log.Println("Begin password encoding randomization...")

	strongChars := []rune(strongPasswordChars)

	log.Println("Shuffling variable of defined strong characters...")
	random.Shuffle(64, func(i, j int) {
		strongChars[i], strongChars[j] = strongChars[j], strongChars[i]
	})

	log.Println("Characters shuffled")

	encoding := string(strongChars)[:64]

	log.Println("End password encoding randomization: SUCCESS")

	return encoding
}

func CryptPassword(instructions InstructPasswordCrypt) (string, error) {
	crypter := instructions.cryptAlgorithm
	password := instructions.passwordInput

	log.Println("Begin password crypt creation...")

	log.Println("Creating cryptAlgorithm crypt handler...")
	log.Printf("Crypting %d bytes...", len(password))
	cryptPassword, err := crypter.Generate(password, []byte(""))
	if err != nil {
		log.Printf("error: %s", err)
		return "", errors.New("failed crypt password bytes")
	}

	log.Println("End password crypt creation: SUCCESS")

	return cryptPassword, err
}

func GeneratePassword(characters int) ([]byte, error) {
	log.Println("Begin datasource password generation...")

	log.Printf("Creating length %d byte arry...", characters)
	passwordBytes := make([]byte, characters)
	log.Printf("Byte array created")

	log.Println("Writing pseudo-random data...")
	_, err := rand.Read(passwordBytes)
	if err != nil {
		log.Printf("error: %d", err)
		return nil, errors.New("failed to write pseudo-random data")
	}
	log.Println("Wrote pseudo-random data")

	log.Println("Encoding data to base64...")
	encodedPasswordBytes := generatedPasswordEncoding.EncodeToString(passwordBytes)[:characters]
	log.Printf("length of password: %d", len(passwordBytes))
	log.Println("Encoded data created")

	if len(encodedPasswordBytes) != characters {
		log.Printf("invalid password size: got %d, expected %d", len(encodedPasswordBytes), characters)
		return nil, errors.New("invalid generated character length")
	}

	log.Println("End datasource password generation: SUCCESS")

	return []byte(encodedPasswordBytes), nil
}

func HashPassword(password []byte, instructions InstructPasswordHash) (string, error) {
	var passwordSumHex string

	hashHandler := instructions.hashAlgorithm
	hashSumSize := instructions.hashSumSize

	log.Println("Begin password hash creation...")

	log.Printf("Creating %s hashAlgorithm hash handler...", hashHandler)
	passwordHasher := hashHandler.New()
	log.Println("Algorithm handler created")

	log.Println("Writing password bytes to handler buffer...")
	passwordHasher.Write(password)
	log.Printf("Wrote %d bytes to handler buffer", len(password))

	log.Println("Calculating hash sum...")
	passwordHashSum := passwordHasher.Sum(nil)
	log.Println("Hash sum calculated")

	log.Println("Encoding sum to hex...")
	passwordSumHex = hex.EncodeToString(passwordHashSum)
	log.Println("Encoded hash created")

	if len(passwordSumHex) != hashSumSize {
		log.Printf("invalid hash sum size: got %d, expected %d", len(passwordSumHex), hashSumSize)
		return "", errors.New("invalid generated hash length")
	}

	log.Println("End password hash creation: SUCCESS")

	return passwordSumHex, nil
}
