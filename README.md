extract_strava
==============

Tool for exporting TCX files from Strava's Android database. This is useful if the Strava app fails to upload a ride for whatever reason.

It is worth noting that Strava's Android database only contains waypoints that have not been successfully uploaded yet.

Dumping your database:
======================
Assuming you already have root access to your Android device and the "adb" command installed:

    % mkdir com.strava && cd com.strava
    % adb shell
    % > su
    % > chmod 775 /data/data; chmod -R 775 /data/data/com.strava
    % adb pull /data/data/com.strava
    % adb shell
    % > su
    % > chmod 771 /data/data; chmod -R 771 /data/data/com.strava


Pre-built binaries:
===================

Look at the bin/ directory. It contains pre-built 64-bit binaries for Mac OS X and Linux.


Building extract_strava:
========================

This will fetch all dependencies and compile the extract_strava binary.

    % export GOPATH=.
    % go get github.com/jmoiron/sqlx
    % go get github.com/mattn/go-sqlite3
    % go build extract_strava.go


Usage:
======

Once the tool has been built and the database extracted:

    % ./extract_strava ./databases/strava .


