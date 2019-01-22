# Basic building blocks of a centralized exchange

 - Wallets (user side)
  - keep track of funds / keys

 - Orders
  - Order types, database

 - Matching engine
  - Taking orders from above, keep track of sides, always be matching

 - Accounts

## Wallets

Of course users have their own wallets, but there needs to be some way to make a deposit to a certain account.
I'm guessing the exchange makes an address when you say you want to deposit, or assigns an address to your account. Assigning the address to the account seems like the best idea, because then you can deposit anytime. The private key for account address could be a child private key (BIP32 style maybe) of a private key the exchange owns.

Basically money always comes in to the exchange, the exchange keeps track of your balance (and holds the funds in their wallets so they could pull a mtgox)

## Interacting with the exchange
The exchange needs to have a few functions that the user interacts with:

Un-permissioned commands (and simple mockups of how it might work):

 - Register account
`ocx register username password`

 - Log in (need a way to keep sessions or something)
`ocx login username password`

 - View orderbook
`ocx vieworderbook`
Or maybe we want dark pools? Could be a feature, probably out of scope, would be difficult to do if you want to match in a decentralized way. Decentralized matching for confidential _orders_ is probably very difficult, aside from the whole decentralized matching problem.
 - Get price (really just getorderbook but with a few more operations)
`ocx getprice`

- Get volume (need to track that server side)
`ocx getvolume`

 - TODO: think of more that you might need

Permissioned commands:

 - Place order
`ocx placeorder price {buy|sell} assetwant assethave amounttobuy`
This will print a description of the order after making it, and prompt the user before actually sending it.

 - Get account's address
`ocx getaddress`
This will return the address that is assigned to the user's account

 - Withdraw
`ocx withdrawtoaddress asset amount recvaddress`
Withdraw will send a transaction to the blockchain.

 - Delete account
`ocx deleteaccount`

For authentication, let's just do some user data storage and send a random token that expires in 30 minutes or something. Server checks token, client stores token and sends it with json.