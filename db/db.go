package main

import (
	"fmt"
	"os"
	"runtime"
)

const perlTemplate = `#!/usr/bin/perl
use strict;
use warnings;
use DBI;

# -------------------------------------------------------------------------
# Load DB config
# -------------------------------------------------------------------------
my $PATHconf = '/etc/astguiclient.conf';
open(my $conf_fh, '<', $PATHconf) or die "Can't open $PATHconf: $!\n";
my %conf;
while (my $line = <$conf_fh>) {
    $line =~ s/[#;].*$//;
    $line =~ s/^\s+|\s+$//g;
    next unless $line =~ /=/;
    my ($key, $val) = split /=/, $line, 2;
    $conf{$key} = $val if defined $val && $val ne '';
}
close $conf_fh;

my $VARDB_server   = $conf{'VARDB_server'}   // 'localhost';
my $VARDB_database = $conf{'VARDB_database'} // 'asterisk';
my $VARDB_user     = $conf{'VARDB_user'}     // '';
my $VARDB_pass     = $conf{'VARDB_pass'}     // '';
my $VARDB_port     = $conf{'VARDB_port'}     // 3306;

# -------------------------------------------------------------------------
# DB Connection
# -------------------------------------------------------------------------
my $dbh = DBI->connect(
    "DBI:mysql:$VARDB_database:$VARDB_server:$VARDB_port",
    $VARDB_user,
    $VARDB_pass,
    { RaiseError => 1, AutoCommit => 1, mysql_enable_utf8 => 1 }
) or die "Couldn't connect to database: " . DBI->errstr;
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: db <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Write the file (Windows ignores perms, Unix honors them)
	err := os.WriteFile(filename, []byte(perlTemplate), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	// If Unix-like, apply chmod +x
	if runtime.GOOS != "windows" {
		err = os.Chmod(filename, 0755)
		if err != nil {
			fmt.Printf("Warning: couldn't set executable bit: %v\n", err)
		}
	}

	fmt.Printf("Created %s\n", filename)
}
