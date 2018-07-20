


package msp 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 




type MSPConfig struct {
	
	
	Type int32 `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	
	Config               []byte   `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MSPConfig) Reset()         { *m = MSPConfig{} }
func (m *MSPConfig) String() string { return proto.CompactTextString(m) }
func (*MSPConfig) ProtoMessage()    {}
func (*MSPConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{0}
}
func (m *MSPConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MSPConfig.Unmarshal(m, b)
}
func (m *MSPConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MSPConfig.Marshal(b, m, deterministic)
}
func (dst *MSPConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MSPConfig.Merge(dst, src)
}
func (m *MSPConfig) XXX_Size() int {
	return xxx_messageInfo_MSPConfig.Size(m)
}
func (m *MSPConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_MSPConfig.DiscardUnknown(m)
}

var xxx_messageInfo_MSPConfig proto.InternalMessageInfo

func (m *MSPConfig) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *MSPConfig) GetConfig() []byte {
	if m != nil {
		return m.Config
	}
	return nil
}









type FabricMSPConfig struct {
	
	
	
	
	
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	
	
	
	RootCerts [][]byte `protobuf:"bytes,2,rep,name=root_certs,json=rootCerts,proto3" json:"root_certs,omitempty"`
	
	
	
	
	
	
	
	
	IntermediateCerts [][]byte `protobuf:"bytes,3,rep,name=intermediate_certs,json=intermediateCerts,proto3" json:"intermediate_certs,omitempty"`
	
	Admins [][]byte `protobuf:"bytes,4,rep,name=admins,proto3" json:"admins,omitempty"`
	
	RevocationList [][]byte `protobuf:"bytes,5,rep,name=revocation_list,json=revocationList,proto3" json:"revocation_list,omitempty"`
	
	
	
	SigningIdentity *SigningIdentityInfo `protobuf:"bytes,6,opt,name=signing_identity,json=signingIdentity" json:"signing_identity,omitempty"`
	
	
	
	OrganizationalUnitIdentifiers []*FabricOUIdentifier `protobuf:"bytes,7,rep,name=organizational_unit_identifiers,json=organizationalUnitIdentifiers" json:"organizational_unit_identifiers,omitempty"`
	
	
	CryptoConfig *FabricCryptoConfig `protobuf:"bytes,8,opt,name=crypto_config,json=cryptoConfig" json:"crypto_config,omitempty"`
	
	
	TlsRootCerts [][]byte `protobuf:"bytes,9,rep,name=tls_root_certs,json=tlsRootCerts,proto3" json:"tls_root_certs,omitempty"`
	
	
	TlsIntermediateCerts [][]byte `protobuf:"bytes,10,rep,name=tls_intermediate_certs,json=tlsIntermediateCerts,proto3" json:"tls_intermediate_certs,omitempty"`
	
	
	FabricNodeOus        *FabricNodeOUs `protobuf:"bytes,11,opt,name=blockchain_node_ous,json=blockchainNodeOus" json:"blockchain_node_ous,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *FabricMSPConfig) Reset()         { *m = FabricMSPConfig{} }
func (m *FabricMSPConfig) String() string { return proto.CompactTextString(m) }
func (*FabricMSPConfig) ProtoMessage()    {}
func (*FabricMSPConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{1}
}
func (m *FabricMSPConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FabricMSPConfig.Unmarshal(m, b)
}
func (m *FabricMSPConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FabricMSPConfig.Marshal(b, m, deterministic)
}
func (dst *FabricMSPConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FabricMSPConfig.Merge(dst, src)
}
func (m *FabricMSPConfig) XXX_Size() int {
	return xxx_messageInfo_FabricMSPConfig.Size(m)
}
func (m *FabricMSPConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_FabricMSPConfig.DiscardUnknown(m)
}

var xxx_messageInfo_FabricMSPConfig proto.InternalMessageInfo

func (m *FabricMSPConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FabricMSPConfig) GetRootCerts() [][]byte {
	if m != nil {
		return m.RootCerts
	}
	return nil
}

func (m *FabricMSPConfig) GetIntermediateCerts() [][]byte {
	if m != nil {
		return m.IntermediateCerts
	}
	return nil
}

func (m *FabricMSPConfig) GetAdmins() [][]byte {
	if m != nil {
		return m.Admins
	}
	return nil
}

func (m *FabricMSPConfig) GetRevocationList() [][]byte {
	if m != nil {
		return m.RevocationList
	}
	return nil
}

func (m *FabricMSPConfig) GetSigningIdentity() *SigningIdentityInfo {
	if m != nil {
		return m.SigningIdentity
	}
	return nil
}

func (m *FabricMSPConfig) GetOrganizationalUnitIdentifiers() []*FabricOUIdentifier {
	if m != nil {
		return m.OrganizationalUnitIdentifiers
	}
	return nil
}

func (m *FabricMSPConfig) GetCryptoConfig() *FabricCryptoConfig {
	if m != nil {
		return m.CryptoConfig
	}
	return nil
}

func (m *FabricMSPConfig) GetTlsRootCerts() [][]byte {
	if m != nil {
		return m.TlsRootCerts
	}
	return nil
}

func (m *FabricMSPConfig) GetTlsIntermediateCerts() [][]byte {
	if m != nil {
		return m.TlsIntermediateCerts
	}
	return nil
}

func (m *FabricMSPConfig) GetFabricNodeOus() *FabricNodeOUs {
	if m != nil {
		return m.FabricNodeOus
	}
	return nil
}




type FabricCryptoConfig struct {
	
	
	
	SignatureHashFamily string `protobuf:"bytes,1,opt,name=signature_hash_family,json=signatureHashFamily" json:"signature_hash_family,omitempty"`
	
	
	
	IdentityIdentifierHashFunction string   `protobuf:"bytes,2,opt,name=identity_identifier_hash_function,json=identityIdentifierHashFunction" json:"identity_identifier_hash_function,omitempty"`
	XXX_NoUnkeyedLiteral           struct{} `json:"-"`
	XXX_unrecognized               []byte   `json:"-"`
	XXX_sizecache                  int32    `json:"-"`
}

func (m *FabricCryptoConfig) Reset()         { *m = FabricCryptoConfig{} }
func (m *FabricCryptoConfig) String() string { return proto.CompactTextString(m) }
func (*FabricCryptoConfig) ProtoMessage()    {}
func (*FabricCryptoConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{2}
}
func (m *FabricCryptoConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FabricCryptoConfig.Unmarshal(m, b)
}
func (m *FabricCryptoConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FabricCryptoConfig.Marshal(b, m, deterministic)
}
func (dst *FabricCryptoConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FabricCryptoConfig.Merge(dst, src)
}
func (m *FabricCryptoConfig) XXX_Size() int {
	return xxx_messageInfo_FabricCryptoConfig.Size(m)
}
func (m *FabricCryptoConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_FabricCryptoConfig.DiscardUnknown(m)
}

var xxx_messageInfo_FabricCryptoConfig proto.InternalMessageInfo

func (m *FabricCryptoConfig) GetSignatureHashFamily() string {
	if m != nil {
		return m.SignatureHashFamily
	}
	return ""
}

func (m *FabricCryptoConfig) GetIdentityIdentifierHashFunction() string {
	if m != nil {
		return m.IdentityIdentifierHashFunction
	}
	return ""
}



type IdemixMSPConfig struct {
	
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	
	Ipk []byte `protobuf:"bytes,2,opt,name=ipk,proto3" json:"ipk,omitempty"`
	
	Signer *IdemixMSPSignerConfig `protobuf:"bytes,3,opt,name=signer" json:"signer,omitempty"`
	
	RevocationPk []byte `protobuf:"bytes,4,opt,name=revocation_pk,json=revocationPk,proto3" json:"revocation_pk,omitempty"`
	
	Epoch                int64    `protobuf:"varint,5,opt,name=epoch" json:"epoch,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdemixMSPConfig) Reset()         { *m = IdemixMSPConfig{} }
func (m *IdemixMSPConfig) String() string { return proto.CompactTextString(m) }
func (*IdemixMSPConfig) ProtoMessage()    {}
func (*IdemixMSPConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{3}
}
func (m *IdemixMSPConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdemixMSPConfig.Unmarshal(m, b)
}
func (m *IdemixMSPConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdemixMSPConfig.Marshal(b, m, deterministic)
}
func (dst *IdemixMSPConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdemixMSPConfig.Merge(dst, src)
}
func (m *IdemixMSPConfig) XXX_Size() int {
	return xxx_messageInfo_IdemixMSPConfig.Size(m)
}
func (m *IdemixMSPConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_IdemixMSPConfig.DiscardUnknown(m)
}

var xxx_messageInfo_IdemixMSPConfig proto.InternalMessageInfo

func (m *IdemixMSPConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *IdemixMSPConfig) GetIpk() []byte {
	if m != nil {
		return m.Ipk
	}
	return nil
}

func (m *IdemixMSPConfig) GetSigner() *IdemixMSPSignerConfig {
	if m != nil {
		return m.Signer
	}
	return nil
}

func (m *IdemixMSPConfig) GetRevocationPk() []byte {
	if m != nil {
		return m.RevocationPk
	}
	return nil
}

func (m *IdemixMSPConfig) GetEpoch() int64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}


