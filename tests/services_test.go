package tests

import (
	"testing"
	"abgo/service"
)

func TestReadUrls(t *testing.T){
	args := &service.Flags{
		UrlFile:"urls.txt"	}
	d := &service.Dispatcher{
		Args:args}

	val1 := d.ReadUrlFromFile();
	val2 := d.ReadUrlFromFile();

	if( val1 == val2){
		t.Errorf("Url values equals  %s == %s", val1, val2)
	}

	val3:= d.ReadUrlFromFile();
	if( val1 !=  val3){
		t.Errorf("Url values must be equals  %s == %s", val1, val3)
	}
}
