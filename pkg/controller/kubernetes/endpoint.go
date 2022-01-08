/*
 * Copyright (c) 2019 Huawei Technologies Co., Ltd.
 * MeshAccelerating is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: LemmyHuang
 * Create: 2021-12-22
 */

package kubernetes

import (
	"fmt"
	apiCoreV1 "k8s.io/api/core/v1"
	"openeuler.io/mesh/pkg/bpf/maps"
	"openeuler.io/mesh/pkg/nets"
)

type endpointKeyAndValue struct {
	key		maps.MapKey
	value	maps.Endpoint
}
type endpointData map[endpointKeyAndValue]objOptionFlag

func (data endpointData) deleteInvalid(kv *endpointKeyAndValue) {
	if data[*kv] == serviceOptionAllFlag {
		delete(data, *kv)
	}
}

func (data endpointData) extractData(flag objOptionFlag, ep *apiCoreV1.Endpoints, nameID uint32) {
	var kv endpointKeyAndValue

	if ep == nil {
		return
	}
	kv.key.NameID = nameID

	for _, sub := range ep.Subsets {
		for _, epPort := range sub.Ports {
			if !nets.GetConfig().IsEnabledProtocol(string(epPort.Protocol)) {
				continue
			}

			kv.value.Address.Protocol = protocolStrToC[epPort.Protocol]
			kv.value.Address.Port = nets.ConvertPortToLittleEndian(epPort.Port)
			kv.key.Port = kv.value.Address.Port

			for _, epAddr := range sub.Addresses {
				kv.value.Address.IPv4 = nets.ConvertIpToUint32(epAddr.IP)
				data[kv] |= flag
				data.deleteInvalid(&kv)
			}
		}
	}
}

func (data endpointData) flushMap(flag objOptionFlag, count objCount, addressToMapKey objAddressToMapKey) int {
	var err error
	var num int

	for kv, f := range data {
		if f != flag {
			continue
		}

		switch flag {
		case serviceOptionDeleteFlag:
			err = kv.deleteMap(count, addressToMapKey)
		case serviceOptionUpdateFlag:
			err = kv.updateMap(count, addressToMapKey)
		default:
			// ignore
		}
		num++

		if err != nil {
			log.Errorln(err)
		}
	}

	return num
}

func (kv *endpointKeyAndValue) updateMap(count objCount, addressToMapKey objAddressToMapKey) error {
	kv.key.Index = count[kv.key.Port]

	if err := kv.value.Update(&kv.key); err != nil {
		return fmt.Errorf("update endpoint failed, %v, %s", kv.key, err)
	}

	// update count
	count[kv.key.Port]++
	addressToMapKey[kv.value.Address] = kv.key

	lb := &maps.Loadbalance{}
	lb.MapKey = kv.key
	if err := lb.Update(&kv.key); err != nil {
		kv.deleteMap(count, addressToMapKey)
		return fmt.Errorf("update loadbalance failed, %v, %s", kv.key, err)
	}

	return nil
}

func (kv *endpointKeyAndValue) deleteMap(count objCount, addressToMapKey objAddressToMapKey) error {
	lb := &maps.Loadbalance{}
	mapKey := addressToMapKey[kv.value.Address]

	kv.key.Index = mapKey.Index
	if kv.key != mapKey {
		return fmt.Errorf("delete endpoint using invalid key, %v != %v", kv.key, mapKey)
	}

	mapKeyTail := mapKey
	mapKeyTail.Index = count[mapKey.Port] - 1

	if mapKey != mapKeyTail {
		if err := kv.value.Lookup(&mapKeyTail); err == nil {
			kv.value.Update(&mapKey)
		}
		if err := lb.Lookup(&mapKeyTail); err == nil {
			lb.Update(&mapKey)
		}
	}
	kv.value.Delete(&mapKeyTail)
	lb.Delete(&mapKeyTail)

	// update count
	delete(addressToMapKey, kv.value.Address)
	count[kv.key.Port]--
	if count[kv.key.Port] <= 0 {
		delete(count, kv.key.Port)
	}

	return nil
}
