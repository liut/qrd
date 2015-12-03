QRD
====

    qrcode generate with fastcgi service


## build

    make

## change nginx host config

````

	location /qr {
		fastcgi_pass   127.0.0.1:9001;
		include		fastcgi_params;
	}


````

## test

browse: http://localhost/qr?c=Hello%20QRcode
