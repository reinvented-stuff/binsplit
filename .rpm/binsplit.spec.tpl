Name:           binsplit
Version:        __VERSION__
Release:        1%{?dist}
Summary:        Splits a binary file into chunks by a boundary sequence

License:        MIT
URL:            https://reinvented-stuff.com/
Source0:        __SOURCE_TARGZ_FILENAME__


%description
binsplit looks up a hex sequence in a binary file and splits the file
into parts.

%prep
%setup -q


%build
make %{?_smp_mflags} build


%install
rm -rf $RPM_BUILD_ROOT
%make_install


%files
%attr(755, root, root) /usr/bin/binsplit
%attr(644, root, root) /usr/share/doc/binsplit-__VERSION__/README.md
%doc


%changelog
