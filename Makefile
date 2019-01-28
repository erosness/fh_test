version="0.10.6"
version_file=VERSION
working_dir=$(shell pwd)
arch="armhf"

clean:
	-rm tpflow

build-go:
	go build -o thingsplex_service_template src/service.go

build-go-arm:
	cd ./src;GOOS=linux GOARCH=arm GOARM=6 go build -o thingsplex_service_template service.go;cd ../

build-go-amd:
	cd ./src;GOOS=linux GOARCH=amd64 go build -o thingsplex_service_template src/service.go;cd ../


configure-arm:
	python ./scripts/config_env.py prod $(version) armhf

configure-amd64:
	python ./scripts/config_env.py prod $(version) amd64


package-tar:
	tar cvzf thingsplex_service_template_$(version).tar.gz thingsplex_service_template VERSION

package-deb-doc-1:
	@echo "Packaging application as debian package"
	chmod a+x package/debian/DEBIAN/*
	cp ./src/thingsplex_service_template package/debian/opt/thingsplex/thingsplex_service_template
	cp VERSION package/debian/opt/thingsplex/thingsplex_service_template
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian
	@echo "Done"

package-deb-doc-2:
	@echo "Packaging application as debian package"
	chmod a+x package/debian/DEBIAN/*
	cp ./src/thingsplex_service_template package/debian/usr/bin/thingsplex_service_template
	cp VERSION package/debian/var/lib/thingsplex/thingsplex_service_template
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian
	@echo "Done"


tar-arm: build-js build-go-arm package-deb-doc-2
	@echo "The application was packaged into tar archive "

deb-arm : clean configure-arm build-go-arm package-deb-doc-2
	mv package/debian.deb package/build/thingsplex_service_template_$(version)_armhf.deb

deb-amd : configure-amd64 build-go-amd package-deb-doc-2
	mv debian.deb thingsplex_service_template_$(version)_amd64.deb

run :
	go run src/service.go -c testdata/var/config.json


.phony : clean
