while
     qcli help &> /dev/null
     rc=$?; if [[ $rc == 0 ]]; then break; fi
do :;  done

balance=`qcli getbalance`
if [ "${balance:0:1}" == "0" ]
then
	qcli generate 600
fi

WALLETFILE=test-wallet
LOCKFILE=${QTUM_DATADIR}/import-test-wallet.lock

if [ ! -e $LOCKFILE ]; then
  while
       qcli getaddressesbyaccount "" &> /dev/null
       rc=$?; if [[ $rc != 0 ]]; then continue; fi

       qcli importwallet "${WALLETFILE}"
       solar prefund qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW 500
       solar prefund qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf 500
       touch $LOCKFILE
       break
  do :;  done
fi
