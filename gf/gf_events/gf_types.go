// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package gf_events

import (
	"time"
)

//-----------------------------------------------------------------
// PROTOCOL_MANAGER

// GF_EVENT__NEW_PEER_LIFECYCLE
type GFeventNewPeerLifecycle struct {
	//--------------------------
	// GF_EVENT
	Id      string  `parquet:"name=gf_id,    type=UTF8" csv:"gf_id"`
	TimeSec float64 `parquet:"name=time_sec, type=UTF8" csv:"time_sec"`
	Module  string  `parquet:"name=module,   type=UTF8" csv:"module"` // protocol_manager
	Type    string  `parquet:"name=type,     type=UTF8" csv:"type"`   // handle_new_peer
	
	//--------------------------
	PeerEnodeID   string `parquet:"name=peer_enode_id,  type=UTF8" csv:"peer_enode_id"`
	Name          string `parquet:"name=name,           type=UTF8" csv:"name"` 
	RemoteAddress string `parquet:"name=remote_address, type=UTF8" csv:"remote_address"` 
	LocalAddress  string `parquet:"name=local_address,  type=UTF8" csv:"local_address"`
}

// GF_EVENT__DROPPING_UNSYNCED_NODE_DURING_FAST_SYNC
type GFeventDroppingUnsyncedNodeDuringFastSync struct {
	//--------------------------
	//GF_EVENT
	Id      string  `parquet:"name=gf_id,    type=UTF8"`
	TimeSec float64 `parquet:"name=time_sec, type=UTF8"`
	Module  string  `parquet:"name=module,   type=UTF8"` // protocol_manager
	Type    string  `parquet:"name=type,     type=UTF8"` // dropping_unsynced_node_during_fast_sync

	//--------------------------
	PeerEnodeID   string `parquet:"name=peer_enode_id,  type=UTF8"`
	Name          string `parquet:"name=name,           type=UTF8"` 
	RemoteAddress string `parquet:"name=remote_address, type=UTF8"` 
	LocalAddress  string `parquet:"name=local_address,  type=UTF8"`
}

//-----------------------------------------------------------------
// DOWNLOADER

// GF_EVENT__NEW_PEER_REGISTER
type GFeventNewPeerRegister struct {
	//--------------------------
	//GF_EVENT
	Id      string  `parquet:"name=gf_id,    type=UTF8"`
	TimeSec float64 `parquet:"name=time_sec, type=UTF8"`
	Module  string  `parquet:"name=module,   type=UTF8"` // downloader
	Type    string  `parquet:"name=type,     type=UTF8"` // register_peer

	//--------------------------
	PeerID        string        `parquet:"name=peer_id, type=UTF8"`
	RoundTripTime time.Duration // FIX!! - serialize to float64
}

// GF_EVENT__NEW_HEADER_FROM_PEER
type GFeventNewHeaderFromPeer struct {
	//--------------------------
	// GF_EVENT
	Id      string  `parquet:"name=gf_id,    type=UTF8"`
	TimeSec float64 `parquet:"name=time_sec, type=UTF8"`
	Module  string  `parquet:"name=module,   type=UTF8"` // downloader
	Type    string  `parquet:"name=type,     type=UTF8"` // new_header_from_peer
	
	//--------------------------
	PeerID           string `parquet:"name=peer_id,           type=UTF8"`
	HeaderNumber     uint64 `parquet:"name=header_number,     type=INT64"`
	HeaderTime       uint64 `parquet:"name=header_time,       type=INT64"`
	HeaderDifficulty uint64 `parquet:"name=header_difficulty, type=INT64"` 
	GasUsed          uint64 `parquet:"name=gas_used,          type=INT64"`
	RootHashHex      string `parquet:"name=root_hash_hex,     type=UTF8"`
	ParentHashHex    string `parquet:"name=parent_hash_hex,   type=UTF8"`
	UncleHashHex     string `parquet:"name=uncle_hash_hex,    type=UTF8"`
}

// GF_EVENT__BLOCK_SYNCHRONISE_WITH_PEER
type GFeventBlockSynchroniseWithPeer struct {
	//--------------------------
	// GF_EVENT
	Id      string  `parquet:"name=gf_id,    type=UTF8"`
	TimeSec float64 `parquet:"name=time_sec, type=UTF8"`
	Module  string  `parquet:"name=module,   type=UTF8"` // downloader
	Type    string  `parquet:"name=type,     type=UTF8"` // block_synchronise_with_peer
	
	//--------------------------
	PeerID string `parquet:"name=peer_id, type=UTF8"`

}

// GF_EVENT__DROPPING_PEER_SYNC_FAILED
type GFeventDroppingPeerSyncFailed struct {
	//--------------------------
	// GF_EVENT
	Id      string  `parquet:"name=gf_id,    type=UTF8"`
	TimeSec float64 `parquet:"name=time_sec, type=UTF8"`
	Module  string  `parquet:"name=module,   type=UTF8"`
	Type    string  `parquet:"name=type,     type=UTF8"`
	
	//--------------------------
	PeerID string `parquet:"name=peer_id, type=UTF8"`
}