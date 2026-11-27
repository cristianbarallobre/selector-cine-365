package main
import "testing"
func TestAppChannelIsolation(t *testing.T){
 if appIsLab&&appDataDirectoryName!="Selector_Cine_365_LAB"{t.Fatal("LAB no está aislada")}
 if !appIsLab&&appDataDirectoryName!="Selector_Cine_365"{t.Fatal("producción ha cambiado de directorio")}
}
