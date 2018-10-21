#!/usr/bin/env bash

# set -x

ELF_NAME=$(grep 'ELF_NAME =' Makefile | awk '{print $3}')
ELF_BUILD_ENV=${3:-"elfBuildEnvMissing"}
ELF_APPENVIRONMENT=${3:-"elfAppEnvMissing"}
ELF_VERSION=${4:-"elfVersionMissing"}
OUTPUT_DIR=${5:-"elfOutputDirMissing"}
BUILDRUNTIME=$(go version | awk '{print $3}')
BUILDUSER="$USER"'@'$(hostname)
BUILDDATE=$(date +%Y%m%d-%H:%M:%S)

if [[ "${ELF_APPENVIRONMENT}" = "dev" ]];
then
  APPVERSIONPRERELEASE="dev"
fi

GITCOMMIT=${GIT_COMMIT:-$(git rev-parse --short HEAD)}
GITBRANCH=$(git rev-parse --abbrev-ref HEAD)
GITCOMMIT_AND_DIRTY=""
if [[ -n "`git status --porcelain`" ]]
then
  GIT_DIRTY="+UNCOMMITEDCHANGES"
else
  GIT_DIRTY=""
fi

# echo "---"
# echo $ELF_NAME
# echo $ELF_BUILD_ENV
# echo $ELF_APPENVIRONMENT
# echo $ELF_VERSION
# echo $OUTPUT_DIR
# echo $APPVERSIONPRERELEASE

# exit 12

function resolve_dependencies(){
  dep ensure
}
function build(){

  if [[ $GIT_DIRTY != "" ]]
  then
    GITCOMMIT_AND_DIRTY=${GITCOMMIT}${GIT_DIRTY}
    LDFLAGS="
            -X github.com/stefancocora/${ELF_NAME}/internal/version.GitCommit=${GITCOMMIT_AND_DIRTY} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Gitbranch=${GITBRANCH} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Buildruntime=${BUILDRUNTIME} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Gitbuilduser=${BUILDUSER} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Gitbuilddate=${BUILDDATE} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.AppEnvironment=${ELF_APPENVIRONMENT} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.VersionPrerelease=${APPVERSIONPRERELEASE} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Version=${ELF_VERSION}
    "
  else
    GITCOMMIT_AND_DIRTY=${GITCOMMIT}
    LDFLAGS="
            -X github.com/stefancocora/${ELF_NAME}/internal/version.GitCommit=${GITCOMMIT_AND_DIRTY} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Gitbranch=${GITBRANCH} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Buildruntime=${BUILDRUNTIME} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Gitbuilduser=${BUILDUSER} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Gitbuilddate=${BUILDDATE} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.AppEnvironment=${ELF_APPENVIRONMENT} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.VersionPrerelease=${APPVERSIONPRERELEASE} \
            -X github.com/stefancocora/${ELF_NAME}/internal/version.Version=${ELF_VERSION}
    "
  fi

  if [[ "${ELF_APPENVIRONMENT}" = "dev" ]];
  then
    ELF_VERSIONED="${ELF_NAME}-${ELF_APPENVIRONMENT}-${ELF_VERSION}-${GITCOMMIT_AND_DIRTY}"
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "./${OUTPUT_DIR}/${ELF_VERSIONED}" -ldflags "${LDFLAGS}" cmd/"${ELF_NAME}"/main.go
  elif [[ "${ELF_APPENVIRONMENT}" = "production" ]] && [[ "${GIT_DIRTY}" = "" ]] ;
  then
    ELF_VERSIONED="${ELF_NAME}-${ELF_VERSION}-${GITCOMMIT_AND_DIRTY}"
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "./${OUTPUT_DIR}/${ELF_VERSIONED}" -ldflags "${LDFLAGS}" cmd/"${ELF_NAME}"/main.go
  elif [[ "${ELF_APPENVIRONMENT}" = "production" ]] && [[ "${GIT_DIRTY}" != "" ]] ;
  then
    date
    echo "=== exception when building code: won't build ${ELF_APPENVIRONMENT} binary with a dirty git tree"
    exit 1
  fi

  if [[ $? -ne 0 ]];
  then
    printf "\n=== build: something went wrong when go build was called\n"
    exit 1
  fi


  strip "${OUTPUT_DIR}/${ELF_VERSIONED}"

  printf "\n=== build: info about the built binary\n"
  echo ""
  ls -lha "${OUTPUT_DIR}/${ELF_VERSIONED}"
  file "${OUTPUT_DIR}/${ELF_VERSIONED}"
  ldd "${OUTPUT_DIR}/${ELF_VERSIONED}"
  printf "md5sum:    "
  md5sum "${OUTPUT_DIR}/${ELF_VERSIONED}"
  printf "sha256sum: "
  sha256sum "${OUTPUT_DIR}/${ELF_VERSIONED}"
  printf "sha512sum: "
  sha512sum "${OUTPUT_DIR}/${ELF_VERSIONED}"
}

resolve_dependencies
build ${OUTPUT_DIR}
