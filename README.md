mlcloud project is used to provided machine learning framework for users include but not limit to tensorflow, MXNet, Caffe
============================

## development details
the project was based on golang. the dependency was managed by dep
you can find it in: https://github.com/golang/dep

## compile the project
```sh
# step 1: download the dependencies
# the packages declare in Gopkg.toml file. except that, only you import that
# pkg in your code. then dep will download the packages.
dep ensure

# step2: build the mlcloud server
cd src && go build -o mlcloud
```

## deploy the project
the project can be deployed on
* baremetal server
* docker container
* kubernetes cluster

see the deploy details under /deploy

## test pull request
## test1
