/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)



type OrganizationalUnitIdentifiersConfiguration struct {
	
	Certificate string `yaml:"Certificate,omitempty"`
	
	OrganizationalUnitIdentifier string `yaml:"OrganizationalUnitIdentifier,omitempty"`
}





type NodeOUs struct {
	
	Enable bool `yaml:"Enable,omitempty"`
	
	ClientOUIdentifier *OrganizationalUnitIdentifiersConfiguration `yaml:"ClientOUIdentifier,omitempty"`
	
	PeerOUIdentifier *OrganizationalUnitIdentifiersConfiguration `yaml:"PeerOUIdentifier,omitempty"`
}



type Configuration struct {
	
	
	OrganizationalUnitIdentifiers []*OrganizationalUnitIdentifiersConfiguration `yaml:"OrganizationalUnitIdentifiers,omitempty"`
	
	
	NodeOUs *NodeOUs `yaml:"NodeOUs,omitempty"`
}

func readFile(file string) ([]byte, error) {
	fileCont, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", file)
	}

	return fileCont, nil
}

func readPemFile(file string) ([]byte, error) {
	bytes, err := readFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "reading from file %s failed", file)
	}

	b, _ := pem.Decode(bytes)
	if b == nil { 
		return nil, errors.Errorf("no pem content for file %s", file)
	}

	return bytes, nil
}

func getPemMaterialFromDir(dir string) ([][]byte, error) {
	mspLogger.Debugf("Reading directory %s", dir)

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil, err
	}

	content := make([][]byte, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read directory %s", dir)
	}

	for _, f := range files {
		fullName := filepath.Join(dir, f.Name())

		f, err := os.Stat(fullName)
		if err != nil {
			mspLogger.Warningf("Failed to stat %s: %s", fullName, err)
			continue
		}
		if f.IsDir() {
			continue
		}

		mspLogger.Debugf("Inspecting file %s", fullName)

		item, err := readPemFile(fullName)
		if err != nil {
			mspLogger.Warningf("Failed reading file %s: %s", fullName, err)
			continue
		}

		content = append(content, item)
	}

	return content, nil
}

const (
	cacerts              = "cacerts"
	admincerts           = "admincerts"
	signcerts            = "signcerts"
	keystore             = "keystore"
	intermediatecerts    = "intermediatecerts"
	crlsfolder           = "crls"
	configfilename       = "config.yaml"
	tlscacerts           = "tlscacerts"
	tlsintermediatecerts = "tlsintermediatecerts"
)

func SetupBCCSPKeystoreConfig(bccspConfig *factory.FactoryOpts, keystoreDir string) *factory.FactoryOpts {
	if bccspConfig == nil {
		bccspConfig = factory.GetDefaultOpts()
	}

	if bccspConfig.ProviderName == "SW" {
		if bccspConfig.SwOpts == nil {
			bccspConfig.SwOpts = factory.GetDefaultOpts().SwOpts
		}

		
		if bccspConfig.SwOpts.FileKeystore == nil ||
			bccspConfig.SwOpts.FileKeystore.KeyStorePath == "" {
			bccspConfig.SwOpts.Ephemeral = false
			bccspConfig.SwOpts.FileKeystore = &factory.FileKeystoreOpts{KeyStorePath: keystoreDir}
		}
	}

	return bccspConfig
}




func GetLocalMspConfigWithType(dir string, bccspConfig *factory.FactoryOpts, ID, mspType string) (*msp.MSPConfig, error) {
	switch mspType {
	case ProviderTypeToString(FABRIC):
		return GetLocalMspConfig(dir, bccspConfig, ID)
	case ProviderTypeToString(IDEMIX):
		return GetIdemixMspConfig(dir, ID)
	default:
		return nil, errors.Errorf("unknown MSP type '%s'", mspType)
	}
}

func GetLocalMspConfig(dir string, bccspConfig *factory.FactoryOpts, ID string) (*msp.MSPConfig, error) {
	signcertDir := filepath.Join(dir, signcerts)
	keystoreDir := filepath.Join(dir, keystore)
	bccspConfig = SetupBCCSPKeystoreConfig(bccspConfig, keystoreDir)

	err := factory.InitFactories(bccspConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "could not initialize BCCSP Factories")
	}

	signcert, err := getPemMaterialFromDir(signcertDir)
	if err != nil || len(signcert) == 0 {
		return nil, errors.Wrapf(err, "could not load a valid signer certificate from directory %s", signcertDir)
	}

	

	sigid := &msp.SigningIdentityInfo{PublicSigner: signcert[0], PrivateSigner: nil}

	return getMspConfig(dir, ID, sigid)
}


