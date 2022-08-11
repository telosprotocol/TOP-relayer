package topsdk

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func newtopsdk() (*TopSdk, error) {
	url := "http://192.168.30.200:8080"
	return NewTopSdk(url)
}

func TestGetTopElectBlockHeadByHeight(t *testing.T) {
	url := "http://192.168.30.200:8080"
	sdk, err := NewTopSdk(url)
	if err != nil {
		t.Fatalf("NewSDK failed,error:%v", err)
	}
	var h uint64 = 1
	result, err := sdk.GetTopElectBlockHeadByHeight(h)
	if err != nil {
		t.Fatalf("GetTopElectBlockHeadByHeight failed,error:%v", err)
	}
	t.Logf("GetTopElectBlockHeadByHeight ok ,result:%v", result)
}

func TestGetLatestTopElectBlockHeight(t *testing.T) {
	sdk, err := newtopsdk()
	if err != nil {
		t.Fatalf("NewSDK failed,error:%v", err)
	}

	h, err := sdk.GetLatestTopElectBlockHeight()
	if err != nil {
		t.Fatalf("GetLatestTopElectBlockHeight failed,error:%v", err)
	}
	t.Logf("GetLatestTopElectBlockHeight:%v", h)
}

func TestParseJson(t *testing.T) {
	var data string = `{"blockRootHash":"0x0000000000000000000000000000000000000000000000000000000000000000","block_type":"election","hash":"0x4ce10ba22c80b977c40cb38e647adfac6c26a6e064f443417145d80f4611c21a","header":"0x00f902b4b86900f866808080a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000f9022480f90220f842a051559e17814b6c892f095bfaaeb11ff0cbf91a2d6ef45d77d61e24aba0d6d3fea05a9a53741dfa2b330e1f19874ec879bbd04a80f0e05e9a45bf205ea98531a135f842a03ba526113685c55c717ea2a62410ee57ab4e38a22ce3c7151548845a85a45184a0720425107d65c4c7364f19153a6c9cc0f91d683243c9e300a045b115192107baf842a0d4477911b0e186674c778d1b162db80b2f4ac95357b4c732cdc0f5124cfb5560a0eb78324385e62b68bc1d6f9fdc868ceca6af707ec93214cf1732ad5a9ec85d6ef842a044d3804bfa1c516dc53542419062a17739336a176f988f420a0eeaa95794f133a04cf2b80f6caecf8a407d0a9fc7000146f8b49dfd44634095cbd6771c97f6f824f842a05eded84b7ceb36373c8276061ff083170a4d3760c5f7eebda1c9dc06d1f7d6d0a0afd8c081434a0b21bda9b6036e053be1e5130d915ca0ac9c1967bb505e12705af842a0ffacf7a8a5bab86b68553b443d946bad2b0fc0a9803c428d6729ebbcaaec3e74a0d207f47a3e6f0468b0a99d86c2811601c5c51ea5254d05bc22a2161c293978a9f842a032f5faef900537f5e7fb2eb0ff091c0af471bac7af3f225d19cb7d9e3b8b3685a0b5cd591a10f6393b8210c50d3b47100b42e19ef67a42191a8a23ee3dbead0b9ff842a05ca103a04fda64db80eb26b10163434f6bf625e7f9fddf99c78e6fe3b56989aba0eaa0e8ecac03bf6e6238c455f2704e82e9efe84500b630ca89dbcc13d624b814c0","innerHeaderHash":"0x17a0a7c18b6b8391dfae586dda532ee7e47494f0efd3b60597636c4d99f8b668","number":"0x1","parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000","receiptsRootHash":"0x0000000000000000000000000000000000000000000000000000000000000000","size":"0x2e2","timestamp":"0x0","txsRootHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`

	var res TopBlock
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		log.Printf("sdk getLatestTopElectBlockHeight data: %v,error:%v", data, err)
		t.Fatalf("TestParseJson error %v", err)
	}
	fmt.Println("type:", res.BlockType, ", number:", res.Number)
}

func TestFastFindAggregate(t *testing.T) {
	url := "http://192.168.50.31:8080"
	sdk, err := NewTopSdk(url)
	if err != nil {
		t.Fatalf("NewSDK failed,error:%v", err)
	}
	for h := 1; ; h += 1 {
		result, err := sdk.GetTopElectBlockHeadByHeight(uint64(h))
		if err != nil {
			t.Fatalf("GetTopElectBlockHeadByHeight failed,error:%v", err)
		}
		if result.BlockType == "aggregate" {
			t.Logf("GetTopElectBlockHeadByHeight ok ,result:%v", result)
			break
		}
	}
}
