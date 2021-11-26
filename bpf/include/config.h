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
 * Create: 2021-09-17
 */

#ifndef _CONFIG_H_
#define _CONFIG_H_

// ************
// options
#define KMESH_MODULE_ON			1
#define KMESH_MODULE_OFF		0

// L3
#define KMESH_ENABLE_IPV4		KMESH_MODULE_ON
#define KMESH_ENABLE_IPV6		KMESH_MODULE_OFF
// L4
#define KMESH_ENABLE_TCP		KMESH_MODULE_ON
#define KMESH_ENABLE_UDP		KMESH_MODULE_OFF
// L7
#define KMESH_ENABLE_HTTP		KMESH_MODULE_OFF

// ************
// map size
#define MAP_SIZE_OF_PER_LISTENER		64
#define MAP_SIZE_OF_PER_FILTER_CHAIN	4
#define MAP_SIZE_OF_PER_FILTER			4
#define MAP_SIZE_OF_PER_VIRTUAL_HOST	4
#define MAP_SIZE_OF_PER_ROUTE			8
#define MAP_SIZE_OF_PER_CLUSTER			32
#define MAP_SIZE_OF_PER_ENDPOINT		128

#define MAP_SIZE_OF_MAX					8192

#define MAP_SIZE_OF_LISTENER		\
	BPF_MIN(MAP_SIZE_OF_MAX, MAP_SIZE_OF_PER_LISTENER)
#define MAP_SIZE_OF_FILTER_CHAIN	\
	BPF_MIN(MAP_SIZE_OF_MAX, MAP_SIZE_OF_PER_FILTER_CHAIN * MAP_SIZE_OF_LISTENER)
#define MAP_SIZE_OF_FILTER			\
	BPF_MIN(MAP_SIZE_OF_MAX, MAP_SIZE_OF_PER_FILTER * MAP_SIZE_OF_FILTER_CHAIN)
#define MAP_SIZE_OF_VIRTUAL_HOST	\
	BPF_MIN(MAP_SIZE_OF_MAX, MAP_SIZE_OF_PER_VIRTUAL_HOST * MAP_SIZE_OF_FILTER)
#define MAP_SIZE_OF_ROUTE			\
	BPF_MIN(MAP_SIZE_OF_MAX, MAP_SIZE_OF_PER_ROUTE * MAP_SIZE_OF_VIRTUAL_HOST)
#define MAP_SIZE_OF_CLUSTER			\
	BPF_MIN(MAP_SIZE_OF_MAX, MAP_SIZE_OF_PER_CLUSTER * MAP_SIZE_OF_ROUTE)
#define MAP_SIZE_OF_ENDPOINT		\
	BPF_MIN(MAP_SIZE_OF_MAX, MAP_SIZE_OF_PER_ENDPOINT * MAP_SIZE_OF_CLUSTER)

// rename map to void truncation when name length exceeds BPF_OBJ_NAME_LEN = 16
#define map_of_listener			listener
#define map_of_filter_chain		filter_chain
#define map_of_filter			filter
#define map_of_virtual_host		virtual_host
#define map_of_route			route
#define map_of_cluster			cluster
#define map_of_loadbanance		loadbanance
#define map_of_endpoint			endpoint
#define map_of_tail_call_prog	tail_call_prog
#define map_of_tail_call_ctx	tail_call_ctx


// ************
// array len
#define KMESH_NAME_LEN				64
#define KMESH_TYPE_LEN				64
#define KMESH_HOST_LEN				128
#define KMESH_FILTER_CHAINS_LEN		64
#define KMESH_HTTP_DOMAIN_NUM		64
#define KMESH_HTTP_DOMAIN_LEN		128

#endif //_CONFIG_H_
