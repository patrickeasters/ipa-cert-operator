package settings

import (
	"io/ioutil"
	"os"
	"time"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type Settings struct {
	IpaHost         string
	IpaRealm        string
	IpaUser         string
	IpaPassword     string
	CertProfileHost string
	CertProfileUser string
	CaChain         string
	RequeuePeriod   time.Duration
}

var Instance Settings
var log = logf.Log.WithName("settings")

func ParseSettings() {
	var ok bool

	Instance.IpaHost = os.Getenv("IPA_HOST")
	Instance.IpaRealm = os.Getenv("IPA_REALM")
	Instance.IpaUser = os.Getenv("IPA_USER")
	Instance.IpaPassword = os.Getenv("IPA_PASSWORD")

	Instance.CertProfileHost, ok = os.LookupEnv("CERT_PROFILE_HOST")
	if !ok {
		Instance.CertProfileHost = "caIPAserviceCert"
	}

	Instance.CertProfileUser, ok = os.LookupEnv("CERT_PROFILE_USER")
	if !ok {
		Instance.CertProfileUser = "IECUserRoles"
	}

	Instance.RequeuePeriod, _ = time.ParseDuration(os.Getenv("REQUEUE_PERIOD"))
	if Instance.RequeuePeriod.Seconds() < 30 {
		Instance.RequeuePeriod = 6 * time.Hour
	}

	chain_file, ok := os.LookupEnv("CA_CHAIN_FILE")
	if ok {
		chain, err := ioutil.ReadFile(chain_file)
		if err != nil {
			Instance.CaChain = string(chain)
		} else {
			log.Error(err, "Unable to read CA_CHAIN_FILE. Ignoring option.")
		}
	}

}
