package main

import (
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/common"
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/99designs/httpsignatures-go"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	UserId         string `mapstructure:"user_id"`
	TenantId       string `mapstructure:"secret_key"`
	RawRegion      string `mapstructure:"region"`
	SkipValidation bool   `mapstructure:"skip_region_validation"`
	CompartmentId  string `mapstructure:"CompartmentId"`
	FingerPrint    string `mapstructure:"FingerPrint"`
	CommConfig     communicator.Config `mapstructure:",squash"`
}

func main() {
	user_id := ""
	tenancy_id := ""
	finger_print := ""
	key := fmt.Sprintf("%s/%s/%s", tenancy_id, user_id, finger_print)
	dat, err := ioutil.ReadFile("/Users/Sriram/.ssh/orc_auth_id_rsa")
	secret := string(dat)

	req, err := http.NewRequest("GET", "https://core.us-az-phoenix-1.oracleiaas.com/v1/20160918/images/", nil)
	req.Header.Set("host", "iaas.us-phoenix-1.oraclecloud.com")
	signer := httpsignatures.NewSigner(
		httpsignatures.AlgorithmHmacSha256,
		httpsignatures.RequestTarget, "date", "host",
	)
	signer.AuthRequest(key, secret, req)

	fmt.Println(req.Header)
	fmt.Println()
	fmt.Println(req.Header.Get("Authorization"))

	client := &http.Client{}

	resp, err := client.Do(req)
	output, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(output[:]))
	fmt.Println(err)
}

