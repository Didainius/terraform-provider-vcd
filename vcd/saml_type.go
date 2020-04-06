package vcd

import "encoding/xml"

type ResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	A       string   `xml:"a,attr"`
	U       string   `xml:"u,attr"`
	Header  struct {
		Text   string `xml:",chardata"`
		Action struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
		} `xml:"Action"`
		Security struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
			O              string `xml:"o,attr"`
			Timestamp      struct {
				Text    string `xml:",chardata"`
				ID      string `xml:"Id,attr"`
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Timestamp"`
		} `xml:"Security"`
	} `xml:"Header"`
	Body struct {
		Text                                   string `xml:",chardata"`
		RequestSecurityTokenResponseCollection struct {
			Text                         string `xml:",chardata"`
			Trust                        string `xml:"trust,attr"`
			RequestSecurityTokenResponse struct {
				Text     string `xml:",chardata"`
				Lifetime struct {
					Text    string `xml:",chardata"`
					Created struct {
						Text string `xml:",chardata"`
						Wsu  string `xml:"wsu,attr"`
					} `xml:"Created"`
					Expires struct {
						Text string `xml:",chardata"`
						Wsu  string `xml:"wsu,attr"`
					} `xml:"Expires"`
				} `xml:"Lifetime"`
				AppliesTo struct {
					Text              string `xml:",chardata"`
					Wsp               string `xml:"wsp,attr"`
					EndpointReference struct {
						Text    string `xml:",chardata"`
						Wsa     string `xml:"wsa,attr"`
						Address string `xml:"Address"`
					} `xml:"EndpointReference"`
				} `xml:"AppliesTo"`
				// RequestedSecurityToken struct {
				// 	Text         string   `xml:",chardata"`
				RequestedSecurityTokenTxt InnerXML `xml:"RequestedSecurityToken"`
				// } `xml:"RequestedSecurityToken"`
				RequestedDisplayToken struct {
					Text         string `xml:",chardata"`
					I            string `xml:"i,attr"`
					DisplayToken struct {
						Text         string `xml:",chardata"`
						Lang         string `xml:"lang,attr"`
						DisplayClaim []struct {
							Text         string `xml:",chardata"`
							URI          string `xml:"Uri,attr"`
							DisplayTag   string `xml:"DisplayTag"`
							Description  string `xml:"Description"`
							DisplayValue string `xml:"DisplayValue"`
						} `xml:"DisplayClaim"`
					} `xml:"DisplayToken"`
				} `xml:"RequestedDisplayToken"`
				RequestedAttachedReference struct {
					Text                   string `xml:",chardata"`
					SecurityTokenReference struct {
						Text          string `xml:",chardata"`
						TokenType     string `xml:"TokenType,attr"`
						Xmlns         string `xml:"xmlns,attr"`
						B             string `xml:"b,attr"`
						KeyIdentifier struct {
							Text      string `xml:",chardata"`
							ValueType string `xml:"ValueType,attr"`
						} `xml:"KeyIdentifier"`
					} `xml:"SecurityTokenReference"`
				} `xml:"RequestedAttachedReference"`
				RequestedUnattachedReference struct {
					Text                   string `xml:",chardata"`
					SecurityTokenReference struct {
						Text          string `xml:",chardata"`
						TokenType     string `xml:"TokenType,attr"`
						Xmlns         string `xml:"xmlns,attr"`
						B             string `xml:"b,attr"`
						KeyIdentifier struct {
							Text      string `xml:",chardata"`
							ValueType string `xml:"ValueType,attr"`
						} `xml:"KeyIdentifier"`
					} `xml:"SecurityTokenReference"`
				} `xml:"RequestedUnattachedReference"`
				TokenType   string `xml:"TokenType"`
				RequestType string `xml:"RequestType"`
				KeyType     string `xml:"KeyType"`
			} `xml:"RequestSecurityTokenResponse"`
		} `xml:"RequestSecurityTokenResponseCollection"`
	} `xml:"Body"`
}

type InnerXML struct {
	Text string `xml:",innerxml"`
}

// type EncryptedAssertion struct {
// 	XMLName       xml.Name `xml:"EncryptedAssertion"`
// 	Xmlns         string   `xml:"xmlns,attr"`
// 	EncryptedData struct {
// 		Type             string `xml:"Type,attr"`
// 		Xenc             string `xml:"xenc,attr"`
// 		EncryptionMethod struct {
// 			Algorithm string `xml:"Algorithm,attr"`
// 		} `xml:"EncryptionMethod"`
// 		KeyInfo struct {
// 			Xmlns        string `xml:"xmlns,attr"`
// 			EncryptedKey struct {
// 				E                string `xml:"e,attr"`
// 				EncryptionMethod struct {
// 					Algorithm    string `xml:"Algorithm,attr"`
// 					DigestMethod struct {
// 						Text      string `xml:",chardata"`
// 						Algorithm string `xml:"Algorithm,attr"`
// 					} `xml:"DigestMethod"`
// 				} `xml:"EncryptionMethod"`
// 				KeyInfo struct {
// 					X509Data struct {
// 						Text             string `xml:",chardata"`
// 						Ds               string `xml:"ds,attr"`
// 						X509IssuerSerial struct {
// 							Text             string `xml:",chardata"`
// 							X509IssuerName   string `xml:"X509IssuerName"`
// 							X509SerialNumber string `xml:"X509SerialNumber"`
// 						} `xml:"X509IssuerSerial"`
// 					} `xml:"X509Data"`
// 				} `xml:"KeyInfo"`
// 				CipherData struct {
// 					CipherValue string `xml:"CipherValue"`
// 				} `xml:"CipherData"`
// 			} `xml:"EncryptedKey"`
// 		} `xml:"KeyInfo"`
// 		CipherData struct {
// 			CipherValue string `xml:"CipherValue"`
// 		} `xml:"CipherData"`
// 	} `xml:"EncryptedData"`
// }

