package factory

import (
	"errors"
	"testing"

	arwenConfig "github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/vm"
	"github.com/ElrondNetwork/elrond-go/vm/mock"
	"github.com/ElrondNetwork/elrond-go/vm/systemSmartContracts/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockNewSystemScFactoryArgs() ArgsNewSystemSCFactory {
	gasMap := arwenConfig.MakeGasMapForTests()
	gasMap = defaults.FillGasMapInternal(gasMap, 1)
	gasSchedule := mock.NewGasScheduleNotifierMock(gasMap)
	return ArgsNewSystemSCFactory{
		SystemEI:            &mock.SystemEIStub{},
		Economics:           &mock.EconomicsHandlerStub{},
		SigVerifier:         &mock.MessageSignVerifierMock{},
		GasSchedule:         gasSchedule,
		NodesConfigProvider: &mock.NodesConfigProviderStub{},
		Marshalizer:         &mock.MarshalizerMock{},
		Hasher:              &mock.HasherMock{},
		SystemSCConfig: &config.SystemSmartContractsConfig{
			ESDTSystemSCConfig: config.ESDTSystemSCConfig{
				BaseIssuingCost: "100000000",
				OwnerAddress:    "aaaaaa",
			},
			GovernanceSystemSCConfig: config.GovernanceSystemSCConfig{
				ProposalCost:     "500",
				NumNodes:         100,
				MinQuorum:        50,
				MinPassThreshold: 50,
				MinVetoThreshold: 50,
			},
			StakingSystemSCConfig: config.StakingSystemSCConfig{
				GenesisNodePrice:                     "1000",
				UnJailValue:                          "10",
				MinStepValue:                         "10",
				MinStakeValue:                        "1",
				UnBondPeriod:                         1,
				StakingV2Epoch:                       1,
				StakeEnableEpoch:                     0,
				NumRoundsWithoutBleed:                1,
				MaximumPercentageToBleed:             1,
				BleedPercentagePerRound:              1,
				MaxNumberOfNodesForStake:             100,
				ActivateBLSPubKeyMessageVerification: false,
				MinUnstakeTokensValue:                "1",
			},
			DelegationSystemSCConfig: config.DelegationSystemSCConfig{
				EnabledEpoch:  0,
				MinServiceFee: 0,
				MaxServiceFee: 10000,
			},
			DelegationManagerSystemSCConfig: config.DelegationManagerSystemSCConfig{
				MinCreationDeposit:  "10",
				EnabledEpoch:        0,
				MinStakeAmount:      "10",
				ConfigChangeAddress: "aabb00",
			},
		},
		EpochNotifier:          &mock.EpochNotifierStub{},
		AddressPubKeyConverter: &mock.PubkeyConverterMock{},
	}
}

func TestNewSystemSCFactory_NilSystemEI(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.SystemEI = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilSystemEnvironmentInterface, err)
}

func TestNewSystemSCFactory_NilSigVerifier(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.SigVerifier = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilMessageSignVerifier, err)
}

func TestNewSystemSCFactory_NilNodesConfigProvider(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.NodesConfigProvider = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilNodesConfigProvider, err)
}

func TestNewSystemSCFactory_NilMarshalizer(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.Marshalizer = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilMarshalizer, err)
}

func TestNewSystemSCFactory_NilHasher(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.Hasher = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilHasher, err)
}

func TestNewSystemSCFactory_NilEconomicsData(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.Economics = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilEconomicsData, err)
}

func TestNewSystemSCFactory_NilSystemScConfig(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.SystemSCConfig = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilSystemSCConfig, err)
}

func TestNewSystemSCFactory_NilEpochNotifier(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.EpochNotifier = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilEpochNotifier, err)
}

func TestNewSystemSCFactory_NilPubKeyConverter(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.AddressPubKeyConverter = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilAddressPubKeyConverter, err)
}

func TestNewSystemSCFactory_Ok(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, err)
	assert.NotNil(t, scFactory)
}

func TestSystemSCFactory_CreateWithBadDelegationManagerConfigChangeAddressShouldError(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.SystemSCConfig.DelegationManagerSystemSCConfig.ConfigChangeAddress = "not a hex string"
	scFactory, _ := NewSystemSCFactory(arguments)

	container, err := scFactory.Create()

	assert.True(t, check.IfNil(container))
	assert.True(t, errors.Is(err, vm.ErrInvalidAddress))
}

func TestSystemSCFactory_Create(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, _ := NewSystemSCFactory(arguments)

	container, err := scFactory.Create()
	assert.Nil(t, err)
	require.NotNil(t, container)
	assert.Equal(t, 6, container.Len())
}

func TestSystemSCFactory_CreateForGenesis(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, _ := NewSystemSCFactory(arguments)

	container, err := scFactory.CreateForGenesis()
	assert.Nil(t, err)
	assert.Equal(t, 4, container.Len())
}

func TestSystemSCFactory_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, _ := NewSystemSCFactory(arguments)
	assert.False(t, scFactory.IsInterfaceNil())

	scFactory = nil
	assert.True(t, check.IfNil(scFactory))
}
