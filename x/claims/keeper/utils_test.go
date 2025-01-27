package keeper_test

import (
	"github.com/twobitEDD/servermint/v12/constants"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/twobitEDD/servermint/v12/app"
	"github.com/twobitEDD/servermint/v12/crypto/ethsecp256k1"
	"github.com/twobitEDD/servermint/v12/testutil"
	utiltx "github.com/twobitEDD/servermint/v12/testutil/tx"
	evertypes "github.com/twobitEDD/servermint/v12/types"
	"github.com/twobitEDD/servermint/v12/x/claims/types"
	evm "github.com/twobitEDD/servermint/v12/x/evm/types"
	feemarkettypes "github.com/twobitEDD/servermint/v12/x/feemarket/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = utiltx.NewSigner(priv)

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	chainID := constants.TestnetFullChainId
	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState(), chainID)
	header := testutil.NewHeader(
		1, time.Now().UTC(), chainID, consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.ClaimsKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)

	params := types.DefaultParams()
	params.AirdropStartTime = suite.ctx.BlockTime().UTC()
	err = suite.app.ClaimsKeeper.SetParams(suite.ctx, params)
	require.NoError(t, err)

	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = params.GetClaimsDenom()
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	// Set Validator
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	err = suite.app.StakingKeeper.Hooks().AfterValidatorCreated(suite.ctx, validator.GetOperator())
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)
	validators := s.app.StakingKeeper.GetValidators(s.ctx, 1)
	suite.validator = validators[0]

	suite.ethSigner = ethtypes.LatestSignerForChainID(s.app.EvmKeeper.ChainID())
}

func (suite *KeeperTestSuite) SetupTestWithEscrow() {
	suite.SetupTest()
	params := suite.app.ClaimsKeeper.GetParams(suite.ctx)

	coins := sdk.NewCoins(sdk.NewCoin(params.ClaimsDenom, sdk.NewInt(10000000)))
	err := testutil.FundModuleAccount(suite.ctx, suite.app.BankKeeper, types.ModuleName, coins)
	suite.Require().NoError(err)
}

// Commit commits and starts a new block with an updated context.
func (suite *KeeperTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

// Commit commits a block at a given time.
func (suite *KeeperTestSuite) CommitAfter(t time.Duration) {
	var err error
	suite.ctx, err = testutil.Commit(suite.ctx, suite.app, t, nil)
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())

	types.RegisterQueryServer(queryHelper, suite.app.ClaimsKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)
}

func newEthAccount(baseAccount *authtypes.BaseAccount) evertypes.EthAccount {
	return evertypes.EthAccount{
		BaseAccount: baseAccount,
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}
}

func getAddr(priv *ethsecp256k1.PrivKey) sdk.AccAddress {
	return sdk.AccAddress(priv.PubKey().Address().Bytes())
}

func govProposal(_ *ethsecp256k1.PrivKey) (uint64, error) {
	panic("broken test on deprecated module")
}

func sendEthToSelf(priv *ethsecp256k1.PrivKey) {
	chainID := s.app.EvmKeeper.ChainID()
	from := common.BytesToAddress(priv.PubKey().Address().Bytes())
	nonce := s.app.EvmKeeper.GetNonce(s.ctx, from)

	ethTxParams := evm.EvmTxArgs{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &from,
		GasLimit:  100000,
		GasFeeCap: s.app.FeeMarketKeeper.GetBaseFee(s.ctx),
		GasTipCap: big.NewInt(1),
		Accesses:  &ethtypes.AccessList{},
	}
	msgEthereumTx := evm.NewTx(&ethTxParams)
	msgEthereumTx.From = from.String()
	_, err := testutil.DeliverEthTx(s.app, priv, msgEthereumTx)
	s.Require().NoError(err)
}

func getEthTxFee() sdk.Coin {
	baseFee := s.app.FeeMarketKeeper.GetBaseFee(s.ctx)
	baseFee.Mul(baseFee, big.NewInt(100000))
	feeAmt := baseFee.Quo(baseFee, big.NewInt(2))
	return sdk.NewCoin(constants.BaseDenom, sdkmath.NewIntFromBigInt(feeAmt))
}
