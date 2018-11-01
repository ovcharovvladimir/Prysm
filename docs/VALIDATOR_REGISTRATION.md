# Voter Registration & Starting Workflow
This doc summarizes the work flow of registering to become a voter in the ESS blockchain and starting supernode (sness)

### Step 1: Deploy voter registration contract 
To deploy voter registration contract, we can use [deployVRC](https://github.com/ovcharovvladimir/Prysm/tree/master/contracts/voter-registration-contract/deployVRC) utility.  
After successfully contract deploying we will get contract address
```
# Deploy contract with keystore UTCJSON and password
> ./deployVRC --UTCPath /path/to/your/keystore/UTCJSON --passwordFile /path/to/your/password.txt
# Deploy contract with private key
> ./deployVRC --privKey 8a6db3b30934439c9f71f1fa777019810fd538c9c1e396809bcf9fd5535e20ca

INFO[0039] New contract deployed at 0x559eDab2b5896C2Bc37951325666ed08CD41099d
```
### Step 2: Launch Supernode 
get 
Launch sness with supernode public key (pubkey) and the voter contract address (vrcaddr) we just deployed
```
./sness --enable-powchain --vrcaddr 0x559eDab2b5896C2Bc37951325666ed08CD41099d --pubkey 
0x6f1b9df48a267576d9c132468071eebbd56263d11f465567a45d0cf71cddeb67

```

### Step 3: Send a transaction to the deposit function with 32 ESS and supernode account holder's public key as argument

### Step 4: Wait for deposit transaction to mine.
After the deposit transaction gets mined, Supernode will report account holder has been registered. Congrats! 
```
INFO[0000] Starting Supernode
INFO[0000] Starting web3 PoW chain service at ws://127.0.0.1:8546
INFO[0152] Voter registered with public key: 0x6f1b9df48a267576d9c132468071eebbd56263d11f465567a45d0cf71cddeb67
```