type IdemixMSPSignerConfig struct {
	
	Cred []byte `protobuf:"bytes,1,opt,name=cred,proto3" json:"cred,omitempty"`
	
	Sk []byte `protobuf:"bytes,2,opt,name=sk,proto3" json:"sk,omitempty"`
	
	OrganizationalUnitIdentifier string `protobuf:"bytes,3,opt,name=organizational_unit_identifier,json=organizationalUnitIdentifier" json:"organizational_unit_identifier,omitempty"`
	
	IsAdmin bool `protobuf:"varint,4,opt,name=is_admin,json=isAdmin" json:"is_admin,omitempty"`
	
	EnrollmentId string `protobuf:"bytes,5,opt,name=enrollment_id,json=enrollmentId" json:"enrollment_id,omitempty"`
	
	CredentialRevocationInformation []byte   `protobuf:"bytes,6,opt,name=credential_revocation_information,json=credentialRevocationInformation,proto3" json:"credential_revocation_information,omitempty"`
	XXX_NoUnkeyedLiteral            struct{} `json:"-"`
	XXX_unrecognized                []byte   `json:"-"`
	XXX_sizecache                   int32    `json:"-"`
}

func (m *IdemixMSPSignerConfig) Reset()         { *m = IdemixMSPSignerConfig{} }
func (m *IdemixMSPSignerConfig) String() string { return proto.CompactTextString(m) }
func (*IdemixMSPSignerConfig) ProtoMessage()    {}
func (*IdemixMSPSignerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{4}
}
func (m *IdemixMSPSignerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdemixMSPSignerConfig.Unmarshal(m, b)
}
func (m *IdemixMSPSignerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdemixMSPSignerConfig.Marshal(b, m, deterministic)
}
func (dst *IdemixMSPSignerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdemixMSPSignerConfig.Merge(dst, src)
}
func (m *IdemixMSPSignerConfig) XXX_Size() int {
	return xxx_messageInfo_IdemixMSPSignerConfig.Size(m)
}
func (m *IdemixMSPSignerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_IdemixMSPSignerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_IdemixMSPSignerConfig proto.InternalMessageInfo

func (m *IdemixMSPSignerConfig) GetCred() []byte {
	if m != nil {
		return m.Cred
	}
	return nil
}

func (m *IdemixMSPSignerConfig) GetSk() []byte {
	if m != nil {
		return m.Sk
	}
	return nil
}

func (m *IdemixMSPSignerConfig) GetOrganizationalUnitIdentifier() string {
	if m != nil {
		return m.OrganizationalUnitIdentifier
	}
	return ""
}

func (m *IdemixMSPSignerConfig) GetIsAdmin() bool {
	if m != nil {
		return m.IsAdmin
	}
	return false
}

func (m *IdemixMSPSignerConfig) GetEnrollmentId() string {
	if m != nil {
		return m.EnrollmentId
	}
	return ""
}

func (m *IdemixMSPSignerConfig) GetCredentialRevocationInformation() []byte {
	if m != nil {
		return m.CredentialRevocationInformation
	}
	return nil
}




type SigningIdentityInfo struct {
	
	
	
	PublicSigner []byte `protobuf:"bytes,1,opt,name=public_signer,json=publicSigner,proto3" json:"public_signer,omitempty"`
	
	
	PrivateSigner        *KeyInfo `protobuf:"bytes,2,opt,name=private_signer,json=privateSigner" json:"private_signer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SigningIdentityInfo) Reset()         { *m = SigningIdentityInfo{} }
func (m *SigningIdentityInfo) String() string { return proto.CompactTextString(m) }
func (*SigningIdentityInfo) ProtoMessage()    {}
func (*SigningIdentityInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{5}
}
func (m *SigningIdentityInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SigningIdentityInfo.Unmarshal(m, b)
}
func (m *SigningIdentityInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SigningIdentityInfo.Marshal(b, m, deterministic)
}
func (dst *SigningIdentityInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigningIdentityInfo.Merge(dst, src)
}
func (m *SigningIdentityInfo) XXX_Size() int {
	return xxx_messageInfo_SigningIdentityInfo.Size(m)
}
func (m *SigningIdentityInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SigningIdentityInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SigningIdentityInfo proto.InternalMessageInfo

func (m *SigningIdentityInfo) GetPublicSigner() []byte {
	if m != nil {
		return m.PublicSigner
	}
	return nil
}

func (m *SigningIdentityInfo) GetPrivateSigner() *KeyInfo {
	if m != nil {
		return m.PrivateSigner
	}
	return nil
}





type KeyInfo struct {
	
	
	
	KeyIdentifier string `protobuf:"bytes,1,opt,name=key_identifier,json=keyIdentifier" json:"key_identifier,omitempty"`
	
	
	KeyMaterial          []byte   `protobuf:"bytes,2,opt,name=key_material,json=keyMaterial,proto3" json:"key_material,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeyInfo) Reset()         { *m = KeyInfo{} }
func (m *KeyInfo) String() string { return proto.CompactTextString(m) }
func (*KeyInfo) ProtoMessage()    {}
func (*KeyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{6}
}
func (m *KeyInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyInfo.Unmarshal(m, b)
}
func (m *KeyInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyInfo.Marshal(b, m, deterministic)
}
func (dst *KeyInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyInfo.Merge(dst, src)
}
func (m *KeyInfo) XXX_Size() int {
	return xxx_messageInfo_KeyInfo.Size(m)
}
func (m *KeyInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyInfo.DiscardUnknown(m)
}

var xxx_messageInfo_KeyInfo proto.InternalMessageInfo

func (m *KeyInfo) GetKeyIdentifier() string {
	if m != nil {
		return m.KeyIdentifier
	}
	return ""
}

func (m *KeyInfo) GetKeyMaterial() []byte {
	if m != nil {
		return m.KeyMaterial
	}
	return nil
}



type FabricOUIdentifier struct {
	
	
	
	
	
	
	
	Certificate []byte `protobuf:"bytes,1,opt,name=certificate,proto3" json:"certificate,omitempty"`
	
	
	OrganizationalUnitIdentifier string   `protobuf:"bytes,2,opt,name=organizational_unit_identifier,json=organizationalUnitIdentifier" json:"organizational_unit_identifier,omitempty"`
	XXX_NoUnkeyedLiteral         struct{} `json:"-"`
	XXX_unrecognized             []byte   `json:"-"`
	XXX_sizecache                int32    `json:"-"`
}

func (m *FabricOUIdentifier) Reset()         { *m = FabricOUIdentifier{} }
func (m *FabricOUIdentifier) String() string { return proto.CompactTextString(m) }
func (*FabricOUIdentifier) ProtoMessage()    {}
func (*FabricOUIdentifier) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{7}
}
func (m *FabricOUIdentifier) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FabricOUIdentifier.Unmarshal(m, b)
}
func (m *FabricOUIdentifier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FabricOUIdentifier.Marshal(b, m, deterministic)
}
func (dst *FabricOUIdentifier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FabricOUIdentifier.Merge(dst, src)
}
func (m *FabricOUIdentifier) XXX_Size() int {
	return xxx_messageInfo_FabricOUIdentifier.Size(m)
}
func (m *FabricOUIdentifier) XXX_DiscardUnknown() {
	xxx_messageInfo_FabricOUIdentifier.DiscardUnknown(m)
}

