1. The logic of implementing ETH to TOP relayer is mainly from `https://github.com/aurora-is-near/rainbow-bridge/blob/master/eth2near/eth2near-block-relay-rs/src/eth2near_relay.rs`
2. For supporting the Ethereum hard fork, following logic may be impacted:
   1. `ethereum.ExecutionPayloadMerkleTreeNew`: the update logic can be borrowed from `https://github.com/prysmaticlabs/prysm/blob/v4.2.1/proto/engine/v1/generated.ssz.go#L1165`.
   2. `ethereum.BeaconBlockBodyMerkleTreeNew`: the update logic can be borrowed from `https://github.com/prysmaticlabs/prysm/blob/v4.2.1/proto/prysm/v1alpha1/generated.ssz.go#L5778`.
   3. `ethereum.BeaconClient.constructFromBeaconBlockBody`.
   4. Update the dependencies of the project.

known issues:
1. For Prysm v4.2.1 consensus node, the return JSON from /eth/v1/beacon/light_client/updates API is not conformed to the schema defined by the standard. The standard is defined at `https://ethereum.github.io/beacon-APIs/#/Beacon/getLightClientUpdatesByRange`. The sample output of v4.2.1 can be found [here](test_data/sepolia/prism-v4.2.1-eth-v1-beacon-lightclient-update-period-518.json).
2. In `ethereum.BeaconClient.GetPrysmLightClientUpdate` API, the HTTP request is sent by http.Client instance. For the Prysm v4.2.1 beacon.Client cannot be used because it cannot handle `http://beacon-node-test-server.com:port/eth/v1/beacon/light_client/updates?start_period=518&count=1` correctly. The request URI will be encoded as `http://beacon-node-test-server.com:port/eth/v1/beacon/light_client/updates%3Fstart_period=518&count=1`, and got 404 error.