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
// IMPORTANT!! - this is the main message interpreted for new peers being processed
type GFeventNewPeerLifecycle struct {
	PeerEnodeID   string `parquet:"name=peer_enode_id,  type=UTF8" csv:"peer_enode_id"  json:"peer_enode_id"`
	Name          string `parquet:"name=name,           type=UTF8" csv:"name"           json:"name"` 
	RemoteAddress string `parquet:"name=remote_address, type=UTF8" csv:"remote_address" json:"remote_address"` 
	LocalAddress  string `parquet:"name=local_address,  type=UTF8" csv:"local_address"  json:"local_address"`
}

// GF_EVENT__DROPPING_UNSYNCED_NODE_DURING_FAST_SYNC
type GFeventDroppingUnsyncedNodeDuringFastSync struct {
	PeerEnodeID   string `parquet:"name=peer_enode_id,  type=UTF8" json:"peer_enode_id"`
	Name          string `parquet:"name=name,           type=UTF8" json:"name"` 
	RemoteAddress string `parquet:"name=remote_address, type=UTF8" json:"remote_address"` 
	LocalAddress  string `parquet:"name=local_address,  type=UTF8" json:"local_address"`
}

//-----------------------------------------------------------------
// DOWNLOADER

// GF_EVENT__NEW_PEER_REGISTER
type GFeventNewPeerRegister struct {
	PeerID        string        `parquet:"name=peer_id, type=UTF8" json:"peer_id"`
	RoundTripTime time.Duration // FIX!! - serialize to float64
}

// GF_EVENT__NEW_HEADER_FROM_PEER
type GFeventNewHeaderFromPeer struct {
	PeerID           string `parquet:"name=peer_id,           type=UTF8"  json:"peer_id"`
	HeaderNumber     uint64 `parquet:"name=header_number,     type=INT64" json:"header_number"`
	HeaderTime       uint64 `parquet:"name=header_time,       type=INT64" json:"header_time"`
	HeaderDifficulty uint64 `parquet:"name=header_difficulty, type=INT64" json:"header_difficulty"` 
	GasUsed          uint64 `parquet:"name=gas_used,          type=INT64" json:"gas_used"`
	RootHashHex      string `parquet:"name=root_hash_hex,     type=UTF8"  json:"root_hash_hex"`
	ParentHashHex    string `parquet:"name=parent_hash_hex,   type=UTF8"  json:"parent_hash_hex"`
	UncleHashHex     string `parquet:"name=uncle_hash_hex,    type=UTF8"  json:"uncle_hash_hex"`
} 

// GF_EVENT__BLOCK_SYNCHRONISE_WITH_PEER
type GFeventBlockSynchroniseWithPeer struct {
	PeerID string `parquet:"name=peer_id, type=UTF8" json:"peer_id"`
}

// GF_EVENT__DROPPING_PEER_SYNC_FAILED
type GFeventDroppingPeerSyncFailed struct {
	PeerID string `parquet:"name=peer_id, type=UTF8"`
}