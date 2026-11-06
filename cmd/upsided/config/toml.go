package config

// DefaultUpsideConfigTemplate defines the configuration template for upside configuration
const DefaultUpsideConfigTemplate = `
###############################################################################
###                             Upside Configuration                       ###
###############################################################################

[upside]

# Set to true to disable tx replacement.
disable-tx-replacement = {{ .Upside.DisableTxReplacement }}

# Set to true to disable optimistic execution (not recommended on validator nodes).
disable-optimistic-execution = {{ .Upside.DisableOptimisticExecution }}
`
