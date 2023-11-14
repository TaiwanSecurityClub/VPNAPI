builder := go
builddir := bin
exe := wgapi
path := /usr/local/bin
instdir := /usr/local/share/wgapi
instrelativedir := ../share/wgapi
systemddir := /etc/systemd/system
config := .env.sample
systemd := wgapi.service
tags := release

all: $(builddir)/$(exe)

$(builddir)/$(exe): main.go go.mod go.sum models router utils
		$(builder) build -o $(builddir)/$(exe) -tags $(tags) $<

install: $(path)/$(exe) $(systemddir)/$(systemd)

$(path)/$(exe): $(instdir)/$(exe) $(instdir)/.env
		ln -s $(instrelativedir)/$(exe) $(path)/$(exe)

$(instdir): 
		mkdir $(instdir)

$(instdir)/$(exe): $(instdir) $(builddir)/$(exe)
		cp $(builddir)/$(exe) $(instdir)/$(exe)
		chown root:root $(instdir)/$(exe)
		chmod 4755 $(instdir)/$(exe)

$(instdir)/.env: $(instdir) $(config)
		cp $(config) $(instdir)/.env

$(systemddir)/$(systemd): $(systemd)
		cp $(systemd) $(systemddir)/$(systemd)

uninstall:
		rm -rf $(path)/$(exe)
		rm -rf $(instdir)
		rm -rf $(systemddir)/$(systemd)
clean: 
		rm -rf $(builddir)
