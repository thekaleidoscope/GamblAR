package atomic_elements

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"

	"github.com/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/fabric-sdk-go/pkg/core/config"
	"github.com/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

type Handeler struct {
	//Channel Info
	ChannelID string

	//OrgInfo
	OrgAdmins     map[string]string
	Organizations []string
	initialized   bool

	//Instance Reference
	client *channel.Client
	admin  *resmgmt.Client
	sdk    *fabsdk.FabricSDK
	event  *event.Client

	//Basic Config
	ConfigFile string
}

//Initializer creates the sdk context from config file and instantiate a sdk instance
func (handel *Handeler) Initializer() error {

	if handel.initialized {
		return fmt.Errorf("SDK already initialized.")
	}

	sdk, err := fabsdk.New(config.FromFile(handel.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed create SDK instance.")
	}

	handel.sdk = sdk

	return nil
}

func (handel *Handeler) CreateAndJoinChannel() error {

	ordererClientContext := handel.sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Orderer"))

	chMgmtClient, err := resmgmt.New(ordererClientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel management client")
	} else {
		fmt.Println("Created new channel management client")
	}
	signId := []msp.SigningIdentity{}
	for _, orgName := range handel.Organizations {

		mspContext, err := mspclient.New(handel.sdk.Context(), mspclient.WithOrg(strings.ToLower(orgName)))
		if err != nil {
			return errors.WithMessage(err, "failed to create MSPContext:")
		}

		signId_t, err := mspContext.GetSigningIdentity(handel.OrgAdmins[orgName])
		if err != nil {
			return errors.WithMessage(err, "failed to GetSigningIdentity:")
		}
		signId = append(signId, signId_t)
	}
	req := resmgmt.SaveChannelRequest{
		ChannelID:         handel.ChannelID,
		ChannelConfigPath: ChannelConfigPath("townchannel.tx"),
		SigningIdentities: signId}

	txID, err := chMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint("orderer.gamblar.com"))
	if err != nil {
		return errors.WithMessage(err, "failed to get SaveChannel.")
	} else {
		fmt.Println("Created Channel", handel.ChannelID, "\n TxID: ", txID)
	}

	for _, orgName := range handel.Organizations {

		resourceContext := handel.sdk.Context(fabsdk.WithUser(handel.OrgAdmins[orgName]), fabsdk.WithOrg(orgName))
		resMgmtClient, err := resmgmt.New(resourceContext)
		if err != nil {
			return errors.WithMessage(err, "failed to get resMgmtClient.")
		} else {

		}

		if err := resMgmtClient.JoinChannel(handel.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.gamblar.com")); err != nil {
			return errors.WithMessage(err, "failed to Join Channel.")
		} else {
			fmt.Println(orgName, "Joined Channel", handel.ChannelID)

		}

	}
	return nil
}

func ChannelConfigPath(filename string) string {
	return path.Join(os.Getenv("GOPATH"), "src", "github.com", "GamblAR", "units", "collection", filename)
}
