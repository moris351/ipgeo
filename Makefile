# This is how we want to name the binary output
base=${shell basename `pwd`}
#base=`basename `pwd``
OUTPUT=${base}
# These are the values we want to pass for Version and BuildTime
GITTAG=`git describe --tags`
BUILD_TIME=`date +%FT%T%z`
# Setup the -ldflags option for go build here, interpolate the variable values

VERTAG=
ifdef vertag
	ifeq (${vertag},fromgit)
		VERTAG=${GITTAG}
	else 
		VERTAG=${vertag}
	endif
endif

lDFLAGS=-ldflags "-X main.VerTag=${VERTAG} -X main.BuildTime=${BUILD_TIME}"
all:
	go build ${LDFLAGS} ./...;\
	cd cmd;\
	go build ${lDFLAGS} -o ipgeo .;\
	cd -;
run:
	cd cmd;\
	./ipgeo;\
	cd -;\
test:
	go test
clean:
	go clean
