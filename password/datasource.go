//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package password

import (
	"crypto"
	"encoding/base64"
	"errors"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/tredoe/crypt/md5_crypt"
	"github.com/tredoe/crypt/sha256_crypt"
	"github.com/tredoe/crypt/sha512_crypt"
	"github.com/zclconf/go-cty/cty"
	"log"
	"strings"
)

type Config struct {
	// Crypt algorithm to crypt provided/generated password value //
	// The corresponding DatasourceOutput.Crypt object is "crypted" with a salt using the provided algorithm //
	// Values: md5, sha256, sha512 //
	// Required: False
	// Default: sha512 //
	Crypt string `mapstructure:"crypt"`
	// Hash algorithm to hash provided/generated password value //
	// The corresponding DatasourceOutput.Hash object strictly hashed and does not provide Crypt //
	// Values: md5, sha256, sha512 //
	// Required: False
	// Default: md5 //
	Hash string `mapstructure:"hash"`
	// Input provided password will be used for crypting - none will be generated //
	// Values: string(Any) //
	// Required: False
	// Default: None //
	Input string `mapstructure:"input"`
	// Length of the generated password value - ignored when Input is provided //
	// The corresponding DatasourceOutput.PlainText object will be of the specified length
	// Values: 8 - 128 //
	// Default: 32 //
	Length int `mapstructure:"length"`

	ctx interpolate.Context
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Base64    string `mapstructure:"base64"`
	Crypt     string `mapstructure:"crypt"`
	Hash      string `mapstructure:"hash"`
	PlainText string `mapstructure:"plaintext"`
}

func (d *Datasource) Configure(raws ...interface{}) error {
	if err := config.Decode(&d.config, nil, raws...); err != nil {
		return err
	}

	log.Printf("datasource crypt provided: %v", d.config.Crypt != "")

	if d.config.Crypt == "" {
		log.Printf("overriding datasource crypt: ==> %s", defaultDatasourceCrypt)
		d.config.Crypt = defaultDatasourceCrypt
	}

	log.Printf("datasource hash provided: %v", d.config.Hash != "")

	if d.config.Hash == "" {
		log.Printf("overriding datasource hash: ==> %s", defaultDatasourceHash)
		d.config.Hash = defaultDatasourceHash
	}

	log.Printf("datasource length provided: %v", d.config.Length != 0)

	log.Printf("datasource input provided: %v", d.config.Input != "")

	if d.config.Input != "" {
		d.config.Length = len(d.config.Input)
	} else if d.config.Length == 0 {
		log.Printf("overriding datasource length: ==> %d", defaultDatasourceLength)
		d.config.Length = defaultDatasourceLength
	}

	return nil
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	null := cty.NullVal(cty.EmptyObject)

	var cryptInstructions InstructPasswordCrypt
	var hashInstructions InstructPasswordHash

	var passwordHash string
	var passwordCrypt string
	var passwordBase64 string
	var passwordPlainText string

	// add passwordLength to InstructPasswordCrypt struct
	switch length := d.config.Length; {

	case length < minimumPasswordLength:
		return null, errors.New("password length below minimum threshold")
	case length > maximumPasswordLength:
		return null, errors.New("password length above maximum threshold")
	default:
		log.Printf("detected datasource password length: %d", length)
		cryptInstructions.passwordLength = length
	}

	// add cryptAlgorithm to InstructPasswordCrypt struct
	switch suppliedCryptAlgo := d.config.Crypt; {

	case suppliedCryptAlgo == "md5":
		{
			cryptInstructions.cryptAlgorithm = md5_crypt.New()
		}
	case suppliedCryptAlgo == "sha256":
		{
			cryptInstructions.cryptAlgorithm = sha256_crypt.New()
		}
	case suppliedCryptAlgo == "sha512":
		{
			cryptInstructions.cryptAlgorithm = sha512_crypt.New()
		}
	default:
		return null, errors.New("crypt algorithm is not supported")
	}

	// add hashAlgorithm, hashSumSize to InstructPasswordHash struct
	switch suppliedHashAlgo := d.config.Hash; {

	case suppliedHashAlgo == "md5":
		{
			hashInstructions.hashAlgorithm = crypto.MD5
			hashInstructions.hashSumSize = 32
		}
	case suppliedHashAlgo == "sha256":
		{
			hashInstructions.hashAlgorithm = crypto.SHA256
			hashInstructions.hashSumSize = 64
		}
	case suppliedHashAlgo == "sha512":
		{
			hashInstructions.hashAlgorithm = crypto.SHA512
			hashInstructions.hashSumSize = 128
		}
	default:
		return null, errors.New("hash algorithm is not supported")
	}

	log.Printf("datasource password will be generated: %v", d.config.Input == "")

	// add passwordInput to InstructPasswordCrypt struct
	switch passwordInput := d.config.Input; {

	case passwordInput != "":
		{
			cryptInstructions.passwordInput = []byte(passwordInput)
		}
	default:
		generatedPassword, err := GeneratePassword(cryptInstructions.passwordLength)
		if err != nil {
			log.Println("End datasource password generation: FAILED")
			return null, err
		}

		cryptInstructions.passwordInput = generatedPassword
	}

	// base64 encode password for DatasourceOutput.Base64 output
	passwordBase64 = base64.RawURLEncoding.EncodeToString(cryptInstructions.passwordInput)

	// convert password to string for DatasourceOutput.PlainText output
	passwordPlainText = string(cryptInstructions.passwordInput)

	log.Printf("detected datasource hash algorithm: %s", strings.ToUpper(d.config.Hash))
	// create password hash for DatasourceOutput.Hash output
	passwordHash, err := HashPassword(cryptInstructions.passwordInput, hashInstructions)
	if err != nil {
		log.Println("End password hash creation: FAILED")
		return null, err
	}

	log.Printf("detected datasource crypt algorithm: %s", strings.ToUpper(d.config.Crypt))
	// create password crypt for DatasourceOutput.Crypt output
	passwordCrypt, err = CryptPassword(cryptInstructions)
	if err != nil {
		log.Println("End password crypt creation: FAILED")
		return null, err
	}

	output := DatasourceOutput{
		Base64:    passwordBase64,
		Crypt:     passwordCrypt,
		Hash:      passwordHash,
		PlainText: passwordPlainText,
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
