# ibc-client-go

Grug light client for Cosmos SDK chains, written in Go.

This is basically a wrapper over [ibc-go](https://github.com/cosmos/ibc-go)'s [07-tendermint](https://github.com/cosmos/ibc-go/tree/main/modules/light-clients/07-tendermint) client, except for the client state's `VerifyMembership` and `VerifyNonMembership` methods are substituted with [Grug's proof format](https://github.com/left-curve/grug/blob/main/crates/jellyfish-merkle/src/proof.rs#L29-L33), which is not compatible with [ICS-23](https://github.com/cosmos/ics23).

## License

TBD
