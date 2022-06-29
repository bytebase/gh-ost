/*
   Copyright 2022 GitHub Inc.
	 See https://github.com/github/gh-ost/blob/master/LICENSE
*/

package base

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	test "github.com/openark/golib/tests"
)

func TestGetTableNames(t *testing.T) {
	{
		context := NewMigrationContext()
		context.OriginalTableName = "some_table"
		test.S(t).ExpectEquals(context.GetOldTableName(), "~some_table_del")
		test.S(t).ExpectEquals(context.GetGhostTableName(), "~some_table_gho")
		test.S(t).ExpectEquals(context.GetChangelogTableName(), "~some_table_ghc")
	}
	{
		context := NewMigrationContext()
		context.OriginalTableName = "a123456789012345678901234567890123456789012345678901234567890"
		test.S(t).ExpectEquals(context.GetOldTableName(), "~a1234567890123456789012345678901234567890123456789012345678_del")
		test.S(t).ExpectEquals(context.GetGhostTableName(), "~a1234567890123456789012345678901234567890123456789012345678_gho")
		test.S(t).ExpectEquals(context.GetChangelogTableName(), "~a1234567890123456789012345678901234567890123456789012345678_ghc")
	}
	{
		context := NewMigrationContext()
		context.OriginalTableName = "a123456789012345678901234567890123456789012345678901234567890123"
		oldTableName := context.GetOldTableName()
		test.S(t).ExpectEquals(oldTableName, "~a1234567890123456789012345678901234567890123456789012345678_del")
	}
	{
		context := NewMigrationContext()
		context.OriginalTableName = "a123456789012345678901234567890123456789012345678901234567890123"
		context.TimestampOldTable = true
		longForm := "Jan 2, 2006 at 3:04pm (MST)"
		context.StartTime, _ = time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
		oldTableName := context.GetOldTableName()
		test.S(t).ExpectEquals(oldTableName, "~a12345678901234567890123456789012345678901234567_1359921240_del")
	}
	{
		context := NewMigrationContext()
		context.OriginalTableName = "a123456789012345678901234567890123456789012345678901234567890123"
		context.TimestampAllTable = true
		longForm := "Jan 2, 2006 at 3:04pm (MST)"
		context.StartTime, _ = time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
		test.S(t).ExpectEquals(context.GetGhostTableName(), "~a12345678901234567890123456789012345678901234567_1359921240_gho")
	}
	{
		context := NewMigrationContext()
		context.OriginalTableName = "foo_bar_baz"
		context.ForceTmpTableName = "tmp"
		test.S(t).ExpectEquals(context.GetOldTableName(), "~tmp_del")
		test.S(t).ExpectEquals(context.GetGhostTableName(), "~tmp_gho")
		test.S(t).ExpectEquals(context.GetChangelogTableName(), "~tmp_ghc")
	}
}

func TestReadConfigFile(t *testing.T) {
	{
		context := NewMigrationContext()
		context.ConfigFile = "/does/not/exist"
		if err := context.ReadConfigFile(); err == nil {
			t.Fatal("Expected .ReadConfigFile() to return an error, got nil")
		}
	}
	{
		f, err := ioutil.TempFile("", t.Name())
		if err != nil {
			t.Fatalf("Failed to create tmp file: %v", err)
		}
		defer os.Remove(f.Name())

		f.Write([]byte("[client]"))
		context := NewMigrationContext()
		context.ConfigFile = f.Name()
		if err := context.ReadConfigFile(); err != nil {
			t.Fatalf(".ReadConfigFile() failed: %v", err)
		}
	}
	{
		f, err := ioutil.TempFile("", t.Name())
		if err != nil {
			t.Fatalf("Failed to create tmp file: %v", err)
		}
		defer os.Remove(f.Name())

		f.Write([]byte("[client]\nuser=test\npassword=123456"))
		context := NewMigrationContext()
		context.ConfigFile = f.Name()
		if err := context.ReadConfigFile(); err != nil {
			t.Fatalf(".ReadConfigFile() failed: %v", err)
		}

		if context.config.Client.User != "test" {
			t.Fatalf("Expected client user %q, got %q", "test", context.config.Client.User)
		} else if context.config.Client.Password != "123456" {
			t.Fatalf("Expected client password %q, got %q", "123456", context.config.Client.Password)
		}
	}
	{
		f, err := ioutil.TempFile("", t.Name())
		if err != nil {
			t.Fatalf("Failed to create tmp file: %v", err)
		}
		defer os.Remove(f.Name())

		f.Write([]byte("[osc]\nmax_load=10"))
		context := NewMigrationContext()
		context.ConfigFile = f.Name()
		if err := context.ReadConfigFile(); err != nil {
			t.Fatalf(".ReadConfigFile() failed: %v", err)
		}

		if context.config.Osc.Max_Load != "10" {
			t.Fatalf("Expected osc 'max_load' %q, got %q", "10", context.config.Osc.Max_Load)
		}
	}
}
