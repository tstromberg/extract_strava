extract_strava
==============

Tool for exporting TCX files from Strava's Android database. This is useful if the Strava app fails to upload a ride for whatever reason.

It is worth noting that Strava's Android database only contains waypoints that have not been successfully uploaded yet.

Dumping your database:
======================
Assuming you already have root access to your Android device and the "adb" command installed:

// % adb shell
// % > su
// % > chmod 775 /data/data; chmod -R 775 /data/data/com.strava
//
// % mkdir com.strava && cd com.strava
// % adb pull /data/data/com.strava
// % adb shell
// % > su
// % > chmod 771 /data/data; chmod -R 771 /data/data/com.strava


Usage:
======

Once extracted:

% extract_strava ./databases/strava .


Using on other platforms:
=========================
extract_strava bundles the go-sqlite3 library, precompiled for Mac OS X. If you need it for other platforms,

% export GOPATH=/path/to/extract_strava_directory
% go get github.com/mattn/go-sqlite3