func GetVerifyingMspConfig(dir, ID, mspType string) (*msp.MSPConfig, error) {
	switch mspType {
	case ProviderTypeToString(FABRIC):
		return getMspConfig(dir, ID, nil)
	case ProviderTypeToString(IDEMIX):
		return GetIdemixMspConfig(dir, ID)
	default:
		return nil, errors.Errorf("unknown MSP type '%s'", mspType)
	}
}

func getMspConfig(dir string, ID string, sigid *msp.SigningIdentityInfo) (*msp.MSPConfig, error) {
	cacertDir := filepath.Join(dir, cacerts)
	admincertDir := filepath.Join(dir, admincerts)
	intermediatecertsDir := filepath.Join(dir, intermediatecerts)
	crlsDir := filepath.Join(dir, crlsfolder)
	configFile := filepath.Join(dir, configfilename)
	tlscacertDir := filepath.Join(dir, tlscacerts)
	tlsintermediatecertsDir := filepath.Join(dir, tlsintermediatecerts)

	cacerts, err := getPemMaterialFromDir(cacertDir)
	if err != nil || len(cacerts) == 0 {
		return nil, errors.WithMessage(err, fmt.Sprintf("could not load a valid ca certificate from directory %s", cacertDir))
	}

	admincert, err := getPemMaterialFromDir(admincertDir)
	if err != nil || len(admincert) == 0 {
		return nil, errors.WithMessage(err, fmt.Sprintf("could not load a valid admin certificate from directory %s", admincertDir))
	}

	intermediatecerts, err := getPemMaterialFromDir(intermediatecertsDir)
	if os.IsNotExist(err) {
		mspLogger.Debugf("Intermediate certs folder not found at [%s]. Skipping. [%s]", intermediatecertsDir, err)
	} else if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("failed loading intermediate ca certs at [%s]", intermediatecertsDir))
	}

	tlsCACerts, err := getPemMaterialFromDir(tlscacertDir)
	tlsIntermediateCerts := [][]byte{}
	if os.IsNotExist(err) {
		mspLogger.Debugf("TLS CA certs folder not found at [%s]. Skipping and ignoring TLS intermediate CA folder. [%s]", tlsintermediatecertsDir, err)
	} else if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("failed loading TLS ca certs at [%s]", tlsintermediatecertsDir))
	} else if len(tlsCACerts) != 0 {
		tlsIntermediateCerts, err = getPemMaterialFromDir(tlsintermediatecertsDir)
		if os.IsNotExist(err) {
			mspLogger.Debugf("TLS intermediate certs folder not found at [%s]. Skipping. [%s]", tlsintermediatecertsDir, err)
		} else if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("failed loading TLS intermediate ca certs at [%s]", tlsintermediatecertsDir))
		}
	} else {
		mspLogger.Debugf("TLS CA certs folder at [%s] is empty. Skipping.", tlsintermediatecertsDir)
	}

	crls, err := getPemMaterialFromDir(crlsDir)
	if os.IsNotExist(err) {
		mspLogger.Debugf("crls folder not found at [%s]. Skipping. [%s]", crlsDir, err)
	} else if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("failed loading crls at [%s]", crlsDir))
	}

	
	
	
	var ouis []*msp.FabricOUIdentifier
	var nodeOUs *msp.FabricNodeOUs
	_, err = os.Stat(configFile)
	if err == nil {
		
		
		raw, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed loading configuration file at [%s]", configFile)
		}

		configuration := Configuration{}
		err = yaml.Unmarshal(raw, &configuration)
		if err != nil {
			return nil, errors.Wrapf(err, "failed unmarshalling configuration file at [%s]", configFile)
		}

		
		if len(configuration.OrganizationalUnitIdentifiers) > 0 {
			for _, ouID := range configuration.OrganizationalUnitIdentifiers {
				f := filepath.Join(dir, ouID.Certificate)
				raw, err = readFile(f)
				if err != nil {
					return nil, errors.Wrapf(err, "failed loading OrganizationalUnit certificate at [%s]", f)
				}

				oui := &msp.FabricOUIdentifier{
					Certificate:                  raw,
					OrganizationalUnitIdentifier: ouID.OrganizationalUnitIdentifier,
				}
				ouis = append(ouis, oui)
			}
		}

		
		if configuration.NodeOUs != nil && configuration.NodeOUs.Enable {
			mspLogger.Debug("Loading NodeOUs")
			if configuration.NodeOUs.ClientOUIdentifier == nil || len(configuration.NodeOUs.ClientOUIdentifier.OrganizationalUnitIdentifier) == 0 {
				return nil, errors.New("Failed loading NodeOUs. ClientOU must be different from nil.")
			}
			if configuration.NodeOUs.PeerOUIdentifier == nil || len(configuration.NodeOUs.PeerOUIdentifier.OrganizationalUnitIdentifier) == 0 {
				return nil, errors.New("Failed loading NodeOUs. PeerOU must be different from nil.")
			}

			nodeOUs = &msp.FabricNodeOUs{
				Enable:             configuration.NodeOUs.Enable,
				ClientOuIdentifier: &msp.FabricOUIdentifier{OrganizationalUnitIdentifier: configuration.NodeOUs.ClientOUIdentifier.OrganizationalUnitIdentifier},
				PeerOuIdentifier:   &msp.FabricOUIdentifier{OrganizationalUnitIdentifier: configuration.NodeOUs.PeerOUIdentifier.OrganizationalUnitIdentifier},
			}

			

			
			f := filepath.Join(dir, configuration.NodeOUs.ClientOUIdentifier.Certificate)
			raw, err = readFile(f)
			if err != nil {
				mspLogger.Infof("Failed loading ClientOU certificate at [%s]: [%s]", f, err)
			} else {
				nodeOUs.ClientOuIdentifier.Certificate = raw
			}

			
			f = filepath.Join(dir, configuration.NodeOUs.PeerOUIdentifier.Certificate)
			raw, err = readFile(f)
			if err != nil {
				mspLogger.Debugf("Failed loading PeerOU certificate at [%s]: [%s]", f, err)
			} else {
				nodeOUs.PeerOuIdentifier.Certificate = raw
			}
		}
	} else {
		mspLogger.Debugf("MSP configuration file not found at [%s]: [%s]", configFile, err)
	}

	
	cryptoConfig := &msp.FabricCryptoConfig{
		SignatureHashFamily:            bccsp.SHA2,
		IdentityIdentifierHashFunction: bccsp.SHA256,
	}

	
	fmspconf := &msp.FabricMSPConfig{
		Admins:                        admincert,
		RootCerts:                     cacerts,
		IntermediateCerts:             intermediatecerts,
		SigningIdentity:               sigid,
		Name:                          ID,
		OrganizationalUnitIdentifiers: ouis,
		RevocationList:                crls,
		CryptoConfig:                  cryptoConfig,
		TlsRootCerts:                  tlsCACerts,
		TlsIntermediateCerts:          tlsIntermediateCerts,
		FabricNodeOus:                 nodeOUs,
	}

	fmpsjs, _ := proto.Marshal(fmspconf)

	mspconf := &msp.MSPConfig{Config: fmpsjs, Type: int32(FABRIC)}

	return mspconf, nil
}

