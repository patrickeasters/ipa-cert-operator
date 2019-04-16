package settings

import (
	"io/ioutil"
	"os"
	"strconv"
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
	IpaTimeout      time.Duration
	RequeuePeriod   time.Duration
	RenewalPeriod   time.Duration
	HostAutoCreate  bool
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

	ipaTimeout, ok := os.LookupEnv("IPA_TIMEOUT")
	Instance.IpaTimeout, _ = time.ParseDuration(ipaTimeout)
	if !ok {
		Instance.IpaTimeout = 30 * time.Second
	}

	Instance.RequeuePeriod, _ = time.ParseDuration(os.Getenv("REQUEUE_PERIOD"))
	if Instance.RequeuePeriod.Minutes() < 1 {
		Instance.RequeuePeriod = 6 * time.Hour
	}

	Instance.RenewalPeriod, _ = time.ParseDuration(os.Getenv("RENEWAL_PERIOD"))
	if Instance.RenewalPeriod.Minutes() < 1 {
		Instance.RenewalPeriod = 30 * 24 * time.Hour
	}

	Instance.HostAutoCreate, _ = strconv.ParseBool(os.Getenv("HOST_AUTO_CREATE"))

	chain_file, ok := os.LookupEnv("CA_CHAIN_FILE")
	if ok {
		chain, err := ioutil.ReadFile(chain_file)
		if err == nil {
			Instance.CaChain = string(chain)
		} else {
			log.Error(err, "Unable to read CA_CHAIN_FILE. Ignoring option.", "Read Error", err)
		}
	}

}
