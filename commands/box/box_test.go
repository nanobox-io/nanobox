// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package box

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/nanobox-io/nanobox-cli/config/mock"
	"github.com/nanobox-io/nanobox-cli/util/mock"
	"github.com/nanobox-io/nanobox-cli/util/vagrant/mock"
	"github.com/spf13/cobra"
	"testing"
)

func newUtil(mockCtrl *gomock.Controller) *mock_util.MockUtil {
	testUtil := mock_util.NewMockUtil(mockCtrl)
	Util = testUtil
	return testUtil
}

func newConfig(mockCtrl *gomock.Controller) *mock_config.MockConfig {
	testConfig := mock_config.NewMockConfig(mockCtrl)
	Config = testConfig
	return testConfig
}

func newVagrant(mockCtrl *gomock.Controller) *mock_vagrant.MockVagrant {
	testVagrant := mock_vagrant.NewMockVagrant(mockCtrl)
	Vagrant = testVagrant
	return testVagrant
}

func TestInstallWithImage(test *testing.T) {
	mockCtrl := gomock.NewController(test)
	defer mockCtrl.Finish()

	testVagrant := newVagrant(mockCtrl)

	testVagrant.EXPECT().HaveImage().Return(true)

	Install(&cobra.Command{}, []string{})
}

func TestInstallWithoutImage(test *testing.T) {
	mockCtrl := gomock.NewController(test)
	defer mockCtrl.Finish()

	testVagrant := newVagrant(mockCtrl)

	testVagrant.EXPECT().HaveImage().Return(false)
	testVagrant.EXPECT().Install()

	Install(&cobra.Command{}, []string{})
}

func TestInstallFail(test *testing.T) {
	mockCtrl := gomock.NewController(test)
	defer mockCtrl.Finish()

	testVagrant := newVagrant(mockCtrl)
	testConfig := newConfig(mockCtrl)
	err := errors.New("something went wrong")

	testVagrant.EXPECT().HaveImage().Return(false)
	testVagrant.EXPECT().Install().Return(err)

	testConfig.EXPECT().Fatal(gomock.Any(), err.Error())

	Install(&cobra.Command{}, []string{})
}

func TestUpdateNotNeeded(test *testing.T) {
	mockCtrl := gomock.NewController(test)
	defer mockCtrl.Finish()

	testVagrant := newVagrant(mockCtrl)
	testConfig := newConfig(mockCtrl)
	testUtil := newUtil(mockCtrl)

	testVagrant.EXPECT().HaveImage().Return(false)
	testVagrant.EXPECT().Install()

	testConfig.EXPECT().Root().Return("")
	testUtil.EXPECT().MD5sMatch(gomock.Any(), gomock.Any()).Return(true, nil)

	Update(&cobra.Command{}, []string{})
}

func TestUpdateNeeded(test *testing.T) {
	mockCtrl := gomock.NewController(test)
	defer mockCtrl.Finish()

	testVagrant := newVagrant(mockCtrl)
	testConfig := newConfig(mockCtrl)
	testUtil := newUtil(mockCtrl)

	testVagrant.EXPECT().HaveImage().Return(false)
	testVagrant.EXPECT().Install()

	testConfig.EXPECT().Root().Return("")
	testUtil.EXPECT().MD5sMatch(gomock.Any(), gomock.Any()).Return(false, nil)

	testVagrant.EXPECT().Update()

	Update(&cobra.Command{}, []string{})
}

func TestUpdateFailed(test *testing.T) {
	mockCtrl := gomock.NewController(test)
	defer mockCtrl.Finish()

	testVagrant := newVagrant(mockCtrl)
	testConfig := newConfig(mockCtrl)
	testUtil := newUtil(mockCtrl)

	testVagrant.EXPECT().HaveImage().Return(false)
	testVagrant.EXPECT().Install()

	testConfig.EXPECT().Root().Return("")
	testUtil.EXPECT().MD5sMatch(gomock.Any(), gomock.Any()).Return(false, nil)
	err := errors.New("something went wrong")
	testVagrant.EXPECT().Update().Return(err)

	testConfig.EXPECT().Fatal(gomock.Any(), err.Error())

	Update(&cobra.Command{}, []string{})
}