const (
	IdemixConfigDirMsp                  = "msp"
	IdemixConfigDirUser                 = "user"
	IdemixConfigFileIssuerPublicKey     = "IssuerPublicKey"
	IdemixConfigFileRevocationPublicKey = "RevocationPublicKey"
	IdemixConfigFileSigner              = "SignerConfig"
)


func GetIdemixMspConfig(dir string, ID string) (*msp.MSPConfig, error) {
	ipkBytes, err := readFile(filepath.Join(dir, IdemixConfigDirMsp, IdemixConfigFileIssuerPublicKey))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read issuer public key file")
	}

	revocationPkBytes, err := readFile(filepath.Join(dir, IdemixConfigDirMsp, IdemixConfigFileRevocationPublicKey))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read revocation public key file")
	}

	idemixConfig := &msp.IdemixMSPConfig{
		Name:         ID,
		Ipk:          ipkBytes,
		RevocationPk: revocationPkBytes,
	}

	signerBytes, err := readFile(filepath.Join(dir, IdemixConfigDirUser, IdemixConfigFileSigner))
	if err == nil {
		signerConfig := &msp.IdemixMSPSignerConfig{}
		err = proto.Unmarshal(signerBytes, signerConfig)
		if err != nil {
			return nil, err
		}
		idemixConfig.Signer = signerConfig
	}

	confBytes, err := proto.Marshal(idemixConfig)
	if err != nil {
		return nil, err
	}

	return &msp.MSPConfig{Config: confBytes, Type: int32(IDEMIX)}, nil
}
