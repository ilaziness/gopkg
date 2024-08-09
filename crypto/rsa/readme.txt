PFX证书转换成PEM格式证书
FPX 使用工具转换会包含CA证书
#PFX转换成crt证书
openssl pkcs12 -in ./server.pfx -clcerts -nokeys -out ./server.crt
1、-clcerts：仅仅输出客户端证书，不输出CA证书。
2、-nokeys：不输出任何私钥信息值。


#PFX 转换 密钥
openssl pkcs12 -in ./server.pfx -nocerts -nodes -out ./server.key
1、-nocerts：不输出任何证书。
2、-nodes：一直对私钥不加密。


如果报错加上-legacy参数，报错下面:
Error outputting keys and certificates
4077214A00710000:error:0308010C:digital envelope routines:inner_evp_generic_fetch:unsupported....................