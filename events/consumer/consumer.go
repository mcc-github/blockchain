/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package consumer

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/comm"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	ehpb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
)

var consumerLogger = flogging.MustGetLogger("eventhub_consumer")


type EventsClient struct {
	sync.RWMutex
	peerAddress string
	regTimeout  time.Duration
	stream      ehpb.Events_ChatClient
	adapter     EventAdapter
}



type RegistrationConfig struct {
	InterestedEvents []*ehpb.Interest
	Timestamp        *timestamp.Timestamp
	TlsCert          *x509.Certificate
}


func NewEventsClient(peerAddress string, regTimeout time.Duration, adapter EventAdapter) (*EventsClient, error) {
	var err error
	if regTimeout < 100*time.Millisecond {
		regTimeout = 100 * time.Millisecond
		err = fmt.Errorf("regTimeout >= 0, setting to 100 msec")
	} else if regTimeout > 60*time.Second {
		regTimeout = 60 * time.Second
		err = fmt.Errorf("regTimeout > 60, setting to 60 sec")
	}
	return &EventsClient{sync.RWMutex{}, peerAddress, regTimeout, nil, adapter}, err
}


func newEventsClientConnectionWithAddress(peerAddress string) (*grpc.ClientConn, error) {
	return comm.NewClientConnectionWithAddress(peerAddress, true, false,
		nil, nil)
}

func (ec *EventsClient) send(emsg *ehpb.Event) error {
	ec.Lock()
	defer ec.Unlock()

	
	localMsp := mspmgmt.GetLocalMSP()
	if localMsp == nil {
		return errors.New("nil local MSP manager")
	}

	signer, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		return fmt.Errorf("could not obtain the default signing identity, err %s", err)
	}

	
	signerCert, err := signer.Serialize()
	if err != nil {
		return fmt.Errorf("fail to serialize the default signing identity, err %s", err)
	}
	emsg.Creator = signerCert

	signedEvt, err := utils.GetSignedEvent(emsg, signer)
	if err != nil {
		return fmt.Errorf("could not sign outgoing event, err %s", err)
	}

	return ec.stream.Send(signedEvt)
}


func (ec *EventsClient) RegisterAsync(config *RegistrationConfig) error {
	creator, err := getCreatorFromLocalMSP()
	if err != nil {
		return fmt.Errorf("error getting creator from MSP: %s", err)
	}
	emsg := &ehpb.Event{Event: &ehpb.Event_Register{Register: &ehpb.Register{Events: config.InterestedEvents}}, Creator: creator, Timestamp: config.Timestamp}

	if config.TlsCert != nil {
		emsg.TlsCertHash = util.ComputeSHA256(config.TlsCert.Raw)
	}
	if err = ec.send(emsg); err != nil {
		consumerLogger.Errorf("error on Register send %s\n", err)
	}
	return err
}


func (ec *EventsClient) register(config *RegistrationConfig) error {
	var err error
	if err = ec.RegisterAsync(config); err != nil {
		return err
	}

	regChan := make(chan struct{})
	go func() {
		defer close(regChan)
		in, inerr := ec.stream.Recv()
		if inerr != nil {
			err = inerr
			return
		}
		switch in.Event.(type) {
		case *ehpb.Event_Register:
		case nil:
			err = fmt.Errorf("invalid nil object for register")
		default:
			err = fmt.Errorf("invalid registration object")
		}
	}()
	select {
	case <-regChan:
	case <-time.After(ec.regTimeout):
		err = fmt.Errorf("timeout waiting for registration")
	}
	return err
}


func (ec *EventsClient) UnregisterAsync(ies []*ehpb.Interest) error {
	creator, err := getCreatorFromLocalMSP()
	if err != nil {
		return fmt.Errorf("error getting creator from MSP: %s", err)
	}
	emsg := &ehpb.Event{Event: &ehpb.Event_Unregister{Unregister: &ehpb.Unregister{Events: ies}}, Creator: creator}

	if err = ec.send(emsg); err != nil {
		err = fmt.Errorf("error on unregister send: %s", err)
	}

	return err
}


func (ec *EventsClient) Recv() (*ehpb.Event, error) {
	in, err := ec.stream.Recv()
	if err == io.EOF {
		
		if ec.adapter != nil {
			ec.adapter.Disconnected(nil)
		}
		return nil, err
	}
	if err != nil {
		if ec.adapter != nil {
			ec.adapter.Disconnected(err)
		}
		return nil, err
	}
	return in, nil
}

func (ec *EventsClient) processEvents() error {
	defer ec.stream.CloseSend()
	for {
		in, err := ec.stream.Recv()
		if err == io.EOF {
			
			if ec.adapter != nil {
				ec.adapter.Disconnected(nil)
			}
			return nil
		}
		if err != nil {
			if ec.adapter != nil {
				ec.adapter.Disconnected(err)
			}
			return err
		}
		if ec.adapter != nil {
			cont, err := ec.adapter.Recv(in)
			if !cont {
				return err
			}
		}
	}
}


func (ec *EventsClient) Start() error {
	conn, err := newEventsClientConnectionWithAddress(ec.peerAddress)
	if err != nil {
		return fmt.Errorf("could not create client conn to %s:%s", ec.peerAddress, err)
	}

	ies, err := ec.adapter.GetInterestedEvents()
	if err != nil {
		return fmt.Errorf("error getting interested events:%s", err)
	}

	if len(ies) == 0 {
		return fmt.Errorf("must supply interested events")
	}

	serverClient := ehpb.NewEventsClient(conn)
	ec.stream, err = serverClient.Chat(context.Background())
	if err != nil {
		return fmt.Errorf("could not create client conn to %s:%s", ec.peerAddress, err)
	}

	regConfig := &RegistrationConfig{InterestedEvents: ies, Timestamp: util.CreateUtcTimestamp()}
	if err = ec.register(regConfig); err != nil {
		return err
	}

	go ec.processEvents()

	return nil
}


func (ec *EventsClient) Stop() error {
	if ec.stream == nil {
		
		return nil
	}
	return ec.stream.CloseSend()
}

func getCreatorFromLocalMSP() ([]byte, error) {
	localMsp := mspmgmt.GetLocalMSP()
	if localMsp == nil {
		return nil, errors.New("nil local MSP manager")
	}
	signer, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		return nil, fmt.Errorf("could not obtain the default signing identity, err %s", err)
	}
	creator, err := signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing the signer: %s", err)
	}
	return creator, nil
}
