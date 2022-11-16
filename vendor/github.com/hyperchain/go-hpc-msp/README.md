MSP
===

> Member Service Provider interface in go.

## Table of Contents

- [Usage](#usage)
- [API](#api)
- [Mockgen](#mockgen)
- [GitCZ](#gitcz)
- [Contribute](#contribute)
- [License](#license)

## usage
### IdentityManager
```
    im := msp.NewIdentityManagerNew(config.MSPConfigMock, logger.MSPLoggerSingleCase,db)
    es, err := im.Sign("ecert", []byte("hello"))
    peerEIdentity, err := im.NewIdentity([]byte(peerEcert), []byte(peerEcertPriv))
    err = im.CheckAndVerify(peerEIdentity, []byte("hello"), peerSign)
    rca, err = im.NewCA(cert, nil)
    sk11, err := im.GetIdentities("ecert")[0].GenerateSessionKeyWith(peerEIdentity, randA, randB)
```
### Identity
```
    peerSign, err := peerEIdentity.Sign([]byte("hello"))
    sk21, err := peerEIdentity.GenerateSessionKeyWith(im.GetIdentities("ecert")[0], randB, randA)
    err = peerEIdentity.VerifySignature(rs, []byte("hello"))
```
### CA
```
    rca.RegisterAuthFunc(f)
    rca.RegisterExpandedVerifyFunc(f)
    _, err = rca.GenIdentity(&idr, false)
    ok, err := lca.CheckIdentity(nil, true)
```
## api
### IdentityManager
```
//IdentityManager provides a major functional interface that represents the current identity and admission control policies.
type IdentityManager interface {
   //get configuration item
   GetConfig() config.EncryptionConfigInterface
   //get log
   GetLogger() logger.Logger
   //Get the key Identity in the current IdentityManager
   GetIdentities(name ...string) []Identity
   //Trusted CA list
   GetCAs(name ...string) []CA 
   //For example, adding new nodes to a distributed CA
   NewIdentity(identity, private []byte) (Identity, error)
   //For example, a new node is added, using a different CA
   NewCA(ca, private []byte) (CA, error)
   RevokeIdentity(name string) error
   DeleteCA(name string) error
   //Verify certificates and signatures, with caching
   CheckAndVerify(cert Identity, msg, sign []byte) error  
   //signature   
   Sign(alias string, msg []byte) ([]byte, error)               
   //To generate ID
   GenerateIdentity(caName string, info []byte) ([]byte, error)
   Close() error
}
```
### Identity
```
//Identity represents the identity, such as the identity of the node, the identity of the SDK, etc., is the object being verified
type Identity interface {
   //get cert name
   GetName() string
   //get ca name
   GetCAName() string
   //Get the DER encoded content of the certificate
   GetIdentityBytes() []byte
   //If there is no private key, an error is returned
   Sign(msg []byte) ([]byte, error)
   //If there is no private key, an error is returned
   GenerateSessionKeyWith(identity Identity, randA, randB []byte) (SessionKey, error)
   VerifySignature(sign, msg []byte) error
   //revoke
   Revoke() bool
}
```
### CA
```
//CA stands for certificate provider
type CA interface {
   GetName() string
   //Register the authentication function callback, 
   //which will be called back when applying for a certificate from the CA.
   RegisterAuthFunc(auth Authentication)
   //Register an additional validation function callback that will be called 
   //when the certificate is validated
   RegisterExpandedVerifyFunc(verifier ExpandedVerify)
   IsRevoked(identity Identity) (bool, error)
   //generate Identity
   GenIdentity(info *IdentityRequest, useAuth bool) ([]byte, error)
   CheckIdentity(identity []byte, useExpandedVerify bool) (bool, error)
   //Revoke revoke identity
   Revoke(identity Identity) error
   //Delete delete self from trust ca list
   Delete() error
}
```
## Mockgen

Install **mockgen** : `go get github.com/golang/mock/mockgen`

How to use?

- source： Specify interface file
- destination: Generated file name
- package:The package name of the generated file
- imports: Dependent package that requires import
- aux_files: Attach a file when there is more than one file in the interface file
- build_flags: Parameters passed to the build tool

Eg.`mockgen -destination mock/mock_msp.go -package mock_msp -source msp.go`

## GitCZ

**Note**: Please use command `npm install` if you are the first time to use `git cz` in this repo.

## Contribute

PRs are welcome!

Small note: If editing the Readme, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

LGPL © Ultramesh