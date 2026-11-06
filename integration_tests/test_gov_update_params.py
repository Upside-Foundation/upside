import pytest
from web3 import exceptions

from .cosmoscli import module_address
from .utils import CONTRACTS, deploy_contract, submit_gov_proposal

pytestmark = pytest.mark.gov


def test_evm_update_param(upside, tmp_path):
    contract = deploy_contract(
        upside.w3,
        CONTRACTS["Random"],
    )
    res = contract.caller.randomTokenId()
    assert res > 0, res
    cli = upside.cosmos_cli()
    p = cli.query_params("evm")
    del p["chain_config"]["merge_netsplit_block"]
    del p["chain_config"]["shanghai_time"]
    del p["chain_config"]["cancun_time"]
    del p["chain_config"]["prague_time"]
    authority = module_address("gov")
    msg = "/ethermint.evm.v1.MsgUpdateParams"
    submit_gov_proposal(
        upside,
        msg,
        messages=[
            {
                "@type": msg,
                "authority": authority,
                "params": p,
            }
        ],
    )
    p = cli.query_params("evm")
    assert not p["chain_config"]["merge_netsplit_block"]
    assert not p["chain_config"]["shanghai_time"]
    assert not p["chain_config"]["cancun_time"]
    assert not p["chain_config"]["prague_time"]
    invalid_msg = "invalid opcode: PUSH0"
    with pytest.raises(exceptions.Web3RPCError) as e_info:
        contract.caller.randomTokenId()
    assert invalid_msg in str(e_info.value)
    with pytest.raises(exceptions.Web3RPCError) as e_info:
        deploy_contract(upside.w3, CONTRACTS["Greeter"])
    assert invalid_msg in str(e_info.value)


def test_gov_update_params(upside):
    params = {
        "upside_admin": "up12luku6uxehhak02py4rcz65zu0swh7wjngq88t",
        "enable_auto_deployment": False,
        "ibc_upside_denom": "ibc/6411AE2ADA1E73DB59DB151"
        "A8988F9B7D5E7E233D8414DB6817F8F1A01600000",
        "ibc_timeout": "96400000000000",
        "max_callback_gas": "400000",
        # Required: the whole param set is replaced, and a half-life under
        # 1 year (including a missing one) is rejected when the proposal runs.
        "mint_half_life_seconds": "63115200",
        # Must equal the stored value; only MsgSetMintPaused changes it.
        "mint_paused": False,
    }
    msg = "/upside.MsgUpdateParams"
    authority = module_address("gov")
    submit_gov_proposal(
        upside,
        msg,
        messages=[
            {
                "@type": msg,
                "authority": authority,
                "params": params,
            }
        ],
    )
    assert upside.cosmos_cli().query_params() == params