var xxx_messageInfo_FabricOUIdentifier proto.InternalMessageInfo

func (m *FabricOUIdentifier) GetCertificate() []byte {
	if m != nil {
		return m.Certificate
	}
	return nil
}

func (m *FabricOUIdentifier) GetOrganizationalUnitIdentifier() string {
	if m != nil {
		return m.OrganizationalUnitIdentifier
	}
	return ""
}




type FabricNodeOUs struct {
	
	Enable bool `protobuf:"varint,1,opt,name=enable" json:"enable,omitempty"`
	
	ClientOuIdentifier *FabricOUIdentifier `protobuf:"bytes,2,opt,name=client_ou_identifier,json=clientOuIdentifier" json:"client_ou_identifier,omitempty"`
	
	PeerOuIdentifier     *FabricOUIdentifier `protobuf:"bytes,3,opt,name=peer_ou_identifier,json=peerOuIdentifier" json:"peer_ou_identifier,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *FabricNodeOUs) Reset()         { *m = FabricNodeOUs{} }
func (m *FabricNodeOUs) String() string { return proto.CompactTextString(m) }
func (*FabricNodeOUs) ProtoMessage()    {}
func (*FabricNodeOUs) Descriptor() ([]byte, []int) {
	return fileDescriptor_msp_config_9843a417ec9ee543, []int{8}
}
func (m *FabricNodeOUs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FabricNodeOUs.Unmarshal(m, b)
}
func (m *FabricNodeOUs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FabricNodeOUs.Marshal(b, m, deterministic)
}
func (dst *FabricNodeOUs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FabricNodeOUs.Merge(dst, src)
}
func (m *FabricNodeOUs) XXX_Size() int {
	return xxx_messageInfo_FabricNodeOUs.Size(m)
}
func (m *FabricNodeOUs) XXX_DiscardUnknown() {
	xxx_messageInfo_FabricNodeOUs.DiscardUnknown(m)
}

var xxx_messageInfo_FabricNodeOUs proto.InternalMessageInfo

func (m *FabricNodeOUs) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

func (m *FabricNodeOUs) GetClientOuIdentifier() *FabricOUIdentifier {
	if m != nil {
		return m.ClientOuIdentifier
	}
	return nil
}

func (m *FabricNodeOUs) GetPeerOuIdentifier() *FabricOUIdentifier {
	if m != nil {
		return m.PeerOuIdentifier
	}
	return nil
}

func init() {
	proto.RegisterType((*MSPConfig)(nil), "msp.MSPConfig")
	proto.RegisterType((*FabricMSPConfig)(nil), "msp.FabricMSPConfig")
	proto.RegisterType((*FabricCryptoConfig)(nil), "msp.FabricCryptoConfig")
	proto.RegisterType((*IdemixMSPConfig)(nil), "msp.IdemixMSPConfig")
	proto.RegisterType((*IdemixMSPSignerConfig)(nil), "msp.IdemixMSPSignerConfig")
	proto.RegisterType((*SigningIdentityInfo)(nil), "msp.SigningIdentityInfo")
	proto.RegisterType((*KeyInfo)(nil), "msp.KeyInfo")
	proto.RegisterType((*FabricOUIdentifier)(nil), "msp.FabricOUIdentifier")
	proto.RegisterType((*FabricNodeOUs)(nil), "msp.FabricNodeOUs")
}

func init() { proto.RegisterFile("msp/msp_config.proto", fileDescriptor_msp_config_9843a417ec9ee543) }

var fileDescriptor_msp_config_9843a417ec9ee543 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x55, 0x5f, 0x6f, 0xe3, 0x44,
	0x10, 0x57, 0x92, 0x36, 0x6d, 0x26, 0x4e, 0x52, 0xf6, 0x7a, 0xc5, 0x20, 0xee, 0x2e, 0x35, 0x20,
	0xf2, 0x42, 0x2a, 0xf5, 0x90, 0x90, 0x10, 0x2f, 0x5c, 0xe1, 0x84, 0x81, 0xd2, 0x6a, 0xab, 0xbe,
	0xf0, 0x62, 0x6d, 0xec, 0x4d, 0xb2, 0xb2, 0xbd, 0x6b, 0xed, 0xae, 0x4f, 0x04, 0xf1, 0x15, 0x78,
	0xe2, 0x3b, 0xf0, 0xcc, 0x2b, 0xdf, 0x0e, 0xed, 0x9f, 0xc6, 0xce, 0x5d, 0x15, 0x78, 0xdb, 0x99,
	0xf9, 0xcd, 0xcf, 0xb3, 0xbf, 0x99, 0x59, 0xc3, 0x69, 0xa9, 0xaa, 0x8b, 0x52, 0x55, 0x49, 0x2a,
	0xf8, 0x92, 0xad, 0xe6, 0x95, 0x14, 0x5a, 0xa0, 0x5e, 0xa9, 0xaa, 0xe8, 0x4b, 0x18, 0x5c, 0xdf,
	0xdd, 0x5e, 0x59, 0x3f, 0x42, 0x70, 0xa0, 0x37, 0x15, 0x0d, 0x3b, 0xd3, 0xce, 0xec, 0x10, 0xdb,
	0x33, 0x3a, 0x83, 0xbe, 0xcb, 0x0a, 0xbb, 0xd3, 0xce, 0x2c, 0xc0, 0xde, 0x8a, 0xfe, 0x3e, 0x80,
	0xc9, 0x6b, 0xb2, 0x90, 0x2c, 0xdd, 0xc9, 0xe7, 0xa4, 0x74, 0xf9, 0x03, 0x6c, 0xcf, 0xe8, 0x19,
	0x80, 0x14, 0x42, 0x27, 0x29, 0x95, 0x5a, 0x85, 0xdd, 0x69, 0x6f, 0x16, 0xe0, 0x81, 0xf1, 0x5c,
	0x19, 0x07, 0xfa, 0x1c, 0x10, 0xe3, 0x9a, 0xca, 0x92, 0x66, 0x8c, 0x68, 0xea, 0x61, 0x3d, 0x0b,
	0x7b, 0xaf, 0x1d, 0x71, 0xf0, 0x33, 0xe8, 0x93, 0xac, 0x64, 0x5c, 0x85, 0x07, 0x16, 0xe2, 0x2d,
	0xf4, 0x19, 0x4c, 0x24, 0x7d, 0x23, 0x52, 0xa2, 0x99, 0xe0, 0x49, 0xc1, 0x94, 0x0e, 0x0f, 0x2d,
	0x60, 0xdc, 0xb8, 0x7f, 0x62, 0x4a, 0xa3, 0x2b, 0x38, 0x51, 0x6c, 0xc5, 0x19, 0x5f, 0x25, 0x2c,
	0xa3, 0x5c, 0x33, 0xbd, 0x09, 0xfb, 0xd3, 0xce, 0x6c, 0x78, 0x19, 0xce, 0x4b, 0x55, 0xcd, 0xef,
	0x5c, 0x30, 0xf6, 0xb1, 0x98, 0x2f, 0x05, 0x9e, 0xa8, 0x5d, 0x27, 0x4a, 0xe0, 0x85, 0x90, 0x2b,
	0xc2, 0xd9, 0x6f, 0x96, 0x98, 0x14, 0x49, 0xcd, 0x99, 0xf6, 0x84, 0x4b, 0x46, 0xa5, 0x0a, 0x8f,
	0xa6, 0xbd, 0xd9, 0xf0, 0xf2, 0x7d, 0xcb, 0xe9, 0x64, 0xba, 0xb9, 0x8f, 0xb7, 0x71, 0xfc, 0x6c,
	0x37, 0xff, 0x9e, 0x33, 0xdd, 0x44, 0x15, 0xfa, 0x1a, 0x46, 0xa9, 0xdc, 0x54, 0x5a, 0xf8, 0x8e,
	0x85, 0xc7, 0xb6, 0xc4, 0x36, 0xdd, 0x95, 0x8d, 0x3b, 0xe1, 0x71, 0x90, 0xb6, 0x2c, 0xf4, 0x09,
	0x8c, 0x75, 0xa1, 0x92, 0x96, 0xec, 0x03, 0xab, 0x45, 0xa0, 0x0b, 0x85, 0xb7, 0xca, 0x7f, 0x01,
	0x67, 0x06, 0xf5, 0x88, 0xfa, 0x60, 0xd1, 0xa7, 0xba, 0x50, 0xf1, 0x3b, 0x0d, 0xf8, 0x0a, 0x26,
	0x4b, 0xfb, 0xfd, 0x84, 0x8b, 0x8c, 0x26, 0xa2, 0x56, 0xe1, 0xd0, 0xd6, 0x86, 0x5a, 0xb5, 0xfd,
	0x2c, 0x32, 0x7a, 0x73, 0xaf, 0xf0, 0x68, 0xd9, 0x98, 0xb5, 0x8a, 0xfe, 0xec, 0x00, 0x7a, 0xb7,
	0x78, 0x74, 0x09, 0x4f, 0x8d, 0xc0, 0x44, 0xd7, 0x92, 0x26, 0x6b, 0xa2, 0xd6, 0xc9, 0x92, 0x94,
	0xac, 0xd8, 0xf8, 0x31, 0x7a, 0xb2, 0x0d, 0x7e, 0x4f, 0xd4, 0xfa, 0xb5, 0x0d, 0xa1, 0x18, 0xce,
	0x1f, 0xda, 0xd7, 0x92, 0xdd, 0x67, 0xd7, 0x3c, 0x35, 0xb2, 0xda, 0x81, 0x1d, 0xe0, 0xe7, 0x0f,
	0xc0, 0x46, 0x60, 0x4b, 0xe4, 0x51, 0xd1, 0x5f, 0x1d, 0x98, 0xc4, 0x19, 0x2d, 0xd9, 0xaf, 0xfb,
	0x07, 0xf9, 0x04, 0x7a, 0xac, 0xca, 0xfd, 0x16, 0x98, 0x23, 0xba, 0x84, 0xbe, 0xa9, 0x8d, 0xca,
	0xb0, 0x67, 0x25, 0xf8, 0xd0, 0x4a, 0xb0, 0xe5, 0xba, 0xb3, 0x31, 0xdf, 0x21, 0x8f, 0x44, 0x1f,
	0xc3, 0xa8, 0x35, 0xa8, 0x55, 0x1e, 0x1e, 0x58, 0xbe, 0xa0, 0x71, 0xde, 0xe6, 0xe8, 0x14, 0x0e,
	0x69, 0x25, 0xd2, 0x75, 0x78, 0x38, 0xed, 0xcc, 0x7a, 0xd8, 0x19, 0xd1, 0x1f, 0x5d, 0x78, 0xfa,
	0x28, 0xb9, 0x29, 0x37, 0x95, 0x34, 0xb3, 0xe5, 0x06, 0xd8, 0x9e, 0xd1, 0x18, 0xba, 0xea, 0xa1,
	0xda, 0xae, 0xca, 0xd1, 0xb7, 0xf0, 0x7c, 0xff, 0xcc, 0xda, 0x4b, 0x0c, 0xf0, 0x47, 0xfb, 0x26,
	0x13, 0x7d, 0x00, 0xc7, 0x4c, 0x25, 0x76, 0xe9, 0x6c, 0xe5, 0xc7, 0xf8, 0x88, 0xa9, 0x6f, 0x8c,
	0x69, 0x6e, 0x46, 0xb9, 0x14, 0x45, 0x51, 0x52, 0x6e, 0x78, 0x6d, 0xf1, 0x03, 0x1c, 0x34, 0xce,
	0x38, 0x43, 0x3f, 0xc0, 0xb9, 0xa9, 0xce, 0xf0, 0x91, 0x22, 0x69, 0x29, 0xc1, 0xf8, 0x52, 0xc8,
	0xd2, 0x9e, 0xed, 0x3e, 0x06, 0xf8, 0x45, 0x03, 0xc4, 0x5b, 0x5c, 0xdc, 0xc0, 0x22, 0x01, 0x4f,
	0x1e, 0xd9, 0x56, 0x53, 0x47, 0x55, 0x2f, 0x0a, 0x96, 0x26, 0xbe, 0x39, 0x4e, 0x95, 0xc0, 0x39,
	0x9d, 0x6e, 0xe8, 0x25, 0x8c, 0x2b, 0xc9, 0xde, 0x98, 0x99, 0xf7, 0xa8, 0xae, 0x6d, 0x61, 0x60,
	0x5b, 0xf8, 0x23, 0x75, 0x8b, 0x3f, 0xf2, 0x18, 0x97, 0x14, 0xdd, 0xc1, 0x91, 0x8f, 0xa0, 0x4f,
	0x61, 0x9c, 0xd3, 0xf6, 0xe8, 0xf9, 0x51, 0x19, 0xe5, 0xb4, 0x35, 0x67, 0xe8, 0x1c, 0x02, 0x03,
	0x2b, 0x89, 0xa6, 0x92, 0x91, 0xc2, 0xb7, 0x63, 0x98, 0xd3, 0xcd, 0xb5, 0x77, 0x45, 0xbf, 0x3f,
	0xec, 0x44, 0xfb, 0x7d, 0x40, 0x53, 0x18, 0x9a, 0x5d, 0x64, 0x4b, 0x96, 0x12, 0x4d, 0xfd, 0x15,
	0xda, 0xae, 0xff, 0xd1, 0xcf, 0xee, 0x7f, 0xf7, 0x33, 0xfa, 0xa7, 0x03, 0xa3, 0x9d, 0x9d, 0x35,
	0x2f, 0x2c, 0xe5, 0x64, 0x51, 0xb8, 0x8f, 0x1e, 0x63, 0x6f, 0xa1, 0x18, 0x4e, 0xd3, 0x82, 0x99,
	0xd6, 0x8a, 0xfa, 0xed, 0xaf, 0xec, 0x79, 0xe8, 0x90, 0x4b, 0xba, 0xa9, 0x5b, 0x97, 0xfb, 0x0e,
	0x50, 0x45, 0xa9, 0x7c, 0x8b, 0xa8, 0xb7, 0x9f, 0xe8, 0xc4, 0xa4, 0xb4, 0x69, 0x5e, 0x25, 0x70,
	0x2e, 0xe4, 0x6a, 0xbe, 0xde, 0x54, 0x54, 0x16, 0x34, 0x5b, 0x51, 0x39, 0x77, 0xef, 0x8d, 0xfb,
	0xbf, 0x29, 0xc3, 0xf4, 0xea, 0xe4, 0x5a, 0x55, 0x6e, 0x4b, 0x6e, 0x49, 0x9a, 0x93, 0x15, 0xfd,
	0x65, 0xb6, 0x62, 0x7a, 0x5d, 0x2f, 0xe6, 0xa9, 0x28, 0x2f, 0x5a, 0xb9, 0x17, 0x2e, 0xf7, 0xc2,
	0xe5, 0x9a, 0xbf, 0xe5, 0xa2, 0x6f, 0xcf, 0x2f, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x91,
	0xb8, 0x49, 0x3f, 0x07, 0x00, 0x00,
}
