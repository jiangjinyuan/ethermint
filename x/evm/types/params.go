package types

import (
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/ethereum/go-ethereum/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/types"
)

var (
	// DefaultEVMDenom defines the default EVM denomination on Ethermint
	DefaultEVMDenom = types.AttoPhoton
	// DefaultMinGasMultiplier is 0.5 or 50%
	DefaultMinGasMultiplier = sdk.NewDecWithPrec(50, 2)
	// DefaultAllowUnprotectedTxs rejects all unprotected txs (i.e false)
	DefaultAllowUnprotectedTxs = false
	// MsgEthTxAmino Amino name
	MsgEthTxAmino = "ethermint/MsgEthereumTx"
)

// Parameter keys
var (
	ParamStoreKeyEVMDenom            = []byte("EVMDenom")
	ParamStoreKeyEnableCreate        = []byte("EnableCreate")
	ParamStoreKeyEnableCall          = []byte("EnableCall")
	ParamStoreKeyExtraEIPs           = []byte("EnableExtraEIPs")
	ParamStoreKeyChainConfig         = []byte("ChainConfig")
	ParamStoreKeyAllowUnprotectedTxs = []byte("AllowUnprotectedTxs")

	// AvailableExtraEIPs define the list of all EIPs that can be enabled by the
	// EVM interpreter. These EIPs are applied in order and can override the
	// instruction sets from the latest hard fork enabled by the ChainConfig. For
	// more info check:
	// https://github.com/ethereum/go-ethereum/blob/master/core/vm/interpreter.go#L97
	AvailableExtraEIPs = []int64{1344, 1884, 2200, 2929, 3198, 3529}
)

// NewParams creates a new Params instance
func NewParams(evmDenom string, allowUnprotectedTxs, enableCreate, enableCall bool, config ChainConfig, extraEIPs ExtraEIPs) Params {
	return Params{
		EvmDenom:            evmDenom,
		AllowUnprotectedTxs: allowUnprotectedTxs,
		EnableCreate:        enableCreate,
		EnableCall:          enableCall,
		ExtraEips:           extraEIPs,
		ChainConfig:         config,
	}
}

// DefaultParams returns default evm parameters
// ExtraEIPs is empty to prevent overriding the latest hard fork instruction set
func DefaultParams() Params {
	return Params{
		EvmDenom:            DefaultEVMDenom,
		EnableCreate:        true,
		EnableCall:          true,
		ChainConfig:         DefaultChainConfig(),
		ExtraEips:           ExtraEIPs{},
		AllowUnprotectedTxs: DefaultAllowUnprotectedTxs,
	}
}

// Validate performs basic validation on evm parameters.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.EvmDenom); err != nil {
		return err
	}

	if err := validateEIPs(p.ExtraEips.ExtraEIPs); err != nil {
		return err
	}

	return p.ChainConfig.Validate()
}

// EIPs returns the ExtraEips as a int slice
func (p Params) EIPs() []int {
	eips := make([]int, len(p.ExtraEips.ExtraEIPs))
	for i, eip := range p.ExtraEips.ExtraEIPs {
		eips[i] = int(eip)
	}
	return eips
}

func validateEVMDenom(i interface{}) error {
	denom, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter EVM denom type: %T", i)
	}

	return sdk.ValidateDenom(denom)
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateEIPs(i interface{}) error {
	eips, ok := i.([]int64)
	if !ok {
		return fmt.Errorf("invalid EIP slice type: %T", i)
	}

	for _, eip := range eips {
		if !vm.ValidEip(int(eip)) {
			return fmt.Errorf("EIP %d is not activateable, valid EIPS are: %s", eip, vm.ActivateableEips())
		}
	}

	return nil
}

func validateChainConfig(i interface{}) error {
	cfg, ok := i.(ChainConfig)
	if !ok {
		return fmt.Errorf("invalid chain config type: %T", i)
	}

	return cfg.Validate()
}

// IsLondon returns if london hardfork is enabled.
func IsLondon(ethConfig *params.ChainConfig, height int64) bool {
	return ethConfig.IsLondon(big.NewInt(height))
}

// RegisterLegacyAminoCodec registers the legacy amino codec for MsgEthereumTx
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgEthereumTx{}, MsgEthTxAmino, nil)
}
