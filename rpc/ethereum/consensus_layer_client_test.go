package ethereum
import (
"os"
"testing"
)

func TestPrismV421LightClientUpdateJsonUnmarshal(t *testing.T) {
	testDataFilePath := "./test_data/sepolia/prism-v4.2.1-eth-v1-beacon-lightclient-update-period-518.json"
	content, err := os.ReadFile(testDataFilePath)
	if err != nil {
		t.Fatal(err)
	}
	jsonString := string(content)

}
