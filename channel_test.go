package main

import "testing"

func TestAppChannelIsolation(t *testing.T) {
 if appIsLab {
  if appDataDirectoryName != "Selector_Cine_365_LAB" || oneDriveFolderName != "Selector_Cine_365_LAB" || cloudDirectoryName != "Selector_Cine_365_LAB_Copias" { t.Fatal("LAB no está aislada") }
 } else {
  if appDataDirectoryName != "Selector_Cine_365" || oneDriveFolderName != "Selector_Cine_365" || cloudDirectoryName != "Selector_Cine_365_Copias" { t.Fatal("producción ha cambiado de directorios") }
 }
 if appDataDirectoryName == oneDriveFolderName && appIsLab && appDataDirectoryName != "Selector_Cine_365_LAB" { t.Fatal("rutas de LAB incorrectas") }
}
