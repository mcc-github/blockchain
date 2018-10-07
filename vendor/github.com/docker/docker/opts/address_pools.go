package opts

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	types "github.com/docker/libnetwork/ipamutils"
)


type PoolsOpt struct {
	values []*types.NetworkToSplit
}


func (p *PoolsOpt) UnmarshalJSON(raw []byte) error {
	return json.Unmarshal(raw, &(p.values))
}


func (p *PoolsOpt) Set(value string) error {
	csvReader := csv.NewReader(strings.NewReader(value))
	fields, err := csvReader.Read()
	if err != nil {
		return err
	}

	poolsDef := types.NetworkToSplit{}

	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid field '%s' must be a key=value pair", field)
		}

		key := strings.ToLower(parts[0])
		value := strings.ToLower(parts[1])

		switch key {
		case "base":
			poolsDef.Base = value
		case "size":
			size, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid size value: %q (must be integer): %v", value, err)
			}
			poolsDef.Size = size
		default:
			return fmt.Errorf("unexpected key '%s' in '%s'", key, field)
		}
	}

	p.values = append(p.values, &poolsDef)

	return nil
}


func (p *PoolsOpt) Type() string {
	return "pool-options"
}


func (p *PoolsOpt) String() string {
	var pools []string
	for _, pool := range p.values {
		repr := fmt.Sprintf("%s %d", pool.Base, pool.Size)
		pools = append(pools, repr)
	}
	return strings.Join(pools, ", ")
}


func (p *PoolsOpt) Value() []*types.NetworkToSplit {
	return p.values
}


func (p *PoolsOpt) Name() string {
	return "default-address-pools"
}
