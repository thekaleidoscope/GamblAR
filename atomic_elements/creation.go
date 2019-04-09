package atomic_elements

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

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

	//Basic Config
	ConfigFile string

	//Chaincode
	ChainCodeID     string
	ChaincodeGoPath string
	ChaincodePath   string
}

type Bid struct {
	GameName  string
	Acc       string
	BidAmount string
}

type Game map[string]Bid

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
			if orgName == "Org1" {
				handel.admin = resMgmtClient
			}
		}

		if err := resMgmtClient.JoinChannel(handel.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.gamblar.com")); err != nil {
			return errors.WithMessage(err, "failed to Join Channel.")
		} else {
			fmt.Println(orgName, "Joined Channel", handel.ChannelID)

		}
		fmt.Println("-------Successful  Creating And Join Channel For each Organizations--------\n")

		fmt.Println("------- Chaincode Installation and Instantiation --------")

		ccPkg, err := packager.NewCCPackage(handel.ChaincodePath, handel.ChaincodeGoPath)
		if err != nil {
			return errors.WithMessage(err, "failed to create chaincode package by "+orgName)
		}
		fmt.Println("ccPkg created by ", orgName)

		// Install example cc to org peers
		installCCReq := resmgmt.InstallCCRequest{Name: handel.ChainCodeID, Path: handel.ChaincodePath, Version: "0", Package: ccPkg}
		_, err = resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return errors.WithMessage(err, "failed to install chaincode by "+orgName)
		}
		fmt.Println("Chaincode installed by ", orgName)
	}
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.gamblar.com", "org2.gamblar.com"})

	resp, err := handel.admin.InstantiateCC(handel.ChannelID, resmgmt.InstantiateCCRequest{Name: handel.ChainCodeID, Path: handel.ChaincodeGoPath, Version: "0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {

		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}

	fmt.Println("Chaincode instantiated")

	return nil
}

func ChannelConfigPath(filename string) string {
	return path.Join(os.Getenv("GOPATH"), "src", "github.com", "GamblAR", "units", "collection", filename)
}
