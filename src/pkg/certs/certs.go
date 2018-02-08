package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	certutil "k8s.io/client-go/util/cert"
	"github.com/deepinsight/mlcloud/src/pkg/common"
	"github.com/deepinsight/mlcloud/src/pkg/certs/pkiutil"
)


// CreateCACertAndKeyfiles create a new self signed CA certificate and key files.
// If the CA certificate and key files already exists in the target folder, they are used only if evaluated equal; otherwise an error is returned.
func CreateCACertAndKeyfiles(cfg *common.AppConfiguration) error {

	caCert, caKey, err := NewCACertAndKey()
	if err != nil {
		return err
	}

	return writeCertificateAuthorithyFilesIfNotExist(
		cfg.CertificatesDir,
		common.CACertAndKeyBaseName,
		caCert,
		caKey,
	)
}

// CreateUserCertAndKeyFiles create a new certificate and key files for the end user.
func CreateUserCertAndKeyFiles(cfg *common.AppConfiguration) error {
    caCert, caKey, err := loadCertificateAuthorithy(cfg.CertificatesDir, common.CACertAndKeyBaseName)
    if err != nil {
        return err
    }

    apiCert, apiKey, err := NewUserCertAndKey(cfg, caCert, caKey)
    if err != nil {
        return err
    }

    baseName := fmt.Sprintf("%s-%s", cfg.GroupName, cfg.UserName)
    return writeCertificateFilesIfNotExist(
        cfg.CertificatesDir,
        baseName,
        caCert,
        apiCert,
        apiKey,
    )
}