type EncryptedAssertion struct {
	XMLName       xml.Name `xml:"EncryptedAssertion"`
	Text          string   `xml:",chardata"`
	Xmlns         string   `xml:"xmlns,attr"`
	Trust         string   `xml:"trust,attr"`
	S             string   `xml:"s,attr"`
	A             string   `xml:"a,attr"`
	U             string   `xml:"u,attr"`
	EncryptedData struct {
		Text             string `xml:",chardata"`
		Xenc             string `xml:"xenc,attr"`
		Type             string `xml:"Type,attr"`
		EncryptionMethod struct {
			Text      string `xml:",chardata"`
			Algorithm string `xml:"Algorithm,attr"`
		} `xml:"EncryptionMethod"`
		KeyInfo struct {
			Text         string `xml:",chardata"`
			Xmlns        string `xml:"xmlns,attr"`
			EncryptedKey struct {
				Text             string `xml:",chardata"`
				E                string `xml:"e,attr"`
				EncryptionMethod struct {
					Text         string `xml:",chardata"`
					Algorithm    string `xml:"Algorithm,attr"`
					DigestMethod struct {
						Text      string `xml:",chardata"`
						Algorithm string `xml:"Algorithm,attr"`
					} `xml:"DigestMethod"`
				} `xml:"EncryptionMethod"`
				KeyInfo struct {
					Text     string `xml:",chardata"`
					X509Data struct {
						Text             string `xml:",chardata"`
						Ds               string `xml:"ds,attr"`
						X509IssuerSerial struct {
							Text             string `xml:",chardata"`
							X509IssuerName   string `xml:"X509IssuerName"`
							X509SerialNumber string `xml:"X509SerialNumber"`
						} `xml:"X509IssuerSerial"`
					} `xml:"X509Data"`
				} `xml:"KeyInfo"`
				CipherData struct {
					Text        string `xml:",chardata"`
					CipherValue string `xml:"CipherValue"`
				} `xml:"CipherData"`
			} `xml:"EncryptedKey"`
		} `xml:"KeyInfo"`
		CipherData struct {
			Text        string `xml:",chardata"`
			CipherValue string `xml:"CipherValue"`
		} `xml:"CipherData"`
	} `xml:"EncryptedData"`
}

type ErrorEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	A       string   `xml:"a,attr"`
	U       string   `xml:"u,attr"`
	Header  struct {
		Text   string `xml:",chardata"`
		Action struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
		} `xml:"Action"`
		Security struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
			O              string `xml:"o,attr"`
			Timestamp      struct {
				Text    string `xml:",chardata"`
				ID      string `xml:"Id,attr"`
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Timestamp"`
		} `xml:"Security"`
	} `xml:"Header"`
	Body struct {
		Text  string `xml:",chardata"`
		Fault struct {
			Text string `xml:",chardata"`
			Code struct {
				Text    string `xml:",chardata"`
				Value   string `xml:"Value"`
				Subcode struct {
					Text  string `xml:",chardata"`
					Value struct {
						Text string `xml:",chardata"`
						A    string `xml:"a,attr"`
					} `xml:"Value"`
				} `xml:"Subcode"`
			} `xml:"Code"`
			Reason struct {
				Chardata string `xml:",chardata"`
				Text     struct {
					Text string `xml:",chardata"`
					Lang string `xml:"lang,attr"`
				} `xml:"Text"`
			} `xml:"Reason"`
		} `xml:"Fault"`
	} `xml:"Body"`
}

type EntityDescriptor struct {
	XMLName         xml.Name `xml:"EntityDescriptor"`
	Text            string   `xml:",chardata"`
	ID              string   `xml:"ID,attr"`
	EntityID        string   `xml:"entityID,attr"`
	Md              string   `xml:"md,attr"`
	SPSSODescriptor struct {
		Text                       string `xml:",chardata"`
		AuthnRequestsSigned        string `xml:"AuthnRequestsSigned,attr"`
		WantAssertionsSigned       string `xml:"WantAssertionsSigned,attr"`
		ProtocolSupportEnumeration string `xml:"protocolSupportEnumeration,attr"`
		KeyDescriptor              []struct {
			Text    string `xml:",chardata"`
			Use     string `xml:"use,attr"`
			KeyInfo struct {
				Text     string `xml:",chardata"`
				Ds       string `xml:"ds,attr"`
				X509Data struct {
					Text            string `xml:",chardata"`
					X509Certificate string `xml:"X509Certificate"`
				} `xml:"X509Data"`
			} `xml:"KeyInfo"`
		} `xml:"KeyDescriptor"`
		SingleLogoutService []struct {
			Text     string `xml:",chardata"`
			Binding  string `xml:"Binding,attr"`
			Location string `xml:"Location,attr"`
		} `xml:"SingleLogoutService"`
		NameIDFormat             []string `xml:"NameIDFormat"`
		AssertionConsumerService []struct {
			Text            string `xml:",chardata"`
			Binding         string `xml:"Binding,attr"`
			Location        string `xml:"Location,attr"`
			Index           string `xml:"index,attr"`
			IsDefault       string `xml:"isDefault,attr"`
			ProtocolBinding string `xml:"ProtocolBinding,attr"`
			Hoksso          string `xml:"hoksso,attr"`
		} `xml:"AssertionConsumerService"`
	} `xml:"SPSSODescriptor"`
}
