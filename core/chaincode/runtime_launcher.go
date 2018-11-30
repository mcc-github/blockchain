/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"strconv"
	"time"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/pkg/errors"
)


type LaunchRegistry interface {
	Launching(cname string) (launchState *LaunchState, started bool)
	Deregister(cname string) error
}


type PackageProvider interface {
	GetChaincodeCodePackage(ccname string, ccversion string) ([]byte, error)
}


type RuntimeLauncher struct {
	Runtime         Runtime
	Registry        LaunchRegistry
	PackageProvider PackageProvider
	StartupTimeout  time.Duration
	Metrics         *LaunchMetrics
}

func (r *RuntimeLauncher) Launch(ccci *ccprovider.ChaincodeContainerInfo) error {
	var startFailCh chan error
	var timeoutCh <-chan time.Time

	startTime := time.Now()
	cname := ccci.Name + ":" + ccci.Version
	launchState, started := r.Registry.Launching(cname)
	if !started {
		startFailCh = make(chan error, 1)
		timeoutCh = time.NewTimer(r.StartupTimeout).C

		codePackage, err := r.getCodePackage(ccci)
		if err != nil {
			return err
		}

		go func() {
			if err := r.Runtime.Start(ccci, codePackage); err != nil {
				startFailCh <- errors.WithMessage(err, "error starting container")
			}
		}()
	}

	var err error
	select {
	case <-launchState.Done():
		err = errors.WithMessage(launchState.Err(), "chaincode registration failed")
	case err = <-startFailCh:
		launchState.Notify(err)
		r.Metrics.LaunchFailures.With(
			"chaincode", cname,
		).Add(1)
	case <-timeoutCh:
		err = errors.Errorf("timeout expired while starting chaincode %s for transaction", cname)
		launchState.Notify(err)
		r.Metrics.LaunchTimeouts.With(
			"chaincode", cname,
		).Add(1)
	}

	success := true
	if err != nil && !started {
		success = false
		chaincodeLogger.Debugf("stopping due to error while launching: %+v", err)
		defer r.Registry.Deregister(cname)
		if err := r.Runtime.Stop(ccci); err != nil {
			chaincodeLogger.Debugf("stop failed: %+v", err)
		}
	}

	r.Metrics.LaunchDuration.With(
		"chaincode", cname,
		"success", strconv.FormatBool(success),
	).Observe(time.Since(startTime).Seconds())

	chaincodeLogger.Debug("launch complete")
	return err
}

func (r *RuntimeLauncher) getCodePackage(ccci *ccprovider.ChaincodeContainerInfo) ([]byte, error) {
	if ccci.ContainerType == inproccontroller.ContainerType {
		return nil, nil
	}

	codePackage, err := r.PackageProvider.GetChaincodeCodePackage(ccci.Name, ccci.Version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chaincode package")
	}

	return codePackage, nil
}
