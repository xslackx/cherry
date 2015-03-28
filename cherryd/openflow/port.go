/*
 * Cherry - An OpenFlow Controller
 *
 * Copyright (C) 2015 Samjung Data Service Co., Ltd.,
 * Kitae Kim <superkkt@sds.co.kr>
 */

package openflow

import (
	"encoding/binary"
	"net"
	"strings"
)

type Port struct {
	Number uint16
	MAC    net.HardwareAddr
	Name   string
	// Bitmap of OFPPC_* flags
	config uint32
	// Bitmap of OFPPS_* flags
	state uint32
	//
	//  Bitmaps of OFPPF_* that describe features. All bits zeroed if unsupported or unavailable.
	//
	current    uint32
	advertised uint32
	supported  uint32
	peer       uint32
}

type PortFeature struct {
	OFPPF_10MB_HD    bool
	OFPPF_10MB_FD    bool
	OFPPF_100MB_HD   bool
	OFPPF_100MB_FD   bool
	OFPPF_1GB_HD     bool
	OFPPF_1GB_FD     bool
	OFPPF_10GB_FD    bool
	OFPPF_COPPER     bool
	OFPPF_FIBER      bool
	OFPPF_AUTONEG    bool
	OFPPF_PAUSE      bool
	OFPPF_PAUSE_ASYM bool
}

// Whether it is administratively down
func (r *Port) IsPortDown() bool {
	if r.config&OFPPC_PORT_DOWN != 0 {
		return true
	}

	return false
}

// Whether physical link is present
func (r *Port) IsLinkDown() bool {
	if r.state&OFPPS_LINK_DOWN != 0 {
		return true
	}

	return false
}

func getFeatures(v uint32) *PortFeature {
	return &PortFeature{
		OFPPF_10MB_HD:    v&OFPPF_10MB_HD != 0,
		OFPPF_10MB_FD:    v&OFPPF_10MB_FD != 0,
		OFPPF_100MB_HD:   v&OFPPF_100MB_HD != 0,
		OFPPF_100MB_FD:   v&OFPPF_100MB_FD != 0,
		OFPPF_1GB_HD:     v&OFPPF_1GB_HD != 0,
		OFPPF_1GB_FD:     v&OFPPF_1GB_FD != 0,
		OFPPF_10GB_FD:    v&OFPPF_10GB_FD != 0,
		OFPPF_COPPER:     v&OFPPF_COPPER != 0,
		OFPPF_FIBER:      v&OFPPF_FIBER != 0,
		OFPPF_AUTONEG:    v&OFPPF_AUTONEG != 0,
		OFPPF_PAUSE:      v&OFPPF_PAUSE != 0,
		OFPPF_PAUSE_ASYM: v&OFPPF_PAUSE_ASYM != 0,
	}
}

func (r *Port) GetCurrentFeatures() *PortFeature {
	return getFeatures(r.current)
}

func (r *Port) GetAdvertisedFeatures() *PortFeature {
	return getFeatures(r.advertised)
}

func (r *Port) GetSupportedFeatures() *PortFeature {
	return getFeatures(r.supported)
}

func (r *Port) UnmarshalBinary(data []byte) error {
	if len(data) < 48 {
		return ErrInvalidPacketLength
	}

	r.Number = binary.BigEndian.Uint16(data[0:2])
	r.MAC = make(net.HardwareAddr, 6)
	copy(r.MAC, data[2:8])
	r.Name = strings.TrimRight(string(data[8:24]), "\x00")
	r.config = binary.BigEndian.Uint32(data[24:28])
	r.state = binary.BigEndian.Uint32(data[28:32])
	r.current = binary.BigEndian.Uint32(data[32:36])
	r.advertised = binary.BigEndian.Uint32(data[36:40])
	r.supported = binary.BigEndian.Uint32(data[40:44])
	r.peer = binary.BigEndian.Uint32(data[44:48])

	return nil
}
