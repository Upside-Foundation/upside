# Integration tests

The Python integration suite starts development networks and exercises the daemon and contracts. Run it on a development machine with disposable state. Use separate homes and ports from any operational node.

## Prerequisites

The test runner uses Nix, Node.js, npm, Python, pytest, and the contract toolchain. The environment is defined in [integration_tests/shell.nix](../integration_tests/shell.nix), and the runner is [scripts/run-integration-tests](../scripts/run-integration-tests).

The Nix build configuration is separate from the pinned Docker release recipe. A successful test build does not establish that its binary matches a release checksum.

## Run the suite

From the repository root:

```bash
make run-integration-tests
```

The target runs three steps:

1. `make gen-bindings-contracts` generates the contract bindings inside a Nix shell.
2. [scripts/run-integration-tests](../scripts/run-integration-tests) runs `npm install` and `npm run typechain` in `integration_tests/contracts`. These npm commands run directly on your machine, outside Nix, and npm runs the install scripts of the packages it downloads.
3. The same script runs pytest inside the Nix shell defined in [integration_tests/shell.nix](../integration_tests/shell.nix).

The default run currently fails. It includes `integration_tests/test_upgrade.py`, which builds `integration_tests/configs/upgrade-test-package.nix`. That file downloads earlier releases from a GitHub repository that doesn't exist, and the test walks upgrade plans that the app no longer registers. Run the tests that have no marker instead:

```bash
TESTS_TO_RUN=unmarked make run-integration-tests
```

To select other pytest markers:

```bash
TESTS_TO_RUN=MARKER make run-integration-tests
```

Replace `MARKER` with a marker used by the suite.

## Run selected tests

First prepare the contract bindings and test dependencies. The npm commands run on your machine, as described in [Run the suite](#run-the-suite):

```bash
make gen-bindings-contracts
cd integration_tests/contracts
npm install
npm run typechain
cd ../..
```

Then enter the environment:

```bash
nix-shell integration_tests/shell.nix
pytest --collect-only -q integration_tests
pytest -s -vv integration_tests -k test_basic
```

Use the collected names and markers to select the desired test cases. Some cases also start a geth instance.

For a manual daemon smoke test with fresh keys and a separate home, follow [Run a localnet](nodes/localnet.md).