// NewAPIServerCertAndKey generate CA certificate for user, signed by the given CA.
func NewUserCertAndKey(cfg *common.AppConfiguration, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {

    org := []string{cfg.GroupName}
    config := certutil.Config{
        CommonName: cfg.UserName,
        Organization: org,
        Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    }
    apiCert, apiKey, err := pkiutil.NewCertAndKey(caCert, caKey, config)
    if err != nil {
        return nil, nil, fmt.Errorf("failure while creating API server key and certificate: %v", err)
    }

    return apiCert, apiKey, nil
}

// NewCACertAndKey will generate a self signed CA.
func NewCACertAndKey() (*x509.Certificate, *rsa.PrivateKey, error) {

	caCert, caKey, err := pkiutil.NewCertificateAuthority()
	if err != nil {
		return nil, nil, fmt.Errorf("failure while generating CA certificate and key: %v", err)
	}

	return caCert, caKey, nil
}


// loadCertificateAuthorithy loads certificate authorithy
func loadCertificateAuthorithy(pkiDir string, baseName string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Checks if certificate authorithy exists in the PKI directory
	if !pkiutil.CertOrKeyExist(pkiDir, baseName) {
		return nil, nil, fmt.Errorf("couldn't load %s certificate authorithy from %s", baseName, pkiDir)
	}

	// Try to load certificate authorithy .crt and .key from the PKI directory
	caCert, caKey, err := pkiutil.TryLoadCertAndKeyFromDisk(pkiDir, baseName)
	if err != nil {
		return nil, nil, fmt.Errorf("failure loading %s certificate authorithy: %v", baseName, err)
	}

	// Make sure the loaded CA cert actually is a CA
	if !caCert.IsCA {
		return nil, nil, fmt.Errorf("%s certificate is not a certificate authorithy", baseName)
	}

	return caCert, caKey, nil
}

// writeCertificateAuthorithyFilesIfNotExist write a new certificate Authorithy to the given path.
// If there already is a certificate file at the given path; kubeadm tries to load it and check if the values in the
// existing and the expected certificate equals. If they do; kubeadm will just skip writing the file as it's up-to-date,
// otherwise this function returns an error.
func writeCertificateAuthorithyFilesIfNotExist(pkiDir string, baseName string, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {

	// If cert or key exists, we should try to load them
	if pkiutil.CertOrKeyExist(pkiDir, baseName) {

		// Try to load .crt and .key from the PKI directory
		caCert, _, err := pkiutil.TryLoadCertAndKeyFromDisk(pkiDir, baseName)
		if err != nil {
			return fmt.Errorf("failure loading %s certificate: %v", baseName, err)
		}

		// Check if the existing cert is a CA
		if !caCert.IsCA {
			return fmt.Errorf("certificate %s is not a CA", baseName)
		}

		// kubeadm doesn't validate the existing certificate Authorithy more than this;
		// Basically, if we find a certificate file with the same path; and it is a CA
		// kubeadm thinks those files are equal and doesn't bother writing a new file
		fmt.Printf("[certificates] Using the existing %s certificate and key.\n", baseName)
	} else {

		// Write .crt and .key files to disk
		if err := pkiutil.WriteCertAndKey(pkiDir, baseName, caCert, caKey); err != nil {
			return fmt.Errorf("failure while saving %s certificate and key: %v", baseName, err)
		}

		fmt.Printf("[certificates] Generated %s certificate and key.\n", baseName)
	}
	return nil
}

// writeCertificateFilesIfNotExist write a new certificate to the given path.
// If there already is a certificate file at the given path; kubeadm tries to load it and check if the values in the
// existing and the expected certificate equals. If they do; kubeadm will just skip writing the file as it's up-to-date,
// otherwise this function returns an error.
func writeCertificateFilesIfNotExist(pkiDir string, baseName string, signingCert *x509.Certificate, cert *x509.Certificate, key *rsa.PrivateKey) error {

	// Checks if the signed certificate exists in the PKI directory
	if pkiutil.CertOrKeyExist(pkiDir, baseName) {
		// Try to load signed certificate .crt and .key from the PKI directory
		signedCert, _, err := pkiutil.TryLoadCertAndKeyFromDisk(pkiDir, baseName)
		if err != nil {
			return fmt.Errorf("failure loading %s certificate: %v", baseName, err)
		}

		// Check if the existing cert is signed by the given CA
		if err := signedCert.CheckSignatureFrom(signingCert); err != nil {
			return fmt.Errorf("certificate %s is not signed by corresponding CA", baseName)
		}

		// kubeadm doesn't validate the existing certificate more than this;
		// Basically, if we find a certificate file with the same path; and it is signed by
		// the expected certificate authorithy, kubeadm thinks those files are equal and
		// doesn't bother writing a new file
		fmt.Printf("[certificates] Using the existing %s certificate and key.\n", baseName)
	} else {

		// Write .crt and .key files to disk
		if err := pkiutil.WriteCertAndKey(pkiDir, baseName, cert, key); err != nil {
			return fmt.Errorf("failure while saving %s certificate and key: %v", baseName, err)
		}

		fmt.Printf("[certificates] Generated %s certificate and key.\n", baseName)
		if pkiutil.HasServerAuth(cert) {
			fmt.Printf("[certificates] %s serving cert is signed for DNS names %v and IPs %v\n", baseName, cert.DNSNames, cert.IPAddresses)
		}
	}

	return nil
}

// writeKeyFilesIfNotExist write a new key to the given path.
// If there already is a key file at the given path; kubeadm tries to load it and check if the values in the
// existing and the expected key equals. If they do; kubeadm will just skip writing the file as it's up-to-date,
// otherwise this function returns an error.
func writeKeyFilesIfNotExist(pkiDir string, baseName string, key *rsa.PrivateKey) error {

	// Checks if the key exists in the PKI directory
	if pkiutil.CertOrKeyExist(pkiDir, baseName) {

		// Try to load .key from the PKI directory
		_, err := pkiutil.TryLoadKeyFromDisk(pkiDir, baseName)
		if err != nil {
			return fmt.Errorf("%s key existed but it could not be loaded properly: %v", baseName, err)
		}

		// kubeadm doesn't validate the existing certificate key more than this;
		// Basically, if we find a key file with the same path kubeadm thinks those files
		// are equal and doesn't bother writing a new file
		fmt.Printf("[certificates] Using the existing %s key.\n", baseName)
	} else {

		// Write .key and .pub files to disk
		if err := pkiutil.WriteKey(pkiDir, baseName, key); err != nil {
			return fmt.Errorf("failure while saving %s key: %v", baseName, err)
		}

		if err := pkiutil.WritePublicKey(pkiDir, baseName, &key.PublicKey); err != nil {
			return fmt.Errorf("failure while saving %s public key: %v", baseName, err)
		}
		fmt.Printf("[certificates] Generated %s key and public key.\n", baseName)
	}

	return nil
}

