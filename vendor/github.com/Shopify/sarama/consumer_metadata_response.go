package sarama

import (
	"net"
	"strconv"
)

type ConsumerMetadataResponse struct {
	Err             KError
	Coordinator     *Broker
	CoordinatorID   int32  
	CoordinatorHost string 
	CoordinatorPort int32  
}

func (r *ConsumerMetadataResponse) decode(pd packetDecoder, version int16) (err error) {
	tmp, err := pd.getInt16()
	if err != nil {
		return err
	}
	r.Err = KError(tmp)

	coordinator := new(Broker)
	if err := coordinator.decode(pd); err != nil {
		return err
	}
	if coordinator.addr == ":0" {
		return nil
	}
	r.Coordinator = coordinator

	
	
	host, portstr, err := net.SplitHostPort(r.Coordinator.Addr())
	if err != nil {
		return err
	}
	port, err := strconv.ParseInt(portstr, 10, 32)
	if err != nil {
		return err
	}
	r.CoordinatorID = r.Coordinator.ID()
	r.CoordinatorHost = host
	r.CoordinatorPort = int32(port)

	return nil
}

func (r *ConsumerMetadataResponse) encode(pe packetEncoder) error {
	pe.putInt16(int16(r.Err))
	if r.Coordinator != nil {
		host, portstr, err := net.SplitHostPort(r.Coordinator.Addr())
		if err != nil {
			return err
		}
		port, err := strconv.ParseInt(portstr, 10, 32)
		if err != nil {
			return err
		}
		pe.putInt32(r.Coordinator.ID())
		if err := pe.putString(host); err != nil {
			return err
		}
		pe.putInt32(int32(port))
		return nil
	}
	pe.putInt32(r.CoordinatorID)
	if err := pe.putString(r.CoordinatorHost); err != nil {
		return err
	}
	pe.putInt32(r.CoordinatorPort)
	return nil
}

func (r *ConsumerMetadataResponse) key() int16 {
	return 10
}

func (r *ConsumerMetadataResponse) version() int16 {
	return 0
}

func (r *ConsumerMetadataResponse) requiredVersion() KafkaVersion {
	return V0_8_2_0
}
