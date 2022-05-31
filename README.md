# SolidityCode.Reviews
Smart Contract Analyzer using Docker, Go, Slither and the Solidity Compiler.      
https://SolidityCode.Reviews   

Docker Image        
https://hub.docker.com/repository/docker/karlpothast/slither-web-interface

---

### 1. Navigate to https://SolidityCode.Reviews

![1](https://raw.githubusercontent.com/karlpothast/SolidityCode.Reviews/master/documentation/SolidityCodeReviewsMainPage.png)

### 2. Upload a solidity file (.sol file)

> For this example, I copied the solidity contract code directly from the [Etherscan page] for the Parity Multi-Sig Wallet app that was [hacked in November of 2017].  I placed it into a .sol file that I've included in this repository for others to use. You can download the sample solidity file [here].
>
> ![2](https://raw.githubusercontent.com/karlpothast/SolidityCode.Reviews/master/documentation/solidityCodeDirectFromEtherscan.png)

[here]: https://raw.githubusercontent.com/karlpothast/SolidityCode.Reviews/master/test-contracts/parity-wallet-hack.sol
[hacked in November of 2017]: https://www.coindesk.com/markets/2017/07/19/30-million-ether-reported-stolen-due-to-parity-wallet-breach/
[Etherscan page]: https://etherscan.io/address/0x863df6bfa4469f3ead0be8f9f2aae51c91a907b4#code

### 3. Slither Scan Results
> __The solidity compiler (solc) will change dynamically based on the file being analyzed__

As you can see below the Parity Wallet code analysis did not fare too well.  90 total issues were found although many were informational or recommendations.

![3](https://raw.githubusercontent.com/karlpothast/SolidityCode.Reviews/master/documentation/parityWalletCodeScanResults.png)
















