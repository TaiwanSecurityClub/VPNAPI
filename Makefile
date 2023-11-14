builder := go
builddir := bin
exe := wgapi
path := /usr/local/bin
instdir := /usr/local/share/wgapi
systemddir := /etc/systemd/system
config := .env.sample
systemd := wgapi.service

all: $(builddir)/$(exe)

$(builddir)/$(exe): main.go go.mod go.sum models router utils
		$(builder) build -o $(builddir)/$(exe) $<

install: $(instdir)/$(exe) $(instdir)/.env $(systemddir)/$(systemd)

$(instdir): 
		mkdir $(instdir)

$(instdir)/$(exe): $(instdir) $(builddir)/$(exe)
		cp $(builddir)/$(exe) $(instdir)/$(exe)
		chown root:root $(instdir)/$(exe)
		chmod 4755 $(instdir)/$(exe)

$(instdir)/.env: $(instdir) $(builddir)/$(config)
		cp $(builddir)/$(config) $(instdir)/.env

$(systemddir)/$(systemd): $(builddir)/$(systemd)
		cp $(builddir)/$(systemd) $(systemddir)/$(systemd)

uninstall:
		rm -rf $(instdir)
		rm -rf $(systemddir)/$(systemd)
clean: 
		rm -rf $(builddir)
