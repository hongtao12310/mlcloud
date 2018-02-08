mlcloud support swagger annotation. we can generate the swagger spec by use the swagger command line
==========================

## generate the swagger spec
```sh
   swagger generate spec -m -o ./swagger.json
```

## run the swagger server
```sh
   swagger generate server -A mlcloud-swagger -f swagger.json
```

