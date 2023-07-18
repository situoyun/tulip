#!/bin/bash
  
rm -rf /nfs/data/tuorong/1/tulip
kubectl delete job tulipdownload -n tuorong
kubectl apply -f /nfs/data/tuorong/1/rollup/tulip_download.yaml

sleep 1s

tar -xzvf /nfs/data/tuorong/1/rollup/mysql.tar.gz -C /nfs/data/tuorong/1/rollup
mv -f /nfs/data/tuorong/1/rollup/test/mysql /nfs/data/tuorong/1/rollup/test/mysql$(date +%Y-%m-%d)
mv -f /nfs/data/tuorong/1/rollup/mysql /nfs/data/tuorong/1/rollup/test/mysql
kubectl rollout restart deployment mysqltest -n tuorong

kubectl delete job tulipupload -n tuorong
kubectl apply -f /nfs/data/tuorong/1/rollup/tulip_upload.yaml