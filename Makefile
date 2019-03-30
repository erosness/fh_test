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

package-deb-doc-tp:
	@echo "Packaging application as Thingsplex debian package"
	chmod a+x package/debian_tp/DEBIAN/*
	cp ./src/thingsplex_service_template package/debian_tp/opt/thingsplex/thingsplex_service_template
	cp VERSION package/debian_tp/opt/thingsplex/thingsplex_service_template
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian_tp
	@echo "Done"

package-deb-doc-fh:
	@echo "Packaging application as Futurehome debian package"
	chmod a+x package/debian_fh/DEBIAN/*
	cp ./src/thingsplex_service_template package/debian_fh/usr/bin/thingsplex_service_template
	cp VERSION package/debian_fh/var/lib/futurehome/thingsplex_service_template
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian_fh
	@echo "Done"


tar-arm: build-js build-go-arm package-deb-doc-2
	@echo "The application was packaged into tar archive "

deb-arm-fh : clean configure-arm build-go-arm package-deb-doc-fh
	mv package/debian_fh.deb package/build/thingsplex_service_template_$(version)_armhf.deb

deb-arm-tp : clean configure-arm build-go-arm package-deb-doc-tp
	mv package/debian_tp.deb package/build/thingsplex_service_template_$(version)_armhf.deb

deb-amd : configure-amd64 build-go-amd package-deb-doc-tp
	mv debian.deb thingsplex_service_template_$(version)_amd64.deb

run :
	go run src/service.go -c testdata/var/config.json


.phony : clean
