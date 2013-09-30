extract_strava
==============
Tool for exporting TCX files from Strava's Android database. This is useful if
the Strava app fails to upload a ride, as it does not give you an option to
retry an upload if the server rejects it.

Dumping your database:
======================
Assuming you already have root access to your Android device and the "adb" command installed:

    % mkdir com.strava && cd com.strava
    % adb shell
    % > su
    % > chmod 775 /data/data; chmod -R 775 /data/data/com.strava
    % adb pull /data/data/com.strava

This will create a file called "databases/strava" among other things. Once the
data has been pulled, return the permissions back to normal.

    % adb shell
    % > su
    % > chmod 771 /data/data; chmod -R 771 /data/data/com.strava


Pre-built binaries:
===================
Look at the bin/ directory. There are pre-built 64-bit binaries for Mac OS X and Linux.

Note that extract_strava looks for "tcx.tmpl" in the current working directory.


Building extract_strava:
========================
Go 1.1 is required for building. Download from https://code.google.com/p/go/downloads/list

These commands will fetch all dependencies and compile the extract_strava binary.

    % export GOPATH=.
    % go get github.com/jmoiron/sqlx
    % go get github.com/mattn/go-sqlite3
    % go build extract_strava.go


Usage:
======
Once the tool has been built and the database extracted, you may run the tool:

    % ./extract_strava ./databases/strava .

extract_strava will dump out one TCX file per ride found in the database. These
files may be uploaded to Strava via the web interface. Thankfully, Strava will
detect and notify you if a ride has already been submitted.
