/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"sync/atomic"
)


type TLSCertificates struct {
	TLSServerCert atomic.Value 
	TLSClientCert atomic.Value 
}
