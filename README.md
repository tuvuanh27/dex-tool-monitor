Crawler transactions and pairs of DEXes
=======================================
- To add new chain:
  - please add `rpc` and `stable_coin` to `config.json` respectively.
  - in file `src/type/chain.go`, add new chain to enum `Chain` and `func GetChain(chain string) Chain`.


- To add new DEX, add data to collection `dexes` in database
