# 명령어 정리

go run main.go는 기본적으로 실행해야 하기 떄문에 붙여주어야 한다.


# option

option으로는 다음과 같은 항목이 존재

`getBalance -address <계정 이름>` : 계정의 밸런스를 가져 온다.


`createBlockChain -address <계정 이름>` : 계정을 생성한다.

최초 생성되는 계정은 100개의 토큰을 소유하고 있다.


`printChain` : 모든 원장 데이터를 가져 온다.

`send -from <계정> -to <계정> -amount` : 토큰을 전송한다.

`createWallet` : 지갑을 생성한다.
 
`listWallets` : 만들어져 있는 지갑 List를 가져온다.